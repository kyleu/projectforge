package view

import (
	"fmt"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/enum"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/project/export/model"
	"projectforge.dev/projectforge/app/util"
)

const viewIndent = "  "

func viewDetailColumn(
	g *golang.Template, ret *golang.Block, models model.Models, m *model.Model, link bool, col *model.Column, modelKey string, indent int, enums enum.Enums,
) {
	ind := util.StringRepeat(viewIndent, indent)
	rels := m.RelationsFor(col)
	if len(rels) == 0 {
		viewString := col.ToGoViewString(modelKey, true, false, enums, util.KeyDetail)
		ret.W(colRow(ind, col, m.LinkURL(modelKey, enums), viewString, link))
		return
	}

	if types.IsString(col.Type) {
		g.AddImport(helper.ImpURL)
	}

	toStrings := strings.Join(viewDetailColumnString(rels, models, m, col), "")

	ret.W(ind + "<td class=\"nowrap\">")
	if col.PK && link {
		ret.W(ind + "  <a href=\"" + m.LinkURL(modelKey, enums) + "\">" + col.ToGoViewString(modelKey, true, false, enums, util.KeyDetail) + toStrings + "</a>")
	} else {
		ret.W(ind + viewIndent + col.ToGoViewString(modelKey, true, false, enums, util.KeyDetail) + toStrings)
	}
	const l = "<a title=%q href=\"{%%%%s %s %%%%}\">{%%%%= components.SVGRef(%q, 18, 18, \"\", ps) %%%%}</a>"
	const msgNotNull = "%s  " + l
	const msg = "%s  {%%%% if %s%s != nil %%%%}" + l + "{%%%% endif %%%%}"
	lo.ForEach(rels, func(rel *model.Relation, _ int) {
		if lo.Contains(rel.Src, col.Name) {
			switch col.Type.Key() {
			case types.KeyBool, types.KeyInt, types.KeyFloat:
				g.AddImport(helper.ImpFmt)
			}
			relModel := models.Get(rel.Table)
			wp := rel.WebPath(m, relModel, modelKey)
			if col.Nullable {
				ret.W(msg, ind, modelKey, col.Proper(), relModel.Title(), wp, relModel.Icon)
			} else {
				ret.W(msgNotNull, ind, relModel.Title(), wp, relModel.Icon)
			}
		}
	})
	ret.W(ind + "</td>")
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
		msg := "{%%%% if p.%sBy%s != nil %%%%} ({%%%%s p.%sBy%s.TitleString() %%%%}){%%%% endif %%%%}"
		return fmt.Sprintf(msg, relModel.Proper(), lNames, relModel.Proper(), lNames), true
	})
}
