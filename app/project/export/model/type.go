package model

import (
	"fmt"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/enum"
)

const (
	keyJSONB          = "jsonb"
	goTypeIntArray    = "[]int"
	goTypeMapArray    = "[]util.ValueMap"
	goTypeStringArray = "[]string"
	goTypeAnyArray    = "[]any"
)

func ToGoType(t types.Type, nullable bool, pkg string, enums enum.Enums) (string, error) {
	var ret string
	switch t.Key() {
	case types.KeyAny:
		ret = "any"
	case types.KeyBool:
		ret = types.KeyBool
	case types.KeyEnum:
		e, err := AsEnumInstance(t, enums)
		if err != nil {
			return "", err
		}
		if e.Package == pkg {
			ret = e.Proper()
		} else {
			ret = e.Package + "." + e.Proper()
		}
	case types.KeyInt:
		ret = types.KeyInt
	case types.KeyFloat:
		ret = "float64"
	case types.KeyList:
		lt := types.TypeAs[*types.List](t)
		switch lt.V.Key() {
		case types.KeyString:
			ret = goTypeStringArray
		case types.KeyInt:
			ret = goTypeIntArray
		case types.KeyEnum:
			e, err := AsEnumInstance(lt.V, enums)
			if err != nil {
				return "", err
			}
			if e.Package == pkg {
				ret = e.ProperPlural()
			} else {
				ret = e.Package + "." + e.ProperPlural()
			}
		case types.KeyMap, types.KeyValueMap:
			ret = goTypeMapArray
		default:
			ret = goTypeAnyArray
		}
	case types.KeyMap, types.KeyValueMap:
		ret = "util.ValueMap"
	case types.KeyReference:
		ref, err := AsRef(t)
		if err != nil {
			return "", err
		}
		if ref.Pkg.Last() == pkg {
			ret = fmt.Sprintf("*%s", ref.K)
		} else {
			ret = fmt.Sprintf("*%s.%s", ref.Pkg.Last(), ref.K)
		}
	case types.KeyString:
		ret = types.KeyString
	case types.KeyDate, types.KeyTimestamp:
		ret = "time.Time"
	case types.KeyUUID:
		ret = "uuid.UUID"
	default:
		return "", errors.Errorf("ERROR:Unhandled[%s]", t.Key())
	}
	if nullable && !t.Scalar() {
		return "*" + ret, nil
	}
	return ret, nil
}

func ToGoString(t types.Type, nullable bool, prop string, alwaysString bool) string {
	switch t.Key() {
	case types.KeyAny, types.KeyBool, types.KeyInt, types.KeyFloat:
		return fmt.Sprintf("fmt.Sprint(%s)", prop)
	case types.KeyList:
		if alwaysString {
			return fmt.Sprintf("util.ToJSON(%s)", prop)
		}
		return prop
	case types.KeyEnum:
		if alwaysString {
			return fmt.Sprintf("%s.String()", prop)
		}
		return prop
	case types.KeyMap, types.KeyValueMap:
		return fmt.Sprintf("util.ToJSON(%s)", prop)
	case types.KeyDate:
		if alwaysString {
			if nullable {
				return fmt.Sprintf("util.TimeToYMD(%s)", prop)
			}
			return fmt.Sprintf("util.TimeToYMD(&%s)", prop)
		}
		return prop
	case types.KeyTimestamp, types.KeyTimestampZoned:
		if alwaysString {
			if nullable {
				return fmt.Sprintf("util.TimeToFull(%s)", prop)
			}
			return fmt.Sprintf("util.TimeToFull(&%s)", prop)
		}
		return prop
	case types.KeyUUID, types.KeyReference:
		if alwaysString && nullable {
			return fmt.Sprintf("util.StringNullable(%s)", prop)
		}
		return fmt.Sprintf("%s.String()", prop)
	default:
		return prop
	}
}
