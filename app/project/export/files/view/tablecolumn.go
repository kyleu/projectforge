package view

import (
	"fmt"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/project/export/golang"
)

func viewTableColumn(
	ret *golang.Block, models model.Models, m *model.Model, link bool, col *model.Column, modelKey string, prefix string, indent int, enums enum.Enums,
) {
	rels := m.RelationsFor(col)
	call, err := getTableColumnString(m, modelKey, rels, models, prefix)
	if err != nil {
		panic(err)
	}
	viewColumn("table", ret, m, col, call, link, modelKey, indent, models, enums, "paths")
}

func getTableColumnString(m *model.Model, modelKey string, rels model.Relations, models model.Models, prefix string) (string, error) {
	if len(rels) == 0 {
		return "", nil
	}
	if len(rels) > 1 {
		return "", errors.Errorf("expected one relation, found [%d]", len(rels))
	}
	rel := rels[0]
	relModel := models.Get(rel.Table)
	if !relModel.CanTraverseRelation() {
		return "", errors.Errorf("can't traverse relation [%s]", rel.Name)
	}
	srcCol := m.Columns.Get(rel.Src[0])
	tgtCol := relModel.Columns.Get(rel.Tgt[0])
	k := relModel.CamelPlural()
	if prefix != "" {
		k = prefix + relModel.ProperPlural()
	}
	relTitles := relModel.Columns.WithTag("title")
	if len(relTitles) == 0 {
		relTitles = relModel.PKs()
	}
	get := fmt.Sprintf("%sBy%s.Get(%s%s)", k, srcCol.Proper(), modelKey, srcCol.Proper())
	if srcCol.Nullable && !srcCol.Type.Scalar() && !tgtCol.Nullable {
		get = fmt.Sprintf("%sBy%s.Get(*%s%s)", k, srcCol.Proper(), modelKey, srcCol.Proper())
	}
	return get, nil
}
