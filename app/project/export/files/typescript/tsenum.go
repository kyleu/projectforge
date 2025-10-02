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
		ret.WF("  %s = %q%s", v.Proper(), v.Key, suffix)
	})
	ret.W("}")
	return ret
}

func tsEnumParse(e *enum.Enum) *golang.Block {
	ret := golang.NewBlock("parseEnum", "ts")
	ret.WF("export function parse%s(value: string): %s | undefined {", e.Proper(), e.Proper())
	ret.WF("  if (Object.values(%s).includes(value as %s)) {", e.Proper(), e.Proper())
	ret.WF("    return value as %s;", e.Proper())
	ret.W("  }")
	ret.W("  return undefined;")
	ret.W("}")
	return ret
}

func tsEnumGet(e *enum.Enum) *golang.Block {
	ret := golang.NewBlock("getEnum", "ts")
	ret.WF("export function get%s(value: string): %s {", e.Proper(), e.Proper())
	ret.WF("  const x = parse%s(value);", e.Proper())
	ret.W("  if (x === undefined) {")
	ret.WF("    throw new Error(`invalid [%s]: ${value}`);", e.Proper())
	ret.W("  }")
	ret.W("  return x;")
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
	ret = append(ret, tsEnum(e), tsEnumParse(e), tsEnumGet(e))
	return ret
}

func EnumContent(e *enum.Enum, imps []string, linebreak string) (*file.File, error) {
	dir := []string{"client", "src"}
	dir = append(dir, e.PackageWithGroup(""))
	filename := e.Kebab()
	g := golang.NewGoTemplate(dir, filename+".ts")
	g.AddBlocks(tsEnumContent(imps, e)...)
	return g.Render(linebreak)
}
