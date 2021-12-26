package export

import (
	"github.com/kyleu/projectforge/app/util"
)

type Model struct {
	Key     string  `json:"key"`
	Pkg     string  `json:"pkg"`
	Columns Columns `json:"columns"`
}

func (m Model) camel() string {
	return util.StringToLowerCamel(m.Key)
}

func (m Model) proper() string {
	return util.StringToCamel(m.Key)
}

func (m Model) properPlural() string {
	return util.StringToPlural(m.proper())
}

func (m Model) firstLetter() string {
	return m.Key[0:1]
}

type Models []*Model
