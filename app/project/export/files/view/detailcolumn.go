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
	ind := util.StringRepeat(ind1, indent)
	rels := m.RelationsFor(col)
	if len(rels) == 0 {
		viewString := col.ToGoViewString(modelKey, true, false, enums, util.KeyDetail)
		ret.W(colRow(ind, col, ModelLinkURL(m, modelKey, enums), viewString, link))
		return
	}

	if types.IsString(col.Type) {
		g.AddImport(helper.ImpURL)
	}

	toStrings := strings.Join(viewDetailColumnString(rels, models, m, col), "")

	ret.W(ind + helper.TextTDStart)
	if col.PK && link {
		cv := col.ToGoViewString(modelKey, true, false, enums, util.KeyDetail)
		ret.W(ind + linkStart + ModelLinkURL(m, modelKey, enums) + `"` + ">" + cv + toStrings + helper.TextEndAnchor)
	} else {
		ret.W(ind + ind1 + col.ToGoViewString(modelKey, true, false, enums, util.KeyDetail) + toStrings)
	}
	lo.ForEach(rels, func(rel *model.Relation, _ int) {
		if lo.Contains(rel.Src, col.Name) {
			switch col.Type.Key() {
			case types.KeyBool, types.KeyInt, types.KeyFloat:
				g.AddImport(helper.ImpFmt)
			}
			relModel := models.Get(rel.Table)
			wp := RelationWebPath(rel, m, relModel, modelKey)
			if col.Nullable {
				ret.W(anchorMsg, ind, modelKey, col.Proper(), relModel.Title(), wp, relModel.Icon)
			} else {
				ret.W(anchorMsgNotNull, ind, relModel.Title(), wp, relModel.Icon)
			}
		}
	})
	ret.W(ind + helper.TextTDEnd)
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
		msg := "{%%%% if p.%sBy%s != nil %%%%} ({%%%%s p.%sBy%s.TitleString() %%%%})" + helper.TextEndIfExtra
		return fmt.Sprintf(msg, relModel.Proper(), lNames, relModel.Proper(), lNames), true
	})
}
