package typescript

import (
	"strings"

	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/files/gohelper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/util"
)

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

func tsFromObject(cols model.Columns, str gohelper.StringProvider, enums enum.Enums, ret *golang.Block) error {
	ret.WB()
	ret.WF("  static fromObject(obj: { [_: string]: unknown }): %s {", str.Proper())
	for _, col := range cols {
		err := tsFromObjectColumn(col, enums, ret)
		if err != nil {
			return err
		}
	}
	ret.WB()
	args := cols.CamelNames()
	argsJoined := strings.Join(args, ", ")
	if len(argsJoined) < 100 {
		ret.WF("    return new %s(%s);", str.Proper(), argsJoined)
	} else {
		ret.WF("    return new %s(", str.Proper())
		for idx, c := range cols.CamelNames() {
			comma := util.Choose(idx+1 < len(cols.CamelNames()), ",", "")
			ret.WF("      %s%s", c, comma)
		}
		ret.W("    );")
	}
	ret.W("  }")
	return nil
}
