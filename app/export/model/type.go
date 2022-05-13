package model

import (
	"fmt"

	"github.com/pkg/errors"
	"projectforge.dev/projectforge/app/lib/types"
)

func ToGoType(t types.Type, nullable bool, pkg string) string {
	var ret string
	switch t.Key() {
	case types.KeyAny:
		ret = "any"
	case types.KeyBool:
		ret = types.KeyBool
	case types.KeyInt:
		ret = types.KeyInt
	case types.KeyMap, types.KeyValueMap:
		ret = "util.ValueMap"
	case types.KeyReference:
		ref, err := AsRef(t)
		if err != nil {
			return "ERROR:" + err.Error()
		}
		if ref.Pkg.Last() == pkg {
			ret = fmt.Sprintf("*%s", ref.K)
		} else {
			ret = fmt.Sprintf("*%s.%s", ref.Pkg.Last(), ref.K)
		}
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

func ToGoDTOType(t types.Type, nullable bool, pkg string) string {
	switch t.Key() {
	case types.KeyAny, types.KeyMap, types.KeyValueMap, types.KeyReference:
		return "json.RawMessage"
	default:
		return ToGoType(t, nullable, pkg)
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
	case types.KeyUUID, types.KeyReference:
		return fmt.Sprintf("%s.String()", prop)
	default:
		return prop
	}
}

func ToGoViewString(t types.Type, prop string, nullable bool, format string) string {
	switch t.Key() {
	case types.KeyAny:
		return "{%%= components.JSON(" + prop + ") %%}"
	case types.KeyBool:
		return "{%%v " + prop + " %%}"
	case types.KeyInt:
		return "{%%d " + prop + " %%}"
	case types.KeyMap, types.KeyValueMap, types.KeyReference:
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
	case types.KeyString:
		switch format {
		case FmtCode:
			return "<pre>{%%s " + ToGoString(t, prop) + " %%}</pre>"
		case FmtURL:
			x := "{%%s " + ToGoString(t, prop) + " %%}"
			return fmt.Sprintf("<a href=%q target=\"_blank\">%s</a>", x, x)
		default:
			return "{%%s " + ToGoString(t, prop) + " %%}"
		}
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
	case types.KeyMap, types.KeyValueMap, types.KeyReference:
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

func AsRef(t types.Type) (*types.Reference, error) {
	w, ok := t.(*types.Wrapped)
	if ok {
		t = w.T
	}
	ref, ok := t.(*types.Reference)
	if !ok {
		return nil, errors.Errorf("InvalidType(%T)", w.T)
	}
	return ref, nil
}

func asRefK(t types.Type) string {
	ref, err := AsRef(t)
	if err != nil {
		return fmt.Sprintf("ERROR: %s", err.Error())
	}
	return ref.K
}
