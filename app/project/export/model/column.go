package model

import (
	"fmt"
	"strings"

	"golang.org/x/exp/slices"

	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/util"
)

const (
	FmtCode   = "code"
	FmtURL    = "url"
	FmtSelect = "select"
)

type Column struct {
	Name       string         `json:"name"`
	Type       *types.Wrapped `json:"type"`
	PK         bool           `json:"pk,omitempty"`
	Nullable   bool           `json:"nullable,omitempty"`
	Search     bool           `json:"search,omitempty"`
	SQLDefault string         `json:"sqlDefault,omitempty"`
	Display    string         `json:"display,omitempty"`
	Format     string         `json:"format,omitempty"`
	Values     []string       `json:"values,omitempty"`
	Tags       []string       `json:"tags,omitempty"`
	HelpString string         `json:"helpString,omitempty"`
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

func (c *Column) Title() string {
	return util.StringToTitle(c.Name)
}

func (c *Column) TitleLower() string {
	return strings.ToLower(c.Title())
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
	return slices.Contains(c.Tags, t)
}

func (c *Column) ToGoString(prefix string) string {
	return ToGoString(c.Type, prefix+c.Proper())
}

func (c *Column) ToGoViewString(prefix string) string {
	return ToGoViewString(c.Type, prefix+c.Proper(), c.Nullable, c.Format)
}

func (c *Column) ToGoType(pkg string) string {
	return ToGoType(c.Type, c.Nullable, pkg)
}

func (c *Column) ToGoDTOType(pkg string) string {
	return ToGoDTOType(c.Type, c.Nullable, pkg)
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
