package gomodel

import (
	"fmt"
	"strings"

	"github.com/samber/lo"

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
	lo.ForEach(helper.ImportsForTypes("go", "", m.PKs().Types()...), func(imp *golang.Import, _ int) {
		g.AddImport(imp)
	})
	lo.ForEach(helper.ImportsForTypes("string", "", m.PKs().Types()...), func(imp *golang.Import, _ int) {
		g.AddImport(imp)
	})
	g.AddImport(helper.ImpSlices)
	g.AddBlocks(modelArray(m))
	ag, err := modelArrayGet(g, m, args.Enums)
	if err != nil {
		return nil, err
	}
	g.AddBlocks(ag)
	lo.ForEach(m.PKs(), func(pk *model.Column, _ int) {
		if pk.Proper() != "Title" {
			if pk.Type.Key() != types.KeyList {
				g.AddBlocks(modelArrayGetBy(m, pk, args.Enums))
			}
			g.AddBlocks(modelArrayCol(m, pk, args.Enums))
			g.AddBlocks(modelArrayColStrings(m, pk))
		}
	})
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
	lo.ForEach(m.PKs(), func(pk *model.Column, _ int) {
		if types.IsList(pk.Type) {
			g.AddImport(helper.ImpSlices)
			comps = append(comps, fmt.Sprintf("slices.Equal(x.%s, %s)", pk.Proper(), pk.Camel()))
		} else {
			comps = append(comps, fmt.Sprintf("x.%s == %s", pk.Proper(), pk.Camel()))
		}
	})
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

func modelArrayGetBy(m *model.Model, col *model.Column, enums enum.Enums) *golang.Block {
	ret := golang.NewBlock(fmt.Sprintf("%sArrayGetBy%s", m.Proper(), col.Proper()), "func")
	name := col.ProperPlural()
	if name == col.Proper() {
		name = col.Proper() + "Set"
	}
	if name[len(name)-1] == 'S' {
		name = name[:len(name)-1] + "s"
	}
	t, _ := col.ToGoType(m.Package, enums)
	ret.W("func (%s %s) GetBy%s(%s ...%s) %s {", m.FirstLetter(), m.ProperPlural(), name, col.CamelPlural(), t, m.ProperPlural())
	ret.W("\tvar ret %s", m.ProperPlural())
	ret.W("\tfor _, x := range %s {", m.FirstLetter())
	ret.W("\t\tif slices.Contains(%s, x.%s) {", col.CamelPlural(), col.Proper())
	ret.W("\t\t\tret = append(ret, x)")
	ret.W("\t\t}")
	ret.W("\t}")
	ret.W("\treturn ret")
	ret.W("}")
	return ret
}

func modelArrayCol(m *model.Model, col *model.Column, enums enum.Enums) *golang.Block {
	ret := golang.NewBlock(fmt.Sprintf("%sArray%s", m.Proper(), col.ProperPlural()), "func")
	name := col.ProperPlural()
	if name == col.Proper() {
		name = col.Proper() + "Set"
	}
	if name[len(name)-1] == 'S' {
		name = name[:len(name)-1] + "s"
	}
	t, _ := col.ToGoType(m.Package, enums)
	ret.W("func (%s %s) %s() []%s {", m.FirstLetter(), m.ProperPlural(), name, t)
	ret.W("\tret := make([]%s, 0, len(%s)+1)", t, m.FirstLetter())
	ret.W("\tfor _, x := range %s {", m.FirstLetter())
	ret.W("\t\tret = append(ret, x.%s)", col.Proper())
	ret.W("\t}")
	ret.W("\treturn ret")
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
