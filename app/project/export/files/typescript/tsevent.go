package typescript

import (
	"strings"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/util"
)

func tsEventConstructor(m *model.Event, enums enum.Enums, ret *golang.Block) {
	s := &util.StringSlice{}
	for _, col := range m.Columns {
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
	for _, col := range m.Columns {
		ret.WF("    this.%s = %s;", col.Camel(), col.Camel())
	}
	ret.W("  }")
}

func tsEvent(m *model.Event, enums enum.Enums) (*golang.Block, error) {
	ret := golang.NewBlock("ts-event-"+m.Name, "ts")
	ret.WF("export class %s {", m.Proper())
	for _, col := range m.Columns {
		optional := util.Choose(col.Nullable || col.HasTag("optional-json"), " | undefined", "")
		ret.WF("  %s: %s%s;", col.Camel(), tsType(col.Type, enums), optional)
	}
	ret.WB()
	tsEventConstructor(m, enums, ret)
	err := tsFromObject(m.Columns, m, enums, ret)
	if err != nil {
		return nil, err
	}
	ret.W("}")
	ret.WB()
	ret.WF("export type %s = Array<%s>;", m.ProperPlural(), m.Proper())
	return ret, nil
}

func tsEventContent(imps []string, m *model.Event, enums enum.Enums) (golang.Blocks, error) {
	var ret golang.Blocks
	if len(imps) > 0 {
		b := golang.NewBlock("imports", "ts")
		for _, l := range imps {
			b.W(l)
		}
		ret = append(ret, b)
	}
	n, err := tsEvent(m, enums)
	if err != nil {
		return nil, err
	}
	ret = append(ret, n)
	return ret, nil
}

func EventContent(m *model.Event, allEnums enum.Enums, imps []string, linebreak string) (*file.File, error) {
	dir := []string{"client", "src"}
	dir = append(dir, m.PackageWithGroup(""))
	filename := m.Camel()
	g := golang.NewGoTemplate(dir, filename+".ts")
	b, err := tsEventContent(imps, m, allEnums)
	if err != nil {
		return nil, err
	}
	g.AddBlocks(b...)
	return g.Render(linebreak)
}
