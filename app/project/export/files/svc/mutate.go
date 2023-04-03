package svc

import (
	"strings"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/project/export/model"
	"projectforge.dev/projectforge/app/util"
)

const serviceAssignmentToken = ":="

func ServiceMutate(m *model.Model, args *model.Args, addHeader bool) (*file.File, error) {
	g := golang.NewFile(m.Package, []string{"app", m.PackageWithGroup("")}, "servicemutate")
	for _, imp := range helper.ImportsForTypes("go", "", m.PKs().Types()...) {
		g.AddImport(imp)
	}
	g.AddImport(helper.ImpAppUtil, helper.ImpContext, helper.ImpSQLx, helper.ImpDatabase)

	if add, err := serviceCreate(m, g, args.Database()); err == nil {
		g.AddBlocks(add)
	} else {
		return nil, err
	}
	if upd, err := serviceUpdate(m, g, args.Database()); err == nil {
		g.AddBlocks(upd)
	} else {
		return nil, err
	}
	if m.IsRevision() || m.IsHistory() {
		if updIN, err := serviceUpdateIfNeeded(m, g, args.Database()); err == nil {
			g.AddBlocks(updIN)
		} else {
			return nil, err
		}
	}
	if save, err := serviceSave(m, g, args.Database()); err == nil {
		g.AddBlocks(save)
	} else {
		return nil, err
	}
	if m.IsRevision() {
		g.AddBlocks(serviceUpsertCore(m, g, args.Database()), serviceInsertRevision(m, args.Database()))
	}
	if m.IsSoftDelete() {
		g.AddImport(helper.ImpTime)
		sdel, err := serviceSoftDelete(m, args.Enums)
		if err != nil {
			return nil, err
		}
		g.AddBlocks(sdel, serviceSoftDeleteWhere(m), serviceAddDeletedClause(m))
	} else {
		del, err := serviceDelete(m, args.Enums)
		if err != nil {
			return nil, err
		}
		g.AddBlocks(del, serviceDeleteWhere(m))
	}
	return g.Render(addHeader)
}

func serviceCreate(m *model.Model, g *golang.File, database string) (*golang.Block, error) {
	ret := golang.NewBlock("Create", "func")
	ret.W("func (s *Service) Create(ctx context.Context, tx *sqlx.Tx, logger util.Logger, models ...*%s) error {", m.Proper())
	ret.W("\tif len(models) == 0 {")
	ret.W("\t\treturn nil")
	ret.W("\t}")

	if m.IsRevision() {
		revCol := m.HistoryColumn()
		ret.W("\trevs, err := s.getCurrent%s(ctx, tx, logger, models...)", revCol.ProperPlural())
		ret.W("\tif err != nil {")
		ret.W("\t\treturn err")
		ret.W("\t}")

		if err := serviceAddCreatedUpdated(m, ret, g, false); err != nil {
			return nil, err
		}

		ret.W("")
		ret.W("\terr = s.upsertCore(ctx, tx, logger, models...)")
		ret.W("\tif err != nil {")
		ret.W("\t\treturn err")
		ret.W("\t}")
		ret.W("\terr = s.insert%s(ctx, tx, logger, models...)", revCol.Proper())
		ret.W("\tif err != nil {")
		ret.W("\t\treturn err")
		ret.W("\t}")
		ret.W("\treturn nil")
	} else {
		if err := serviceAddCreatedUpdated(m, ret, g, false); err != nil {
			return nil, err
		}
		ret.W("\tq := database.SQLInsert(tableQuoted, columnsQuoted, len(models), s.db.Placeholder())")
		ret.W("\tvals := make([]any, 0, len(models)*len(columnsQuoted))")
		ret.W("\tfor _, arg := range models {")
		ret.W("\t\tvals = append(vals, arg.ToData()...)")
		ret.W("\t}")
		ret.W("\treturn s.db.Insert(ctx, q, tx, logger, vals...)")
	}
	ret.W("}")
	return ret, nil
}

func serviceUpdate(m *model.Model, g *golang.File, database string) (*golang.Block, error) {
	ret := golang.NewBlock("Update", "func")
	ret.W("func (s *Service) Update(ctx context.Context, tx *sqlx.Tx, model *%s, logger util.Logger) error {", m.Proper())
	if m.IsRevision() {
		revCol := m.HistoryColumn()
		ret.W("\trevs, err := s.getCurrent%s(ctx, tx, logger, model)", revCol.ProperPlural())
		ret.W("\tif err != nil {")
		ret.W("\t\treturn err")
		ret.W("\t}")
		ret.W("\tmodel.%s = revs[model.String()] + 1", revCol.Proper())
	}

	if cc := m.Columns.WithTag("created"); len(cc) > 0 {
		g.AddImport(helper.ImpErrors)
		ret.W("\tcurr, err := s.Get(ctx, tx, %s%s, logger)", m.PKs().ToRefs("model."), m.SoftDeleteSuffix())
		ret.W("\tif err != nil {")
		ret.W("\t\treturn errors.Wrapf(err, \"can't get original %s [%%%%s]\", model.String())", m.TitleLower())
		ret.W("\t}")
		for _, c := range cc {
			ret.W("\tmodel.%s = curr.%s", c.Proper(), c.Proper())
		}
	}

	for _, updated := range m.Columns.WithTag("updated") {
		err := serviceSetVal(updated, g, ret, 1)
		if err != nil {
			return nil, errors.Wrap(err, "")
		}
	}

	if m.IsHistory() {
		ret.W("")
		ret.W("\t_, hErr := s.SaveHistory(ctx, tx, curr, model, logger)")
		ret.W("\tif hErr != nil {")
		ret.W("\t\treturn errors.Wrap(hErr, \"unable to save history\")")
		ret.W("\t}")
	}

	pks := m.PKs()
	pkVals := make([]string, 0, len(pks))
	for _, pk := range pks {
		pkVals = append(pkVals, "model."+pk.Proper())
	}
	if m.IsRevision() {
		revCol := m.HistoryColumn()
		ret.W("")
		ret.W("\terr = s.upsertCore(ctx, tx, logger, model)")
		ret.W("\tif err != nil {")
		ret.W("\t\treturn err")
		ret.W("\t}")
		ret.W("\terr = s.insert%s(ctx, tx, logger, model)", revCol.Proper())
		ret.W("\tif err != nil {")
		ret.W("\t\treturn err")
		ret.W("\t}")
		ret.W("\treturn nil")
	} else {
		placeholder := ""
		if database == model.SQLServer {
			placeholder = "@"
		}
		ret.W("\tq := database.SQLUpdate(tableQuoted, columnsQuoted, %q, s.db.Placeholder())", pks.WhereClause(len(m.Columns), placeholder))
		ret.W("\tdata := model.ToData()")
		ret.W("\tdata = append(data, %s)", strings.Join(pkVals, ", "))
		token := "="
		if len(m.Columns.WithTag("created")) == 0 {
			token = serviceAssignmentToken
		}
		ret.W("\t_, err %s s.db.Update(ctx, q, tx, 1, logger, data...)", token)
		ret.W("\tif err != nil {")
		ret.W("\t\treturn err")
		ret.W("\t}")
		ret.W("\treturn nil")
	}
	ret.W("}")
	return ret, nil
}

func serviceUpdateIfNeeded(m *model.Model, g *golang.File, database string) (*golang.Block, error) {
	ret := golang.NewBlock("UpdateIfNeeded", "func")
	ret.W("func (s *Service) UpdateIfNeeded(ctx context.Context, tx *sqlx.Tx, model *%s, logger util.Logger) error {", m.Proper())
	if m.IsRevision() {
		revCol := m.HistoryColumn()
		ret.W("\trevs, err := s.getCurrent%s(ctx, tx, logger, model)", revCol.ProperPlural())
		ret.W("\tif err != nil {")
		ret.W("\t\treturn err")
		ret.W("\t}")
		ret.W("\tmodel.%s = revs[model.String()] + 1", revCol.Proper())
	}

	if cc := m.Columns.WithTag("created"); len(cc) > 0 {
		g.AddImport(helper.ImpErrors)
		ret.W("\tcurr, err := s.Get(ctx, tx, %s%s, logger)", m.PKs().ToRefs("model."), m.SoftDeleteSuffix())
		ret.W("\tif curr == nil || err != nil {")
		ret.W("\t\treturn s.Create(ctx, tx, logger, model)")
		ret.W("\t}")
		for _, c := range cc {
			ret.W("\tmodel.%s = curr.%s", c.Proper(), c.Proper())
		}
	}

	for _, updated := range m.Columns.WithTag("updated") {
		err := serviceSetVal(updated, g, ret, 1)
		if err != nil {
			return nil, errors.Wrap(err, "")
		}
	}

	if m.IsHistory() {
		ret.W("")
		ret.W("\th, hErr := s.SaveHistory(ctx, tx, curr, model, logger)")
		ret.W("\tif hErr != nil {")
		ret.W("\t\treturn errors.Wrap(hErr, \"unable to save history\")")
		ret.W("\t}")
		ret.W("\tif h == nil || len(h.Changes) == 0 {")
		ret.W("\t\treturn nil")
		ret.W("\t}")
	}

	pks := m.PKs()
	pkVals := make([]string, 0, len(pks))
	for _, pk := range pks {
		pkVals = append(pkVals, "model."+pk.Proper())
	}
	if m.IsRevision() {
		revCol := m.HistoryColumn()
		ret.W("")
		ret.W("\terr = s.upsertCore(ctx, tx, logger, model)")
		ret.W("\tif err != nil {")
		ret.W("\t\treturn err")
		ret.W("\t}")
		ret.W("\terr = s.insert%s(ctx, tx, logger, model)", revCol.Proper())
		ret.W("\tif err != nil {")
		ret.W("\t\treturn err")
		ret.W("\t}")
		ret.W("\treturn nil")
	} else {
		placeholder := ""
		if database == model.SQLServer {
			placeholder = "@"
		}
		ret.W("\tq := database.SQLUpdate(tableQuoted, columnsQuoted, %q, s.db.Placeholder())", pks.WhereClause(len(m.Columns), placeholder))
		ret.W("\tdata := model.ToData()")
		ret.W("\tdata = append(data, %s)", strings.Join(pkVals, ", "))
		token := "="
		if len(m.Columns.WithTag("created")) == 0 {
			token = serviceAssignmentToken
		}
		ret.W("\t_, err %s s.db.Update(ctx, q, tx, 1, logger, data...)", token)
		ret.W("\tif err != nil {")
		ret.W("\t\treturn err")
		ret.W("\t}")
		ret.W("\treturn nil")
	}
	ret.W("}")
	return ret, nil
}

func serviceSave(m *model.Model, g *golang.File, database string) (*golang.Block, error) {
	ret := golang.NewBlock("Save", "func")
	ret.W("func (s *Service) Save(ctx context.Context, tx *sqlx.Tx, logger util.Logger, models ...*%s) error {", m.Proper())
	ret.W("\tif len(models) == 0 {")
	ret.W("\t\treturn nil")
	ret.W("\t}")

	if m.IsRevision() {
		ret.W("\trevs, err := s.getCurrent%s(ctx, tx, logger, models...)", m.HistoryColumns(true).Col.ProperPlural())
		ret.W("\tif err != nil {")
		ret.W("\t\treturn err")
		ret.W("\t}")
	}

	if err := serviceAddCreatedUpdated(m, ret, g, false); err != nil {
		return nil, err
	}
	if m.IsRevision() {
		ret.W("")
		ret.W("\terr = s.upsertCore(ctx, tx, logger, models...)")
		ret.W("\tif err != nil {")
		ret.W("\t\treturn err")
		ret.W("\t}")
		ret.W("\terr = s.insert%s(ctx, tx, logger, models...)", m.HistoryColumn().Proper())
		ret.W("\tif err != nil {")
		ret.W("\t\treturn err")
		ret.W("\t}")
		ret.W("\treturn nil")
	} else {
		q := strings.Join(m.PKs().NamesQuoted(), ", ")
		ret.W("\tq := database.SQLUpsert(tableQuoted, columnsQuoted, len(models), []string{%s}, columnsQuoted, s.db.Placeholder())", q)
		ret.W("\tvar data []any")
		ret.W("\tfor _, model := range models {")
		ret.W("\t\tdata = append(data, model.ToData()...)")
		ret.W("\t}")
		ret.W("\treturn s.db.Insert(ctx, q, tx, logger, data...)")
	}
	ret.W("}")
	return ret, nil
}

func serviceUpsertCore(m *model.Model, g *golang.File, database string) *golang.Block {
	g.AddImport(helper.ImpAppUtil)
	ret := golang.NewBlock("UpsertCore", "func")
	ret.W("func (s *Service) upsertCore(ctx context.Context, tx *sqlx.Tx, logger util.Logger, models ...*%s) error {", m.Proper())
	ret.W("\tconflicts := util.StringArrayQuoted([]string{%s})", strings.Join(m.PKs().NamesQuoted(), ", "))
	ret.W("\tq := database.SQLUpsert(tableQuoted, columnsCore, len(models), conflicts, columnsCore, s.db.Placeholder())")
	ret.W("\tdata := make([]any, 0, len(columnsCore)*len(models))")
	ret.W("\tfor _, model := range models {")
	ret.W("\t\tdata = append(data, model.ToDataCore()...)")
	ret.W("\t}")
	ret.W("\t_, err := s.db.Update(ctx, q, tx, 1, logger, data...)")
	ret.W("\treturn err")
	ret.W("}")
	return ret
}

func serviceInsertRevision(m *model.Model, database string) *golang.Block {
	revCol := m.HistoryColumn()
	ret := golang.NewBlock("InsertRev", "func")
	ret.W("func (s *Service) insert%s(ctx context.Context, tx *sqlx.Tx, logger util.Logger, models ...*%s) error {", m.HistoryColumn().Proper(), m.Proper())
	ret.W("\tq := database.SQLInsert(table%sQuoted, columns%s, len(models), s.db.Placeholder())", revCol.Proper(), revCol.Proper())
	ret.W("\tdata := make([]any, 0, len(columns%s)*len(models))", revCol.Proper())
	ret.W("\tfor _, model := range models {")
	ret.W("\t\tdata = append(data, model.ToData%s()...)", m.HistoryColumn().Proper())
	ret.W("\t}")
	ret.W("\treturn s.db.Insert(ctx, q, tx, logger, data...)")
	ret.W("}")
	return ret
}

func serviceAddCreatedUpdated(m *model.Model, ret *golang.Block, g *golang.File, loadCurr bool, additional ...string) error {
	createdCols := m.Columns.WithTag("created")
	updatedCols := m.Columns.WithTag("updated")
	if len(createdCols) > 0 || len(updatedCols) > 0 || m.IsRevision() {
		ret.W("\tfor _, model := range models {")
		err := serviceLoadCreated(g, ret, m, createdCols, loadCurr)
		if err != nil {
			return err
		}
		if m.IsRevision() {
			ret.W("\t\tmodel.%s = revs[model.String()] + 1", m.HistoryColumn().Proper())
		}
		for _, updated := range updatedCols {
			err := serviceSetVal(updated, g, ret, 2)
			if err != nil {
				return err
			}
		}
		if m.IsHistory() && loadCurr {
			ret.W("")
			ret.W("\t\t_, hErr := s.SaveHistory(ctx, tx, curr, model)")
			ret.W("\t\tif hErr != nil {")
			ret.W("\t\t\treturn errors.Wrap(hErr, \"unable to save history\")")
			ret.W("\t\t}")
		}
		ret.W("\t}")
	}
	return nil
}

func serviceLoadCreated(g *golang.File, ret *golang.Block, m *model.Model, createdCols model.Columns, loadCurr bool) error {
	if len(createdCols) > 0 {
		if loadCurr {
			ret.W("\t\tcurr, e := s.Get(ctx, tx, %s%s)", m.PKs().ToRefs("model."), m.SoftDeleteSuffix())
			ret.W("\t\tif e == nil && curr != nil {")
			for _, created := range createdCols {
				ret.W("\t\t\tmodel.%s = curr.%s", created.Proper(), created.Proper())
			}
			ret.W("\t\t} else {")
			for _, created := range createdCols {
				err := serviceSetVal(created, g, ret, 3)
				if err != nil {
					return err
				}
			}
			ret.W("\t\t}")
		} else {
			for _, created := range createdCols {
				err := serviceSetVal(created, g, ret, 2)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func serviceSetVal(c *model.Column, g *golang.File, ret *golang.Block, indent int) error {
	ind := util.StringRepeat("\t", indent)
	if c.Type.Key() == types.KeyTimestamp || c.Type.Key() == types.KeyDate {
		if c.Nullable {
			g.AddImport(helper.ImpAppUtil)
			ret.W(ind+"model.%s = util.NowPointer()", c.Proper())
		} else {
			g.AddImport(helper.ImpTime)
			ret.W(ind+"model.%s = time.Now()", c.Proper())
		}
	} else {
		return errors.New("unhandled type [" + c.Type.Key() + "]")
	}
	return nil
}
