package typescript

import (
	"github.com/samber/lo"
	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/util"
	"strings"
)

func tsContent(enums enum.Enums, models model.Models) golang.Blocks {
	ret := make(golang.Blocks, 0, len(enums)+len(models))
	for _, e := range enums {
		ret = append(ret, tsEnum(e))
	}
	for _, m := range models {
		ret = append(ret, tsModel(m, enums))
	}
	return ret
}

func tsEnum(e *enum.Enum) *golang.Block {
	ret := golang.NewBlock("tsenum-"+e.Name, "ts")
	ret.W("// eslint-disable-next-line no-shadow")
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
	case types.KeyReference:
		r, _ := model.AsRef(t)
		return strings.TrimPrefix(r.K, "*")
	case types.KeyValueMap, types.KeyMap:
		return "{ [key: string]: unknown }"
	default:
		return t.String()
	}
}