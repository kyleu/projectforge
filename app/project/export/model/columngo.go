package model

import (
	"fmt"
	"strings"

	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/util"
)

func (c *Column) ToSQLType() string {
	ret := ToSQLType(c.Type)
	if !c.Nullable {
		ret += " not null"
	}
	if c.SQLDefault != "" {
		ret += " default " + c.SQLDefault
	}
	return ret
}

func (c *Column) ToGoEditString(prefix string, format string) string {
	switch c.Type.Key() {
	case types.KeyAny:
		return fmt.Sprintf(`{%%%%= components.TableTextarea(%q, %q, 8, util.ToJSON(%s), 5, %q) %%%%}`, c.Camel(), c.Title(), c.ToGoString(prefix), c.Help())
	case types.KeyBool:
		return fmt.Sprintf(`{%%%%= components.TableBoolean(%q, %q, %s, 5, %q) %%%%}`, c.Camel(), c.Title(), prefix+c.Proper(), c.Help())
	case types.KeyInt:
		return fmt.Sprintf(`{%%%%= components.TableInputNumber(%q, %q, %s, 5, %q) %%%%}`, c.Camel(), c.Title(), prefix+c.Proper(), c.Help())
	case types.KeyFloat:
		return fmt.Sprintf(`{%%%%= components.TableInputFloat(%q, %q, %s, 5, %q) %%%%}`, c.Camel(), c.Title(), prefix+c.Proper(), c.Help())
	case types.KeyList, types.KeyMap, types.KeyValueMap:
		return fmt.Sprintf(`{%%%%= components.TableTextarea(%q, %q, 8, util.ToJSON(%s), 5, %q) %%%%}`, c.Camel(), c.Title(), c.ToGoString(prefix), c.Help())
	case types.KeyReference:
		return fmt.Sprintf(`{%%%%= components.TableTextarea(%q, %q, 8, util.ToJSON(%s), 5, %q) %%%%}`, c.Camel(), c.Title(), prefix+c.Proper(), c.Help())
	case types.KeyTimestamp:
		gs := c.ToGoString(prefix)
		if !c.Nullable {
			gs = "&" + gs
		}
		return fmt.Sprintf(`{%%%%= components.TableInputTimestamp(%q, %q, %s, 5, %q) %%%%}`, c.Camel(), c.Title(), gs, c.Help())
	case types.KeyUUID:
		gs := prefix + c.Proper()
		if !c.Nullable {
			gs = "&" + gs
		}
		return fmt.Sprintf(`{%%%%= components.TableInputUUID(%q, %q, %s, 5, %q) %%%%}`, c.Camel(), c.Title(), gs, c.Help())
	case types.KeyString:
		switch format {
		case FmtCode:
			return fmt.Sprintf(`{%%%%= components.TableTextarea(%q, %q, 8, %s, 5, %q) %%%%}`, c.Camel(), c.Title(), c.ToGoString(prefix), c.Help())
		case FmtSelect:
			if len(c.Values) == 0 {
				return fmt.Sprintf(`{%%%%= components.TableInput(%q, %q, %s, 5, %q) %%%%}`, c.Camel(), c.Title(), c.ToGoString(prefix), c.Help())
			}
			sel := `{%%%%= components.TableSelect(%q, %q, %s, %s, nil, 5, %q) %%%%}`
			opts := "[]string{" + strings.Join(util.StringArrayQuoted(c.Values), ", ") + "}"
			return fmt.Sprintf(sel, c.Camel(), c.Title(), c.ToGoString(prefix), opts, c.Help())
		default:
			return fmt.Sprintf(`{%%%%= components.TableInput(%q, %q, %s, 5, %q) %%%%}`, c.Camel(), c.Title(), c.ToGoString(prefix), c.Help())
		}
	default:
		return fmt.Sprintf(`{%%%%= components.TableInput(%q, %q, %s, 5, %q) %%%%}`, c.Camel(), c.Title(), c.ToGoString(prefix), c.Help())
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
		// l := types.TypeAs[*types.List](t)
		// if l == nil {
		// 	return fmt.Sprintf("ERROR:invalid list type [%T]", t)
		// }
		return "Array" // + toGoMapParse(l.V)
	case types.KeyMap, types.KeyValueMap:
		return "Map"
	case types.KeyReference:
		return asRefK(t)
	case types.KeyString:
		return "String"
	case types.KeyTimestamp:
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
	case types.KeyTimestamp:
		return "time.Time{}"
	case types.KeyUUID:
		return "uuid.UUID{}"
	default:
		return "ERROR:unhandled zero value for type [" + c.Type.Key() + "]"
	}
}
