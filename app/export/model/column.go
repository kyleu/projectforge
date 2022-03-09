package model

import (
	"fmt"
	"strings"

	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/util"
)

type Column struct {
	Name       string         `json:"name"`
	Type       *types.Wrapped `json:"type"`
	PK         bool           `json:"pk,omitempty"`
	Nullable   bool           `json:"nullable,omitempty"`
	Search     bool           `json:"search,omitempty"`
	SQLDefault string         `json:"sqlDefault,omitempty"`
	Display    string         `json:"display,omitempty"`
	Tags       []string       `json:"tags,omitempty"`
}

func (c *Column) Clone() *Column {
	return &Column{Name: c.Name, Type: c.Type, PK: c.PK, Nullable: c.Nullable, Search: c.Search, SQLDefault: c.SQLDefault, Display: c.Display, Tags: c.Tags}
}

func (c *Column) NameQuoted() string {
	return fmt.Sprintf("%q", c.Name)
}

func (c *Column) Camel() string {
	return util.StringToLowerCamel(c.Name)
}

func (c *Column) CamelPlural() string {
	return util.StringToPlural(c.Camel())
}

func (c *Column) Proper() string {
	return util.StringToCamel(c.Name)
}

func (c *Column) Plural() string {
	ret := util.StringToPlural(c.Name)
	if ret == c.Name {
		return ret + "Set"
	}
	return ret
}

func (c *Column) ProperPlural() string {
	return util.StringToPlural(c.Proper())
}

func (c *Column) HasTag(t string) bool {
	return util.StringArrayContains(c.Tags, t)
}

func (c *Column) ToGoString(prefix string) string {
	return ToGoString(c.Type, prefix+c.Proper())
}

func (c *Column) ToGoViewString(prefix string) string {
	return ToGoViewString(c.Type, prefix+c.Proper(), c.Nullable)
}

func (c *Column) ToGoType() string {
	return ToGoType(c.Type, c.Nullable)
}

func (c *Column) ToGoDTOType() string {
	return ToGoDTOType(c.Type, c.Nullable)
}

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

func (c *Column) ToGoEditString(prefix string) string {
	switch c.Type.Key() {
	case types.KeyAny:
		return fmt.Sprintf(`{%%%%= components.TableTextarea(%q, %q, 8, util.ToJSON(%s), 5, %q) %%%%}`, c.Camel(), c.Title(), c.ToGoString(prefix), c.Help())
	case types.KeyBool:
		return fmt.Sprintf(`{%%%%= components.TableBoolean(%q, %q, %s, 5, %q) %%%%}`, c.Camel(), c.Title(), prefix+c.Proper(), c.Help())
	case types.KeyInt:
		return fmt.Sprintf(`{%%%%= components.TableInputNumber(%q, %q, %s, 5, %q) %%%%}`, c.Camel(), c.Title(), prefix+c.Proper(), c.Help())
	case types.KeyMap:
		return fmt.Sprintf(`{%%%%= components.TableTextarea(%q, %q, 8, util.ToJSON(%s), 5, %q) %%%%}`, c.Camel(), c.Title(), c.ToGoString(prefix), c.Help())
	case types.KeyTimestamp:
		gs := c.ToGoString(prefix)
		if !c.Nullable {
			gs = "&" + gs
		}
		return fmt.Sprintf(`{%%%%= components.TableInputTimestamp(%q, %q, %s, 5, %q) %%%%}`, c.Camel(), c.Title(), gs, c.Help())
	default:
		return fmt.Sprintf(`{%%%%= components.TableInput(%q, %q, %s, 5, %q) %%%%}`, c.Camel(), c.Title(), c.ToGoString(prefix), c.Help())
	}
}

func (c *Column) ToGoMapParse() string {
	switch c.Type.Key() {
	case types.KeyAny:
		return "Interface"
	case types.KeyBool:
		return "Bool"
	case types.KeyInt:
		return "Int"
	case types.KeyMap:
		return "Map"
	case types.KeyString:
		return "String"
	case types.KeyTimestamp:
		return "Time"
	case types.KeyUUID:
		return "UUID"
	default:
		return "ERROR:unhandled map parse for type [" + c.Type.Key() + "]"
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
	case types.KeyInt:
		return "0"
	case types.KeyMap:
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

func (c *Column) Title() string {
	return util.StringToTitle(c.Name)
}

func (c *Column) TitleLower() string {
	return strings.ToLower(c.Title())
}

func (c *Column) ShouldDisplay(k string) bool {
	switch c.Display {
	case "detail":
		return k == c.Display
	default:
		return true
	}
}

func (c *Column) BC() string {
	if c == nil {
		return ""
	}
	return fmt.Sprintf(", %q", c.Camel())
}
