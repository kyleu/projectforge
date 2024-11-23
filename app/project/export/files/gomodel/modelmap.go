package gomodel

import (
	"fmt"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
)

func ModelMap(m *model.Model, args *model.Args, linebreak string) (*file.File, error) {
	g := golang.NewFile(m.Package, []string{"app", m.PackageWithGroup("")}, strings.ToLower(m.Camel())+"map")
	g.AddImport(helper.ImpAppUtil)
	imps, err := helper.SpecialImports(m.Columns, m.PackageWithGroup(""), args.Models, args.Enums)
	if err != nil {
		return nil, err
	}
	g.AddImport(imps...)
	imps, err = helper.EnumImports(m.Columns.Types(), m.PackageWithGroup(""), args.Models, args.Enums)
	if err != nil {
		return nil, err
	}
	g.AddImport(imps...)
	g.AddImport(m.Imports.Supporting("map")...)
	g.AddBlocks(modelToMap(g, m, args.Enums, args.Database))
	if b, e := modelFromMap(g, m, args.Models, args.Enums, args.Database); e == nil {
		g.AddBlocks(b)
	} else {
		return nil, e
	}
	g.AddBlocks(modelToOrderedMap(g, m, args.Enums, args.Database))
	return g.Render(linebreak)
}

func modelToMap(g *golang.File, m *model.Model, enums enum.Enums, database string) *golang.Block {
	ret := golang.NewBlock(m.Package+"ToMap", "func")
	ret.WF("func (%s *%s) ToMap() util.ValueMap {", m.FirstLetter(), m.Proper())
	content := strings.Join(lo.Map(m.Columns, func(col *model.Column, _ int) string {
		return fmt.Sprintf(`%q: %s.%s`, col.Camel(), m.FirstLetter(), col.Proper())
	}), ", ")
	ret.W("\treturn util.ValueMap{" + content + "}")
	ret.W("}")

	_ = ` {
		return util.ValueMap{"id": m.ID, "t": m.T, "data": m.Data, "occurred": m.Occurred}
	}`
	return ret
}

func modelFromMap(g *golang.File, m *model.Model, models model.Models, enums enum.Enums, database string) (*golang.Block, error) {
	cols := m.Columns.NotDerived().WithoutTags("created", "updated")
	pks := cols.PKs()
	nonPKs := cols.NonPKs()

	ret := golang.NewBlock(m.Package+"FromMap", "func")
	ret.WF("func %sFromMap(m util.ValueMap, setPK bool) (*%s, util.ValueMap, error) {", m.Proper(), m.Proper())
	ret.WF("\tret := &%s{}", m.Proper())
	ret.W("\textra := util.ValueMap{}")
	ret.W("\tfor k, v := range m {")
	ret.W("\t\tvar err error")
	ret.W("\t\tswitch k {")
	for _, col := range pks {
		ret.WF("\t\tcase %q:", col.CamelNoReplace())
		ret.W("\t\t\tif setPK {")
		if err := forMapCol(g, ret, 4, m, models, enums, col); err != nil {
			return nil, err
		}
		ret.W("\t\t\t}")
	}
	for _, col := range nonPKs {
		ret.WF("\t\tcase %q:", col.CamelNoReplace())
		if err := forMapCol(g, ret, 3, m, models, enums, col); err != nil {
			return nil, err
		}
	}
	ret.W("\t\tdefault:")
	ret.W("\t\t\textra[k] = v")
	ret.W("\t\t}")
	ret.W("\t\tif err != nil {")
	ret.W("\t\t\treturn nil, nil, err")
	ret.W("\t\t}")
	ret.W("\t}")
	ret.W("\t// $PF_SECTION_START(extrachecks)$")
	ret.W("\t// $PF_SECTION_END(extrachecks)$")
	ret.W("\treturn ret, extra, nil")
	ret.W("}")

	return ret, nil
}

func modelToOrderedMap(g *golang.File, m *model.Model, enums enum.Enums, database string) *golang.Block {
	ret := golang.NewBlock(m.Package+"ToOrderedMap", "func")
	ret.WF("func (%s *%s) ToOrderedMap() *util.OrderedMap[any] {", m.FirstLetter(), m.Proper())
	content := strings.Join(lo.Map(m.Columns, func(col *model.Column, _ int) string {
		return fmt.Sprintf(`{K: %q, V: %s.%s}`, col.Camel(), m.FirstLetter(), col.Proper())
	}), ", ")
	ret.W("\tpairs := util.OrderedPairs[any]{" + content + "}")
	ret.W("\treturn util.NewOrderedMap[any](false, 4, pairs...)")
	ret.W("}")

	_ = ` {
		return util.ValueMap{"id": m.ID, "t": m.T, "data": m.Data, "occurred": m.Occurred}
	}`
	return ret
}
