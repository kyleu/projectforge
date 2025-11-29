package view

import (
	"fmt"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/util"
)

const (
	ind1       = "  "
	anchorIcon = "<a title=%q href=%q>%s</a>"
	anchorMsg  = `%s ` + anchorIcon
	titleStr   = "{%%s x.TitleString() %%}"
)

func colRow(ind string, col *model.Column, u string, viewString string, link bool) string {
	ret := viewString
	if col.HasTag("title") {
		ret = "<strong>" + ret + "</strong>"
	}
	if link {
		if col.PK || col.HasTag("link") || col.HasTag("title") {
			ret = fmt.Sprintf("<a href=%q>%s"+helper.TextEndAnchor, u, ret)
		}
	}
	return ind + "<td>" + ret + helper.TextTDEnd
}

func viewColumn(
	key string, ret *golang.Block, m *model.Model, col *model.Column, call string,
	link bool, modelKey string, indent int, models model.Models, enums enum.Enums, chk bool, pathKey string,
) {
	ind := util.StringRepeat(ind1, indent)
	modelLinkStart := `<a href="{%%s ` + modelKey + `WebPath() %%}">`
	rels := m.RelationsFor(col)
	cv := col.ToGoViewString(modelKey, true, false, enums, key)
	if len(rels) == 0 {
		ret.W(colRow(ind, col, ModelLinkURL(modelKey), cv, link))
		return
	}

	ret.W(ind + helper.TextTDStart)
	if chk {
		ret.W(ind + ind1 + "{%% if " + modelKey + col.Proper() + " != " + col.ZeroVal() + " %%}{%% if x := " + call + "; x != nil %%}")
	} else {
		ret.W(ind + ind1 + "{%% if x := " + call + "; x != nil %%}")
	}
	strs := []string{}
	if col.PK {
		strs = append(strs, modelLinkStart, titleStr, "</a>")
	} else {
		strs = append(strs, titleStr)
	}

	lo.ForEach(rels, func(rel *model.Relation, _ int) {
		if lo.Contains(rel.Src, col.Name) {
			relModel := models.Get(rel.Table)
			icon := fmt.Sprintf("{%%%%= components.SVGLink(`%s`, ps) %%%%}", relModel.Icon)
			if icons := relModel.Columns.WithFormat("icon"); len(icons) == 1 {
				icon = "{%%= components.SVGLink(x." + icons[0].IconDerived() + ", ps) %%}"
			}
			wp := `{%%s x.WebPath(` + pathKey + `...) %%}`
			strs = append(strs, fmt.Sprintf(anchorMsg, "", relModel.Title(), wp, icon))
		}
	})
	ret.W(ind + ind1 + util.StringJoin(strs, ""))
	ret.W(ind + ind1 + "{%% else %%}")
	if col.PK {
		ret.W(ind + ind1 + modelLinkStart + cv + "</a>")
	} else {
		ret.W(ind + ind1 + cv)
	}
	if chk {
		ret.W(ind + ind1 + "{%% endif %%}{%% endif %%}")
	} else {
		ret.W(ind + ind1 + "{%% endif %%}")
	}
	ret.W(ind + helper.TextTDEnd)
}
