package svc

import (
	"strings"


	"projectforge.dev/projectforge/app/project/export/enum"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/project/export/model"
)

const argString = "ctx context.Context, tx *sqlx.Tx, wc string, expected int, logger util.Logger, values ...any"

func serviceDelete(m *model.Model, enums enum.Enums) (*golang.Block, error) {
	pks := m.PKs()
	ret := golang.NewBlock("Delete", "func")
	args, err := pks.Args(m.Package, enums)
	if err != nil {
		return nil, err
	}
	ret.W("func (s *Service) Delete(ctx context.Context, tx *sqlx.Tx, %s, logger util.Logger) error {", args)
	ret.W("\tq := database.SQLDelete(tableQuoted, defaultWC(0))")
	ret.W("\t_, err := s.db.Delete(ctx, q, tx, 1, logger, %s)", strings.Join(pks.CamelNames(), ", "))
	ret.W("\treturn err")
	ret.W("}")
	return ret, nil
}

func serviceSoftDelete(m *model.Model, enums enum.Enums) (*golang.Block, error) {
	pks := m.PKs()
	delCols := m.Columns.WithTag("deleted")
	ret := golang.NewBlock("Delete", "func")
	ret.W("// Delete doesn't actually delete, it only sets [" + strings.Join(delCols.Names(), ", ") + "].")
	args, err := pks.Args(m.Package, enums)
	if err != nil {
		return nil, err
	}
	ret.W("func (s *Service) Delete(ctx context.Context, tx *sqlx.Tx, %s, logger util.Logger) error {", args)
	ret.W("\tcols := []string{%s}", strings.Join(delCols.NamesQuoted(), ", "))
	ret.W("\tq := database.SQLUpdate(tableQuoted, cols, defaultWC(len(cols)), \"\")")
	ret.W("\t_, err := s.db.Update(ctx, q, tx, 1, logger, time.Now(), %s)", strings.Join(pks.CamelNames(), ", "))
	ret.W("\treturn err")
	ret.W("}")
	return ret, nil
}

func serviceDeleteWhere(m *model.Model) *golang.Block {
	ret := golang.NewBlock("Delete", "func")
	ret.W("func (s *Service) DeleteWhere(%s) error {", argString)
	ret.W("\tq := database.SQLDelete(tableQuoted, wc)")
	ret.W("\t_, err := s.db.Delete(ctx, q, tx, expected, logger, values...)")
	ret.W("\treturn err")
	ret.W("}")
	return ret
}

func serviceSoftDeleteWhere(m *model.Model) *golang.Block {
	delCols := m.Columns.WithTag("deleted")
	ret := golang.NewBlock("Delete", "func")
	ret.W("// Delete doesn't actually delete, it only sets [" + strings.Join(delCols.Names(), ", ") + "].")
	ret.W("func (s *Service) DeleteWhere(%s) error {", argString)
	ret.W("\tcols := []string{%s}", strings.Join(delCols.NamesQuoted(), ", "))
	ret.W("\tq := database.SQLUpdate(tableQuoted, cols, wc, \"\")")
	ret.W("\t_, err := s.db.Update(ctx, q, tx, expected, logger, append([]any{time.Now()}, values...)...)")
	ret.W("\treturn err")
	ret.W("}")
	return ret
}

func serviceAddDeletedClause(m *model.Model) *golang.Block {
	delCols := m.Columns.WithTag("deleted")
	ret := golang.NewBlock("Delete", "func")
	ret.W("func addDeletedClause(wc string, includeDeleted bool) string {")
	ret.W("\tif includeDeleted {")
	ret.W("\t\treturn wc")
	ret.W("\t}")
	ret.W("\treturn wc + \" and \\\"%s\\\" is null\"", delCols[0].Name)
	ret.W("}")
	return ret
}
