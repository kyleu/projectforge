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

func viewTableColumn(
	g *golang.Template, ret *golang.Block, models model.Models, m *model.Model, link bool,
	col *model.Column, modelKey string, prefix string, indent int, enums enum.Enums,
) {
	ind := util.StringRepeat(viewIndent, indent)
	rels := m.RelationsFor(col)
	if len(rels) == 0 {
		viewString := col.ToGoViewString(modelKey, true, false, enums, "table")
		ret.W(colRow(ind, col, m.LinkURL(modelKey, enums), viewString, link))
		return
	}

	if types.IsString(col.Type) {
		g.AddImport(helper.ImpURL)
	}

	toStrings := getTableColumnStrings(m, modelKey, rels, models, prefix)

	ret.W(ind + "<td class=\"nowrap\">")
	if col.PK && link {
		ret.W(ind + "  <a href=\"" + m.LinkURL(modelKey, enums) + "\">" + col.ToGoViewString(modelKey, true, false, enums, "table") + toStrings + "</a>")
	} else {
		ret.W(ind + viewIndent + col.ToGoViewString(modelKey, true, false, enums, "table") + toStrings)
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
			if col.Nullable {
				ret.W(msg, ind, modelKey, col.Proper(), relModel.Title(), rel.WebPath(m, relModel, modelKey), relModel.Icon)
			} else {
				ret.W(msgNotNull, ind, relModel.Title(), rel.WebPath(m, relModel, modelKey), relModel.Icon)
			}
		}
	})
	ret.W(ind + "</td>")
}

func getTableColumnStrings(m *model.Model, modelKey string, rels model.Relations, models model.Models, prefix string) string {
	ret := lo.FilterMap(rels, func(rel *model.Relation, idx int) (string, bool) {
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
		if srcCol.Nullable && !tgtCol.Nullable {
			get := fmt.Sprintf("%sBy%s.Get(*%s%s)", k, srcCol.Proper(), modelKey, srcCol.Proper())
			return "{%% if " + modelKey + srcCol.Proper() + " != nil %%}" +
				"{%% if x := " + get + "; x != nil %%} ({%%s x.TitleString() %%}){%% endif %%}" +
				"{%% endif %%}", true
		}
		get := fmt.Sprintf("%sBy%s.Get(%s%s)", k, srcCol.Proper(), modelKey, srcCol.Proper())
		return "{%% if x := " + get + "; x != nil %%} ({%%s x.TitleString() %%}){%% endif %%}", true
	})
	return strings.Join(ret, "")
}
