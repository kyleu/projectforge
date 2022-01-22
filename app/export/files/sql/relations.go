package sql

import (
	"github.com/kyleu/projectforge/app/export/golang"
	"github.com/kyleu/projectforge/app/export/model"
)

func sqlRelations(ret *golang.Block, m *model.Model) {
	for _, rel := range m.Relations {
		ret.W("  foreign key (%s) references %q (%s),", rel.SrcQuoted(), rel.Table, rel.TgtQuoted())
	}
}
