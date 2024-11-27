package model

import (
	"fmt"

	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/util"
)

const keyText, keyNVarcharMax = "text", "nvarchar(max)"

func ToSQLType(t types.Type, database string) (string, error) {
	switch t.Key() {
	case types.KeyAny:
		if database == util.DatabaseSQLServer {
			return keyText, nil
		}
		return keyJSONB, nil
	case types.KeyBool:
		if database == util.DatabaseSQLServer {
			return "bit", nil
		}
		return "boolean", nil
	case types.KeyEnum:
		if database == util.DatabaseSQLServer {
			return "nvarchar(255)", nil
		}
		e, err := AsEnum(t)
		if err != nil {
			return "", err
		}
		return e.Ref, nil
	case types.KeyInt:
		if types.Bits(t) == 64 {
			return "bigint", nil
		}
		return "int", nil
	case types.KeyFloat:
		return "double precision", nil
	case types.KeyList, types.KeyMap, types.KeyOrderedMap, types.KeyValueMap, types.KeyReference, types.KeyNumeric:
		if database == util.DatabaseSQLServer {
			return keyNVarcharMax, nil
		}
		return keyJSONB, nil
	case types.KeyString:
		if database == util.DatabaseSQLServer {
			s, ok := types.Wrap(t).T.(*types.String)
			if !ok {
				return keyNVarcharMax, nil
			}
			if s.MaxLength > 0 {
				return fmt.Sprintf("nvarchar(%d)", s.MaxLength), nil
			}
			return keyNVarcharMax, nil
		}
		return keyText, nil
	case types.KeyTimestampZoned:
		return types.KeyTimestamp, nil
	case types.KeyDate, types.KeyTimestamp:
		if database == util.DatabaseSQLServer {
			return "datetime", nil
		}
		return types.KeyTimestamp, nil
	case types.KeyUUID:
		if database == util.DatabaseSQLServer {
			return "uniqueidentifier", nil
		}
		return "uuid", nil
	default:
		return "sql-error-invalid-type", nil
	}
}
