package typescript

import (
	"github.com/samber/lo"
	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/util"
)

func tsContent(models model.Models, enums enum.Enums) *golang.Block {
	ret := golang.NewBlock("TypeScript", "ts")
	for _, e := range enums {
		tsEnum(e, ret)
	}
	for _, m := range models {
		tsModel(m, enums, ret)
	}
	ret.W("// $PF_SECTION_START(models)$")
	ret.W("// $PF_SECTION_END(models)$")
	return ret
}

func tsEnum(e *enum.Enum, ret *golang.Block) {
	ret.W("// eslint-disable-next-line no-shadow")
	ret.WF("export enum %s {", e.Proper())
	lo.ForEach(e.Values, func(v *enum.Value, idx int) {
		suffix := util.Choose(idx == len(e.Values)-1, "", ",")
		ret.WF("  %s = %q%s", v.Name, v.Key, suffix)
	})
	ret.W("}")
	ret.WB()
}

func tsModel(m *model.Model, enums enum.Enums, ret *golang.Block) {
	ret.WF("export class %s {", m.Proper())
	for _, col := range m.Columns {
		optional := util.Choose(col.Nullable || col.HasTag("optional-json"), "?", "")
		ret.WF("  %s%s: %s;", col.Camel(), optional, tsType(col.Type, enums))
	}
	ret.WF("  // $PF_SECTION_START(model-%s)$", m.Proper())
	ret.WF("  // $PF_SECTION_END(model-%s)$", m.Proper())
	ret.W("}")
	ret.WB()
}

func tsType(t *types.Wrapped, enums enum.Enums) string {
	switch t.Key() {
	case types.KeyUUID:
		return "string"
	case types.KeyInt, types.KeyFloat:
		return "number"
	case types.KeyTimestamp:
		return "Date"
	case types.KeyEnum:
		e := enums.Get(t.EnumKey())
		return e.Proper()
	case types.KeyValueMap, types.KeyMap:
		return "{ [key: string]: unknown }"
	default:
		return t.String()
	}
}
