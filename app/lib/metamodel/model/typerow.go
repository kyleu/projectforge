package model

import (
	"fmt"

	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/util"
)

func ToGoRowType(t types.Type, nullable bool, pkg string, enums enum.Enums, database string) (string, error) {
	switch t.Key() {
	case types.KeyAny, types.KeyList, types.KeyMap, types.KeyOrderedMap, types.KeyValueMap, types.KeyReference:
		if database == util.DatabaseSQLite || database == util.DatabaseSQLServer {
			return types.KeyString, nil
		}
		return "[]byte", nil
	case types.KeyJSON:
		return "util.NilJSON", nil
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
				return "util.NilString", nil
			case types.KeyInt:
				if i := types.TypeAs[*types.Int](t); i != nil {
					b := util.Choose(i.Bits == 0, 64, i.Bits)
					return fmt.Sprintf("util.NilInt%d", b), nil
				}
				return "util.NilInt64", nil
			case types.KeyFloat:
				return "util.NilFloat64", nil
			case types.KeyBool:
				return "util.NilBool", nil
			}
		}
		return ToGoType(t, nullable, pkg, enums)
	}
}
