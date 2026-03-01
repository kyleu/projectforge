package gomodel

import (
	"fmt"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/util"
)

func modelBreadcrumb(m *model.Model) *golang.Block {
	ret := golang.NewBlock("Breadcrumb", "func")
	ret.WF("func (%s *%s) Breadcrumb(extra ...string) string {", m.FirstLetter(), m.Proper())
	l := fmt.Sprintf("\t"+`return %s.TitleString() + "||" + %s.WebPath(extra...)`, m.FirstLetter(), m.FirstLetter())
	if m.Icon != "" {
		l += fmt.Sprintf(` + "**%s"`, m.Icon)
	}
	ret.W(l)
	ret.W("}")
	return ret
}

func modelWebPath(g *golang.File, m *model.Model) *golang.Block {
	ret := golang.NewBlock("WebPath", "type")
	ret.WF("func (%s *%s) WebPath(paths ...string) string {", m.FirstLetter(), m.Proper())
	ret.W("\tif len(paths) == 0 {")
	ret.W("\t\tpaths = []string{DefaultRoute}")
	ret.W("\t}")
	keys := make([]string, 0, len(m.PKs()))
	lo.ForEach(m.PKs(), func(pk *model.Column, _ int) {
		g.AddImport(helper.ImpURL)
		const fn = "url.QueryEscape"
		goStr := pk.ToGoString(m.FirstLetter() + ".")
		switch {
		case types.IsStringList(pk.Type):
			keys = append(keys, fmt.Sprintf(fn+`(util.StringJoin(%s, ","))`, goStr))
		case types.IsString(pk.Type) && pk.HasTag("path"):
			g.AddImport(helper.ImpStrings)
			keys = append(keys, fn+"(strings.ReplaceAll("+goStr+`, "/", "||"))`)
		case pk.Type.Key() == types.KeyTimestamp:
			keys = append(keys, fn+"(util.TimeToJS("+goStr+"))")
		default:
			if _, err := model.AsEnum(pk.Type); err == nil {
				keys = append(keys, fn+"("+goStr+".String())")
			} else {
				keys = append(keys, fn+"("+goStr+")")
			}
		}
	})
	ret.WF("\treturn util.StringPath(append(paths, %s)...)", util.StringJoin(keys, ", "))
	ret.W("}")
	return ret
}
