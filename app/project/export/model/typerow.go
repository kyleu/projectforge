package model

import (
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/enum"
	"projectforge.dev/projectforge/app/util"
)

func ToGoRowType(t types.Type, nullable bool, pkg string, enums enum.Enums, database string) (string, error) {
	switch t.Key() {
	case types.KeyAny, types.KeyList, types.KeyMap, types.KeyValueMap, types.KeyReference:
		if database == util.DatabaseSQLite {
			return types.KeyString, nil
		}
		return "json.RawMessage", nil
	default:
		if t.Key() == types.KeyUUID && database == util.DatabaseSQLServer {
			if nullable {
				return "*any", nil
			}
			return "mssql.UniqueIdentifier", nil
		}
		if t.Scalar() && nullable {
			switch t.Key() {
			case types.KeyString:
				return "sql.NullString", nil
			case types.KeyInt:
				return "sql.NullInt64", nil
			case types.KeyFloat:
				return "sql.NullFloat64", nil
			case types.KeyBool:
				return "sql.NullBool", nil
			}
		}
		return ToGoType(t, nullable, pkg, enums)
	}
}
