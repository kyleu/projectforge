package svc

import (
	"fmt"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/project/export/model"
	"strings"
)

func serviceGetMultipleSinglePK(m *model.Model, dbRef string) *golang.Block {
	name := "GetMultiple"
	ret := golang.NewBlock(name, "func")
	pk := m.PKs()[0]
	t := model.ToGoType(pk.Type, pk.Nullable, m.Package)
	msg := "func (s *Service) %s(ctx context.Context, tx *sqlx.Tx%s, logger util.Logger, %s ...%s) (%s, error) {"
	ret.W(msg, name, getSuffix(m), pk.Plural(), t, m.ProperPlural())
	ret.W("\tif len(%s) == 0 {", pk.Plural())
	ret.W("\t\treturn %s{}, nil", m.ProperPlural())
	ret.W("\t}")
	ret.W("\twc := database.SQLInClause(%q, len(%s), 0)", pk.Name, pk.Plural())
	if m.IsSoftDelete() {
		ret.W("\twc = addDeletedClause(wc, includeDeleted)")
	}
	ret.W("\tret := dtos{}")
	ret.W("\tq := database.SQLSelectSimple(columnsString, %s, wc)", tableClauseFor(m))
	ret.W("\tvals := make([]any, 0, len(%s))", pk.Plural())
	ret.W("\tfor _, x := range %s {", pk.Plural())
	ret.W("\t\tvals = append(vals, x)")
	ret.W("\t}")
	ret.W("\terr := s.%s.Select(ctx, &ret, q, tx, logger, vals...)", dbRef)
	ret.W("\tif err != nil {")
	ret.W("\t\treturn nil, errors.Wrapf(err, \"unable to get %s for [%%%%d] %s\", len(%s))", m.ProperPlural(), pk.Plural(), pk.Plural())
	ret.W("\t}")
	ret.W("\treturn ret.To%s(), nil", m.ProperPlural())
	ret.W("}")
	return ret
}

func serviceGetMultipleManyPKs(m *model.Model, dbRef string) *golang.Block {
	name := "GetMultiple"
	ret := golang.NewBlock(name, "func")
	pks := m.PKs()
	//t := model.ToGoType(pk.Type, pk.Nullable, m.Package)
	msg := "func (s *Service) %s(ctx context.Context, tx *sqlx.Tx%s, logger util.Logger, pks ...*PK) (%s, error) {"
	ret.W(msg, name, getSuffix(m), m.ProperPlural())
	ret.W("\tif len(pks) == 0 {")
	ret.W("\t\treturn %s{}, nil", m.ProperPlural())
	ret.W("\t}")
	ret.W("\twc := \"(\"")
	ret.W("\tfor idx := range pks {")
	ret.W("\t\tif idx > 0 {")
	ret.W("\t\t\twc += \" or \"")
	ret.W("\t\t}")
	var tags []string
	var idxs []string
	for idx, pk := range pks {
		tags = append(tags, fmt.Sprintf("%s = $%%%%d", pk.Name))
		idxs = append(idxs, fmt.Sprintf("(idx*%d)+%d", len(pks), idx+1))
	}
	ret.W("\t\twc += fmt.Sprintf(\"(%s)\", %s)", strings.Join(tags, " and "), strings.Join(idxs, ", "))
	ret.W("\t}")
	ret.W("\twc += \")\"")
	if m.IsSoftDelete() {
		ret.W("\twc = addDeletedClause(wc, includeDeleted)")
	}
	ret.W("\tret := dtos{}")
	ret.W("\tq := database.SQLSelectSimple(columnsString, %s, wc)", tableClauseFor(m))
	ret.W("\tvals := make([]any, 0, len(pks)*%d)", len(pks))
	ret.W("\tfor _, x := range pks {")
	var refs []string
	for _, pk := range pks {
		refs = append(refs, fmt.Sprintf("x.%s", pk.Proper()))
	}
	ret.W("\t\tvals = append(vals, %s)", strings.Join(refs, ", "))
	ret.W("\t}")
	ret.W("\terr := s.%s.Select(ctx, &ret, q, tx, logger, vals...)", dbRef)
	ret.W("\tif err != nil {")
	ret.W("\t\treturn nil, errors.Wrapf(err, \"unable to get %s for [%%%%d] pks\", len(pks))", m.ProperPlural())
	ret.W("\t}")
	ret.W("\treturn ret.To%s(), nil", m.ProperPlural())
	ret.W("}")
	return ret
}
