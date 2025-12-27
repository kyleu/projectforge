package gohelper

import (
	"fmt"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/metamodel"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/util"
)

func ToMap(m metamodel.StringProvider, cols model.Columns) *golang.Block {
	ret := golang.NewBlock(m.PackageName()+"ToMap", "func")
	ret.WF("func (%s *%s) ToMap() util.ValueMap {", m.FirstLetter(), m.Proper())
	content := util.StringJoin(lo.Map(cols, func(col *model.Column, _ int) string {
		return fmt.Sprintf(`%q: %s.%s`, col.Camel(), m.FirstLetter(), col.Proper())
	}), ", ")
	ret.W("\treturn util.ValueMap{" + content + "}")
	ret.W("}")
	return ret
}

func FromMap(g *golang.File, m metamodel.StringProvider, cols model.Columns, args *metamodel.Args) (*golang.Block, error) {
	c := cols.NotDerived().WithoutTags("created", "updated")
	pks := c.PKs()
	nonPKs := c.NonPKs()

	ret := golang.NewBlock(m.PackageName()+"FromMap", "func")
	ret.WF("func %sFromMap(m util.ValueMap, setPK bool) (*%s, util.ValueMap, error) {", m.Proper(), m.Proper())
	ret.WF("\tret := &%s{}", m.Proper())
	ret.W("\textra := util.ValueMap{}")
	ret.W("\tfor k, v := range m {")
	ret.W("\t\tvar err error")
	ret.W("\t\tswitch k {")
	for _, col := range pks {
		ret.WF("\t\tcase %q:", col.CamelNoReplace())
		ret.W("\t\t\tif setPK {")
		if err := forMapCol(g, ret, 4, m, args, col); err != nil {
			return nil, err
		}
		ret.W("\t\t\t}")
	}
	for _, col := range nonPKs {
		ret.WF("\t\tcase %q:", col.CamelNoReplace())
		if err := forMapCol(g, ret, 3, m, args, col); err != nil {
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

func ToOrderedMap(m metamodel.StringProvider, cols model.Columns) *golang.Block {
	ret := golang.NewBlock(m.PackageName()+"ToOrderedMap", "func")
	ret.WF("func (%s *%s) ToOrderedMap() *util.OrderedMap[any] {", m.FirstLetter(), m.Proper())
	ret.WF("\tif %s == nil {", m.FirstLetter())
	ret.W("\t\treturn nil")
	ret.W("\t}")

	content := util.StringJoin(lo.Map(cols, func(col *model.Column, _ int) string {
		return fmt.Sprintf(`{K: %q, V: %s.%s}`, col.Camel(), m.FirstLetter(), col.Proper())
	}), ", ")
	ret.W("\tpairs := util.OrderedPairs[any]{" + content + "}")
	ret.W("\treturn util.NewOrderedMap(false, 4, pairs...)")
	ret.W("}")
	return ret
}
