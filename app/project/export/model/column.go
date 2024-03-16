package model

import (
	"fmt"
	"slices"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/enum"
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
	// {Key: "helpstring", Title: "Help String", Description: "The X of the column"},
}

type Column struct {
	Name       string         `json:"name"`
	Type       *types.Wrapped `json:"type"`
	PK         bool           `json:"pk,omitempty"`
	Nullable   bool           `json:"nullable,omitempty"`
	Search     bool           `json:"search,omitempty"`
	SQLDefault string         `json:"sqlDefault,omitempty"`
	Indexed    bool           `json:"indexed,omitempty"`
	Display    string         `json:"display,omitempty"`
	Format     string         `json:"format,omitempty"`
	JSON       string         `json:"json,omitempty"`
	Example    string         `json:"example,omitempty"`
	Validation string         `json:"validation,omitempty"`
	Values     []string       `json:"values,omitempty"`
	Tags       []string       `json:"tags,omitempty"`
	HelpString string         `json:"helpString,omitempty"`
}

func (c *Column) Clone() *Column {
	return &Column{
		Name: c.Name, Type: c.Type, PK: c.PK, Nullable: c.Nullable, Search: c.Search, SQLDefault: c.SQLDefault, Display: c.Display,
		Format: c.Format, Example: c.Example, JSON: c.JSON, Validation: c.Validation, Values: c.Values, Tags: c.Tags, HelpString: c.HelpString,
	}
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
	return lo.Contains(c.Tags, t)
}

func (c *Column) AddTag(t string) {
	if !lo.Contains(c.Tags, t) {
		c.Tags = append(c.Tags, t)
		slices.Sort(c.Tags)
	}
}

func (c *Column) ToGoString(prefix string) string {
	return ToGoString(c.Type, c.Nullable, prefix+c.Proper(), false)
}

func (c *Column) ToGoViewString(prefix string, verbose bool, url bool, enums enum.Enums, src string) string {
	return ToGoViewString(c.Type, prefix+c.Proper(), c.Nullable, c.Format, verbose, url, enums, src)
}

func (c *Column) ToGoType(pkg string, enums enum.Enums) (string, error) {
	return ToGoType(c.Type, c.Nullable, pkg, enums)
}

func (c *Column) ToGoRowType(pkg string, enums enum.Enums, database string) (string, error) {
	return ToGoRowType(c.Type, c.Nullable, pkg, enums, database)
}

func (c *Column) ShouldDisplay(k string) bool {
	switch c.Display {
	case util.KeyDetail:
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

func (c *Column) NeedsErr(_ string, database string) bool {
	if database == util.DatabaseSQLServer && c.Type.Key() == types.KeyUUID && c.Nullable {
		return true
	}
	if c.Type.Scalar() {
		return true
	}
	if c.Nullable && (c.Type.Key() == types.KeyDate || c.Type.Key() == types.KeyTimestamp || c.Type.Key() == types.KeyReference) {
		return true
	}
	return false
}

func (c *Column) RemoveTag(t string) {
	if idx := slices.Index(c.Tags, t); idx > -1 {
		c.Tags = append(c.Tags[:idx], c.Tags[idx+1:]...)
	}
}
