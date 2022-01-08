package model

import (
	"fmt"

	"github.com/kyleu/projectforge/app/util"
)

type Column struct {
	Name       string   `json:"name"`
	Type       *Type    `json:"type"`
	PK         bool     `json:"pk,omitempty"`
	Nullable   bool     `json:"nullable,omitempty"`
	Search     bool     `json:"search,omitempty"`
	SQLDefault string   `json:"sqlDefault,omitempty"`
	Tags       []string `json:"tags,omitempty"`
}

func (c *Column) Camel() string {
	return util.StringToLowerCamel(c.Name)
}

func (c *Column) Proper() string {
	return util.StringToCamel(c.Name)
}

func (c *Column) Plural() string {
	return util.StringToPlural(c.Name)
}

func (c *Column) ProperPlural() string {
	return util.StringToPlural(c.Proper())
}

func (c *Column) HasTag(t string) bool {
	return util.StringArrayContains(c.Tags, t)
}

func (c *Column) ToGoString(prefix string) string {
	return c.Type.ToGoString(prefix + c.Proper())
}

func (c *Column) ToGoViewString(prefix string) string {
	return c.Type.ToGoViewString(prefix+c.Proper(), c.Nullable)
}

func (c *Column) ToGoType() string {
	return c.Type.ToGoType(c.Nullable)
}

func (c *Column) ToGoDTOType() string {
	return c.Type.ToGoDTOType(c.Nullable)
}

func (c *Column) ToSQLType() string {
	ret := c.Type.ToSQLType()
	if !c.Nullable {
		ret += " not null"
	}
	if c.SQLDefault != "" {
		ret += " default " + c.SQLDefault
	}
	return ret
}

func (c *Column) ToGoEditString(prefix string) string {
	switch c.Type.Key {
	case TypeInt.Key:
		return fmt.Sprintf(`{%%%%= components.TableInputNumber(%q, %q, %s, 5, %q) %%%%}`, c.Camel(), c.Title(), prefix+c.Proper(), c.Help())
	case TypeMap.Key:
		return fmt.Sprintf(`{%%%%= components.TableTextarea(%q, %q, 8, util.ToJSON(%s), 5, %q) %%%%}`, c.Camel(), c.Title(), c.ToGoString(prefix), c.Help())
	case TypeTimestamp.Key:
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
	suffix := ""
	if c.Nullable {
		suffix = "Opt"
	}
	switch c.Type.Key {
	case TypeInt.Key:
		return "Int" + suffix
	case TypeMap.Key:
		return "Map" + suffix
	case TypeString.Key:
		return "String" + suffix
	case TypeTimestamp.Key:
		return "Time" + suffix
	case TypeUUID.Key:
		return "UUID" + suffix
	default:
		return "ERROR:unhandled map parse for type [" + c.Type.Key + "]"
	}
}

const nilStr = "nil"

func (c *Column) ZeroVal() string {
	if c.Nullable {
		return nilStr
	}
	switch c.Type.Key {
	case TypeInt.Key:
		return "0"
	case TypeMap.Key:
		return nilStr
	case TypeString.Key:
		return "\"\""
	case TypeTimestamp.Key:
		return "time.Time{}"
	case TypeUUID.Key:
		return "uuid.UUID{}"
	default:
		return "ERROR:unhandled zero value for type [" + c.Type.Key + "]"
	}
}

func (c *Column) Title() string {
	return util.StringToTitle(c.Name)
}

func (c *Column) BC() string {
	if c == nil {
		return ""
	}
	return fmt.Sprintf(", %q", c.Camel())
}
