package model

import (
	"fmt"
)

type Type struct {
	Key   string
	Names []string
}

func newType(names ...string) *Type {
	return &Type{Key: names[0], Names: names}
}

func (t *Type) ToGoType(nullable bool) string {
	var ret string
	switch t.Key {
	case TypeInt.Key:
		ret = TypeInt.Key
	case TypeMap.Key:
		ret = "util.ValueMap"
	case TypeString.Key:
		ret = TypeString.Key
	case TypeTimestamp.Key:
		ret = "time.Time"
	case TypeUUID.Key:
		ret = "uuid.UUID"
	default:
		return "ERROR:Unhandled[" + t.Key + "]"
	}
	if nullable {
		return "*" + ret
	}
	return ret
}

func (t *Type) ToGoDTOType(nullable bool) string {
	switch t.Key {
	case TypeMap.Key:
		return "json.RawMessage"
	default:
		return t.ToGoType(nullable)
	}
}

func (t *Type) ToGoString(prop string) string {
	switch t.Key {
	case TypeInt.Key:
		return fmt.Sprintf("fmt.Sprint(%s)", prop)
	case TypeUUID.Key:
		return fmt.Sprintf("%s.String()", prop)
	default:
		return prop
	}
}

func (t *Type) ToGoViewString(prop string, nullable bool) string {
	switch t.Key {
	case TypeInt.Key:
		return "{%%d " + prop + " %%}"
	case TypeMap.Key:
		return "{%%= components.JSON(" + prop + ") %%}"
	case TypeTimestamp.Key:
		if nullable {
			return "{%%= components.DisplayTimestamp(" + prop + ") %%}"
		}
		return "{%%= components.DisplayTimestamp(&" + prop + ") %%}"
	default:
		return "{%%s " + t.ToGoString(prop) + " %%}"
	}
}
