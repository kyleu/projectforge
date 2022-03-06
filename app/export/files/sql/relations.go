package sql

import (
	"projectforge.dev/app/export/golang"
	"projectforge.dev/app/export/model"
)

func sqlRelations(ret *golang.Block, m *model.Model) {
	for _, rel := range m.Relations {
		ret.W("  foreign key (%s) references %q (%s),", rel.SrcQuoted(), rel.Table, rel.TgtQuoted())
	}
}
