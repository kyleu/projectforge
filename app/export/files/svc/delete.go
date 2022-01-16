package svc

import (
	"strings"

	"github.com/kyleu/projectforge/app/export/golang"
	"github.com/kyleu/projectforge/app/export/model"
	"github.com/pkg/errors"
)

func serviceDelete(m *model.Model) *golang.Block {
	pks := m.PKs()
	ret := golang.NewBlock("Delete", "func")
	ret.W("func (s *Service) Delete(ctx context.Context, tx *sqlx.Tx, %s) error {", pks.Args())
	ret.W("\tq := database.SQLDelete(tableQuoted, %q)", pks.WhereClause(0))
	ret.W("\t_, err := s.db.Delete(ctx, q, tx, 1, %s)", strings.Join(pks.CamelNames(), ", "))
	ret.W("\treturn err")
	ret.W("}")
	return ret
}

func serviceSoftDelete(m *model.Model) (*golang.Block, error) {
	pks := m.PKs()
	delCols := m.Columns.WithTag("deleted")
	if len(delCols) == 0 {
		return nil, errors.New("when [softDelete] is set, exactly one column must be tagged [deleted]")
	}
	if len(delCols) > 1 {
		return nil, errors.New("when [softDelete] is set, no more than one column may be tagged [deleted]")
	}
	ret := golang.NewBlock("Delete", "func")
	ret.W("// Delete doesn't actually delete, it only sets [" + strings.Join(delCols.Names(), ", ") + "].")
	ret.W("func (s *Service) Delete(ctx context.Context, tx *sqlx.Tx, %s) error {", pks.Args())
	ret.W("\tq := database.SQLUpdate(tableQuoted, []string{%s}, %q, \"\")", strings.Join(delCols.NamesQuoted(), ", "), pks.WhereClause(1))
	ret.W("\t_, err := s.db.Update(ctx, q, tx, 1, time.Now(), %s)", strings.Join(pks.CamelNames(), ", "))
	ret.W("\treturn err")
	ret.W("}")
	ret.W("")
	ret.W("func addDeletedClause(wc string, includeDeleted bool) string {")
	ret.W("\tif includeDeleted {")
	ret.W("\t\treturn wc")
	ret.W("\t}")
	ret.W("\treturn wc + \" and \\\"%s\\\" is null\"", delCols[0].Name)
	ret.W("}")
	return ret, nil
}
