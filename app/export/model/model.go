package model

import (
	"golang.org/x/exp/slices"

	"projectforge.dev/projectforge/app/lib/filter"
	"projectforge.dev/projectforge/app/util"
)

type Model struct {
	Name           string           `json:"name"`
	Package        string           `json:"package"`
	Description    string           `json:"description,omitempty"`
	Icon           string           `json:"icon,omitempty"`
	Ordering       filter.Orderings `json:"ordering,omitempty"`
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
	historyMap     *HistoryMap
	historyMapDB   *HistoryMap
}

func (m *Model) HasTag(t string) bool {
	return slices.Contains(m.Tags, t)
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
	for _, pk := range pks {
		linkURL += "/" + pk.ToGoViewString(prefix)
	}
	return linkURL
}

func (m *Model) RelationsFor(col *Column) Relations {
	var ret Relations
	for _, r := range m.Relations {
		if slices.Contains(r.Src, col.Name) {
			ret = append(ret, r)
		}
	}
	return ret
}

func (m *Model) CanTraverseRelation() bool {
	return len(m.PKs()) == 1 && len(m.Columns.WithTag("title")) > 0
}
