package typescript

import (
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/util"
)

func tsEnum(e *enum.Enum) *golang.Block {
	ret := golang.NewBlock("ts-enum-"+e.Name, "ts")
	// ret.W("// eslint-disable-next-line no-shadow")
	ret.WF("export enum %s {", e.Proper())
	lo.ForEach(e.Values, func(v *enum.Value, idx int) {
		suffix := util.Choose(idx == len(e.Values)-1, "", ",")
		ret.WF("  %s = %q%s", v.Name, v.Key, suffix)
	})
	ret.W("}")
	return ret
}

func tsEnumContent(imps []string, e *enum.Enum) golang.Blocks {
	var ret golang.Blocks
	if len(imps) > 0 {
		b := golang.NewBlock("imports", "ts")
		for _, l := range imps {
			b.W(l)
		}
		ret = append(ret, b)
	}
	ret = append(ret, tsEnum(e))
	return ret
}

func EnumContent(e *enum.Enum, imps []string, linebreak string) (*file.File, error) {
	dir := []string{"client", "src"}
	dir = append(dir, e.PackageWithGroup(""))
	filename := e.Camel()
	g := golang.NewGoTemplate(dir, filename+".ts")
	g.AddBlocks(tsEnumContent(imps, e)...)
	return g.Render(linebreak)
}
