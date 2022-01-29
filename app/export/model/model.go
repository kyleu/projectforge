package model

import (
	"strings"

	"github.com/kyleu/projectforge/app/lib/filter"
	"github.com/kyleu/projectforge/app/util"
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

func (m *Model) Camel() string {
	return util.StringToLowerCamel(m.Name)
}

func (m *Model) CamelPlural() interface{} {
	return util.StringToPlural(m.Camel())
}

func (m *Model) Title() string {
	if m.TitleOverride == "" {
		return m.Proper()
	}
	return m.TitleOverride
}

func (m *Model) Proper() string {
	if m.ProperOverride == "" {
		return util.StringToCamel(m.Name)
	}
	return util.StringToCamel(m.ProperOverride)
}

func (m *Model) Route() string {
	if m.RouteOverride == "" {
		return m.Package
	}
	return m.RouteOverride
}

func (m *Model) TitleLower() string {
	return strings.ToLower(m.Title())
}

func (m *Model) TitlePlural() string {
	return util.StringToPlural(m.Title())
}

func (m *Model) TitlePluralLower() string {
	return util.StringToPlural(m.TitleLower())
}

func (m *Model) Plural() string {
	ret := util.StringToPlural(m.Name)
	if ret == m.Name {
		return ret + "Set"
	}
	return ret
}

func (m *Model) ProperPlural() string {
	ret := util.StringToPlural(m.Proper())
	if ret == m.Proper() {
		return ret + "Set"
	}
	return ret
}

func (m *Model) FirstLetter() string {
	return strings.ToLower(m.Name[0:1])
}

func (m *Model) IconSafe() string {
	_, ok := util.SVGLibrary[m.Icon]
	if !ok {
		return "star"
	}
	return m.Icon
}

func (m *Model) URLPath(prefix string) string {
	url := "\"/" + m.Route() + "\""
	for _, pk := range m.PKs() {
		url += "+\"/\"+" + pk.ToGoString(prefix)
	}
	return url
}

func (m *Model) ClassRef() string {
	return m.Package + "." + m.Proper()
}

func (m *Model) LinkURL(prefix string) string {
	pks := m.PKs()
	linkURL := "/" + m.Route()
	for _, pk := range pks {
		linkURL += "/" + pk.ToGoViewString(prefix)
	}
	return linkURL
}

func (m *Model) HasTag(t string) bool {
	return util.StringArrayContains(m.Tags, t)
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

func (m *Model) PKs() Columns {
	return m.Columns.PKs()
}

func (m *Model) GroupedColumns() Columns {
	return m.Columns.WithTag("grouped")
}

func (m *Model) IsRevision() bool {
	return m.History == RevisionType
}

func (m *Model) IsHistory() bool {
	return m.History == HistoryType
}

func (m *Model) RelationsFor(col *Column) Relations {
	var ret Relations
	for _, r := range m.Relations {
		if util.StringArrayContains(r.Src, col.Name) {
			ret = append(ret, r)
		}
	}
	return ret
}
