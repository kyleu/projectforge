package svc

import (
	"fmt"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/project/export/enum"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/project/export/model"
	"projectforge.dev/projectforge/app/util"
)

func serviceSearch(m *model.Model, grp *model.Column, dbRef string, enums enum.Enums, database string) (*golang.Block, error) {
	prefix := ""
	if grp != nil {
		prefix = "By" + grp.Proper()
	}
	var clauses []string
	hasEqual, hasLike := false, false
	lo.ForEach(m.AllSearches(database), func(s string, _ int) {
		if strings.HasPrefix(s, "=") {
			hasEqual = true
		} else {
			hasLike = true
		}
	})
	eq, like := "$1", "$2"
	if database == util.DatabaseSQLServer {
		eq, like = "@p1", "@p2"
	}
	var params []string
	if hasEqual {
		params = append(params, "strings.ToLower(query)")
	}
	if hasLike {
		params = append(params, `"%%" + strings.ToLower(query) + "%%"`)
		if !hasEqual {
			like = eq
		}
	}
	lo.ForEach(m.AllSearches(database), func(s string, _ int) {
		if strings.HasPrefix(s, "=") {
			clauses = append(clauses, s+" = "+eq)
		} else {
			clauses = append(clauses, s+" like "+like)
		}
	})

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
	defSC := "searchClause"
	if m.IsSoftDelete() {
		defSC = "addDeletedClause(searchClause, includeDeleted)"
	}
	ret.W("\tvar wc string")
	ret.W("\tvar vals []any")
	ret.W("\tvar err error")
	ret.W("\tif strings.Contains(query, \":\") {")
	ret.W("\t\twc, vals, err = database.QueryFieldDescs(FieldDescs, query, 0)")
	ret.W("\t} else {")
	ret.W("\t\twc = %s", defSC)
	ret.W("\t\tvals = []any{%s}", strings.Join(params, ", "))
	ret.W("\t}")
	ret.WE(1, "nil")

	ret.W("\tq := database.SQLSelect(columnsString, %s, wc, params.OrderByString(), params.Limit, params.Offset, s.db.Type)", tableClause)
	ret.W("\tret := rows{}")
	ret.W("\terr = s.%s.Select(ctx, &ret, q, tx, logger, vals...)", dbRef)
	ret.WE(1, "nil")

	ret.W("\treturn ret.To%s(), nil", m.ProperPlural())
	ret.W("}")
	return ret, nil
}
