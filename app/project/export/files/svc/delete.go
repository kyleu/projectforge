package svc

import (
	"strings"

	"projectforge.dev/projectforge/app/project/export/enum"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/project/export/model"
)

const (
	argString = "ctx context.Context, tx *sqlx.Tx, wc string, expected int, logger util.Logger, values ...any"
	delMsg    = "// Delete doesn't actually delete, it only sets [%s]."
)

func serviceDelete(m *model.Model, enums enum.Enums) (*golang.Block, error) {
	pks := m.PKs()
	ret := golang.NewBlock("Delete", "func")
	args, err := pks.Args(m.Package, enums)
	if err != nil {
		return nil, err
	}
	ret.W("func (s *Service) Delete(ctx context.Context, tx *sqlx.Tx, %s, logger util.Logger) error {", args)
	ret.W("\tq := database.SQLDelete(tableQuoted, defaultWC(0), s.db.Type)")
	ret.W("\t_, err := s.db.Delete(ctx, q, tx, 1, logger, %s)", strings.Join(pks.CamelNames(), ", "))
	ret.W("\treturn err")
	ret.W("}")
	return ret, nil
}

func serviceSoftDelete(m *model.Model, enums enum.Enums) (*golang.Block, error) {
	pks := m.PKs()
	delCols := m.Columns.WithTag("deleted")
	ret := golang.NewBlock("Delete", "func")
	ret.W(delMsg, strings.Join(delCols.Names(), ", "))
	args, err := pks.Args(m.Package, enums)
	if err != nil {
		return nil, err
	}
	ret.W("func (s *Service) Delete(ctx context.Context, tx *sqlx.Tx, %s, logger util.Logger) error {", args)
	ret.W("\tcols := []string{%s}", strings.Join(delCols.NamesQuoted(), ", "))
	ret.W("\tq := database.SQLUpdate(tableQuoted, cols, defaultWC(len(cols)), s.db.Type)")
	ret.W("\t_, err := s.db.Update(ctx, q, tx, 1, logger, util.TimeCurrent(), %s)", strings.Join(pks.CamelNames(), ", "))
	ret.W("\treturn err")
	ret.W("}")
	return ret, nil
}

func serviceDeleteWhere(_ *model.Model) *golang.Block {
	ret := golang.NewBlock("Delete", "func")
	ret.W("func (s *Service) DeleteWhere(%s) error {", argString)
	ret.W("\tq := database.SQLDelete(tableQuoted, wc, s.db.Type)")
	ret.W("\t_, err := s.db.Delete(ctx, q, tx, expected, logger, values...)")
	ret.W("\treturn err")
	ret.W("}")
	return ret
}

func serviceSoftDeleteWhere(m *model.Model) *golang.Block {
	delCols := m.Columns.WithTag("deleted")
	ret := golang.NewBlock("Delete", "func")
	ret.W(delMsg, strings.Join(delCols.Names(), ", "))
	ret.W("func (s *Service) DeleteWhere(%s) error {", argString)
	ret.W("\tcols := []string{%s}", strings.Join(delCols.NamesQuoted(), ", "))
	ret.W("\tq := database.SQLUpdate(tableQuoted, cols, wc, s.db.Type)")
	ret.W("\t_, err := s.db.Update(ctx, q, tx, expected, logger, append([]any{util.TimeCurrent()}, values...)...)")
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
	ret.W("\treturn wc + \" and \\\"%s\\\""+helper.TextIsNull+"\"", delCols[0].SQL())
	ret.W("}")
	return ret
}
