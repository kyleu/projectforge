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
	ret.W("func (%s *%s) String() string {", m.FirstLetter(), m.Proper())
	if pks := m.PKs(); len(pks) == 1 {
		ret.W("\treturn %s", model.ToGoString(pks[0].Type, pks[0].Nullable, fmt.Sprintf("%s.%s", m.FirstLetter(), pks[0].Proper()), false))
	} else {
		g.AddImport(helper.ImpFmt)
		s := "\treturn fmt.Sprintf(\""
		lo.ForEach(m.PKs(), func(_ *model.Column, idx int) {
			if idx > 0 {
				s += "::"
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
	ret.W("func (%s *%s) TitleString() string {", m.FirstLetter(), m.Proper())
	titles := m.Columns.WithTag("title")
	var toStrings []string
	switch len(titles) {
	case 0:
		ret.W("\treturn %s.String()", m.FirstLetter())
	case 1:
		title := titles[0]
		x := model.ToGoString(title.Type, title.Nullable, fmt.Sprintf("%s.%s", m.FirstLetter(), title.Proper()), true)
		if strings.HasPrefix(x, "fmt.") {
			g.AddImport(helper.ImpFmt)
		}
		ret.W("\tif xx := %s; xx != \"\" {", x)
		ret.W("\t\treturn xx")
		ret.W("\t}")
		ret.W("\treturn %s.String()", m.FirstLetter())
	default:
		toStrings = lo.Map(titles, func(title *model.Column, _ int) string {
			x := model.ToGoString(title.Type, title.Nullable, fmt.Sprintf("%s.%s", m.FirstLetter(), title.Proper()), true)
			if strings.HasPrefix(x, "fmt.") {
				g.AddImport(helper.ImpFmt)
			}
			return x
		})
		ret.W("\treturn %s", strings.Join(toStrings, " + \" / \" + "))
	}
	ret.W("}")
	return ret
}

func modelWebPath(g *golang.File, m *model.Model) *golang.Block {
	ret := golang.NewBlock("WebPath", "type")
	ret.W("func (%s *%s) WebPath(paths ...string) string {", m.FirstLetter(), m.Proper())
	ret.W("\tif len(paths) == 0 {")
	ret.W("\t\tpaths = []string{DefaultRoute}")
	ret.W("\t}")
	keys := make([]string, 0, len(m.PKs()))
	lo.ForEach(m.PKs(), func(pk *model.Column, _ int) {
		g.AddImport(helper.ImpURL)
		const fn = "url.QueryEscape"
		switch {
		case types.IsStringList(pk.Type):
			g.AddImport(helper.ImpStrings)
			keys = append(keys, fmt.Sprintf(fn+`(strings.Join(%s, ","))`, pk.ToGoString(m.FirstLetter()+".")))
		default:
			keys = append(keys, fn+"("+pk.ToGoString(m.FirstLetter()+".")+")")
		}
	})
	ret.W("\treturn path.Join(append(paths, %s)...)", strings.Join(keys, ", "))
	ret.W("}")
	return ret
}
