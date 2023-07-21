package svc

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/project/export/model"
	"projectforge.dev/projectforge/app/util"
)

const serviceAssignmentToken = ":="

func ServiceMutate(m *model.Model, args *model.Args, addHeader bool, linebreak string) (*file.File, error) {
	g := golang.NewFile(m.Package, []string{"app", m.PackageWithGroup("")}, "servicemutate")
	lo.ForEach(helper.ImportsForTypes("go", "", m.PKs().Types()...), func(imp *golang.Import, _ int) {
		g.AddImport(imp)
	})
	g.AddImport(helper.ImpAppUtil, helper.ImpContext, helper.ImpSQLx, helper.ImpDatabase, helper.ImpLo)

	if add, err := serviceCreate(g, m, args.Audit(m)); err == nil {
		g.AddBlocks(add)
	} else {
		return nil, err
	}
	if upd, err := serviceUpdate(g, m, args.Audit(m), args.Database()); err == nil {
		g.AddBlocks(upd)
	} else {
		return nil, err
	}
	if m.IsRevision() || m.IsHistory() {
		if updIN, err := serviceUpdateIfNeeded(g, m, args.Database()); err == nil {
			g.AddBlocks(updIN)
		} else {
			return nil, err
		}
	}
	if save, err := serviceSave(g, m); err == nil {
		g.AddBlocks(save)
	} else {
		return nil, err
	}
	if m.IsRevision() {
		g.AddBlocks(serviceUpsertCore(g, m), serviceInsertRevision(m))
	}
	if m.IsSoftDelete() {
		g.AddImport(helper.ImpAppUtil)
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
	return g.Render(addHeader, linebreak)
}

func serviceCreate(g *golang.File, m *model.Model, audit bool) (*golang.Block, error) {
	ret := golang.NewBlock("Create", "func")
	ret.W("func (s *Service) Create(ctx context.Context, tx *sqlx.Tx, logger util.Logger, models ...*%s) error {", m.Proper())
	ret.W("\tif len(models) == 0 {")
	ret.W("\t\treturn nil")
	ret.W("\t}")

	if m.IsRevision() {
		revCol := m.HistoryColumn()
		ret.W("\trevs, err := s.getCurrent%s(ctx, tx, logger, models...)", revCol.ProperPlural())
		ret.WE(1)

		if err := serviceAddCreatedUpdated(g, m, ret, false); err != nil {
			return nil, err
		}

		ret.WB()
		ret.W("\terr = s.upsertCore(ctx, tx, logger, models...)")
		ret.WE(1)
		ret.W("\terr = s.insert%s(ctx, tx, logger, models...)", revCol.Proper())
		ret.WE(1)
		ret.W("\treturn nil")
	} else {
		if err := serviceAddCreatedUpdated(g, m, ret, false); err != nil {
			return nil, err
		}
		ret.W("\tq := database.SQLInsert(tableQuoted, columnsQuoted, len(models), s.db.Placeholder())")

		if audit {
			ret.W("\tvals := make([]any, 0, len(models)*len(columnsQuoted))")
			ret.W("\tfor _, arg := range models {")
			msg := "\t\t_, _, err := s.audit.ApplyObjSimple(ctx, \"%s.create\", \"created new %s\", nil, arg, %q, nil, logger)"
			ret.W(msg, m.Proper(), m.TitleLower(), m.Proper())
			ret.WE(2)
			ret.W("\t\tvals = append(vals, arg.ToData()...)")
			ret.W("\t}")
		} else {
			ret.W("\tvals := lo.FlatMap(models, func(arg *%s, _ int) []any {", m.Proper())
			ret.W("\t\treturn arg.ToData()")
			ret.W("\t})")
		}
		ret.W("\treturn s.db.Insert(ctx, q, tx, logger, vals...)")
	}
	ret.W("}")
	return ret, nil
}

func serviceUpdate(g *golang.File, m *model.Model, audit bool, database string) (*golang.Block, error) {
	ret := golang.NewBlock("Update", "func")
	ret.W("func (s *Service) Update(ctx context.Context, tx *sqlx.Tx, model *%s, logger util.Logger) error {", m.Proper())
	if m.IsRevision() {
		revCol := m.HistoryColumn()
		ret.W("\trevs, err := s.getCurrent%s(ctx, tx, logger, model)", revCol.ProperPlural())
		ret.WE(1)
		ret.W("\tmodel.%s = revs[model.String()] + 1", revCol.Proper())
	}

	if cc := m.Columns.WithTag("created"); len(cc) > 0 || audit {
		g.AddImport(helper.ImpErrors)
		ret.W("\tcurr, err := s.Get(ctx, tx, %s%s, logger)", m.PKs().ToRefs("model."), m.SoftDeleteSuffix())
		ret.W("\tif err != nil {")
		ret.W("\t\treturn errors.Wrapf(err, \"can't get original %s [%%%%s]\", model.String())", m.TitleLower())
		ret.W("\t}")
		lo.ForEach(cc, func(c *model.Column, _ int) {
			ret.W("\tmodel.%s = curr.%s", c.Proper(), c.Proper())
		})
	}

	for _, updated := range m.Columns.WithTag("updated") {
		err := serviceSetVal(g, updated, ret, 1)
		if err != nil {
			return nil, errors.Wrap(err, "")
		}
	}

	if m.IsHistory() {
		ret.WB()
		ret.W("\t_, hErr := s.SaveHistory(ctx, tx, curr, model, logger)")
		ret.W("\tif hErr != nil {")
		ret.W("\t\treturn errors.Wrap(hErr, \"unable to save history\")")
		ret.W("\t}")
	}

	pks := m.PKs()
	pkVals := make([]string, 0, len(pks))
	lo.ForEach(pks, func(pk *model.Column, _ int) {
		pkVals = append(pkVals, "model."+pk.Proper())
	})
	if m.IsRevision() {
		revCol := m.HistoryColumn()
		ret.WB()
		ret.W("\terr = s.upsertCore(ctx, tx, logger, model)")
		ret.WE(1)
		ret.W("\terr = s.insert%s(ctx, tx, logger, model)", revCol.Proper())
		ret.WE(1)
		if audit {
			serviceAuditApply(g, m, ret)
		}
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
		if len(m.Columns.WithTag("created")) == 0 && (!audit) {
			token = serviceAssignmentToken
		}
		ret.W("\t_, err %s s.db.Update(ctx, q, tx, 1, logger, data...)", token)
		ret.WE(1)
		if audit {
			serviceAuditApply(g, m, ret)
		}
		ret.W("\treturn nil")
	}
	ret.W("}")
	return ret, nil
}

func serviceAuditApply(g *golang.File, m *model.Model, ret *golang.Block) {
	g.AddImport(helper.ImpFmt)
	ret.W("\tmsg := fmt.Sprintf(\"updated %s [%%%%s]\", model.String())", m.Title())
	ret.W("\t_, _, err = s.audit.ApplyObjSimple(ctx, \"%s.update\", msg, curr, model, %q, nil, logger)", m.Proper(), m.Proper())
	ret.WE(1)
}

func serviceUpdateIfNeeded(g *golang.File, m *model.Model, database string) (*golang.Block, error) {
	ret := golang.NewBlock("UpdateIfNeeded", "func")
	ret.W("func (s *Service) UpdateIfNeeded(ctx context.Context, tx *sqlx.Tx, model *%s, logger util.Logger) error {", m.Proper())
	if m.IsRevision() {
		revCol := m.HistoryColumn()
		ret.W("\trevs, err := s.getCurrent%s(ctx, tx, logger, model)", revCol.ProperPlural())
		ret.WE(1)
		ret.W("\tmodel.%s = revs[model.String()] + 1", revCol.Proper())
	}

	if cc := m.Columns.WithTag("created"); len(cc) > 0 {
		g.AddImport(helper.ImpErrors)
		ret.W("\tcurr, err := s.Get(ctx, tx, %s%s, logger)", m.PKs().ToRefs("model."), m.SoftDeleteSuffix())
		ret.W("\tif curr == nil || err != nil {")
		ret.W("\t\treturn s.Create(ctx, tx, logger, model)")
		ret.W("\t}")
		lo.ForEach(cc, func(c *model.Column, _ int) {
			ret.W("\tmodel.%s = curr.%s", c.Proper(), c.Proper())
		})
	}

	for _, updated := range m.Columns.WithTag("updated") {
		err := serviceSetVal(g, updated, ret, 1)
		if err != nil {
			return nil, errors.Wrap(err, "")
		}
	}

	if m.IsHistory() {
		ret.WB()
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
	lo.ForEach(pks, func(pk *model.Column, _ int) {
		pkVals = append(pkVals, "model."+pk.Proper())
	})
	if m.IsRevision() {
		revCol := m.HistoryColumn()
		ret.WB()
		ret.W("\terr = s.upsertCore(ctx, tx, logger, model)")
		ret.WE(1)
		ret.W("\terr = s.insert%s(ctx, tx, logger, model)", revCol.Proper())
		ret.WE(1)
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
		ret.WE(1)
		ret.W("\treturn nil")
	}
	ret.W("}")
	return ret, nil
}

func serviceSave(g *golang.File, m *model.Model) (*golang.Block, error) {
	ret := golang.NewBlock("Save", "func")
	ret.W("func (s *Service) Save(ctx context.Context, tx *sqlx.Tx, logger util.Logger, models ...*%s) error {", m.Proper())
	ret.W("\tif len(models) == 0 {")
	ret.W("\t\treturn nil")
	ret.W("\t}")

	if m.IsRevision() {
		ret.W("\trevs, err := s.getCurrent%s(ctx, tx, logger, models...)", m.HistoryColumns(true).Col.ProperPlural())
		ret.WE(1)
	}

	if err := serviceAddCreatedUpdated(g, m, ret, false); err != nil {
		return nil, err
	}
	if m.IsRevision() {
		ret.WB()
		ret.W("\terr = s.upsertCore(ctx, tx, logger, models...)")
		ret.WE(1)
		ret.W("\terr = s.insert%s(ctx, tx, logger, models...)", m.HistoryColumn().Proper())
		ret.WE(1)
		ret.W("\treturn nil")
	} else {
		q := strings.Join(m.PKs().NamesQuoted(), ", ")
		ret.W("\tq := database.SQLUpsert(tableQuoted, columnsQuoted, len(models), []string{%s}, columnsQuoted, s.db.Placeholder())", q)
		ret.W("\tdata := lo.FlatMap(models, func(model *%s, _ int) []any {", m.Proper())
		ret.W("\t\treturn model.ToData()")
		ret.W("\t})")
		ret.W("\treturn s.db.Insert(ctx, q, tx, logger, data...)")
	}
	ret.W("}")
	return ret, nil
}

func serviceUpsertCore(g *golang.File, m *model.Model) *golang.Block {
	g.AddImport(helper.ImpAppUtil)
	ret := golang.NewBlock("UpsertCore", "func")
	ret.W("func (s *Service) upsertCore(ctx context.Context, tx *sqlx.Tx, logger util.Logger, models ...*%s) error {", m.Proper())
	ret.W("\tconflicts := util.StringArrayQuoted([]string{%s})", strings.Join(m.PKs().NamesQuoted(), ", "))
	ret.W("\tq := database.SQLUpsert(tableQuoted, columnsCore, len(models), conflicts, columnsCore, s.db.Placeholder())")
	ret.W("\tdata := lo.FlatMap(models, func(model *%s, _ int) []any {", m.Proper())
	ret.W("\t\treturn model.ToDataCore()")
	ret.W("\t})")
	ret.W("\t_, err := s.db.Update(ctx, q, tx, 1, logger, data...)")
	ret.W("\treturn err")
	ret.W("}")
	return ret
}

func serviceInsertRevision(m *model.Model) *golang.Block {
	revCol := m.HistoryColumn()
	ret := golang.NewBlock("InsertRev", "func")
	ret.W("func (s *Service) insert%s(ctx context.Context, tx *sqlx.Tx, logger util.Logger, models ...*%s) error {", m.HistoryColumn().Proper(), m.Proper())
	ret.W("\tq := database.SQLInsert(table%sQuoted, columns%s, len(models), s.db.Placeholder())", revCol.Proper(), revCol.Proper())
	ret.W("\tdata := lo.FlatMap(models, func(model *%s, _ int) []any {", m.Proper())
	ret.W("\t\treturn model.ToData%s()", revCol.Proper())
	ret.W("\t})")
	ret.W("\treturn s.db.Insert(ctx, q, tx, logger, data...)")
	ret.W("}")
	return ret
}

func serviceAddCreatedUpdated(g *golang.File, m *model.Model, ret *golang.Block, loadCurr bool) error {
	createdCols := m.Columns.WithTag("created")
	updatedCols := m.Columns.WithTag("updated")
	if len(createdCols) > 0 || len(updatedCols) > 0 || m.IsRevision() {
		ret.W("\tlo.ForEach(models, func(model *%s, _ int) {", m.Proper())
		err := serviceLoadCreated(g, ret, m, createdCols, loadCurr)
		if err != nil {
			return err
		}
		if m.IsRevision() {
			ret.W("\t\tmodel.%s = revs[model.String()] + 1", m.HistoryColumn().Proper())
		}
		for _, updated := range updatedCols {
			err := serviceSetVal(g, updated, ret, 2)
			if err != nil {
				return err
			}
		}
		if m.IsHistory() && loadCurr {
			ret.WB()
			ret.W("\t\t_, hErr := s.SaveHistory(ctx, tx, curr, model)")
			ret.W("\t\tif hErr != nil {")
			ret.W("\t\t\treturn errors.Wrap(hErr, \"unable to save history\")")
			ret.W("\t\t}")
		}
		ret.W("\t})")
	}
	return nil
}

func serviceLoadCreated(g *golang.File, ret *golang.Block, m *model.Model, createdCols model.Columns, loadCurr bool) error {
	if len(createdCols) > 0 {
		if loadCurr {
			ret.W("\t\tcurr, e := s.Get(ctx, tx, %s%s)", m.PKs().ToRefs("model."), m.SoftDeleteSuffix())
			ret.W("\t\tif e == nil && curr != nil {")
			lo.ForEach(createdCols, func(created *model.Column, _ int) {
				ret.W("\t\t\tmodel.%s = curr.%s", created.Proper(), created.Proper())
			})
			ret.W("\t\t} else {")
			for _, created := range createdCols {
				err := serviceSetVal(g, created, ret, 3)
				if err != nil {
					return err
				}
			}
			ret.W("\t\t}")
		} else {
			for _, created := range createdCols {
				err := serviceSetVal(g, created, ret, 2)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func serviceSetVal(g *golang.File, c *model.Column, ret *golang.Block, indent int) error {
	ind := util.StringRepeat("\t", indent)
	if c.Type.Key() == types.KeyTimestamp || c.Type.Key() == types.KeyDate {
		if c.Nullable {
			g.AddImport(helper.ImpAppUtil)
			ret.W(ind+"model.%s = util.TimeCurrentP()", c.Proper())
		} else {
			g.AddImport(helper.ImpAppUtil)
			ret.W(ind+"model.%s = util.TimeCurrent()", c.Proper())
		}
	} else {
		return errors.New("unhandled type [" + c.Type.Key() + "]")
	}
	return nil
}
