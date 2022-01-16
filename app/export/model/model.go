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
	Config         util.ValueMap    `json:"config,omitempty"`
	Columns        Columns          `json:"columns"`
	historyMap     *HistoryMap
	historyMapDB   *HistoryMap
}

func (m *Model) Camel() string {
	return util.StringToLowerCamel(m.Name)
}

func (m *Model) Proper() string {
	if m.ProperOverride == "" {
		return util.StringToCamel(m.Name)
	}
	return util.StringToCamel(m.ProperOverride)
}

func (m *Model) Title() string {
	if m.TitleOverride == "" {
		return m.Proper()
	}
	return m.TitleOverride
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
	return util.StringToPlural(m.Name)
}

func (m *Model) ProperPlural() string {
	return util.StringToPlural(m.Proper())
}

func (m *Model) FirstLetter() string {
	return m.Name[0:1]
}

func (m *Model) IconSafe() string {
	_, ok := util.SVGLibrary[m.Icon]
	if !ok {
		return "star"
	}
	return m.Icon
}

func (m *Model) URLPath(prefix string) string {
	url := "\"/" + m.Package + "\""
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
	linkURL := "/" + m.Package
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

func (m *Model) PKs() Columns {
	return m.Columns.PKs()
}

func (m *Model) GroupedColumns() Columns {
	return m.Columns.WithTag("grouped")
}

func (m *Model) IsRevision() bool {
	return m.History == RevisionType
}

func (m *Model) Validate() error {
	if m.IsRevision() {
		hc := m.HistoryColumns(true)
		if hc.Err != nil {
			return hc.Err
		}
		hc = m.HistoryColumns(false)
		if hc.Err != nil {
			return hc.Err
		}
	}
	return nil
}

type Models []*Model
