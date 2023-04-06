package svc

import (
	"fmt"
	"strings"

	"projectforge.dev/projectforge/app/project/export/enum"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/project/export/model"
)

const serviceGetMultipleName = "GetMultiple"

func serviceGetMultipleSinglePK(m *model.Model, dbRef string, enums enum.Enums) (*golang.Block, error) {
	ret := golang.NewBlock(serviceGetMultipleName, "func")
	pk := m.PKs()[0]
	t, err := model.ToGoType(pk.Type, pk.Nullable, m.Package, enums)
	if err != nil {
		return nil, err
	}
	msg := "func (s *Service) %s(ctx context.Context, tx *sqlx.Tx%s, logger util.Logger, %s ...%s) (%s, error) {"
	ret.W(msg, serviceGetMultipleName, getSuffix(m), pk.Plural(), t, m.ProperPlural())
	ret.W("\tif len(%s) == 0 {", pk.Plural())
	ret.W("\t\treturn %s{}, nil", m.ProperPlural())
	ret.W("\t}")
	ret.W("\twc := database.SQLInClause(%q, len(%s), 0, s.db.Placeholder())", pk.Name, pk.Plural())
	if m.IsSoftDelete() {
		ret.W("\twc = addDeletedClause(wc, includeDeleted)")
	}
	ret.W("\tret := rows{}")
	ret.W("\tq := database.SQLSelectSimple(columnsString, %s, s.db.Placeholder(), wc)", tableClauseFor(m))
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
	return ret, nil
}

func serviceGetMultipleManyPKs(m *model.Model, dbRef string) *golang.Block {
	ret := golang.NewBlock(serviceGetMultipleName, "func")
	pks := m.PKs()

	tags := make([]string, 0, len(pks))
	idxs := make([]string, 0, len(pks))
	refs := make([]string, 0, len(pks))
	for idx, pk := range pks {
		tags = append(tags, fmt.Sprintf("%s = $%%%%d", pk.Name))
		idxs = append(idxs, fmt.Sprintf("(idx*%d)+%d", len(pks), idx+1))
		refs = append(refs, fmt.Sprintf("x.%s", pk.Proper()))
	}

	msg := "func (s *Service) %s(ctx context.Context, tx *sqlx.Tx%s, logger util.Logger, pks ...*PK) (%s, error) {"
	ret.W(msg, serviceGetMultipleName, getSuffix(m), m.ProperPlural())
	ret.W("\tif len(pks) == 0 {")
	ret.W("\t\treturn %s{}, nil", m.ProperPlural())
	ret.W("\t}")
	ret.W("\twc := \"(\"")
	ret.W("\tfor idx := range pks {")
	ret.W("\t\tif idx > 0 {")
	ret.W("\t\t\twc += \" or \"")
	ret.W("\t\t}")
	ret.W("\t\twc += fmt.Sprintf(\"(%s)\", %s)", strings.Join(tags, " and "), strings.Join(idxs, ", "))
	ret.W("\t}")
	ret.W("\twc += \")\"")
	if m.IsSoftDelete() {
		ret.W("\twc = addDeletedClause(wc, includeDeleted)")
	}
	ret.W("\tret := rows{}")
	ret.W("\tq := database.SQLSelectSimple(columnsString, %s, s.db.Placeholder(), wc)", tableClauseFor(m))
	ret.W("\tvals := make([]any, 0, len(pks)*%d)", len(pks))
	ret.W("\tfor _, x := range pks {")
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
