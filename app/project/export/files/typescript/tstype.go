package typescript

import (
	"strings"

	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/lib/types"
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
