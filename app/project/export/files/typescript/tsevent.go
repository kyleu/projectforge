package typescript

import (
	"strings"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/util"
)

func tsEventConstructor(evt *model.Event, enums enum.Enums, ret *golang.Block) {
	s := &util.StringSlice{}
	cols := evt.Columns.NotDerived()
	for _, col := range cols {
		optional := util.Choose(col.Nullable || col.HasTag("optional-json"), " | undefined", "")
		s.Pushf("%s: %s%s", col.Camel(), tsType(col.Type, enums), optional)
	}
	args := strings.Join(s.Slice, ", ")
	if len(args) < 106 {
		ret.WF("  constructor(%s) {", args)
	} else {
		ret.WF("  constructor(")
		for idx, col := range s.Slice {
			comma := util.Choose(idx+1 < len(s.Slice), ",", "")
			ret.WF("    %s%s", col, comma)
		}
		ret.W("  ) {")
	}
	for _, col := range cols {
		ret.WF("    this.%s = %s;", col.Camel(), col.Camel())
	}
	ret.W("  }")
}

func tsEvent(evt *model.Event, enums enum.Enums) (*golang.Block, error) {
	ret := golang.NewBlock("ts-event-"+evt.Name, "ts")
	ret.WF("export class %s {", evt.Proper())
	cols := evt.Columns.NotDerived()
	for _, col := range cols {
		optional := util.Choose(col.Nullable || col.HasTag("optional-json"), " | undefined", "")
		ret.WF("  %s: %s%s;%s", col.Camel(), tsType(col.Type, enums), optional, col.CommentString())
	}
	ret.WB()
	tsEventConstructor(evt, enums, ret)
	err := tsFromObject(cols, evt, enums, ret)
	if err != nil {
		return nil, err
	}
	ret.W("}")
	ret.WB()
	ret.WF("export type %s = %s[];", evt.ProperPlural(), evt.Proper())
	return ret, nil
}

func tsEventContent(imps []string, evt *model.Event, enums enum.Enums) (golang.Blocks, error) {
	var ret golang.Blocks
	if len(imps) > 0 {
		b := golang.NewBlock("imports", "ts")
		for _, l := range imps {
			b.W(l)
		}
		ret = append(ret, b)
	}
	n, err := tsEvent(evt, enums)
	if err != nil {
		return nil, err
	}
	ret = append(ret, n)
	return ret, nil
}

func EventContent(evt *model.Event, allEnums enum.Enums, imps []string, linebreak string) (*file.File, error) {
	dir := []string{"client", "src"}
	dir = append(dir, evt.PackageWithGroup(""))
	filename := evt.Kebab()
	g := golang.NewGoTemplate(dir, filename+".ts")
	b, err := tsEventContent(imps, evt, allEnums)
	if err != nil {
		return nil, err
	}
	g.AddBlocks(b...)
	return g.Render(linebreak)
}
