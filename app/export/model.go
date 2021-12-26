package export

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

func (m Model) camel() string {
	return util.StringToLowerCamel(m.Name)
}

func (m Model) proper() string {
	return util.StringToCamel(m.Name)
}

func (m Model) properPlural() string {
	return util.StringToPlural(m.proper())
}

func (m Model) firstLetter() string {
	return m.Name[0:1]
}

func (m Model) packageProper() string {
	return util.StringToCamel(m.Package)
}

func (m Model) IconSafe() string {
	if m.Icon == "" {
		return "star"
	}
	return m.Icon
}

type Models []*Model
