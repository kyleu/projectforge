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
	case TypeBool.Key:
		ret = TypeBool.Key
	case TypeInt.Key:
		ret = TypeInt.Key
	case TypeInterface.Key:
		ret = "interface{}"
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
	if nullable && !t.IsScalar() {
		return "*" + ret
	}
	return ret
}

func (t *Type) ToGoDTOType(nullable bool) string {
	switch t.Key {
	case TypeInterface.Key:
		return "json.RawMessage"
	case TypeMap.Key:
		return "json.RawMessage"
	default:
		return t.ToGoType(nullable)
	}
}

func (t *Type) ToGoString(prop string) string {
	switch t.Key {
	case TypeBool.Key:
		return fmt.Sprintf("fmt.Sprint(%s)", prop)
	case TypeInt.Key:
		return fmt.Sprintf("fmt.Sprint(%s)", prop)
	case TypeInterface.Key:
		return fmt.Sprintf("fmt.Sprint(%s)", prop)
	case TypeUUID.Key:
		return fmt.Sprintf("%s.String()", prop)
	default:
		return prop
	}
}

func (t *Type) ToGoViewString(prop string, nullable bool) string {
	switch t.Key {
	case TypeBool.Key:
		return "{%%v " + prop + " %%}"
	case TypeInt.Key:
		return "{%%d " + prop + " %%}"
	case TypeInterface.Key:
		return "{%%= components.JSON(" + prop + ") %%}"
	case TypeMap.Key:
		return "{%%= components.JSON(" + prop + ") %%}"
	case TypeTimestamp.Key:
		if nullable {
			return "{%%= components.DisplayTimestamp(" + prop + ") %%}"
		}
		return "{%%= components.DisplayTimestamp(&" + prop + ") %%}"
	case TypeUUID.Key:
		if nullable {
			return "{%%= components.DisplayUUID(" + prop + ") %%}"
		}
		return "{%%= components.DisplayUUID(&" + prop + ") %%}"
	default:
		return "{%%s " + t.ToGoString(prop) + " %%}"
	}
}

func (t *Type) ToSQLType() string {
	switch t.Key {
	case TypeBool.Key:
		return "boolean"
	case TypeInt.Key:
		return "int"
	case TypeInterface.Key:
		return "jsonb"
	case TypeMap.Key:
		return "jsonb"
	case TypeString.Key:
		return "text"
	case TypeTimestamp.Key:
		return "timestamp"
	case TypeUUID.Key:
		return "uuid"
	default:
		return "sql-error-invalid-type"
	}
}

func (t *Type) IsScalar() bool {
	switch t.Key {
	case TypeBool.Key, TypeInt.Key, TypeInterface.Key, TypeMap.Key, TypeString.Key:
		return true
	default:
		return false
	}
}
