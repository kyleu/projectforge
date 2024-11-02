package gomodel

import (
	"fmt"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
)

const tSet = "Set"

func Models(m *model.Model, args *model.Args, goVersion string, linebreak string) (*file.File, error) {
	name := strings.ToLower(m.CamelPlural())
	if name == strings.ToLower(m.Camel()) {
		name += "_array"
	}
	g := golang.NewFile(m.Package, []string{"app", m.PackageWithGroup("")}, name)
	lo.ForEach(helper.ImportsForTypes("go", "", m.IndexedColumns(true).Types()...), func(imp *model.Import, _ int) {
		g.AddImport(imp)
	})
	lo.ForEach(helper.ImportsForTypes(types.KeyString, "", m.PKs().Types()...), func(imp *model.Import, _ int) {
		g.AddImport(imp)
	})
	g.AddImport(helper.ImpLo, helper.ImpAppUtil)
	g.AddImport(m.Imports.Supporting("array")...)
	g.AddBlocks(modelArray(m))

	imps, err := helper.SpecialImports(m.IndexedColumns(true), m.PackageWithGroup(""), args.Models, args.Enums)
	if err != nil {
		return nil, err
	}
	g.AddImport(imps...)
	ag, err := modelArrayGet(g, m, m.PKs(), args.Enums, goVersion)
	if err != nil {
		return nil, err
	}
	g.AddBlocks(ag)
	lo.ForEach(m.PKs(), func(pk *model.Column, _ int) {
		if pk.Proper() != "Title" {
			g.AddBlocks(modelArrayCol(m, pk, args.Enums))
			g.AddBlocks(modelArrayColStrings(m, pk))
		}
	})
	g.AddBlocks(modelArrayTitleStrings(m))
	if len(m.PKs()) > 1 {
		g.AddBlocks(modelArrayToPKs(m))
	}
	lo.ForEach(m.IndexedColumns(true), func(col *model.Column, _ int) {
		if col.Type.Key() != types.KeyList {
			if !col.PK {
				g.AddBlocks(modelArrayCol(m, col, args.Enums))
			}
			g.AddBlocks(modelArrayGetBySingle(m, col, args.Enums), modelArrayGetByMulti(m, col, args.Enums))
		}
	})
	g.AddBlocks(modelArrayToCSV(m), modelArrayRandom(m), modelArrayClone(m))
	return g.Render(linebreak)
}

func modelArray(m *model.Model) *golang.Block {
	ret := golang.NewBlock(m.Proper()+"Array", "type")
	ret.WF("type %s []*%s", m.ProperPlural(), m.Proper())
	return ret
}

func modelArrayGet(g *golang.File, m *model.Model, cols model.Columns, enums enum.Enums, goVersion string) (*golang.Block, error) {
	ret := golang.NewBlock(m.Proper()+"ArrayGet", "func")
	args, err := cols.Args(m.Package, enums)
	if err != nil {
		return nil, err
	}
	ret.WF("func (%s %s) Get(%s) *%s {", m.FirstLetter(), m.ProperPlural(), args, m.Proper())

	comps := lo.Map(m.PKs(), func(pk *model.Column, _ int) string {
		if types.IsList(pk.Type) {
			g.AddImport(helper.ImpSlicesForGo(goVersion))
			return fmt.Sprintf("slices.Equal(x.%s, %s)", pk.Proper(), pk.Camel())
		}
		return fmt.Sprintf("x.%s == %s", pk.Proper(), pk.Camel())
	})

	ret.WF("\treturn lo.FindOrElse(%s, nil, func(x *%s) bool {", m.FirstLetter(), m.Proper())
	ret.WF("\t\treturn %s", strings.Join(comps, " && "))
	ret.W("\t})")
	ret.W("}")
	return ret, nil
}

func modelArrayRandom(m *model.Model) *golang.Block {
	ret := golang.NewBlock(m.Proper()+"ArrayRandom", "func")
	ret.WF("func (%s %s) Random() *%s {", m.FirstLetter(), m.ProperPlural(), m.Proper())
	ret.WF("\treturn util.RandomElement(%s)", m.FirstLetter())
	ret.W("}")
	return ret
}

func modelArrayClone(m *model.Model) *golang.Block {
	ret := golang.NewBlock(m.Proper()+"ArrayClone", "func")
	ret.WF("func (%s %s) Clone() %s {", m.FirstLetter(), m.ProperPlural(), m.ProperPlural())
	ret.WF("\treturn lo.Map(%s, func(xx *%s, _ int) *%s {", m.FirstLetter(), m.Proper(), m.Proper())
	ret.W("\t\treturn xx.Clone()")
	ret.W("\t})")
	ret.W("}")
	return ret
}

func modelArrayGetByMulti(m *model.Model, col *model.Column, enums enum.Enums) *golang.Block {
	ret := golang.NewBlock(fmt.Sprintf("%sArrayGetBy%s", m.Proper(), col.Proper()), "func")
	name := col.ProperPlural()
	if name == col.Proper() {
		name = col.Proper() + tSet
	}
	if name[len(name)-1] == 'S' {
		name = name[:len(name)-1] + "s"
	}
	t, _ := col.ToGoType(m.Package, enums)
	ret.WF("func (%s %s) GetBy%s(%s ...%s) %s {", m.FirstLetter(), m.ProperPlural(), name, col.CamelPlural(), t, m.ProperPlural())
	ret.WF("\treturn lo.Filter(%s, func(xx *%s, _ int) bool {", m.FirstLetter(), m.Proper())
	ret.WF("\t\treturn lo.Contains(%s, xx.%s)", col.CamelPlural(), col.Proper())
	ret.W("\t})")
	ret.W("}")
	return ret
}

func modelArrayGetBySingle(m *model.Model, col *model.Column, enums enum.Enums) *golang.Block {
	ret := golang.NewBlock(fmt.Sprintf("%sArrayGetBy%s", m.Proper(), col.Proper()), "func")
	t, _ := col.ToGoType(m.Package, enums)
	ret.WF("func (%s %s) GetBy%s(%s %s) %s {", m.FirstLetter(), m.ProperPlural(), col.Proper(), col.Camel(), t, m.ProperPlural())
	ret.WF("\treturn lo.Filter(%s, func(xx *%s, _ int) bool {", m.FirstLetter(), m.Proper())
	ret.WF("\t\treturn xx.%s == %s", col.Proper(), col.Camel())
	ret.W("\t})")
	ret.W("}")
	return ret
}

func modelArrayCol(m *model.Model, col *model.Column, enums enum.Enums) *golang.Block {
	ret := golang.NewBlock(fmt.Sprintf("%sArray%s", m.Proper(), col.ProperPlural()), "func")
	name := col.ProperPlural()
	if name == col.Proper() {
		name = col.Proper() + tSet
	}
	if name[len(name)-1] == 'S' {
		name = name[:len(name)-1] + "s"
	}
	t, _ := col.ToGoType(m.Package, enums)
	ret.WF("func (%s %s) %s() []%s {", m.FirstLetter(), m.ProperPlural(), name, t)
	ret.WF("\treturn lo.Map(%s, func(xx *%s, _ int) %s {", m.FirstLetter(), m.Proper(), t)
	ret.WF("\t\treturn xx.%s", col.Proper())
	ret.W("\t})")
	ret.W("}")
	return ret
}

func modelArrayColStrings(m *model.Model, col *model.Column) *golang.Block {
	ret := golang.NewBlock(fmt.Sprintf("%sArray%sStrings", m.Proper(), col.Proper()), "func")
	ret.WF("func (%s %s) %sStrings(includeNil bool) []string {", m.FirstLetter(), m.ProperPlural(), col.Proper())
	ret.WF("\tret := make([]string, 0, len(%s)+1)", m.FirstLetter())
	ret.W("\tif includeNil {")
	ret.W("\t\tret = append(ret, \"\")")
	ret.W("\t}")
	ret.WF("\tlo.ForEach(%s, func(x *%s, _ int) {", m.FirstLetter(), m.Proper())
	ret.WF("\t\tret = append(ret, %s)", model.ToGoString(col.Type, col.Nullable, "x."+col.Proper(), true))
	ret.W("\t})")
	ret.W("\treturn ret")
	ret.W("}")
	return ret
}

func modelArrayToPKs(m *model.Model) *golang.Block {
	ret := golang.NewBlock(fmt.Sprintf("%sArrayToPKs", m.Proper()), "func")
	ret.WF("func (%s %s) ToPKs() []*PK {", m.FirstLetter(), m.ProperPlural())
	ret.WF("\treturn lo.Map(%s, func(x *%s, _ int) *PK {", m.FirstLetter(), m.Proper())
	ret.W("\t\treturn x.ToPK()")
	ret.W("\t})")
	ret.W("}")
	return ret
}

func modelArrayTitleStrings(m *model.Model) *golang.Block {
	ret := golang.NewBlock(m.Proper()+"ArrayTitleStrings", "func")
	ret.WF("func (%s %s) TitleStrings(nilTitle string) []string {", m.FirstLetter(), m.ProperPlural())
	ret.WF("\tret := make([]string, 0, len(%s)+1)", m.FirstLetter())
	ret.W("\tif nilTitle != \"\" {")
	ret.W("\t\tret = append(ret, nilTitle)")
	ret.W("\t}")
	ret.WF("\tlo.ForEach(%s, func(x *%s, _ int) {", m.FirstLetter(), m.Proper())
	ret.W("\t\tret = append(ret, x.TitleString())")
	ret.W("\t})")
	ret.W("\treturn ret")
	ret.W("}")
	return ret
}
