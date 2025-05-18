package types

import "fmt"

const (
	emptyList = "[]"
	emptyMap  = "{}"
)

func invalidInput(key string, v any) string {
	return fmt.Sprintf("unable to parse [%s] from [%v] (%T)", key, v, v)
}

func Bits(t Type) int {
	if i := TypeAs[*Int](t); i != nil {
		return i.Bits
	}
	if i := TypeAs[*Float](t); i != nil {
		return i.Bits
	}
	return 0
}

func FromJSONType(s string, ref string) *Wrapped {
	switch s {
	case "", "null":
		if ref != "" {

		}
		return NewAny()
	case "nil", "<nil>":
		return NewNil()
	case "array":
		return NewList(NewAny())
	case "boolean":
		return NewBool()
	case "enum":
		return NewEnum("?")
	case "integer":
		return NewInt(0)
	case "number":
		return NewFloat(0)
	case "object":
		return NewMap(NewString(), NewAny())
	case "string":
		return NewString()
	default:
		return NewError(fmt.Sprintf("unhandled type [%s] for FromJSONType", s))
	}
}

func ToJSONType(v Type) string {
	switch t := v.(type) {
	case *Wrapped:
		return ToJSONType(t.T)
	case *Float:
		return "number"
	case *Int:
		return "integer"
	case *Map:
		return "object"
	case *String:
		return "string"
	default:
		return fmt.Sprintf("unhandled type [%T] for ToJSONType", t)
	}
}
