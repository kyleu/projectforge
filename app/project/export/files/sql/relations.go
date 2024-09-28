package sql

import (
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/project/export/golang"
)

func sqlRelations(ret *golang.Block, m *model.Model, _ model.Models) {
	lo.ForEach(m.Relations, func(rel *model.Relation, _ int) {
		ret.WF("  foreign key (%s) references %q (%s),", rel.SrcQuoted(), rel.Table, rel.TgtQuoted())
	})
}
