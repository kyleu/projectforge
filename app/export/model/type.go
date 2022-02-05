package model

import (
	"fmt"

	"github.com/kyleu/projectforge/app/lib/types"
)

func ToGoType(t types.Type, nullable bool) string {
	var ret string
	switch t.Key() {
	case types.KeyAny:
		ret = "interface{}"
	case types.KeyBool:
		ret = types.KeyBool
	case types.KeyInt:
		ret = types.KeyInt
	case types.KeyMap:
		ret = "util.ValueMap"
	case types.KeyString:
		ret = types.KeyString
	case types.KeyTimestamp:
		ret = "time.Time"
	case types.KeyUUID:
		ret = "uuid.UUID"
	default:
		return "ERROR:Unhandled[" + t.Key() + "]"
	}
	if nullable && !t.Scalar() {
		return "*" + ret
	}
	return ret
}

func ToGoDTOType(t types.Type, nullable bool) string {
	switch t.Key() {
	case types.KeyAny, types.KeyMap:
		return "json.RawMessage"
	default:
		return ToGoType(t, nullable)
	}
}

func ToGoString(t types.Type, prop string) string {
	switch t.Key() {
	case types.KeyAny:
		return fmt.Sprintf("fmt.Sprint(%s)", prop)
	case types.KeyBool:
		return fmt.Sprintf("fmt.Sprint(%s)", prop)
	case types.KeyInt:
		return fmt.Sprintf("fmt.Sprint(%s)", prop)
	case types.KeyUUID:
		return fmt.Sprintf("%s.String()", prop)
	default:
		return prop
	}
}

func ToGoViewString(t types.Type, prop string, nullable bool) string {
	switch t.Key() {
	case types.KeyAny:
		return "{%%= components.JSON(" + prop + ") %%}"
	case types.KeyBool:
		return "{%%v " + prop + " %%}"
	case types.KeyInt:
		return "{%%d " + prop + " %%}"
	case types.KeyMap:
		return "{%%= components.JSON(" + prop + ") %%}"
	case types.KeyTimestamp:
		if nullable {
			return "{%%= components.DisplayTimestamp(" + prop + ") %%}"
		}
		return "{%%= components.DisplayTimestamp(&" + prop + ") %%}"
	case types.KeyUUID:
		if nullable {
			return "{%%= components.DisplayUUID(" + prop + ") %%}"
		}
		return "{%%= components.DisplayUUID(&" + prop + ") %%}"
	default:
		return "{%%s " + ToGoString(t, prop) + " %%}"
	}
}

const keyJSONB = "jsonb"

func ToSQLType(t types.Type) string {
	switch t.Key() {
	case types.KeyAny:
		return keyJSONB
	case types.KeyBool:
		return "boolean"
	case types.KeyInt:
		return "int"
	case types.KeyMap:
		return keyJSONB
	case types.KeyString:
		return "text"
	case types.KeyTimestamp:
		return "timestamp"
	case types.KeyUUID:
		return "uuid"
	default:
		return "sql-error-invalid-type"
	}
}
