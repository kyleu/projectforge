package gomodel

import (
	"fmt"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/metamodel"
	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/lib/types"
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
		gt, err := c.ToGoType(m.Package, args.Enums)
		if err != nil {
			return err.Error()
		}
		if ref, mdl, _ := helper.LoadRef(c, args.Models, args.ExtraTypes); ref != nil && !strings.Contains(gt, ".") {
			if mdl != nil && mdl.Package != m.Package {
				gt = mdl.Package + "." + gt
				if x := len(gt); maxTypeLength < x {
					maxTypeLength = x
				}
			}
		}
		if gt == "*"+types.KeyAny {
			gt = types.KeyAny
		}
		return gt
	})

	for idx, c := range cols {
		goType := util.StringPad(gts[idx], maxTypeLength)
		var tag string
		if c.JSON == "" {
			tag = fmt.Sprintf("json:%q", c.CamelNoReplace()+gohelper.JSONSuffix(c))
		} else {
			tag = fmt.Sprintf("json:%q", c.JSON+gohelper.JSONSuffix(c))
		}
		if c.Validation != "" {
			tag += fmt.Sprintf(",validate:%q", c.Validation)
		}
		if c.Example != "" {
			tag += fmt.Sprintf(",fake:%q", c.Example)
		}
		ret.WF("\t%s %s `%s`", util.StringPad(c.Proper(), maxColLength), goType, tag)
	}
	ret.W("}")
	return ret, nil
}

func modelConstructor(m *model.Model, enums enum.Enums) (*golang.Block, error) {
	ret := golang.NewBlock("New"+m.Proper(), "func")
	args, err := m.PKs().Args(m.Package, enums)
	if err != nil {
		return nil, err
	}
	ret.WF("func New%s(%s) *%s {", m.Proper(), args, m.Proper())
	ret.WF("\treturn &%s{%s}", m.Proper(), m.PKs().Refs())
	ret.W("}")
	return ret, nil
}
