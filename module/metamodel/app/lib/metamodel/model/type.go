package model

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"{{{ .Package }}}/app/lib/metamodel/enum"
	"{{{ .Package }}}/app/lib/types"
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
	case types.KeyAny, types.KeyJSON:
		ret = types.KeyAny
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
		ret = t.String()
	case types.KeyFloat:
		ret = "float64"
	case types.KeyList:
		lt := types.Wrap(t).ListType()
		switch lt.Key() {
		case types.KeyString:
			ret = goTypeStringArray
		case types.KeyInt:
			ret = goTypeIntArray
		case types.KeyEnum:
			e, err := AsEnumInstance(lt, enums)
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
	case types.KeyOrderedMap:
		ret = "*util.OrderedMap[any]"
	case types.KeyNumeric:
		ret = "numeric.Numeric"
	case types.KeyNumericMap:
		ret = "numeric.NumericMap"
	case types.KeyReference:
		ref, err := AsRef(t)
		if err != nil {
			return "", err
		}
		ret = ref.LastRef(ref.Pkg.Last() != pkg)
	case types.KeyString:
		ret = types.KeyString
	case types.KeyDate, types.KeyTimestamp, types.KeyTimestampZoned:
		ret = "time.Time"
	case types.KeyUUID:
		ret = "uuid.UUID"
	default:
		return "", errors.Errorf("ERROR:Unhandled[%s]", t.Key())
	}
	if nullable && !t.Scalar() && t.Key() != types.KeyReference {
		return "*" + strings.TrimPrefix(ret, "*"), nil
	}
	return ret, nil
}

func ToGoString(t types.Type, nullable bool, prop string, alwaysString bool) string {
	switch t.Key() {
	case types.KeyAny, types.KeyBool, types.KeyInt, types.KeyFloat:
		return fmt.Sprintf("fmt.Sprint(%s)", prop)
	case types.KeyList, types.KeyJSON:
		if alwaysString {
			return fmt.Sprintf("util.ToJSONCompact(%s)", prop)
		}
		return prop
	case types.KeyEnum:
		if alwaysString {
			return fmt.Sprintf("%s.String()", prop)
		}
		return prop
	case types.KeyMap, types.KeyOrderedMap, types.KeyValueMap, types.KeyNumericMap:
		return fmt.Sprintf("util.ToJSONCompact(%s)", prop)
	case types.KeyDate:
		if alwaysString {
			if nullable {
				return fmt.Sprintf("util.TimeToYMD(%s)", prop)
			}
			return fmt.Sprintf("util.TimeToYMD(&%s)", prop)
		}
		return prop
	case types.KeyNumeric:
		return fmt.Sprintf("%s.String()", prop)
	case types.KeyTimestamp, types.KeyTimestampZoned:
		if alwaysString {
			if nullable {
				return fmt.Sprintf("util.TimeToFull(%s)", prop)
			}
			return fmt.Sprintf("util.TimeToFull(&%s)", prop)
		}
		return prop
	case types.KeyUUID:
		if alwaysString && nullable {
			return fmt.Sprintf("util.StringNullable(%s)", prop)
		}
		return fmt.Sprintf("%s.String()", prop)
	case types.KeyReference:
		return fmt.Sprintf("util.ToJSONCompact(%s)", prop)
	default:
		return prop
	}
}
