package view

import (
	"fmt"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/util"
)

var (
	ind1             = "  "
	anchorIcon       = "<a title=%q href=%q>%s</a>"
	anchorMsgNotNull = `%s ` + anchorIcon
	anchorMsg        = "%s {%%%% if %s%s != nil %%%%}" + anchorIcon + helper.TextEndIfExtra
	linkStart        = "  <a href=\""
)

func colRow(ind string, col *model.Column, u string, viewString string, pathKey string, link bool) string {
	ret := viewString
	if col.HasTag("title") {
		ret = "<strong>" + ret + "</strong>"
	}
	if (col.PK || col.HasTag("link")) && link {
		ret = fmt.Sprintf("<a href=%q>%s"+helper.TextEndAnchor, u, ret)
	}
	return ind + "<td>" + ret + helper.TextTDEnd
}

func viewColumn(
	key string, ret *golang.Block, m *model.Model, col *model.Column, call string,
	link bool, modelKey string, indent int, models model.Models, enums enum.Enums, pathKey string,
) {
	ind := util.StringRepeat(ind1, indent)
	rels := m.RelationsFor(col)
	cv := col.ToGoViewString(modelKey, true, false, enums, key)
	if len(rels) == 0 {
		ret.W(colRow(ind, col, ModelLinkURL(m, modelKey, enums), cv, pathKey, link))
		return
	}

	ret.W(ind + helper.TextTDStart)
	ret.W(ind + ind1 + "{%% if x := " + call + "; x != nil %%}")
	strs := []string{"{%%s x.TitleString() %%}"}
	lo.ForEach(rels, func(rel *model.Relation, _ int) {
		if lo.Contains(rel.Src, col.Name) {
			relModel := models.Get(rel.Table)
			icon := fmt.Sprintf("{%%%%= components.SVGLink(`%s`, ps) %%%%}", relModel.Icon)
			if icons := relModel.Columns.WithFormat("icon"); len(icons) == 1 {
				icon = "{%%= components.SVGLink(x." + icons[0].IconDerived() + ", ps) %%}"
			}
			wp := `{%%s x.WebPath(` + pathKey + `...) %%}`
			if col.Nullable && !col.Type.Scalar() {
				strs = append(strs, fmt.Sprintf(anchorMsg, "", modelKey, col.Proper(), relModel.Title(), wp, icon))
			} else {
				strs = append(strs, fmt.Sprintf(anchorMsgNotNull, "", relModel.Title(), wp, icon))
			}
		}
	})
	ret.W(ind + ind1 + strings.Join(strs, ""))
	ret.W(ind + ind1 + "{%% else %%}")
	ret.W(ind + ind1 + cv)
	ret.W(ind + ind1 + "{%% endif %%}")

	ret.W(ind + helper.TextTDEnd)
}
