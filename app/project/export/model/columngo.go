package model

import (
	"fmt"
	"strings"

	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/enum"
	"projectforge.dev/projectforge/app/util"
)

const msgTextarea, stringsSuffix = `{%%%%= edit.TextareaTable(%q, %q, %q, 8, util.ToJSON(%s), 5, %s) %%%%}`, ".Strings()"

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
	switch c.Type.Key() {
	case types.KeyAny:
		return fmt.Sprintf(msgTextarea, c.Camel(), id, c.Title(), c.ToGoString(prefix), h), nil
	case types.KeyBool:
		return fmt.Sprintf(`{%%%%= edit.BoolTable(%q, %q, %s, 5, %s) %%%%}`, c.Camel(), c.Title(), prefix+c.Proper(), h), nil
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
			call = fmt.Sprintf("string(%s)", c.ToGoString(prefix))
		}
		msg := `{%%%%= edit.SelectTable(%q, %q, %q, %s, %s, %s, 5, %s) %%%%}`
		return fmt.Sprintf(msg, c.Camel(), id, c.Title(), call, eKeys, eTitles, h), nil
	case types.KeyInt:
		return fmt.Sprintf(`{%%%%= edit.IntTable(%q, %q, %q, %s, 5, %s) %%%%}`, c.Camel(), id, c.Title(), prefix+c.Proper(), h), nil
	case types.KeyFloat:
		return fmt.Sprintf(`{%%%%= edit.FloatTable(%q, %q, %q, %s, 5, %s) %%%%}`, c.Camel(), id, c.Title(), prefix+c.Proper(), h), nil
	case types.KeyList:
		lt := types.TypeAs[*types.List](c.Type)
		e, _ := AsEnumInstance(lt.V, enums)
		if e != nil {
			return fmt.Sprintf(
				`{%%%%= edit.CheckboxTable(%q, %q, %s.Keys(), %s.All%s.Keys(), %s.All%s%s, 5, %s.All%s.Help()) %%%%}`,
				c.Camel(), c.Title(), c.ToGoString(prefix), e.Package, e.ProperPlural(), e.Package, e.ProperPlural(), stringsSuffix, e.Package, e.ProperPlural(),
			), nil
		}
		if c.Display == FmtTags.Key && lt.V.Key() == types.KeyString {
			msg := `{%%%%= edit.TagsTable(%q, util.StringToTitle(%q), %q, %s, ps, 5, %s) %%%%}`
			return fmt.Sprintf(msg, c.Camel(), id, c.Title(), c.ToGoString(prefix), h), nil
		}
		return fmt.Sprintf(msgTextarea, c.Camel(), id, c.Title(), c.ToGoString(prefix), h), nil
	case types.KeyMap, types.KeyValueMap:
		return fmt.Sprintf(msgTextarea, c.Camel(), id, c.Title(), prefix+c.Proper(), h), nil
	case types.KeyReference:
		return fmt.Sprintf(msgTextarea, c.Camel(), id, c.Title(), prefix+c.Proper(), h), nil
	case types.KeyDate:
		gs := c.ToGoString(prefix)
		if !c.Nullable {
			gs = "&" + gs
		}
		return fmt.Sprintf(`{%%%%= edit.TimestampDayTable(%q, %q, %q, %s, 5, %s) %%%%}`, c.Camel(), id, c.Title(), gs, h), nil
	case types.KeyTimestamp:
		gs := c.ToGoString(prefix)
		if !c.Nullable {
			gs = "&" + gs
		}
		return fmt.Sprintf(`{%%%%= edit.TimestampTable(%q, %q, %q, %s, 5, %s) %%%%}`, c.Camel(), id, c.Title(), gs, h), nil
	case types.KeyUUID:
		gs := prefix + c.Proper()
		if !c.Nullable {
			gs = "&" + gs
		}
		return fmt.Sprintf(`{%%%%= edit.UUIDTable(%q, %q, %q, %s, 5, %s) %%%%}`, c.Camel(), id, c.Title(), gs, h), nil
	case types.KeyString:
		switch format {
		case FmtCode.Key, FmtCodeHidden.Key, FmtHTML.Key, FmtJSON.Key, FmtSQL.Key:
			return fmt.Sprintf(`{%%%%= edit.TextareaTable(%q, %q, %q, 8, %s, 5, %s) %%%%}`, c.Camel(), id, c.Title(), c.ToGoString(prefix), h), nil
		case FmtSelect.Key:
			if len(c.Values) == 0 {
				return fmt.Sprintf(`{%%%%= edit.Table(%q, %q, %q, %s, 5, %s) %%%%}`, c.Camel(), id, c.Title(), c.ToGoString(prefix), h), nil
			}
			sel := `{%%%%= edit.SelectTable(%q, "", %q, %s, %s, nil, 5, %s) %%%%}`
			opts := "[]string{" + strings.Join(util.StringArrayQuoted(c.Values), ", ") + "}"
			return fmt.Sprintf(sel, c.Camel(), c.Title(), c.ToGoString(prefix), opts, h), nil
		default:
			return fmt.Sprintf(`{%%%%= edit.StringTable(%q, %q, %q, %s, 5, %s) %%%%}`, c.Camel(), id, c.Title(), c.ToGoString(prefix), h), nil
		}
	default:
		return fmt.Sprintf(`{%%%%= edit.StringTable(%q, %q, %q, %s, 5, %s) %%%%}`, c.Camel(), id, c.Title(), c.ToGoString(prefix), h), nil
	}
}

func (c *Column) ToGoMapParse() string {
	return toGoMapParse(c.Type)
}

func toGoMapParse(t types.Type) string {
	switch t.Key() {
	case types.KeyAny:
		return ""
	case types.KeyBool:
		return "Bool"
	case types.KeyInt:
		return "Int"
	case types.KeyFloat:
		return "Float"
	case types.KeyList:
		l := types.TypeAs[*types.List](t)
		if l == nil {
			return fmt.Sprintf("ERROR:invalid list type [%T]", t)
		}
		return "Array" + toGoMapParse(l.V)
	case types.KeyMap, types.KeyValueMap:
		return "Map"
	case types.KeyReference:
		return asRefK(t)
	case types.KeyString, types.KeyEnum:
		return "String"
	case types.KeyDate, types.KeyTimestamp:
		return "Time"
	case types.KeyUUID:
		return "UUID"
	default:
		return "ERROR:unhandled map parse for type [" + t.Key() + "]"
	}
}

func (c *Column) ZeroVal() string {
	if c.Nullable {
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
	case types.KeyMap, types.KeyValueMap, types.KeyReference:
		return types.KeyNil
	case types.KeyString:
		return `""`
	case types.KeyDate, types.KeyTimestamp:
		return "time.Time{}"
	case types.KeyUUID:
		return "uuid.UUID{}"
	default:
		return "ERROR:unhandled zero value for type [" + c.Type.Key() + "]"
	}
}
