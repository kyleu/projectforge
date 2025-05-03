package model

import (
	"fmt"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/filter"
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/util"
)

type Model struct {
	Name           string           `json:"name"`
	Package        string           `json:"package"`
	Group          []string         `json:"group,omitempty"`
	Schema         string           `json:"schema,omitempty"`
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
	Imports        Imports          `json:"imports,omitempty"`
	acronyms       []string
}

func (m *Model) HasTag(t string) bool {
	return lo.Contains(m.Tags, t)
}

func (m *Model) AddTag(t string) {
	if !m.HasTag(t) {
		m.Tags = util.ArraySorted(append(m.Tags, t))
	}
}

func (m *Model) RemoveTag(t string) {
	m.Tags = lo.Filter(m.Tags, func(x string, _ int) bool {
		return x != t
	})
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
	if x := m.Config.GetStringOpt("pkg-" + prefix); x != "" {
		return x
	}
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

func (m *Model) HasSearches() bool {
	return len(m.AllSearches("")) > 0
}

func (m *Model) AllSearches(db string) []string {
	if !m.HasTag("search") {
		return m.Search
	}
	ret := util.NewStringSlice(util.ArrayCopy(m.Search))
	lo.ForEach(m.Columns, func(c *Column, _ int) {
		if c.Search {
			x := fmt.Sprintf("%q", c.SQL())
			if !types.IsString(c.Type) {
				switch db {
				case util.DatabaseSQLServer:
					x = fmt.Sprintf("cast(%q as nvarchar(2048))", c.SQL())
				case util.DatabaseSQLite:
					x = c.SQL()
				default:
					x = fmt.Sprintf("%q::text", c.SQL())
				}
			}
			ret.Push(fmt.Sprintf("lower(%s)", x))
		}
	})
	return ret.Slice
}

func (m *Model) SetAcronyms(acronyms ...string) {
	m.acronyms = acronyms
	for _, col := range m.Columns {
		col.SetAcronyms(acronyms...)
	}
}

func (m *Model) Cleanup() {
	m.Relations = lo.UniqBy(m.Relations, func(x *Relation) string {
		return x.Uniq()
	})
	for _, rel := range m.Relations {
		rel.Src = lo.Uniq(rel.Src)
		rel.Tgt = lo.Uniq(rel.Tgt)
	}
}
