package view

import (
	"fmt"
	"slices"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/util"
)

func viewTableColumn(
	ret *golang.Block, models model.Models, m *model.Model, link bool, col *model.Column, modelKey string, prefix string, indent int, enums enum.Enums,
) {
	rels := m.RelationsFor(col)
	call, err := getTableColumnString(m, modelKey, rels, models, prefix)
	if err != nil {
		panic(err)
	}
	viewColumn("table", ret, m, col, call, link, modelKey, indent, models, enums, col.Nullable, "paths")
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

	srcCols := lo.Filter(m.Columns, func(x *model.Column, _ int) bool {
		return slices.Contains(rel.Src, x.Name)
	})
	tgtCols := lo.Filter(relModel.Columns, func(x *model.Column, _ int) bool {
		return slices.Contains(rel.Tgt, x.Name)
	})
	if len(srcCols) != len(tgtCols) {
		return "", errors.Errorf("invalid column size for relation [%s]", rel.Name)
	}
	k := relModel.CamelPlural()
	if prefix != "" {
		k = prefix + relModel.ProperPlural()
	}

	calls := make([]string, 0, len(srcCols))
	for i, srcCol := range srcCols {
		tgtCol := tgtCols[i]
		if srcCol.Nullable && !srcCol.Type.Scalar() && !tgtCol.Nullable {
			calls = append(calls, fmt.Sprintf("*%s%s", modelKey, srcCol.Proper()))
		} else {
			calls = append(calls, fmt.Sprintf("%s%s", modelKey, srcCol.Proper()))
		}
	}

	return fmt.Sprintf("%sBy%s.Get(%s)", k, util.StringJoin(srcCols.ProperNames(), ""), util.StringJoin(calls, ", ")), nil
}
