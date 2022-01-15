package svc

import (
	"strings"

	"github.com/kyleu/projectforge/app/export/files/helper"
	"github.com/kyleu/projectforge/app/export/golang"
	"github.com/kyleu/projectforge/app/export/model"
	"github.com/kyleu/projectforge/app/file"
)

func ServiceGet(m *model.Model, args *model.Args) (*file.File, error) {
	g := golang.NewFile(m.Package, []string{"app", m.Package}, "serviceget")
	for _, imp := range helper.ImportsForTypes("go", m.PKs().Types()...) {
		g.AddImport(imp)
	}
	g.AddImport(helper.ImpContext, helper.ImpErrors, helper.ImpSQLx, helper.ImpFilter, helper.ImpDatabase)
	g.AddBlocks(serviceList(m))
	if len(m.PKs()) > 0 {
		g.AddBlocks(serviceGetByPK(m))
	}
	if len(m.PKs()) > 1 {
		for _, pkCol := range m.PKs() {
			g.AddBlocks(serviceGetBy("GetBy"+pkCol.Proper(), m, model.Columns{pkCol}, true))
		}
	}
	for _, grp := range m.GroupedColumns() {
		g.AddBlocks(serviceGrouped(m, grp))
		g.AddImport(helper.ImpAppUtil)
		if !grp.PK {
			g.AddBlocks(serviceGetBy("GetBy"+grp.Proper(), m, model.Columns{grp}, true))
		}
	}
	if len(m.Search) > 0 {
		g.AddImport(helper.ImpStrings)
		g.AddBlocks(serviceSearch(m, nil))
	}
	for _, grp := range m.GroupedColumns() {
		if !grp.PK {
			if len(m.Search) > 0 {
				g.AddBlocks(serviceSearch(m, grp))
			}
		}
	}
	return g.Render()
}

func serviceGrouped(m *model.Model, grp *model.Column) *golang.Block {
	name := "Get" + grp.ProperPlural()
	ret := golang.NewBlock(name, "func")
	ret.W("func (s *Service) %s(ctx context.Context, tx *sqlx.Tx%s) ([]*util.KeyValInt, error) {", name, getSuffix(m))
	ret.W("\twc := \"\"")
	if m.IsSoftDelete() {
		delCols := m.Columns.WithTag("deleted")
		ret.W("\tif !includeDeleted {")
		ret.W("\t\twc = %q", delCols[0].Name+" is null")
		ret.W("\t}")
	}
	cols := "vendor as key, count(*) as val"
	ret.W("\tsql := database.SQLSelectGrouped(%q, %s, wc, %q, %q, 0, 0)", cols, tableClauseFor(m), grp.Camel(), grp.Camel())
	ret.W("\tvar ret []*util.KeyValInt")
	ret.W("\terr := s.db.Select(ctx, &ret, sql, tx)")
	ret.W("\tif err != nil {")
	ret.W("\t\treturn nil, errors.Wrap(err, \"unable to get %s\")", m.ProperPlural())
	ret.W("\t}")
	ret.W("\treturn ret, nil")
	ret.W("}")
	return ret
}

func serviceList(m *model.Model) *golang.Block {
	ret := golang.NewBlock("List", "func")
	ret.W("func (s *Service) List(ctx context.Context, tx *sqlx.Tx, params *filter.Params%s) (%s, error) {", getSuffix(m), m.ProperPlural())
	ret.W("\tparams = filters(params)")
	ret.W("\twc := \"\"")
	if m.IsSoftDelete() {
		delCols := m.Columns.WithTag("deleted")
		ret.W("\tif !includeDeleted {")
		ret.W("\t\twc = %q", delCols[0].Name+" is null")
		ret.W("\t}")
	}
	ret.W("\tsql := database.SQLSelect(columnsString, %s, wc, params.OrderByString(), params.Limit, params.Offset)", tableClauseFor(m))
	ret.W("\tret := dtos{}")
	ret.W("\terr := s.db.Select(ctx, &ret, sql, tx)")
	ret.W("\tif err != nil {")
	ret.W("\t\treturn nil, errors.Wrap(err, \"unable to get %s\")", m.ProperPlural())
	ret.W("\t}")
	ret.W("\treturn ret.To%s(), nil", m.ProperPlural())
	ret.W("}")
	return ret
}

func serviceGetByPK(m *model.Model) *golang.Block {
	return serviceGetBy("Get", m, m.PKs(), false)
}

func serviceGetBy(key string, m *model.Model, cols model.Columns, returnMultiple bool) *golang.Block {
	if returnMultiple {
		return serviceGetMultiple(key, m, cols)
	}
	return serviceGetOne(key, m, cols)
}

func serviceGetOne(key string, m *model.Model, cols model.Columns) *golang.Block {
	if key == "" {
		key = "GetBy" + cols.Smushed()
	}
	ret := golang.NewBlock(key, "func")
	ret.W("func (s *Service) %s(ctx context.Context, tx *sqlx.Tx, %s%s) (*%s, error) {", key, cols.Args(), getSuffix(m), m.Proper())
	ret.W("\twc := %q", cols.WhereClause(0))
	if m.IsSoftDelete() {
		ret.W("\twc = addDeletedClause(wc, includeDeleted)")
	}
	ret.W("\tret := &dto{}")
	ret.W("\tsql := database.SQLSelectSimple(columnsString, %s, wc)", tableClauseFor(m))
	ret.W("\terr := s.db.Get(ctx, ret, sql, tx, %s)", strings.Join(cols.CamelNames(), ", "))
	ret.W("\tif err != nil {")
	sj := strings.Join(cols.CamelNames(), ", ")
	decls := make([]string, 0, len(cols))
	for _, c := range cols {
		decls = append(decls, c.Camel()+" [%%s]")
	}
	ret.W("\t\treturn nil, errors.Wrapf(err, \"unable to get %s by %s\", %s)", m.Camel(), strings.Join(decls, ", "), sj)
	ret.W("\t}")
	ret.W("\treturn ret.To%s(), nil", m.Proper())
	ret.W("}")
	return ret
}

func serviceGetMultiple(key string, m *model.Model, cols model.Columns) *golang.Block {
	if key == "" {
		key = "GetBy" + cols.Smushed()
	}
	ret := golang.NewBlock(key, "func")
	ret.W("func (s *Service) %s(ctx context.Context, tx *sqlx.Tx, %s, params *filter.Params%s) (%s, error) {", key, cols.Args(), getSuffix(m), m.ProperPlural())
	ret.W("\tparams = filters(params)")
	ret.W("\twc := %q", cols.WhereClause(0))
	if m.IsSoftDelete() {
		ret.W("\twc = addDeletedClause(wc, includeDeleted)")
	}
	ret.W("\tsql := database.SQLSelect(columnsString, %s, wc, params.OrderByString(), params.Limit, params.Offset)", tableClauseFor(m))
	ret.W("\tret := dtos{}")
	ret.W("\terr := s.db.Select(ctx, &ret, sql, tx, %s)", strings.Join(cols.CamelNames(), ", "))
	ret.W("\tif err != nil {")
	sj := strings.Join(cols.CamelNames(), ", ")
	decls := make([]string, 0, len(cols))
	for _, c := range cols {
		decls = append(decls, c.Camel()+" [%%s]")
	}
	ret.W("\t\treturn nil, errors.Wrapf(err, \"unable to get %s by %s\", %s)", m.Plural(), strings.Join(decls, ", "), sj)
	ret.W("\t}")
	ret.W("\treturn ret.To%s(), nil", m.ProperPlural())
	ret.W("}")
	return ret
}

func tableClauseFor(m *model.Model) string {
	if m.IsRevision() {
		return "tablesJoined"
	}
	return "table"
}
