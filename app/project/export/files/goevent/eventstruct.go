package goevent

import (
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/metamodel"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/project/export/files/gohelper"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/util"
)

func eventStruct(evt *model.Event, args *metamodel.Args) (*golang.Block, error) {
	ret := golang.NewBlock(evt.Proper(), "struct")
	ret.WF("type %s struct {", evt.Proper())
	cols := evt.Columns
	maxColLength := cols.MaxCamelLength()
	maxTypeLength := cols.MaxGoTypeLength(evt.Package, args.Enums)

	gts := lo.Map(cols, func(c *model.Column, _ int) string {
		gt := helper.GoTypeWithRef(c, evt.Package, args)
		if x := len(gt); maxTypeLength < x {
			maxTypeLength = x
		}
		return gt
	})

	for idx, c := range cols {
		goType := util.StringPad(gts[idx], maxTypeLength)
		tag := gohelper.ColumnTag(c)
		ret.WF("\t%s %s %s", util.StringPad(c.Proper(), maxColLength), goType, tag)
	}
	ret.W("}")
	return ret, nil
}
