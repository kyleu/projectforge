package gomodel

import (
	"fmt"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
)

func modelString(g *golang.File, m *model.Model) *golang.Block {
	ret := golang.NewBlock("String", "func")
	ret.WF("func (%s *%s) String() string {", m.FirstLetter(), m.Proper())
	if pks := m.PKs(); len(pks) == 1 {
		ret.WF("\treturn %s", model.ToGoString(pks[0].Type, pks[0].Nullable, fmt.Sprintf("%s.%s", m.FirstLetter(), pks[0].Proper()), false))
	} else {
		g.AddImport(helper.ImpFmt)
		s := "\treturn fmt.Sprintf(\""
		lo.ForEach(m.PKs(), func(_ *model.Column, idx int) {
			if idx > 0 {
				s += " • "
			}
			s += "%%s"
		})
		s += "\""
		lo.ForEach(m.PKs(), func(c *model.Column, _ int) {
			s += helper.TextCommaSpace + c.ToGoString(m.FirstLetter()+".")
		})
		ret.W(s + ")")
	}
	ret.W("}")
	return ret
}

func modelTitle(g *golang.File, m *model.Model) *golang.Block {
	ret := golang.NewBlock("Title", "func")
	ret.WF("func (%s *%s) TitleString() string {", m.FirstLetter(), m.Proper())
	titles := m.Columns.WithTag("title")
	var toStrings []string
	switch len(titles) {
	case 0:
		ret.WF("\treturn %s.String()", m.FirstLetter())
	case 1:
		title := titles[0]
		x := model.ToGoString(title.Type, title.Nullable, fmt.Sprintf("%s.%s", m.FirstLetter(), title.Proper()), true)
		if strings.HasPrefix(x, "fmt.") {
			g.AddImport(helper.ImpFmt)
		}
		ret.WF("\tif xx := %s; xx != \"\" {", x)
		ret.W("\t\treturn xx")
		ret.W("\t}")
		ret.WF("\treturn %s.String()", m.FirstLetter())
	default:
		toStrings = lo.Map(titles, func(title *model.Column, _ int) string {
			x := model.ToGoString(title.Type, title.Nullable, fmt.Sprintf("%s.%s", m.FirstLetter(), title.Proper()), true)
			if strings.HasPrefix(x, "fmt.") {
				g.AddImport(helper.ImpFmt)
			}
			return x
		})
		ret.WF("\treturn %s", strings.Join(toStrings, " + \" / \" + "))
	}
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
			g.AddImport(helper.ImpStrings)
			keys = append(keys, fmt.Sprintf(fn+`(strings.Join(%s, ","))`, goStr))
		case types.IsString(pk.Type) && pk.HasTag("path"):
			g.AddImport(helper.ImpStrings)
			keys = append(keys, fn+"(strings.ReplaceAll("+goStr+`, "/", "||"))`)
		default:
			keys = append(keys, fn+"("+goStr+")")
		}
	})
	ret.WF("\treturn path.Join(append(paths, %s)...)", strings.Join(keys, ", "))
	ret.W("}")
	return ret
}
