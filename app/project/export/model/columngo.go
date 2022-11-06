package model

import (
	"fmt"
	"strings"

	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/enum"
	"projectforge.dev/projectforge/app/util"
)

func (c *Column) ToSQLType() (string, error) {
	ret, err := ToSQLType(c.Type)
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

func (c *Column) ToGoEditString(prefix string, format string, enums enum.Enums) (string, error) {
	h, err := c.Help(enums)
	if err != nil {
		return "", err
	}
	switch c.Type.Key() {
	case types.KeyAny:
		return fmt.Sprintf(`{%%%%= components.TableTextarea(%q, %q, 8, util.ToJSON(%s), 5, %q) %%%%}`, c.Camel(), c.Title(), c.ToGoString(prefix), h), nil
	case types.KeyBool:
		return fmt.Sprintf(`{%%%%= components.TableBoolean(%q, %q, %s, 5, %q) %%%%}`, c.Camel(), c.Title(), prefix+c.Proper(), h), nil
	// case types.KeyEnum:
	//	return fmt.Sprintf(`{%%%%= components.TableInput(%q, %q, string(%s), 5, %q) %%%%}`, c.Camel(), c.Title(), c.ToGoString(prefix), h), nil
	case types.KeyEnum:
		e, err := AsEnumInstance(c.Type, enums)
		if err != nil {
			return "", err
		}
		eRef := strings.Join(util.StringArrayQuoted(e.Values), ", ")
		msg := `{%%%%= components.TableSelect(%q, %q, string(%s), []string{%s}, []string{%s}, 5, %q) %%%%}`
		return fmt.Sprintf(msg, c.Camel(), c.Title(), c.ToGoString(prefix), eRef, eRef, h), nil
	case types.KeyInt:
		return fmt.Sprintf(`{%%%%= components.TableInputNumber(%q, %q, %s, 5, %q) %%%%}`, c.Camel(), c.Title(), prefix+c.Proper(), h), nil
	case types.KeyFloat:
		return fmt.Sprintf(`{%%%%= components.TableInputFloat(%q, %q, %s, 5, %q) %%%%}`, c.Camel(), c.Title(), prefix+c.Proper(), h), nil
	case types.KeyList, types.KeyMap, types.KeyValueMap:
		return fmt.Sprintf(`{%%%%= components.TableTextarea(%q, %q, 8, util.ToJSON(%s), 5, %q) %%%%}`, c.Camel(), c.Title(), c.ToGoString(prefix), h), nil
	case types.KeyReference:
		return fmt.Sprintf(`{%%%%= components.TableTextarea(%q, %q, 8, util.ToJSON(%s), 5, %q) %%%%}`, c.Camel(), c.Title(), prefix+c.Proper(), h), nil
	case types.KeyDate:
		gs := c.ToGoString(prefix)
		if !c.Nullable {
			gs = "&" + gs
		}
		return fmt.Sprintf(`{%%%%= components.TableInputTimestampDay(%q, %q, %s, 5, %q) %%%%}`, c.Camel(), c.Title(), gs, h), nil
	case types.KeyTimestamp:
		gs := c.ToGoString(prefix)
		if !c.Nullable {
			gs = "&" + gs
		}
		return fmt.Sprintf(`{%%%%= components.TableInputTimestamp(%q, %q, %s, 5, %q) %%%%}`, c.Camel(), c.Title(), gs, h), nil
	case types.KeyUUID:
		gs := prefix + c.Proper()
		if !c.Nullable {
			gs = "&" + gs
		}
		return fmt.Sprintf(`{%%%%= components.TableInputUUID(%q, %q, %s, 5, %q) %%%%}`, c.Camel(), c.Title(), gs, h), nil
	case types.KeyString:
		switch format {
		case FmtCode:
			return fmt.Sprintf(`{%%%%= components.TableTextarea(%q, %q, 8, %s, 5, %q) %%%%}`, c.Camel(), c.Title(), c.ToGoString(prefix), h), nil
		case FmtSelect:
			if len(c.Values) == 0 {
				return fmt.Sprintf(`{%%%%= components.TableInput(%q, %q, %s, 5, %q) %%%%}`, c.Camel(), c.Title(), c.ToGoString(prefix), h), nil
			}
			sel := `{%%%%= components.TableSelect(%q, %q, %s, %s, nil, 5, %q) %%%%}`
			opts := "[]string{" + strings.Join(util.StringArrayQuoted(c.Values), ", ") + "}"
			return fmt.Sprintf(sel, c.Camel(), c.Title(), c.ToGoString(prefix), opts, h), nil
		default:
			return fmt.Sprintf(`{%%%%= components.TableInput(%q, %q, %s, 5, %q) %%%%}`, c.Camel(), c.Title(), c.ToGoString(prefix), h), nil
		}
	default:
		return fmt.Sprintf(`{%%%%= components.TableInput(%q, %q, %s, 5, %q) %%%%}`, c.Camel(), c.Title(), c.ToGoString(prefix), h), nil
	}
}

func (c *Column) ToGoMapParse() string {
	return toGoMapParse(c.Type)
}

func toGoMapParse(t types.Type) string {
	switch t.Key() {
	case types.KeyAny:
		return "Interface"
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
		return "false"
	case types.KeyList:
		return types.KeyNil
	case types.KeyInt:
		return "0"
	case types.KeyFloat:
		return "0.0"
	case types.KeyMap, types.KeyValueMap, types.KeyReference:
		return types.KeyNil
	case types.KeyString:
		return "\"\""
	case types.KeyDate, types.KeyTimestamp:
		return "time.Time{}"
	case types.KeyUUID:
		return "uuid.UUID{}"
	default:
		return "ERROR:unhandled zero value for type [" + c.Type.Key() + "]"
	}
}
