package gomodel

import (
	"fmt"

	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/project/export/model"
)

func modelDiff(m *model.Model, g *golang.File) *golang.Block {
	ret := golang.NewBlock("Diff"+m.Proper(), "func")

	var complexity int
	for _, col := range m.Columns {
		switch col.Type.Key() {
		case types.KeyAny, types.KeyList, types.KeyMap, types.KeyValueMap, types.KeyReference:
			complexity += 2
		case types.KeyBool, types.KeyInt, types.KeyFloat:
			complexity += 4
		case types.KeyEnum:
			complexity += 4
		case types.KeyString:
			complexity += 2
		case types.KeyDate, types.KeyTimestamp, types.KeyUUID:
			complexity += 6
		default:
			complexity += 2
		}
	}
	if complexity > 42 {
		ret.W("//nolint:gocognit")
	}
	ret.W("func (%s *%s) Diff(%sx *%s) util.Diffs {", m.FirstLetter(), m.Proper(), m.FirstLetter(), m.Proper())
	ret.W("\tvar diffs util.Diffs")
	for _, col := range m.Columns {
		if col.HasTag("updated") {
			continue
		}
		l := fmt.Sprintf("%s.%s", m.FirstLetter(), col.Proper())
		r := fmt.Sprintf("%sx.%s", m.FirstLetter(), col.Proper())
		switch col.Type.Key() {
		case types.KeyAny, types.KeyList, types.KeyMap, types.KeyValueMap, types.KeyReference:
			ret.W("\tdiffs = append(diffs, util.DiffObjects(%s, %s, %q)...)", l, r, col.Camel())
		case types.KeyBool, types.KeyInt, types.KeyFloat:
			g.AddImport(helper.ImpFmt)
			ret.W("\tif %s != %s {", l, r)
			ret.W("\t\tdiffs = append(diffs, util.NewDiff(%q, fmt.Sprint(%s), fmt.Sprint(%s)))", col.Camel(), l, r)
			ret.W("\t}")
		case types.KeyEnum:
			ret.W("\tif %s != %s {", l, r)
			ret.W("\t\tdiffs = append(diffs, util.NewDiff(%q, string(%s), string(%s)))", col.Camel(), l, r)
			ret.W("\t}")
		case types.KeyString:
			ret.W("\tif %s != %s {", l, r)
			ret.W("\t\tdiffs = append(diffs, util.NewDiff(%q, %s, %s))", col.Camel(), l, r)
			ret.W("\t}")
		case types.KeyDate, types.KeyTimestamp, types.KeyUUID:
			if col.Nullable {
				msg := "\tif (%s == nil && %s != nil) || (%s != nil && %s == nil) || (%s != nil && %s != nil && *%s != *%s) {"
				line := fmt.Sprintf(msg, l, r, l, r, l, r, l, r)
				if len(line) > 152 {
					line += " //nolint:lll"
				}
				ret.W(line)
				g.AddImport(helper.ImpFmt)
				ret.W("\t\tdiffs = append(diffs, util.NewDiff(%q, fmt.Sprint(%s), fmt.Sprint(%s))) //nolint:gocritic // it's nullable", col.Camel(), l, r)
			} else {
				ret.W("\tif %s != %s {", l, r)
				ret.W("\t\tdiffs = append(diffs, util.NewDiff(%q, %s.String(), %s.String()))", col.Camel(), l, r)
			}
			ret.W("\t}")
		default:
			ret.W("\tTODO: %s", col.Type.Key())
		}
	}
	ret.W("\treturn diffs")
	ret.W("}")
	return ret
}
