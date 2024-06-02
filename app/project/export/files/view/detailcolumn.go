package view

import (
	"fmt"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/util"
)

func viewDetailColumn(
	g *golang.Template, ret *golang.Block, models model.Models, m *model.Model, link bool, col *model.Column, modelKey string, indent int, enums enum.Enums,
) {
	rels := m.RelationsFor(col)
	toStrings := strings.Join(viewDetailColumnString(rels, models, m, col), "")
	viewColumn(util.KeyDetail, g, ret, m, col, toStrings, link, modelKey, indent, models, enums)
}

func ModelLinkURL(m *model.Model, prefix string, enums enum.Enums) string {
	pks := m.PKs()
	linkURL := "/" + m.Route()
	lo.ForEach(pks, func(pk *model.Column, _ int) {
		linkURL += "/" + pk.ToGoViewString(prefix, false, true, enums, util.KeySimple)
	})
	return linkURL
}

func RelationWebPath(rel *model.Relation, src *model.Model, tgt *model.Model, prefix string) interface{} {
	url := "`/" + tgt.Route() + "`"
	lo.ForEach(rel.Src, func(s string, _ int) {
		c := src.Columns.Get(s)
		x := c.ToGoString(prefix)
		if types.IsString(c.Type) {
			x = "url.QueryEscape(" + x + ")"
		}
		url += "+`/`+" + x
	})
	return url
}

func viewDetailColumnString(rels model.Relations, models model.Models, m *model.Model, col *model.Column) []string {
	return lo.FilterMap(rels, func(rel *model.Relation, _ int) (string, bool) {
		relModel := models.Get(rel.Table)
		lCols := rel.SrcColumns(m)
		lNames := strings.Join(lCols.ProperNames(), "")

		relTitles := relModel.Columns.WithTag("title")
		if len(relTitles) == 0 {
			relTitles = relModel.PKs()
		}
		if len(relTitles) == 1 && relTitles[0].Name == col.Name {
			return "", false
		}
		x := fmt.Sprintf("p.%sBy%s", relModel.Proper(), lNames)
		msg := "%s||{%%%% if %s != nil %%%%} ({%%%%s %s.TitleString() %%%%})" + helper.TextEndIfExtra
		return fmt.Sprintf(msg, x, x, x), true
	})
}
