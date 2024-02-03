package svc

import (
	"fmt"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/project/export/enum"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/project/export/model"
)

func serviceGetMultipleSingleCol(m *model.Model, name string, col *model.Column, dbRef string, enums enum.Enums) (*golang.Block, error) {
	ret := golang.NewBlock(name, "func")
	t, err := model.ToGoType(col.Type, col.Nullable, m.Package, enums)
	if err != nil {
		return nil, err
	}
	msg := "func (s *Service) %s(ctx context.Context, tx *sqlx.Tx, params *filter.Params%s, logger util.Logger, %s ...%s) (%s, error) {"
	ret.W(msg, name, getSuffix(m), col.CamelPlural(), t, m.ProperPlural())
	ret.W("\tif len(%s) == 0 {", col.CamelPlural())
	ret.W("\t\treturn %s{}, nil", m.ProperPlural())
	ret.W("\t}")
	ret.W("\tparams = filters(params)")
	ret.W("\twc := database.SQLInClause(%q, len(%s), 0, s.db.Type)", col.Name, col.CamelPlural())
	if m.IsSoftDelete() {
		ret.W("\twc = addDeletedClause(wc, includeDeleted)")
	}
	ret.W("\tq := database.SQLSelect(columnsString, %s, wc, params.OrderByString(), params.Limit, params.Offset, s.db.Type)", tableClause)
	ret.W("\tret := rows{}")
	ret.W("\terr := s.%s.Select(ctx, &ret, q, tx, logger, lo.ToAnySlice(%s)...)", dbRef, col.CamelPlural())
	ret.W("\tif err != nil {")
	ret.W("\t\treturn nil, errors.Wrapf(err, \"unable to get %s for [%%%%d] %s\", len(%s))", m.ProperPlural(), col.CamelPlural(), col.CamelPlural())
	ret.W("\t}")
	ret.W("\treturn ret.To%s(), nil", m.ProperPlural())
	ret.W("}")
	return ret, nil
}

func serviceGetMultipleManyCols(m *model.Model, name string, cols model.Columns, dbRef string) *golang.Block {
	ret := golang.NewBlock(name, "func")

	tags := make([]string, 0, len(cols))
	idxs := make([]string, 0, len(cols))
	refs := make([]string, 0, len(cols))
	lo.ForEach(cols, func(pk *model.Column, idx int) {
		tags = append(tags, fmt.Sprintf("%s = $%%%%d", pk.Name))
		idxs = append(idxs, fmt.Sprintf("(idx*%d)+%d", len(cols), idx+1))
		refs = append(refs, fmt.Sprintf("x.%s", pk.Proper()))
	})

	msg := "func (s *Service) %s(ctx context.Context, tx *sqlx.Tx%s, logger util.Logger, pks ...*PK) (%s, error) {"
	ret.W(msg, name, getSuffix(m), m.ProperPlural())
	ret.W("\tif len(pks) == 0 {")
	ret.W("\t\treturn %s{}, nil", m.ProperPlural())
	ret.W("\t}")
	ret.W("\twc := \"(\"")
	ret.W("\tlo.ForEach(pks, func(_ *PK, idx int) {")
	ret.W("\t\tif idx > 0 {")
	ret.W("\t\t\twc += \" or \"")
	ret.W("\t\t}")
	ret.W("\t\twc += fmt.Sprintf(\"(%s)\", %s)", strings.Join(tags, " and "), strings.Join(idxs, ", "))
	ret.W("\t})")
	ret.W("\twc += \")\"")
	if m.IsSoftDelete() {
		ret.W("\twc = addDeletedClause(wc, includeDeleted)")
	}
	ret.W("\tret := rows{}")
	ret.W("\tq := database.SQLSelectSimple(columnsString, %s, s.db.Type, wc)", tableClause)

	ret.W("\tvals := lo.FlatMap(pks, func(x *PK, _ int) []any {")
	ret.W("\t\treturn []any{%s}", strings.Join(refs, ", "))
	ret.W("\t})")

	ret.W("\terr := s.%s.Select(ctx, &ret, q, tx, logger, vals...)", dbRef)
	ret.W("\tif err != nil {")
	ret.W("\t\treturn nil, errors.Wrapf(err, \"unable to get %s for [%%%%d] pks\", len(pks))", m.ProperPlural())
	ret.W("\t}")
	ret.W("\treturn ret.To%s(), nil", m.ProperPlural())
	ret.W("}")
	return ret
}
