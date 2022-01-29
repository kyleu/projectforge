package model

import (
	"github.com/kyleu/projectforge/app/util"
	"github.com/pkg/errors"
)

var goKeywords = []string{
	"break", "case", "chan", "const", "continue", "default", "defer", "else",
	"fallthrough", "for", "func", "go", "goto", "if", "import", "interface",
	"map", "package", "range", "return", "select", "struct", "switch", "type", "var",
}

var reservedNames = []string{"audit", "audit_record"}

func (m *Model) Validate() error {
	if len(m.PKs()) == 0 {
		return errors.Errorf("model [%s] has no primary key", m.Name)
	}
	if util.StringArrayContains(reservedNames, m.Name) {
		return errors.Errorf("model [%s] uses a reserved name", m.Name)
	}
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
	for _, col := range m.Columns {
		if util.StringArrayContains(goKeywords, col.Name) {
			return errors.Errorf("column [%s] uses reserved keyword", col.Name)
		}
	}
	return nil
}
