package gohelper

import (
	"fmt"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/metamodel"
	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/util"
)

const tSet = "Set"

func Array(
	g *golang.File, m model.StringProvider, cols model.Columns, importantCols model.Columns, args *metamodel.Args, goVersion string, linebreak string,
) (*file.File, error) {
	g.AddImport(helper.ImpLo, helper.ImpAppUtil)
	imps, err := helper.SpecialImports(importantCols, m.PackageWithGroup(""), args)
	if err != nil {
		return nil, err
	}
	g.AddImport(imps...)

	g.AddBlocks(ModelArray(m))
	if len(cols.PKs()) > 0 {
		ag, err := ModelArrayGet(g, m, cols.PKs(), args, goVersion)
		if err != nil {
			return nil, err
		}
		g.AddBlocks(ag)
	}
	lo.ForEach(cols.PKs(), func(pk *model.Column, _ int) {
		if pk.Proper() != "Title" {
			g.AddBlocks(modelArrayCol(m, pk, args.Enums))
			g.AddBlocks(modelArrayColStrings(m, pk, args.Enums))
		}
	})
	g.AddBlocks(modelArrayTitleStrings(m))
	if len(cols.PKs()) > 1 {
		g.AddBlocks(modelArrayToPKs(m))
	}
	lo.ForEach(importantCols, func(col *model.Column, _ int) {
		if col.Type.Key() != types.KeyList {
			if !col.PK {
				g.AddBlocks(modelArrayCol(m, col, args.Enums))
			}
			g.AddBlocks(modelArrayGetBySingle(m, col, args.Enums), modelArrayGetByMulti(m, col, args.Enums))
		}
	})
	g.AddBlocks(ModelArrayToMaps(m), ModelArrayToOrderedMaps(m), BlockArrayToCSV(m), BlockArrayRandom(m), BlockArrayClone(m))
	return g.Render(linebreak)
}

func ModelArray(m model.StringProvider) *golang.Block {
	ret := golang.NewBlock(m.Proper()+"Array", "type")
	ret.WF("type %s []*%s", m.ProperPlural(), m.Proper())
	return ret
}

func ModelArrayGet(g *golang.File, m model.StringProvider, cols model.Columns, x *metamodel.Args, goVersion string) (*golang.Block, error) {
	ret := golang.NewBlock(m.Proper()+"ArrayGet", "func")
	argsString, err := helper.GoArgsWithRef(cols, m.PackageName(), x)
	if err != nil {
		return nil, err
	}
	ret.WF("func (%s %s) Get(%s) *%s {", m.FirstLetter(), m.ProperPlural(), argsString, m.Proper())

	comps := lo.Map(cols.PKs(), func(pk *model.Column, _ int) string {
		if types.IsList(pk.Type) {
			g.AddImport(helper.ImpSlices)
			return fmt.Sprintf("slices.Equal(x.%s, %s)", pk.Proper(), pk.Camel())
		}
		return fmt.Sprintf("x.%s == %s", pk.Proper(), pk.Camel())
	})

	ret.WF("\treturn lo.FindOrElse(%s, nil, func(x *%s) bool {", m.FirstLetter(), m.Proper())
	ret.WF("\t\treturn %s", util.StringJoin(comps, " && "))
	ret.W("\t})")
	ret.W("}")
	return ret, nil
}

func ModelArrayToMaps(m model.StringProvider) *golang.Block {
	ret := golang.NewBlock(m.Proper()+"ArrayToMaps", "func")
	ret.WF("func (%s %s) ToMaps() []util.ValueMap {", m.FirstLetter(), m.ProperPlural())
	ret.WF("\treturn lo.Map(%s, func(x *%s, _ int) util.ValueMap {", m.FirstLetter(), m.Proper())
	ret.W("\t\treturn x.ToMap()")
	ret.W("\t})")
	ret.W("}")
	return ret
}

func ModelArrayToOrderedMaps(m model.StringProvider) *golang.Block {
	ret := golang.NewBlock(m.Proper()+"ArrayToOrderedMaps", "func")
	ret.WF("func (%s %s) ToOrderedMaps() util.OrderedMaps[any] {", m.FirstLetter(), m.ProperPlural())
	ret.WF("\treturn lo.Map(%s, func(x *%s, _ int) *util.OrderedMap[any] {", m.FirstLetter(), m.Proper())
	ret.W("\t\treturn x.ToOrderedMap()")
	ret.W("\t})")
	ret.W("}")
	return ret
}

func modelArrayGetByMulti(m model.StringProvider, col *model.Column, enums enum.Enums) *golang.Block {
	ret := golang.NewBlock(fmt.Sprintf("%sArrayGetBy%s", m.Proper(), col.Proper()), "func")
	name := col.ProperPlural()
	if name == col.Proper() {
		name = col.Proper() + tSet
	}
	if name[len(name)-1] == 'S' {
		name = name[:len(name)-1] + "s"
	}
	t, _ := col.ToGoType(m.PackageName(), enums)
	ret.WF("func (%s %s) GetBy%s(%s ...%s) %s {", m.FirstLetter(), m.ProperPlural(), name, col.CamelPlural(), t, m.ProperPlural())
	ret.WF("\treturn lo.Filter(%s, func(xx *%s, _ int) bool {", m.FirstLetter(), m.Proper())
	ret.WF("\t\treturn lo.Contains(%s, xx.%s)", col.CamelPlural(), col.Proper())
	ret.W("\t})")
	ret.W("}")
	return ret
}

func modelArrayGetBySingle(m model.StringProvider, col *model.Column, enums enum.Enums) *golang.Block {
	ret := golang.NewBlock(fmt.Sprintf("%sArrayGetBy%s", m.Proper(), col.Proper()), "func")
	t, _ := col.ToGoType(m.PackageName(), enums)
	ret.WF("func (%s %s) GetBy%s(%s %s) %s {", m.FirstLetter(), m.ProperPlural(), col.Proper(), col.Camel(), t, m.ProperPlural())
	ret.WF("\treturn lo.Filter(%s, func(xx *%s, _ int) bool {", m.FirstLetter(), m.Proper())
	ret.WF("\t\treturn xx.%s == %s", col.Proper(), col.Camel())
	ret.W("\t})")
	ret.W("}")
	return ret
}

func modelArrayCol(m model.StringProvider, col *model.Column, enums enum.Enums) *golang.Block {
	ret := golang.NewBlock(fmt.Sprintf("%sArray%s", m.Proper(), col.ProperPlural()), "func")
	name := col.ProperPlural()
	if name == col.Proper() {
		name = col.Proper() + tSet
	}
	if name[len(name)-1] == 'S' {
		name = name[:len(name)-1] + "s"
	}
	t, _ := col.ToGoType(m.PackageName(), enums)
	ret.WF("func (%s %s) %s() []%s {", m.FirstLetter(), m.ProperPlural(), name, t)
	ret.WF("\treturn lo.Map(%s, func(xx *%s, _ int) %s {", m.FirstLetter(), m.Proper(), t)
	ret.WF("\t\treturn xx.%s", col.Proper())
	ret.W("\t})")
	ret.W("}")
	return ret
}

func modelArrayColStrings(m model.StringProvider, col *model.Column, enums enum.Enums) *golang.Block {
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

func modelArrayToPKs(m model.StringProvider) *golang.Block {
	ret := golang.NewBlock(fmt.Sprintf("%sArrayToPKs", m.Proper()), "func")
	ret.WF("func (%s %s) ToPKs() []*PK {", m.FirstLetter(), m.ProperPlural())
	ret.WF("\treturn lo.Map(%s, func(x *%s, _ int) *PK {", m.FirstLetter(), m.Proper())
	ret.W("\t\treturn x.ToPK()")
	ret.W("\t})")
	ret.W("}")
	return ret
}

func modelArrayTitleStrings(m model.StringProvider) *golang.Block {
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
