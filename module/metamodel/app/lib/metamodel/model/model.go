package model

import (
	"fmt"
	"slices"
	"strings"

	"github.com/samber/lo"

	"{{{ .Package }}}/app/lib/filter"
	"{{{ .Package }}}/app/lib/types"
	"{{{ .Package }}}/app/util"
)

type Model struct {
	Name           string           `json:"name"`
	Package        string           `json:"package"`
	Group          []string         `json:"group,omitempty"`
	Description    string           `json:"description,omitempty"`
	Icon           string           `json:"icon,omitempty"`
	Ordering       filter.Orderings `json:"ordering,omitempty"`
	SortIndex      int              `json:"sortIndex,omitempty"`
	View           string           `json:"view,omitempty"`
	Search         []string         `json:"search,omitempty"`
	Tags           []string         `json:"tags,omitempty"`
	TitleOverride  string           `json:"title,omitempty"`
	PluralOverride string           `json:"plural,omitempty"`
	ProperOverride string           `json:"proper,omitempty"`
	TableOverride  string           `json:"table,omitempty"`
	RouteOverride  string           `json:"route,omitempty"`
	Config         util.ValueMap    `json:"config,omitempty"`
	Columns        Columns          `json:"columns"`
	Relations      Relations        `json:"relations,omitempty"`
	Indexes        Indexes          `json:"indexes,omitempty"`
	SeedData       [][]any          `json:"seedData,omitempty"`
	Links          Links            `json:"links,omitempty"`
}

func (m *Model) HasTag(t string) bool {
	return lo.Contains(m.Tags, t)
}

func (m *Model) AddTag(t string) {
	if !m.HasTag(t) {
		m.Tags = append(m.Tags, t)
		slices.Sort(m.Tags)
	}
}

func (m *Model) PKs() Columns {
	return m.Columns.PKs()
}

func (m *Model) IsSoftDelete() bool {
	return m.HasTag("softDelete")
}

func (m *Model) SoftDeleteSuffix() string {
	if m.IsSoftDelete() {
		return ", true"
	}
	return ""
}

func (m *Model) RelationsFor(col *Column) Relations {
	return lo.Filter(m.Relations, func(r *Relation, _ int) bool {
		return lo.Contains(r.Src, col.Name)
	})
}

func (m *Model) CanTraverseRelation() bool {
	return len(m.PKs()) == 1
}

func (m *Model) PackageWithGroup(prefix string) string {
	if len(m.Group) == 0 {
		return prefix + m.Package
	}
	x := lo.Map(m.Group, func(g string, _ int) string {
		return strings.ToLower(prefix + g)
	})
	x = append(x, prefix+m.Package)
	return strings.Join(x, "/")
}

func (m *Model) GroupString(prefix string, dflt string) string {
	if len(m.Group) == 0 {
		return dflt
	}
	x := lo.Map(m.Group, func(g string, _ int) string {
		return strings.ToLower(prefix + g)
	})
	return strings.Join(x, "/")
}

func (m *Model) Breadcrumbs() string {
	ret := util.NewStringSlice(lo.Map(m.Group, func(g string, _ int) string {
		return fmt.Sprintf("%q", strings.ToLower(g))
	}))
	ret.Pushf("%q", m.Package)
	return ret.Join(", ")
}

func (m *Model) IndexedColumns(includePK bool) Columns {
	var ret Columns
	a := func(c *Column) {
		for _, x := range ret {
			if x.Name == c.Name {
				return
			}
		}
		ret = append(ret, c)
	}
	lo.ForEach(m.Columns, func(c *Column, _ int) {
		if c.Indexed || (includePK && c.PK) || m.Relations.ContainsSource(c.Name) {
			a(c)
		}
	})
	return ret
}

func (m *Model) AllSearches(db string) []string {
	if !m.HasTag("search") {
		return m.Search
	}
	ret := util.NewStringSlice(slices.Clone(m.Search))
	lo.ForEach(m.Columns, func(c *Column, _ int) {
		if c.Search {
			x := c.Name
			if !types.IsString(c.Type) {
				switch db {
				case dbSQLServer:
					x = fmt.Sprintf("cast(%s as nvarchar(2048))", c.SQL())
				case dbSQLite:
					x = c.SQL()
				default:
					x = fmt.Sprintf("%s::text", c.SQL())
				}
			}
			ret.Push("lower(" + x + ")")
		}
	})
	return ret.Slice
}

func (m *Model) HasSearches() bool {
	return len(m.AllSearches("")) > 0
}

func (m *Model) ProperWithGroup(extraAcronyms []string) string {
	if len(m.Group) > 0 {
		return util.StringToCamel(m.Group[len(m.Group)-1], extraAcronyms...) + util.StringToCamel(m.Package, extraAcronyms...)
	}
	return m.Proper()
}
