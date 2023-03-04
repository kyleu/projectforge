package sql

import (
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/project/export/model"
)

func sqlRelations(ret *golang.Block, m *model.Model, models model.Models) {
	for _, rel := range m.Relations {
		ret.W("  foreign key (%s) references %q (%s),", rel.SrcQuoted(), rel.Table, rel.TgtQuoted())
	}
}
