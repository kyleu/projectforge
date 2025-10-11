package gomodel

import (
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/metamodel"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/project/export/files/gohelper"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/util"
)

func modelStruct(m *model.Model, args *metamodel.Args) (*golang.Block, error) {
	ret := golang.NewBlock(m.Proper(), "struct")
	ret.WF("type %s struct {", m.Proper())
	cols := m.Columns.NotDerived()
	maxColLength := cols.MaxCamelLength()
	maxTypeLength := cols.MaxGoTypeLength(m.Package, args.Enums)

	gts := lo.Map(cols, func(c *model.Column, _ int) string {
		gt := helper.GoTypeWithRef(c, m.Package, args)
		if x := len(gt); maxTypeLength < x {
			maxTypeLength = x
		}
		return gt
	})

	for idx, c := range cols {
		if c.HelpString != "" {
			ret.WF("\t// %s", c.HelpString)
		}
		goType := util.StringPad(gts[idx], maxTypeLength)
		tag := gohelper.ColumnTag(c)
		ret.WF("\t%s %s %s", util.StringPad(c.Proper(), maxColLength), goType, tag)
	}
	ret.W("}")
	return ret, nil
}

func modelConstructor(m *model.Model, x *metamodel.Args) (*golang.Block, error) {
	ret := golang.NewBlock("New"+m.Proper(), "func")
	argsString, err := helper.GoArgsWithRef(m.PKs(), m.Package, x)
	if err != nil {
		return nil, err
	}
	ret.WF("func New%s(%s) *%s {", m.Proper(), argsString, m.Proper())
	ret.WF("\treturn &%s{%s}", m.Proper(), m.PKs().Refs())
	ret.W("}")
	return ret, nil
}
