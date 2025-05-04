package svc

import (
	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/util"
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
	ret.WF("func (s *Service) Delete(ctx context.Context, tx *sqlx.Tx, %s, logger util.Logger) error {", args)
	ret.W("\tq := database.SQLDelete(tableQuoted, defaultWC(0), s.db.Type)")
	ret.WF("\t_, err := s.db.Delete(ctx, q, tx, 1, logger, %s)", util.StringJoin(pks.CamelNames(), ", "))
	if m.HasTag("events") {
		ret.W("\tif s.Events != nil {")
		ret.W("\t\tif e := s.Events.Delete(ctx, tx, logger, id); e != nil {")
		ret.WF("\t\t\treturn errors.Wrap(e, \"error processing [delete] events for %s\")", m.Proper())
		ret.W("\t\t}")
		ret.W("\t}")
	}
	ret.W("\treturn err")
	ret.W("}")
	return ret, nil
}

func serviceSoftDelete(m *model.Model, enums enum.Enums) (*golang.Block, error) {
	pks := m.PKs()
	delCols := m.Columns.WithTag("deleted")
	ret := golang.NewBlock("Delete", "func")
	ret.WF(delMsg, util.StringJoin(delCols.Names(), ", "))
	args, err := pks.Args(m.Package, enums)
	if err != nil {
		return nil, err
	}
	ret.WF("func (s *Service) Delete(ctx context.Context, tx *sqlx.Tx, %s, logger util.Logger) error {", args)
	ret.WF("\tcols := []string{%s}", util.StringJoin(delCols.SQLQuoted(), ", "))
	ret.W("\tq := database.SQLUpdate(tableQuoted, cols, defaultWC(len(cols)), s.db.Type)")
	ret.WF("\t_, err := s.db.Update(ctx, q, tx, 1, logger, util.TimeCurrent(), %s)", util.StringJoin(pks.CamelNames(), ", "))
	if m.HasTag("events") {
		ret.W("\tif s.Events != nil {")
		ret.W("\t\tif e := s.Events.Delete(ctx, tx, logger, id); e != nil {")
		ret.WF("\t\t\treturn errors.Wrap(e, \"error processing [delete] events for %s\")", m.Proper())
		ret.W("\t\t}")
		ret.W("\t}")
	}
	ret.W("\treturn err")
	ret.W("}")
	return ret, nil
}

func serviceDeleteWhere(_ *model.Model) *golang.Block {
	ret := golang.NewBlock("Delete", "func")
	ret.WF("func (s *Service) DeleteWhere(%s) error {", argString)
	ret.W("\tq := database.SQLDelete(tableQuoted, wc, s.db.Type)")
	ret.W("\t_, err := s.db.Delete(ctx, q, tx, expected, logger, values...)")
	ret.W("\treturn err")
	ret.W("}")
	return ret
}

func serviceSoftDeleteWhere(m *model.Model) *golang.Block {
	delCols := m.Columns.WithTag("deleted")
	ret := golang.NewBlock("Delete", "func")
	ret.WF(delMsg, util.StringJoin(delCols.Names(), ", "))
	ret.WF("func (s *Service) DeleteWhere(%s) error {", argString)
	ret.WF("\tcols := []string{%s}", util.StringJoin(delCols.SQLQuoted(), ", "))
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
	ret.WF("\treturn wc + \" and \\\"%s\\\""+helper.TextIsNull+"\"", delCols[0].SQL())
	ret.W("}")
	return ret
}
