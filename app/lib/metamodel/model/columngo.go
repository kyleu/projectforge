package model

import (
	"fmt"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/util"
)

const msgTextarea, stringsSuffix = `{%%%%= edit.TextareaTable(%q, %q, %q, 8, util.ToJSON(%s), 5, %s) %%%%}`, ".Strings()"

func (c *Column) ToGoString(prefix string) string {
	return ToGoString(c.Type, c.Nullable, prefix+c.Proper(), false)
}

func (c *Column) ToGoViewString(prefix string, verbose bool, url bool, enums enum.Enums, src string) string {
	prop := prefix + c.ProperDerived()
	return ToGoViewString(c.Type, prop, c.Nullable, c.Format, c.Display, verbose, url, enums, src)
}

func (c *Column) ToGoType(pkg string, enums enum.Enums) (string, error) {
	return ToGoType(c.Type, c.Nullable, pkg, enums)
}

func (c *Column) ToGoRowType(pkg string, enums enum.Enums, database string) (string, error) {
	return ToGoRowType(c.Type, c.Nullable, pkg, enums, database)
}

func (c *Column) ToSQLType(database string) (string, error) {
	ret, err := ToSQLType(c.Type, database)
	if err != nil {
		return "", err
	}

	if !c.Nullable {
		ret += " not null"
	}
	if c.SQLDefault != "" {
		ret += " default " + c.SQLDefault
	}
	return ret, nil
}

func (c *Column) ToGoEditString(prefix string, format string, id string, enums enum.Enums) (string, error) {
	h, err := c.Help(enums)
	if err != nil {
		return "", err
	}
	prop := c.ToGoString(prefix)
	preprop := prefix + c.Proper()
	key := c.CamelNoReplace()
	switch c.Type.Key() {
	case types.KeyAny:
		return fmt.Sprintf(msgTextarea, key, id, c.Title(), prop, h), nil
	case types.KeyBool:
		return fmt.Sprintf(`{%%%%= edit.BoolTable(%q, %q, %s, 5, %s) %%%%}`, key, c.Title(), preprop, h), nil
	case types.KeyEnum:
		e, err := AsEnumInstance(c.Type, enums)
		if err != nil {
			return "", err
		}
		ePrefix := fmt.Sprintf("All%s", e.ProperPlural())
		if e.Package != "" {
			ePrefix = e.Package + "." + ePrefix
		}
		eKeys := ePrefix + ".Keys()"
		eTitles := ePrefix + stringsSuffix
		call := fmt.Sprintf("%s%s.Key", prefix, c.Proper())
		if e.Simple() {
			call = fmt.Sprintf("string(%s)", prop)
		}
		msg := `{%%%%= edit.SelectTable(%q, %q, %q, %s, %s, %s, 5, %s) %%%%}`
		return fmt.Sprintf(msg, key, id, c.Title(), call, eKeys, eTitles, h), nil
	case types.KeyFloat:
		return fmt.Sprintf(`{%%%%= edit.FloatTable(%q, %q, %q, %s, 5, %s) %%%%}`, key, id, c.Title(), preprop, h), nil
	case types.KeyInt:
		if types.Bits(c.Type) != 0 {
			preprop = "int(" + preprop + ")"
		}
		return fmt.Sprintf(`{%%%%= edit.IntTable(%q, %q, %q, %s, 5, %s) %%%%}`, key, id, c.Title(), preprop, h), nil
	case types.KeyJSON:
		return fmt.Sprintf(msgTextarea, key, id, c.Title(), prop, h), nil
	case types.KeyList:
		lt := c.Type.ListType()
		e, _ := AsEnumInstance(lt, enums)
		if e != nil {
			return fmt.Sprintf(
				`{%%%%= edit.CheckboxTable(%q, %q, %s.Keys(), %s.All%s.Keys(), %s.All%s%s, 5, %s.All%s.Help()) %%%%}`,
				key, c.Title(), prop, e.Package, e.ProperPlural(), e.Package, e.ProperPlural(), stringsSuffix, e.Package, e.ProperPlural(),
			), nil
		}
		if c.Display == FmtTags.Key && lt.Key() == types.KeyString {
			msg := `{%%%%= edit.TagsTable(%q, util.StringToTitle(%q), %q, %s, ps, 5, %s) %%%%}`
			return fmt.Sprintf(msg, key, id, c.Title(), prop, h), nil
		}
		return fmt.Sprintf(msgTextarea, key, id, c.Title(), prop, h), nil
	case types.KeyMap, types.KeyValueMap, types.KeyOrderedMap:
		return fmt.Sprintf(msgTextarea, key, id, c.Title(), preprop, h), nil
	case types.KeyNumeric:
		return fmt.Sprintf(`{%%%%= edit.NumericTable(%q, %q, %q, %s, 5, %s) %%%%}`, key, id, c.Title(), preprop, h), nil
	case types.KeyReference:
		return fmt.Sprintf(msgTextarea, key, id, c.Title(), preprop, h), nil
	case types.KeyDate:
		gs := prop
		if !c.Nullable {
			gs = "&" + gs
		}
		return fmt.Sprintf(`{%%%%= edit.TimestampDayTable(%q, %q, %q, %s, 5, %s) %%%%}`, key, id, c.Title(), gs, h), nil
	case types.KeyTimestamp, types.KeyTimestampZoned:
		gs := prop
		if !c.Nullable {
			gs = "&" + gs
		}
		return fmt.Sprintf(`{%%%%= edit.TimestampTable(%q, %q, %q, %s, 5, %s) %%%%}`, key, id, c.Title(), gs, h), nil
	case types.KeyUUID:
		gs := prefix + c.Proper()
		if !c.Nullable {
			gs = "&" + gs
		}
		return fmt.Sprintf(`{%%%%= edit.UUIDTable(%q, %q, %q, %s, 5, %s) %%%%}`, key, id, c.Title(), gs, h), nil
	case types.KeyString:
		switch format {
		case FmtCode.Key, FmtCodeHidden.Key, FmtHTML.Key, FmtJSON.Key, FmtMarkdown.Key, FmtSQL.Key:
			return fmt.Sprintf(`{%%%%= edit.TextareaTable(%q, %q, %q, 8, %s, 5, %s) %%%%}`, key, id, c.Title(), prop, h), nil
		case FmtColor.Key:
			return fmt.Sprintf(`{%%%%= edit.ColorTable(%q, %q, %q, %s, 5, %s) %%%%}`, key, id, c.Title(), prop, h), nil
		case FmtSelect.Key:
			if len(c.Values) == 0 {
				return fmt.Sprintf(`{%%%%= edit.Table(%q, %q, %q, %s, 5, %s) %%%%}`, key, id, c.Title(), prop, h), nil
			}
			sel := `{%%%%= edit.SelectTable(%q, "", %q, %s, %s, nil, 5, %s) %%%%}`
			opts := "[]string{" + util.StringJoin(util.StringArrayQuoted(c.Values), ", ") + "}"
			return fmt.Sprintf(sel, key, c.Title(), prop, opts, h), nil
		default:
			return fmt.Sprintf(`{%%%%= edit.StringTable(%q, %q, %q, %s, 5, %s) %%%%}`, key, id, c.Title(), prop, h), nil
		}
	default:
		return fmt.Sprintf(`{%%%%= edit.StringTable(%q, %q, %q, %s, 5, %s) %%%%}`, key, id, c.Title(), prop, h), nil
	}
}

func (c *Column) ToGoMapParse() (string, error) {
	return toGoMapParse(c.Type)
}

func toGoMapParse(t types.Type) (string, error) {
	switch t.Key() {
	case types.KeyAny:
		return "Unknown", nil
	case types.KeyBool:
		return "Bool", nil
	case types.KeyFloat:
		return "Float", nil
	case types.KeyInt:
		if types.Bits(t) == 16 {
			return "Int16", nil
		}
		if types.Bits(t) == 32 {
			return "Int32", nil
		}
		if types.Bits(t) == 64 {
			return "Int64", nil
		}
		return "Int", nil
	case types.KeyJSON:
		return "JSON", nil
	case types.KeyList:
		l := types.TypeAs[*types.List](t)
		if l == nil {
			return fmt.Sprintf("ERROR:invalid list type [%T]", t), nil
		}
		x, err := toGoMapParse(l.V)
		if err != nil {
			return "", err
		}
		return "Array" + x, nil
	case types.KeyMap, types.KeyValueMap:
		return "Map", nil
	case types.KeyOrderedMap:
		return "OrderedMap", nil
	case types.KeyNumeric:
		return "Numeric", nil
	case types.KeyReference:
		return asRefK(t), nil
	case types.KeyString, types.KeyEnum:
		return "String", nil
	case types.KeyDate, types.KeyTimestamp, types.KeyTimestampZoned:
		return "Time", nil
	case types.KeyUUID:
		return "UUID", nil
	default:
		return "", errors.Errorf("ERROR:unhandled map parse for type [%s]", t.Key())
	}
}

func (c *Column) ZeroVal() string {
	if c.Nullable && !c.Type.Scalar() {
		return types.KeyNil
	}
	switch c.Type.Key() {
	case types.KeyAny:
		return types.KeyNil
	case types.KeyBool:
		return util.BoolFalse
	case types.KeyList:
		return types.KeyNil
	case types.KeyInt:
		return "0"
	case types.KeyFloat:
		return "0.0"
	case types.KeyMap, types.KeyValueMap, types.KeyOrderedMap, types.KeyReference:
		return types.KeyNil
	case types.KeyNumeric:
		return "numeric.Zero"
	case types.KeyString:
		return `""`
	case types.KeyDate, types.KeyTimestamp, types.KeyTimestampZoned:
		return "time.Time{}"
	case types.KeyUUID:
		return "uuid.UUID{}"
	default:
		return "ERROR:unhandled zero value for type [" + c.Type.Key() + "]"
	}
}
