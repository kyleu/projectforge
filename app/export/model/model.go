package model

import (
	"github.com/kyleu/projectforge/app/util"
)

type Model struct {
	Name        string   `json:"name"`
	Package     string   `json:"package"`
	Description string   `json:"description"`
	Icon        string   `json:"icon"`
	Ordering    string   `json:"ordering"`
	Search      []string `json:"search"`
	Columns     Columns  `json:"columns"`
}

func (m *Model) Camel() string {
	return util.StringToLowerCamel(m.Name)
}

func (m *Model) Proper() string {
	return util.StringToCamel(m.Name)
}

func (m *Model) ProperPlural() string {
	return util.StringToPlural(m.Proper())
}

func (m *Model) FirstLetter() string {
	return m.Name[0:1]
}

func (m *Model) PackageProper() string {
	return util.StringToCamel(m.Package)
}

func (m *Model) IconSafe() string {
	if m.Icon == "" {
		return "star"
	}
	return m.Icon
}

func (m *Model) URLPath(prefix string) string {
	url := "\"/" + m.Package + "\""
	for _, pk := range m.Columns.PKs() {
		url += "+\"/\"+" + pk.ToGoString(prefix)
	}
	return url
}

func (m *Model) ClassRef() string {
	return m.Package + "." + m.Proper()
}

func (m *Model) LinkURL(prefix string) string {
	pks := m.Columns.PKs()
	linkURL := "/" + m.Package
	for _, pk := range pks {
		linkURL += "/" + pk.ToGoViewString(prefix)
	}
	return linkURL
}

type Models []*Model
