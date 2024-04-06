package svc

import (
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/project/export/model"
)

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
