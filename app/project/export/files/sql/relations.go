package sql

import (
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/project/export/model"
)

func sqlRelations(ret *golang.Block, m *model.Model, _ model.Models) {
	lo.ForEach(m.Relations, func(rel *model.Relation, _ int) {
		ret.W("  foreign key (%s) references %q (%s),", rel.SrcQuoted(), rel.Table, rel.TgtQuoted())
	})
}
