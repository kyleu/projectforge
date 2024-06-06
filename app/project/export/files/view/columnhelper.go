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

var (
	ind1             = "  "
	anchorIcon       = "<a title=%q href=\"{%%%%s %s %%%%}\">%s</a>"
	anchorMsgNotNull = `%s  ` + anchorIcon
	anchorMsg        = "%s  {%%%% if %s%s != nil %%%%}" + anchorIcon + helper.TextEndIfExtra
	anchorSVG        = "{%%%%= components.SVGSimple(%s, 18, ps) %%%%}"
	linkStart        = "  <a href=\""
)

func colRow(ind string, col *model.Column, u string, viewString string, link bool) string {
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
	indent int, models model.Models, enums enum.Enums,
) {
	prefix := ""
	if idx := strings.Index(toStrings, "||"); idx > -1 {
		prefix = toStrings[:idx]
		toStrings = toStrings[idx+2:]
	}
	ind := util.StringRepeat(ind1, indent)
	rels := m.RelationsFor(col)
	if len(rels) == 0 {
		viewString := col.ToGoViewString(modelKey, true, false, enums, key)
		ret.W(colRow(ind, col, ModelLinkURL(m, modelKey, enums), viewString, link))
		return
	}

	if types.IsString(col.Type) {
		g.AddImport(helper.ImpURL)
	}

	ret.W(ind + helper.TextTDStart)
	if col.PK && link {
		cv := col.ToGoViewString(modelKey, true, false, enums, key)
		ret.W(ind + linkStart + ModelLinkURL(m, modelKey, enums) + "\">" + cv + toStrings + helper.TextEndAnchor)
	} else {
		ret.W(ind + ind1 + col.ToGoViewString(modelKey, true, false, enums, key) + toStrings)
	}

	lo.ForEach(rels, func(rel *model.Relation, _ int) {
		if lo.Contains(rel.Src, col.Name) {
			switch col.Type.Key() {
			case types.KeyBool, types.KeyInt, types.KeyFloat:
				g.AddImport(helper.ImpFmt)
			}
			relModel := models.Get(rel.Table)
			icon := fmt.Sprintf(anchorSVG, `"`+relModel.Icon+`"`)
			if icons := relModel.Columns.WithFormat("icon"); len(icons) == 1 {
				msg := `{%%%% if x := %s; x != nil %%%%}{%%%%= components.SVGSimple(x.%s, 18, ps) %%%%}{%%%% else %%%%}%s{%%%% endif %%%%}`
				icon = fmt.Sprintf(msg, prefix, icons[0].ProperDerived(), icon)
			}
			wp := RelationWebPath(rel, m, relModel, modelKey)
			if col.Nullable {
				ret.W(anchorMsg, ind, modelKey, col.Proper(), relModel.Title(), wp, icon)
			} else {
				ret.W(anchorMsgNotNull, ind, relModel.Title(), wp, icon)
			}
		}
	})
	ret.W(ind + helper.TextTDEnd)
}
