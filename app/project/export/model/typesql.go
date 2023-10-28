package model

import "projectforge.dev/projectforge/app/lib/types"

func ToSQLType(t types.Type) (string, error) {
	switch t.Key() {
	case types.KeyAny:
		return keyJSONB, nil
	case types.KeyBool:
		return "boolean", nil
	case types.KeyEnum:
		e, err := AsEnum(t)
		if err != nil {
			return "", err
		}
		return e.Ref, nil
	case types.KeyInt:
		i := types.TypeAs[*types.Int](t)
		if i != nil && i.Bits == 64 {
			return "bigint", nil
		}
		return "int", nil
	case types.KeyFloat:
		return "double precision", nil
	case types.KeyList, types.KeyMap, types.KeyValueMap, types.KeyReference:
		return keyJSONB, nil
	case types.KeyString:
		return "text", nil
	case types.KeyDate, types.KeyTimestamp:
		return "timestamp", nil
	case types.KeyUUID:
		return "uuid", nil
	default:
		return "sql-error-invalid-type", nil
	}
}
