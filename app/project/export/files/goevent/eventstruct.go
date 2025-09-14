package goevent

import (
	"fmt"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/metamodel"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/files/gohelper"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/util"
)

func eventStruct(evt *model.Event, args *metamodel.Args) (*golang.Block, error) {
	ret := golang.NewBlock(evt.Proper(), "struct")
	ret.WF("type %s struct {", evt.Proper())
	cols := evt.Columns.NotDerived()
	maxColLength := cols.MaxCamelLength()
	maxTypeLength := cols.MaxGoTypeLength(evt.Package, args.Enums)

	gts := lo.Map(cols, func(c *model.Column, _ int) string {
		gt, err := c.ToGoType(evt.Package, args.Enums)
		if err != nil {
			return err.Error()
		}
		if ref, mdl, _ := helper.LoadRef(c, args.Models, args.ExtraTypes); ref != nil && !strings.Contains(gt, ".") {
			if mdl != nil && mdl.Package != evt.Package {
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
