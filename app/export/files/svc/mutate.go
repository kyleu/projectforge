package svc

import (
	"strings"

	"github.com/kyleu/projectforge/app/export/files/helper"
	"github.com/kyleu/projectforge/app/export/golang"
	"github.com/kyleu/projectforge/app/export/model"
	"github.com/kyleu/projectforge/app/file"
	"github.com/kyleu/projectforge/app/util"
	"github.com/pkg/errors"
)

func ServiceMutate(m *model.Model, args *model.Args) (*file.File, error) {
	g := golang.NewFile(m.Package, []string{"app", m.Package}, "servicemutate")
	for _, imp := range helper.ImportsForTypes("go", m.PKs().Types()...) {
		g.AddImport(imp)
	}
	g.AddImport(helper.ImpContext, helper.ImpSQLx, helper.ImpDatabase)

	if add, err := serviceCreate(m, g); err == nil {
		g.AddBlocks(add)
	} else {
		return nil, err
	}
	if upd, err := serviceUpdate(m, g); err == nil {
		g.AddBlocks(upd)
	} else {
		return nil, err
	}
	if save, err := serviceSave(m, g); err == nil {
		g.AddBlocks(save)
	} else {
		return nil, err
	}
	if m.IsRevision() {
		g.AddBlocks(serviceUpsertCore(m, g), serviceInsertRevision(m))
	}
	if m.IsSoftDelete() {
		g.AddImport(helper.ImpTime)
		softDel, err := serviceSoftDelete(m)
		if err != nil {
			return nil, err
		}
		g.AddBlocks(softDel)
	} else {
		g.AddBlocks(serviceDelete(m))
	}
	return g.Render()
}

func serviceCreate(m *model.Model, g *golang.File) (*golang.Block, error) {
	ret := golang.NewBlock("Create", "func")
	ret.W("func (s *Service) Create(ctx context.Context, tx *sqlx.Tx, models ...*%s) error {", m.Proper())
	ret.W("\tif len(models) == 0 {")
	ret.W("\t\treturn nil")
	ret.W("\t}")

	if m.IsRevision() {
		revCol := m.HistoryColumn()
		ret.W("\trevs, err := s.getCurrent%s(ctx, tx, models...)", revCol.ProperPlural())
		ret.W("\tif err != nil {")
		ret.W("\t\treturn err")
		ret.W("\t}")

		if err := serviceAddCreatedUpdated(m, ret, g); err != nil {
			return nil, err
		}

		ret.W("")
		ret.W("\terr = s.upsertCore(ctx, tx, models...)")
		ret.W("\tif err != nil {")
		ret.W("\t\treturn err")
		ret.W("\t}")
		ret.W("\terr = s.insert%s(ctx, tx, models...)", revCol.Proper())
		ret.W("\tif err != nil {")
		ret.W("\t\treturn err")
		ret.W("\t}")
		ret.W("\treturn nil")
	} else {
		if err := serviceAddCreatedUpdated(m, ret, g); err != nil {
			return nil, err
		}

		ret.W("\tq := database.SQLInsert(tableQuoted, columnsQuoted, len(models), \"\")")
		ret.W("\tvals := make([]interface{}, 0, len(models)*len(columnsQuoted))")
		ret.W("\tfor _, arg := range models {")
		ret.W("\t\tvals = append(vals, arg.ToData()...)")
		ret.W("\t}")
		ret.W("\treturn s.db.Insert(ctx, q, tx, vals...)")
	}
	ret.W("}")
	return ret, nil
}

func serviceUpdate(m *model.Model, g *golang.File) (*golang.Block, error) {
	ret := golang.NewBlock("Update", "func")
	ret.W("func (s *Service) Update(ctx context.Context, tx *sqlx.Tx, model *%s) error {", m.Proper())
	if m.IsRevision() {
		revCol := m.HistoryColumn()
		ret.W("\trevs, err := s.getCurrent%s(ctx, tx, model)", revCol.ProperPlural())
		ret.W("\tif err != nil {")
		ret.W("\t\treturn err")
		ret.W("\t}")
		ret.W("\tmodel.%s = revs[model.String()] + 1", revCol.Proper())
	}
	for _, updated := range m.Columns.WithTag("updated") {
		err := serviceSetVal(updated, g, ret, 1)
		if err != nil {
			return nil, errors.Wrap(err, "")
		}
	}
	pks := m.PKs()
	pkVals := make([]string, 0, len(pks))
	for _, pk := range pks {
		pkVals = append(pkVals, "model."+pk.Proper())
	}
	if m.IsRevision() {
		revCol := m.HistoryColumn()
		ret.W("")
		ret.W("\terr = s.upsertCore(ctx, tx, model)")
		ret.W("\tif err != nil {")
		ret.W("\t\treturn err")
		ret.W("\t}")
		ret.W("\terr = s.insert%s(ctx, tx, model)", revCol.Proper())
		ret.W("\tif err != nil {")
		ret.W("\t\treturn err")
		ret.W("\t}")
		ret.W("\treturn nil")
	} else {
		ret.W("\tq := database.SQLUpdate(tableQuoted, columnsQuoted, %q, \"\")", pks.WhereClause(len(m.Columns)))
		ret.W("\tdata := model.ToData()")
		ret.W("\tdata = append(data, %s)", strings.Join(pkVals, ", "))
		ret.W("\t_, ret := s.db.Update(ctx, q, tx, 1, data...)")
		ret.W("\treturn ret")
	}
	ret.W("}")
	return ret, nil
}

func serviceSave(m *model.Model, g *golang.File) (*golang.Block, error) {
	ret := golang.NewBlock("Save", "func")
	ret.W("func (s *Service) Save(ctx context.Context, tx *sqlx.Tx, models ...*%s) error {", m.Proper())
	ret.W("\tif len(models) == 0 {")
	ret.W("\t\treturn nil")
	ret.W("\t}")

	if m.IsRevision() {
		ret.W("\trevs, err := s.getCurrent%s(ctx, tx, models...)", m.HistoryColumns(true).Col.ProperPlural())
		ret.W("\tif err != nil {")
		ret.W("\t\treturn err")
		ret.W("\t}")
	}

	if err := serviceAddCreatedUpdated(m, ret, g); err != nil {
		return nil, err
	}
	if m.IsRevision() {
		ret.W("")
		ret.W("\terr = s.upsertCore(ctx, tx, models...)")
		ret.W("\tif err != nil {")
		ret.W("\t\treturn err")
		ret.W("\t}")
		ret.W("\terr = s.insert%s(ctx, tx, models...)", m.HistoryColumn().Proper())
		ret.W("\tif err != nil {")
		ret.W("\t\treturn err")
		ret.W("\t}")
		ret.W("\treturn nil")
	} else {
		q := strings.Join(m.PKs().NamesQuoted(), ", ")
		ret.W("\tq := database.SQLUpsert(tableQuoted, columnsQuoted, len(models), []string{%s}, columns, \"\")", q)
		ret.W("\tvar data []interface{}")
		ret.W("\tfor _, model := range models {")
		ret.W("\t\tdata = append(data, model.ToData()...)")
		ret.W("\t}")
		ret.W("\treturn s.db.Insert(ctx, q, tx, data...)")
	}
	ret.W("}")
	return ret, nil
}

func serviceUpsertCore(m *model.Model, g *golang.File) *golang.Block {
	g.AddImport(helper.ImpAppUtil)
	ret := golang.NewBlock("UpsertCore", "func")
	ret.W("func (s *Service) upsertCore(ctx context.Context, tx *sqlx.Tx, models ...*%s) error {", m.Proper())
	ret.W("\tconflicts := util.StringArrayQuoted([]string{%s})", strings.Join(m.PKs().NamesQuoted(), ", "))
	ret.W("\tq := database.SQLUpsert(tableQuoted, columnsCore, len(models), conflicts, columnsCore, \"\")")
	ret.W("\tdata := make([]interface{}, 0, len(columnsCore)*len(models))")
	ret.W("\tfor _, model := range models {")
	ret.W("\t\tdata = append(data, model.ToDataCore()...)")
	ret.W("\t}")
	ret.W("\t_, err := s.db.Update(ctx, q, tx, 1, data...)")
	ret.W("\treturn err")
	ret.W("}")
	return ret
}

func serviceInsertRevision(m *model.Model) *golang.Block {
	revCol := m.HistoryColumn()
	ret := golang.NewBlock("InsertRev", "func")
	ret.W("func (s *Service) insert%s(ctx context.Context, tx *sqlx.Tx, models ...*%s) error {", m.HistoryColumn().Proper(), m.Proper())
	ret.W("\tq := database.SQLInsert(table%sQuoted, columns%s, len(models), \"\")", revCol.Proper(), revCol.Proper())
	ret.W("\tdata := make([]interface{}, 0, len(columns%s)*len(models))", revCol.Proper())
	ret.W("\tfor _, model := range models {")
	ret.W("\t\tdata = append(data, model.ToData%s()...)", m.HistoryColumn().Proper())
	ret.W("\t}")
	ret.W("\treturn s.db.Insert(ctx, q, tx, data...)")
	ret.W("}")
	return ret
}

func serviceAddCreatedUpdated(m *model.Model, ret *golang.Block, g *golang.File, additional ...string) error {
	createdCols := m.Columns.WithTag("created")
	updatedCols := m.Columns.WithTag("updated")
	if len(createdCols) > 0 || len(updatedCols) > 0 || m.IsRevision() {
		ret.W("\tfor _, model := range models {")
		if m.IsRevision() {
			ret.W("\t\tmodel.%s = revs[model.String()] + 1", m.HistoryColumn().Proper())
		}
		for _, created := range createdCols {
			err := serviceSetVal(created, g, ret, 2)
			if err != nil {
				return err
			}
		}
		for _, updated := range updatedCols {
			err := serviceSetVal(updated, g, ret, 2)
			if err != nil {
				return err
			}
		}
		ret.W("\t}")
	}
	return nil
}

func serviceSetVal(c *model.Column, g *golang.File, ret *golang.Block, indent int) error {
	ind := util.StringRepeat("\t", indent)
	if c.Type.Key == model.TypeTimestamp.Key {
		if c.Nullable {
			g.AddImport(helper.ImpAppUtil)
			ret.W(ind+"model.%s = util.NowPointer()", c.Proper())
		} else {
			g.AddImport(helper.ImpTime)
			ret.W(ind+"model.%s = time.Now()", c.Proper())
		}
	} else {
		return errors.New("unhandled type [" + c.Type.Key + "]")
	}
	return nil
}
