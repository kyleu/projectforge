package view

import (
	"fmt"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
)

func viewTableColumn(
	g *golang.Template, ret *golang.Block, models model.Models, m *model.Model, link bool,
	col *model.Column, modelKey string, prefix string, indent int, enums enum.Enums,
) {
	rels := m.RelationsFor(col)
	toStrings := getTableColumnStrings(m, modelKey, rels, models, prefix)
	viewColumn("table", g, ret, m, col, toStrings, link, modelKey, indent, models, enums)
}

func getTableColumnStrings(m *model.Model, modelKey string, rels model.Relations, models model.Models, prefix string) string {
	ret := lo.FilterMap(rels, func(rel *model.Relation, _ int) (string, bool) {
		relModel := models.Get(rel.Table)
		if !relModel.CanTraverseRelation() {
			return "", false
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
		if len(relTitles) == 1 && relTitles[0].Name == tgtCol.Name {
			return "", false
		}
		st, end := "{%% if x := ", "; x != nil %%} ({%%s x.TitleString() %%})"
		if srcCol.Nullable && !tgtCol.Nullable {
			get := fmt.Sprintf("%sBy%s.Get(*%s%s)", k, srcCol.Proper(), modelKey, srcCol.Proper())
			return "{%% if " + modelKey + srcCol.Proper() + " != nil %%}" +
				st + get + end + helper.TextEndIf + helper.TextEndIf, true
		}
		get := fmt.Sprintf("%sBy%s.Get(%s%s)", k, srcCol.Proper(), modelKey, srcCol.Proper())
		return get + "||" + st + get + end + helper.TextEndIf, true
	})
	return strings.Join(ret, "")
}
