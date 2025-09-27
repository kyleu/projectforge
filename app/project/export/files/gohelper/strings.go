package gohelper

import (
	"fmt"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/util"
)

func BlockString(g *golang.File, cols model.Columns, str model.StringProvider) *golang.Block {
	ret := golang.NewBlock("String", "func")
	ret.WF("func (%s *%s) String() string {", str.FirstLetter(), str.Proper())
	if pks := cols.PKs(); len(pks) == 1 {
		ret.WF("\treturn %s", model.ToGoString(pks[0].Type, pks[0].Nullable, fmt.Sprintf("%s.%s", str.FirstLetter(), pks[0].Proper()), false))
	} else {
		g.AddImport(helper.ImpFmt)
		s := "\treturn fmt.Sprintf(\""
		lo.ForEach(cols.PKs(), func(_ *model.Column, idx int) {
			if idx > 0 {
				s += " • "
			}
			s += "%%s"
		})
		s += "\""
		lo.ForEach(cols.PKs(), func(c *model.Column, _ int) {
			s += helper.TextCommaSpace + c.ToGoString(str.FirstLetter()+".")
		})
		ret.W(s + ")")
	}
	ret.W("}")
	return ret
}

func BlockTitle(g *golang.File, cols model.Columns, str model.StringProvider) *golang.Block {
	ret := golang.NewBlock("Title", "func")
	ret.WF("func (%s *%s) TitleString() string {", str.FirstLetter(), str.Proper())
	titles := cols.WithTag("title")
	var toStrings []string
	switch len(titles) {
	case 0:
		ret.WF("\treturn %s.String()", str.FirstLetter())
	case 1:
		title := titles[0]
		x := model.ToGoString(title.Type, title.Nullable, fmt.Sprintf("%s.%s", str.FirstLetter(), title.Proper()), true)
		if strings.HasPrefix(x, "fmt.") {
			g.AddImport(helper.ImpFmt)
		}
		ret.WF("\tif xx := %s; xx != \"\" {", x)
		ret.W("\t\treturn xx")
		ret.W("\t}")
		ret.WF("\treturn %s.String()", str.FirstLetter())
	default:
		toStrings = lo.Map(titles, func(title *model.Column, _ int) string {
			x := model.ToGoString(title.Type, title.Nullable, fmt.Sprintf("%s.%s", str.FirstLetter(), title.Proper()), true)
			if strings.HasPrefix(x, "fmt.") {
				g.AddImport(helper.ImpFmt)
			}
			return x
		})
		ret.WF("\treturn %s", util.StringJoin(toStrings, " + \" / \" + "))
	}
	ret.W("}")
	return ret
}

func BlockWebPath(g *golang.File, cols model.Columns, str model.StringProvider) *golang.Block {
	ret := golang.NewBlock("WebPath", "type")
	ret.WF("func (%s *%s) WebPath(paths ...string) string {", str.FirstLetter(), str.Proper())
	ret.W("\tif len(paths) == 0 {")
	ret.W("\t\tpaths = []string{DefaultRoute}")
	ret.W("\t}")
	keys := make([]string, 0, len(cols.PKs()))
	lo.ForEach(cols.PKs(), func(pk *model.Column, _ int) {
		g.AddImport(helper.ImpURL)
		const fn = "url.QueryEscape"
		goStr := pk.ToGoString(str.FirstLetter() + ".")
		switch {
		case types.IsStringList(pk.Type):
			keys = append(keys, fmt.Sprintf(fn+`(util.StringJoin(%s, ","))`, goStr))
		case types.IsString(pk.Type) && pk.HasTag("path"):
			g.AddImport(helper.ImpStrings)
			keys = append(keys, fn+"(strings.ReplaceAll("+goStr+`, "/", "||"))`)
		default:
			keys = append(keys, fn+"("+goStr+")")
		}
	})
	ret.WF("\treturn util.StringPath(append(paths, %s)...)", util.StringJoin(keys, ", "))
	ret.W("}")
	return ret
}

func BlockStrings(g *golang.File, cols model.Columns, str model.StringProvider) *golang.Block {
	ret := golang.NewBlock("Strings", "func")
	ret.WF("func (%s *%s) Strings() []string {", str.FirstLetter(), str.Proper())
	x := cols.NotDerived().ToGoStrings(str.FirstLetter()+".", true, 160)
	if strings.Contains(x, "fmt.") {
		g.AddImport(helper.ImpFmt)
	}
	ret.WF("\treturn []string{%s}", x)
	ret.W("}")
	return ret
}

func JSONSuffix(col *model.Column) string {
	if col.HasTag("force-json") {
		return ""
	}
	if types.IsList(col.Type) {
		return ",omitempty"
	}
	if col.Type.Key() == types.KeyReference && util.StringToPlural(col.Name) == col.Name {
		return ",omitempty"
	}
	return ",omitzero"
}

func ColumnTag(col *model.Column) string {
	ret := JSONSuffix(col)
	var tag string
	switch {
	case col.HasTag("ignore-json") || col.Derived():
		tag += "json:\"-\""
	case col.JSON == "":
		tag += fmt.Sprintf("json:%q", col.CamelNoReplace()+ret)
	default:
		tag += fmt.Sprintf("json:%q", col.JSON+ret)
	}
	if col.Validation != "" {
		tag += fmt.Sprintf(",validate:%q", col.Validation)
	}
	if col.Example != "" {
		tag += fmt.Sprintf(",fake:%q", col.Example)
	}
	return "`" + tag + "`"
}
