package svc

import (
	"fmt"
	"strings"

	"projectforge.dev/projectforge/app/project/export/enum"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/project/export/model"
)

func serviceSearch(m *model.Model, grp *model.Column, dbRef string, enums enum.Enums, database string) (*golang.Block, error) {
	prefix := ""
	if grp != nil {
		prefix = "By" + grp.Proper()
	}
	var clauses []string
	hasEqual, hasLike := false, false
	for _, s := range m.AllSearches(database) {
		if strings.HasPrefix(s, "=") {
			hasEqual = true
		} else {
			hasLike = true
		}
	}
	eq, like := "$1", "$2"
	if database == model.SQLServer {
		eq, like = "@p1", "@p2"
	}
	var params []string
	if hasEqual {
		params = append(params, "strings.ToLower(query)")
	}
	if hasLike {
		params = append(params, `"%%"+strings.ToLower(query)+"%%"`)
		if !hasEqual {
			like = eq
		}
	}
	for _, s := range m.AllSearches(database) {
		if strings.HasPrefix(s, "=") {
			clauses = append(clauses, s+" = "+eq)
		} else {
			clauses = append(clauses, s+" like "+like)
		}
	}

	ret := golang.NewBlock("search", "func")
	grpTxt := ""
	if grp == nil {
		wc := "(" + strings.Join(clauses, " or ") + ")"
		ret.W("const searchClause = %q", wc)
		ret.WB()
	} else {
		gt, err := grp.ToGoType(m.Package, enums)
		if err != nil {
			return nil, err
		}
		grpTxt = fmt.Sprintf(", %s %s", grp.Camel(), gt)
	}
	const decl = "func (s *Service) Search%s(ctx context.Context%s, query string, tx *sqlx.Tx, params *filter.Params%s, logger util.Logger) (%s, error) {"
	ret.W(decl, prefix, grpTxt, getSuffix(m), m.ProperPlural())
	ret.W("\tparams = filters(params)")
	if m.IsSoftDelete() {
		ret.W("\twc := addDeletedClause(searchClause, includeDeleted)")
	} else {
		ret.W("\twc := searchClause")
	}
	ret.W("\tq := database.SQLSelect(columnsString, %s, wc, params.OrderByString(), params.Limit, params.Offset, s.db.Placeholder())", tableClauseFor(m))
	ret.W("\tret := rows{}")
	ret.W("\terr := s.%s.Select(ctx, &ret, q, tx, logger, %s)", dbRef, strings.Join(params, ", "))
	ret.W("\tif err != nil {")
	ret.W("\t\treturn nil, err")
	ret.W("\t}")
	ret.W("\treturn ret.To%s(), nil", m.ProperPlural())
	ret.W("}")
	return ret, nil
}
