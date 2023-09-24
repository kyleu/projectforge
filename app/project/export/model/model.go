package model

import (
	"fmt"
	"slices"
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
	Description    string           `json:"description,omitempty"`
	Icon           string           `json:"icon,omitempty"`
	Ordering       filter.Orderings `json:"ordering,omitempty"`
	SortIndex      int              `json:"sortIndex,omitempty"`
	Search         []string         `json:"search,omitempty"`
	History        string           `json:"history,omitempty"`
	Tags           []string         `json:"tags,omitempty"`
	TitleOverride  string           `json:"title,omitempty"`
	ProperOverride string           `json:"proper,omitempty"`
	RouteOverride  string           `json:"route,omitempty"`
	Config         util.ValueMap    `json:"config,omitempty"`
	Columns        Columns          `json:"columns"`
	Relations      Relations        `json:"relations,omitempty"`
	Indexes        Indexes          `json:"indexes,omitempty"`
	SeedData       [][]any          `json:"seedData,omitempty"`
	historyMap     *HistoryMap
	historyMapDB   *HistoryMap
}

func (m *Model) HasTag(t string) bool {
	return lo.Contains(m.Tags, t)
}

func (m *Model) AddTag(t string) {
	if !lo.Contains(m.Tags, t) {
		m.Tags = append(m.Tags, t)
		slices.Sort(m.Tags)
	}
}

func (m *Model) PKs() Columns {
	return m.Columns.PKs()
}

func (m *Model) GroupedColumns() Columns {
	return m.Columns.WithTag("grouped")
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

func (m *Model) IsRevision() bool {
	return m.History == RevisionType
}

func (m *Model) IsHistory() bool {
	return m.History == HistoryType
}

func (m *Model) LinkURL(prefix string) string {
	pks := m.PKs()
	linkURL := "/" + m.Route()
	lo.ForEach(pks, func(pk *Column, _ int) {
		linkURL += "/" + pk.ToGoViewString(prefix, false, true)
	})
	return linkURL
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
	ret := lo.Map(m.Group, func(g string, _ int) string {
		return fmt.Sprintf("%q", strings.ToLower(g))
	})
	ret = append(ret, fmt.Sprintf("%q", m.Package))
	return strings.Join(ret, ", ")
}

func (m *Model) IndexedColumns() Columns {
	ret := m.GroupedColumns()
	a := func(c *Column) {
		for _, x := range ret {
			if x.Name == c.Name {
				return
			}
		}
		ret = append(ret, c)
	}
	lo.ForEach(m.Columns, func(c *Column, _ int) {
		if c.Indexed || c.PK {
			a(c)
		}
	})
	return ret
}

func (m *Model) AllSearches(database string) []string {
	if !m.HasTag("search") {
		return m.Search
	}
	ret := slices.Clone(m.Search)
	lo.ForEach(m.Columns, func(c *Column, _ int) {
		if c.Search {
			x := c.Name
			if !types.IsString(c.Type) {
				switch database {
				case util.DatabaseSQLServer:
					x = fmt.Sprintf("cast(%s as nvarchar(2048))", c.Name)
				case util.DatabaseSQLite:
					x = c.Name
				default:
					x = fmt.Sprintf("%s::text", c.Name)
				}
			}
			ret = append(ret, "lower("+x+")")
		}
	})
	return ret
}

func (m *Model) HasSearches() bool {
	return len(m.AllSearches("")) > 0
}
