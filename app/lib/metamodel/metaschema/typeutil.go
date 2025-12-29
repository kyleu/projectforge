package metaschema

import (
	"fmt"

	"projectforge.dev/projectforge/app/lib/types"
)

func FromJSONType(s string, ref string) *types.Wrapped {
	switch s {
	case "", KeyNull:
		return types.NewAny()
	case KeyNil, KeyNilString:
		return types.NewNil()
	case KeyArray:
		return types.NewList(types.NewAny())
	case KeyBoolean:
		return types.NewBool()
	case KeyEnum:
		return types.NewEnum("?")
	case KeyInteger:
		return types.NewInt(0)
	case KeyNumber:
		return types.NewFloat(0)
	case KeyObject:
		return types.NewStringKeyedMap()
	case KeyString:
		return types.NewString()
	default:
		return types.NewError(fmt.Sprintf("unhandled type [%s] for FromJSONType", s))
	}
}

func ToJSONType(v types.Type) string {
	switch t := v.(type) {
	case *types.Wrapped:
		return ToJSONType(t.T)
	case *types.Float:
		return KeyNumber
	case *types.Int:
		return KeyInteger
	case *types.Map:
		return KeyObject
	case *types.String:
		return KeyString
	default:
		return fmt.Sprintf("unhandled type [%T] for ToJSONType", t)
	}
}
