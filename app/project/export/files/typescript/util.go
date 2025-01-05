package typescript

import (
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/util"
)

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

func tsModelContent(imps []string, m *model.Model, enums enum.Enums) golang.Blocks {
	var ret golang.Blocks
	if len(imps) > 0 {
		b := golang.NewBlock("imports", "ts")
		for _, l := range imps {
			b.W(l)
		}
		ret = append(ret, b)
	}
	ret = append(ret, tsModel(m, enums))
	return ret
}

func tsEnum(e *enum.Enum) *golang.Block {
	ret := golang.NewBlock("tsenum-"+e.Name, "ts")
	// ret.W("// eslint-disable-next-line no-shadow")
	ret.WF("export enum %s {", e.Proper())
	lo.ForEach(e.Values, func(v *enum.Value, idx int) {
		suffix := util.Choose(idx == len(e.Values)-1, "", ",")
		ret.WF("  %s = %q%s", v.Name, v.Key, suffix)
	})
	ret.W("}")
	return ret
}

func tsModel(m *model.Model, enums enum.Enums) *golang.Block {
	ret := golang.NewBlock("tsmodel-"+m.Name, "ts")
	ret.WF("export class %s {", m.Proper())
	for _, col := range m.Columns {
		optional := util.Choose(col.Nullable || col.HasTag("optional-json"), "?", "")
		ret.WF("  %s%s: %s;", col.Camel(), optional, tsType(col.Type, enums))
	}
	ret.W("}")
	ret.WB()
	ret.WF("export type %s = Array<%s>;", m.ProperPlural(), m.Proper())
	return ret
}

func tsType(t *types.Wrapped, enums enum.Enums) string {
	switch t.Key() {
	case types.KeyBool:
		return "boolean"
	case types.KeyUUID:
		return "string"
	case types.KeyInt, types.KeyFloat:
		return "number"
	case types.KeyTimestamp:
		return "Date"
	case types.KeyEnum:
		e := enums.Get(t.EnumKey())
		return e.Proper()
	case types.KeyList:
		lt := t.ListType()
		return tsType(lt, enums) + "[]"
	case types.KeyNumeric:
		return "Numeric"
	case types.KeyReference:
		r, _ := model.AsRef(t)
		return strings.TrimPrefix(r.K, "*")
	case types.KeyMap, types.KeyOrderedMap, types.KeyValueMap:
		return "{ [key: string]: unknown }"
	default:
		return t.String()
	}
}
