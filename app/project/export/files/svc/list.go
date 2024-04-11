package svc

import (
	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/project/export/model"
)

func ServiceList(m *model.Model, args *model.Args, addHeader bool, linebreak string) (*file.File, error) {
	dbRef := args.DBRef()
	g := golang.NewFile(m.Package, []string{"app", m.PackageWithGroup("")}, "servicelist")
	g.AddImport(helper.ImpAppUtil, helper.ImpContext, helper.ImpErrors, helper.ImpSQLx, helper.ImpFilter, helper.ImpAppDatabase)
	g.AddBlocks(serviceList(m, args.DBRef()), serviceListSQL(m, args.DBRef()), serviceListWhere(m), serviceCount(g, m, args.DBRef()))
	if m.HasSearches() {
		g.AddImport(helper.ImpStrings)
		ss, err := serviceSearch(m, nil, dbRef, args.Enums, args.Database)
		if err != nil {
			return nil, err
		}
		g.AddBlocks(ss)
	}
	return g.Render(addHeader, linebreak)
}

func serviceList(m *model.Model, dbRef string) *golang.Block {
	ret := golang.NewBlock("List", "func")
	ret.W("func (s *Service) List(ctx context.Context, tx *sqlx.Tx, params *filter.Params%s, logger util.Logger) (%s, error) {", getSuffix(m), m.ProperPlural())
	ret.W("\tparams = filters(params)")
	ret.W("\twc := \"\"")
	if m.IsSoftDelete() {
		delCols := m.Columns.WithTag("deleted")
		ret.W("\tif !includeDeleted {")
		ret.W("\t\twc = %q", delCols[0].NameQuoted()+helper.TextIsNull)
		ret.W("\t}")
	}
	ret.W("\tq := database.SQLSelect(columnsString, %s, wc, params.OrderByString(), params.Limit, params.Offset, s.db.Type)", tableClause)
	ret.W("\tret := rows{}")
	ret.W("\terr := s.%s.Select(ctx, &ret, q, tx, logger)", dbRef)
	ret.W("\tif err != nil {")
	ret.W("\t\treturn nil, errors.Wrap(err, \"unable to get %s\")", m.TitlePluralLower())
	ret.W("\t}")
	ret.W("\treturn ret.To%s(), nil", m.ProperPlural())
	ret.W("}")
	return ret
}

func serviceListSQL(m *model.Model, dbRef string) *golang.Block {
	ret := golang.NewBlock("ListSQL", "func")
	ret.W("func (s *Service) ListSQL(ctx context.Context, tx *sqlx.Tx, sql string, logger util.Logger, values ...any) (%s, error) {", m.ProperPlural())
	ret.W("\tret := rows{}")
	ret.W("\terr := s.%s.Select(ctx, &ret, sql, tx, logger, values...)", dbRef)
	ret.W("\tif err != nil {")
	ret.W("\t\treturn nil, errors.Wrap(err, \"unable to get %s using custom SQL\")", m.TitlePluralLower())
	ret.W("\t}")
	ret.W("\treturn ret.To%s(), nil", m.ProperPlural())
	ret.W("}")
	return ret
}

func serviceListWhere(m *model.Model) *golang.Block {
	ret := golang.NewBlock("ListWhere", "func")
	decl := "func (s *Service) ListWhere(ctx context.Context, tx *sqlx.Tx, where string, params *filter.Params, logger util.Logger, values ...any) (%s, error) {"
	ret.W(decl, m.ProperPlural())
	ret.W("\tparams = filters(params)")
	ret.W("\tsql := database.SQLSelect(columnsString, tableQuoted, where, params.OrderByString(), params.Limit, params.Offset, s.db.Type)")
	ret.W("\treturn s.ListSQL(ctx, tx, sql, logger, values...)")
	ret.W("}")
	return ret
}

func serviceCount(g *golang.File, m *model.Model, dbRef string) *golang.Block {
	g.AddImport(helper.ImpStrings)
	ret := golang.NewBlock("Count", "func")
	ret.W("func (s *Service) Count(ctx context.Context, tx *sqlx.Tx, whereClause string%s, logger util.Logger, args ...any) (int, error) {", getSuffix(m))
	ret.W("\tif strings.Contains(whereClause, \"'\") || strings.Contains(whereClause, \";\") {")
	ret.W("\t\treturn 0, errors.Errorf(\"invalid where clause [%%s]\", whereClause)")
	ret.W("\t}")
	if m.IsSoftDelete() {
		delCols := m.Columns.WithTag("deleted")
		ret.W("\tif !includeDeleted {")
		ret.W("\t\tif whereClause == \"\" {")
		ret.W("\t\t\twhereClause = %q", delCols[0].NameQuoted()+helper.TextIsNull)
		ret.W("\t\t} else {")
		ret.W("\t\t\twhereClause += \" and \" + %q", delCols[0].NameQuoted()+helper.TextIsNull)
		ret.W("\t\t}")
		ret.W("\t}")
	}
	ret.W("\tq := database.SQLSelectSimple(\"count(*) as x\", %s, s.db.Type, whereClause)", tableClause)
	ret.W("\tret, err := s.%s.SingleInt(ctx, q, tx, logger, args...)", dbRef)
	ret.W("\tif err != nil {")
	ret.W("\t\treturn 0, errors.Wrap(err, \"unable to get count of %s\")", m.TitlePluralLower())
	ret.W("\t}")
	ret.W("\treturn int(ret), nil")
	ret.W("}")
	return ret
}
