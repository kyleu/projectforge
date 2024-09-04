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
	anchorMsgNotNull = `%s  ` + anchorIcon
	anchorMsg        = "%s  {%%%% if %s%s != nil %%%%}" + anchorIcon + helper.TextEndIfExtra
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
	key string, g *golang.Template, ret *golang.Block,
	m *model.Model, col *model.Column,
	toStrings string, link bool, modelKey string,
	indent int, models model.Models, enums enum.Enums, pathKey string,
) {
	prefix := ""
	if idx := strings.Index(toStrings, "||"); idx > -1 {
		prefix = toStrings[:idx]
		toStrings = toStrings[idx+2:]
	}
	ind := util.StringRepeat(ind1, indent)
	rels := m.RelationsFor(col)
	if len(rels) == 0 {
		cv := col.ToGoViewString(modelKey, true, false, enums, key)
		ret.W(colRow(ind, col, ModelLinkURL(m, modelKey, enums), cv, pathKey, link))
		return
	}

	ret.W(ind + helper.TextTDStart)
	if col.PK && link {
		cv := col.ToGoViewString(modelKey, true, false, enums, key)
		ret.W(ind + linkStart + ModelLinkURL(m, modelKey, enums) + "\">" + cv + toStrings + helper.TextEndAnchor)
	} else {
		cv := col.ToGoViewString(modelKey, true, false, enums, key)
		ret.W(ind + ind1 + cv + toStrings)
	}

	lo.ForEach(rels, func(rel *model.Relation, _ int) {
		if lo.Contains(rel.Src, col.Name) {
			relModel := models.Get(rel.Table)
			icon := fmt.Sprintf("{%%%%= components.SVGLink(`%s`, ps) %%%%}", relModel.Icon)
			if icons := relModel.Columns.WithFormat("icon"); len(icons) == 1 {
				msg := `{%%%% if x := %s; x != nil %%%%}{%%%%= components.SVGLink(x.%s, ps) %%%%}{%%%% else %%%%}%s{%%%% endif %%%%}`
				icon = fmt.Sprintf(msg, prefix, icons[0].ProperDerived(), icon)
			}
			if prefix == "" {
				println()
			}
			wp := `{%% if x := ` + prefix + `; x != nil %%}{%%s x.WebPath(` + pathKey + `...) %%}{%% endif %%}`
			if col.Nullable {
				ret.W(anchorMsg, ind, modelKey, col.Proper(), relModel.Title(), wp, icon)
			} else {
				ret.W(anchorMsgNotNull, ind, relModel.Title(), wp, icon)
			}
		}
	})
	ret.W(ind + helper.TextTDEnd)
}
