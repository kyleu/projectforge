package model

import (
	"fmt"
	"slices"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/util"
)

var ColumnFieldDescs = util.FieldDescs{
	{Key: "name", Title: "Name", Description: "The name of the column"},
	{Key: "type", Title: "Type", Description: "The type of the column", Type: "type"},
	{Key: "pk", Title: "PK", Description: "Indicates if this column is a primary key", Type: "bool"},
	{Key: "nullable", Title: "Nullable", Description: "Indicates if this column is nullable", Type: "bool"},
	{Key: "search", Title: "Search", Description: "Indicates if this column is included in search", Type: "bool"},
	// {Key: "sqlDefault", Title: "SQL Default", Description: "The X of the column"},
	{Key: "indexed", Title: "Indexed", Description: "Indicates if this column is indexed", Type: "bool"},
	{Key: "display", Title: "Display", Description: "The display setting of the column's value"},
	{Key: "format", Title: "Format", Description: "The formatting applied to the column's value"},
	{Key: "example", Title: "Example", Description: "Example annotation for the column's value"},
	{Key: "json", Title: "JSON", Description: "JSON field name to use instead of [name]"},
	{Key: "validation", Title: "Validation", Description: "Validation annotation for the column's value"},
	// {Key: "values", Title: "Values", Description: "The X of the column"},
	// {Key: "tags", Title: "Tags", Description: "The X of the column"},
	// {Key: "comment", Title: "Comment", Description: "A comment added to generated code"},
	// {Key: "help", Title: "Help String", Description: "The X of the column"},
}

type Column struct {
	Name           string         `json:"name"`
	Type           *types.Wrapped `json:"type"`
	PK             bool           `json:"pk,omitzero"`
	Nullable       bool           `json:"nullable,omitzero"`
	Search         bool           `json:"search,omitzero"`
	SQLDefault     string         `json:"sqlDefault,omitzero"`
	Indexed        bool           `json:"indexed,omitzero"`
	Display        string         `json:"display,omitzero"`
	Format         string         `json:"format,omitzero"`
	JSON           string         `json:"json,omitzero"`
	SQLOverride    string         `json:"sql,omitzero"`
	TitleOverride  string         `json:"title,omitzero"`
	PluralOverride string         `json:"plural,omitzero"`
	ProperOverride string         `json:"proper,omitzero"`
	Example        string         `json:"example,omitzero"`
	Validation     string         `json:"validation,omitzero"`
	Values         []string       `json:"values,omitempty"`
	Tags           []string       `json:"tags,omitempty"`
	Comment        string         `json:"comment,omitzero"`
	Help           string         `json:"help,omitzero"`
	Metadata       util.ValueMap  `json:"metadata,omitzero"`
	acronyms       []string
}

func (c *Column) Clone() *Column {
	return &Column{
		Name: c.Name, Type: c.Type, PK: c.PK, Nullable: c.Nullable, Search: c.Search, SQLDefault: c.SQLDefault,
		Display: c.Display, Format: c.Format, Example: c.Example, JSON: c.JSON, Validation: c.Validation,
		Values: c.Values, Tags: c.Tags, Comment: c.Comment, Help: c.Help, Metadata: c.Metadata, acronyms: c.acronyms,
	}
}

func (c *Column) NameQuoted() string {
	return fmt.Sprintf("%q", c.Name)
}

func (c *Column) Initials() string {
	return util.StringToInitials(c.Name, c.acronyms...)
}

func (c *Column) Camel() string {
	return util.StringToCamel(c.Name, c.acronyms...)
}

func (c *Column) CamelNoReplace() string {
	return util.StringToCamel(c.Name)
}

func (c *Column) CamelPlural() string {
	return util.StringToPlural(c.Camel())
}

func (c *Column) Proper() string {
	if c.ProperOverride == "" {
		return util.StringToProper(c.Name, c.acronyms...)
	}
	return c.ProperOverride
}

func (c *Column) ProperDerived() string {
	if c.Derived() {
		return c.Proper() + "()"
	}
	return c.Proper()
}

func (c *Column) IconDerived() string {
	ret := c.ProperDerived()
	if c.Type.Key() == "enum" {
		ret += ".Icon"
	}
	return ret
}

func (c *Column) Title() string {
	if c.TitleOverride == "" {
		return util.StringToTitle(c.Name, c.acronyms...)
	}
	return c.TitleOverride
}

func (c *Column) TitleLower() string {
	return strings.ToLower(c.Title())
}

func (c *Column) Plural() string {
	if c.PluralOverride != "" {
		return c.PluralOverride
	}
	ret := util.StringToPlural(c.Name)
	if ret == c.Name {
		return ret + tSet
	}
	return ret
}

func (c *Column) ProperPlural() string {
	if len(c.Name) == 1 {
		return c.Proper() + tSet
	}
	return util.StringToPlural(c.Proper())
}

func (c *Column) HasTag(t string) bool {
	return lo.Contains(c.Tags, t)
}

func (c *Column) AddTag(t string) {
	if !lo.Contains(c.Tags, t) {
		c.Tags = util.ArraySorted(append(c.Tags, t))
	}
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

func (c *Column) SQL() string {
	return util.OrDefault(c.SQLOverride, c.Name)
}

var needsErr = []string{types.KeyDate, types.KeyTimestamp, types.KeyTimestampZoned, types.KeyReference}

func (c *Column) NeedsErr(_ string, db string) bool {
	if db == util.DatabaseSQLServer && c.Type.Key() == types.KeyUUID && c.Nullable {
		return true
	}
	if c.Type.Scalar() {
		return true
	}
	if c.Nullable && (!slices.Contains(needsErr, c.Type.Key())) {
		return true
	}
	return false
}

func (c *Column) RemoveTag(t string) {
	if idx := slices.Index(c.Tags, t); idx > -1 {
		c.Tags = append(c.Tags[:idx], c.Tags[idx+1:]...)
	}
}

func (c *Column) Derived() bool {
	return c.HasTag("derived")
}

func (c *Column) SetAcronyms(acronyms ...string) {
	c.acronyms = acronyms
}

func (c *Column) String() string {
	return fmt.Sprintf("%s: %s%s", c.Name, util.Choose(c.Nullable, "*", ""), c.Type.String())
}
