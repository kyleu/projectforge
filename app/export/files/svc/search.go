package svc

import (
	"fmt"
	"strings"

	"github.com/kyleu/projectforge/app/export/golang"
	"github.com/kyleu/projectforge/app/export/model"
)

func serviceSearch(m *model.Model, grp *model.Column) *golang.Block {
	prefix := ""
	if grp != nil {
		prefix = "By" + grp.Proper()
	}
	var clauses []string
	hasEqual, hasLike := false, false
	for _, s := range m.Search {
		if strings.HasPrefix(s, "=") {
			hasEqual = true
		} else {
			hasLike = true
		}
	}
	eq := "$1"
	like := "$2"
	var params []string
	if hasEqual {
		params = append(params, "strings.ToLower(q)")
	}
	if hasLike {
		params = append(params, `"%%"+strings.ToLower(q)+"%%"`)
		if !hasEqual {
			like = eq
		}
	}
	for _, s := range m.Search {
		if strings.HasPrefix(s, "=") {
			clauses = append(clauses, s+" = "+eq)
		} else {
			clauses = append(clauses, s+" like "+like)
		}
	}
	wc := "(" + strings.Join(clauses, " or ") + ")"

	ret := golang.NewBlock("search", "func")
	grpTxt := ""
	if grp == nil {
		ret.W("const searchClause = %q", wc)
		ret.W("")
	} else {
		grpTxt = fmt.Sprintf(", %s %s", grp.Camel(), grp.ToGoType())
	}
	decl := "func (s *Service) Search%s(ctx context.Context%s, q string, tx *sqlx.Tx, params *filter.Params%s) (%s, error) {"
	ret.W(decl, prefix, grpTxt, getSuffix(m), m.ProperPlural())
	ret.W("\tparams = filters(params)")
	if m.IsSoftDelete() {
		ret.W("\twc := addDeletedClause(searchClause, includeDeleted)")
	} else {
		ret.W("\twc := searchClause")
	}
	ret.W("\tsql := database.SQLSelect(columnsString, %s, wc, params.OrderByString(), params.Limit, params.Offset)", tableClauseFor(m))
	ret.W("\tret := dtos{}")
	ret.W("\terr := s.db.Select(ctx, &ret, sql, tx, " + strings.Join(params, ", ") + ")")
	ret.W("\tif err != nil {")
	ret.W("\t\treturn nil, err")
	ret.W("\t}")
	ret.W("\treturn ret.To%s(), nil", m.ProperPlural())
	ret.W("}")
	return ret
}
