package model

import (
	"github.com/pkg/errors"
	"projectforge.dev/projectforge/app/util"
)

var goKeywords = []string{
	"break", "case", "chan", "const", "continue", "default", "defer", "else",
	"fallthrough", "for", "func", "go", "goto", "if", "import", "interface",
	"map", "package", "range", "return", "select", "struct", "switch", "type", "var",
}

var reservedNames = map[string][]string{"audit": {"audit", "audit_record"}}

func (m *Model) Validate(mods []string) error {
	if len(m.PKs()) == 0 {
		return errors.Errorf("model [%s] has no primary key", m.Name)
	}
	for _, mod := range mods {
		if util.StringArrayContains(reservedNames[mod], m.Name) {
			return errors.Errorf("model [%s] uses name which is reserved by [%s]", m.Name, mod)
		}
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
