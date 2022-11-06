package gomodel

import (
	"fmt"
	"strings"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/enum"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/project/export/model"
)

func Models(m *model.Model, args *model.Args, addHeader bool) (*file.File, error) {
	name := strings.ToLower(m.CamelPlural())
	if name == strings.ToLower(m.Camel()) {
		name += "_array"
	}
	g := golang.NewFile(m.Package, []string{"app", m.PackageWithGroup("")}, name)
	for _, imp := range helper.ImportsForTypes("go", m.PKs().Types()...) {
		g.AddImport(imp)
	}
	for _, imp := range helper.ImportsForTypes("string", m.PKs().Types()...) {
		g.AddImport(imp)
	}
	g.AddImport(helper.ImpSlices)
	g.AddBlocks(modelArray(m))
	ag, err := modelArrayGet(g, m, args.Enums)
	if err != nil {
		return nil, err
	}
	g.AddBlocks(ag)
	for _, pk := range m.PKs() {
		if pk.Proper() != "Title" {
			g.AddBlocks(modelArrayColStrings(m, pk))
		}
	}
	g.AddBlocks(modelArrayTitleStrings(m), modelArrayClone(m))
	return g.Render(addHeader)
}

func modelArray(m *model.Model) *golang.Block {
	ret := golang.NewBlock(m.Proper()+"Array", "type")
	ret.W("type %s []*%s", m.ProperPlural(), m.Proper())
	return ret
}

func modelArrayGet(g *golang.File, m *model.Model, enums enum.Enums) (*golang.Block, error) {
	ret := golang.NewBlock(m.Proper()+"ArrayGet", "func")
	args, err := m.PKs().Args(m.Package, enums)
	if err != nil {
		return nil, err
	}
	ret.W("func (%s %s) Get(%s) *%s {", m.FirstLetter(), m.ProperPlural(), args, m.Proper())
	ret.W("\tfor _, x := range %s {", m.FirstLetter())
	comps := make([]string, 0, len(m.PKs()))
	for _, pk := range m.PKs() {
		if types.IsList(pk.Type) {
			g.AddImport(helper.ImpSlices)
			comps = append(comps, fmt.Sprintf("slices.Equal(x.%s, %s)", pk.Proper(), pk.Camel()))
		} else {
			comps = append(comps, fmt.Sprintf("x.%s == %s", pk.Proper(), pk.Camel()))
		}
	}
	ret.W("\t\tif %s {", strings.Join(comps, " && "))
	ret.W("\t\t\treturn x")
	ret.W("\t\t}")
	ret.W("\t}")
	ret.W("\treturn nil")
	ret.W("}")
	return ret, nil
}

func modelArrayClone(m *model.Model) *golang.Block {
	ret := golang.NewBlock(m.Proper()+"ArrayClone", "func")
	ret.W("func (%s %s) Clone() %s {", m.FirstLetter(), m.ProperPlural(), m.ProperPlural())
	ret.W("\treturn slices.Clone(%s)", m.FirstLetter())
	ret.W("}")
	return ret
}

func modelArrayColStrings(m *model.Model, col *model.Column) *golang.Block {
	ret := golang.NewBlock(fmt.Sprintf("%sArray%sStrings", m.Proper(), col.Proper()), "func")
	ret.W("func (%s %s) %sStrings(includeNil bool) []string {", m.FirstLetter(), m.ProperPlural(), col.Proper())
	ret.W("\tret := make([]string, 0, len(%s)+1)", m.FirstLetter())
	ret.W("\tif includeNil {")
	ret.W("\t\tret = append(ret, \"\")")
	ret.W("\t}")
	ret.W("\tfor _, x := range %s {", m.FirstLetter())
	ret.W("\t\tret = append(ret, %s)", model.ToGoString(col.Type, "x."+col.Proper(), true))
	ret.W("\t}")
	ret.W("\treturn ret")
	ret.W("}")
	return ret
}

func modelArrayTitleStrings(m *model.Model) *golang.Block {
	ret := golang.NewBlock(m.Proper()+"ArrayTitleStrings", "func")
	ret.W("func (%s %s) TitleStrings(nilTitle string) []string {", m.FirstLetter(), m.ProperPlural())
	ret.W("\tret := make([]string, 0, len(%s)+1)", m.FirstLetter())
	ret.W("\tif nilTitle != \"\" {")
	ret.W("\t\tret = append(ret, nilTitle)")
	ret.W("\t}")
	ret.W("\tfor _, x := range %s {", m.FirstLetter())
	ret.W("\t\tret = append(ret, x.TitleString())")
	ret.W("\t}")
	ret.W("\treturn ret")
	ret.W("}")
	return ret
}
