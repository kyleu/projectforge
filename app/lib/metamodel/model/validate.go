package model

import (
	"github.com/pkg/errors"
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/util"
)

var goKeywords = []string{
	"break", "case", "chan", "const", "continue", "default", "defer", "else",
	"fallthrough", "for", "func", "go", "goto", "if", "import", "interface",
	"map", "package", "range", "return", "select", "struct", "switch", "type", "var",
}

var reservedNames = map[string][]string{"audit": {"audit", "audit_record"}}

func (m *Model) Validate(mods []string, models Models, groups Groups) error {
	if err := validateBasic(m); err != nil {
		return err
	}
	for _, mod := range mods {
		if lo.Contains(reservedNames[mod], m.Name) {
			return errors.Errorf("model [%s] uses name which is reserved by [%s]", m.Name, mod)
		}
	}
	if (!m.HasTag("menu-hidden")) && len(m.Group) > 0 && groups.Get(m.Group...) == nil {
		msg := "model [%s] references undefined group [%s], and no model matches"
		if len(m.Group) == 1 && models.Get(m.Group[0]) == nil {
			return errors.Errorf(msg, m.Name, util.StringJoin(m.Group, "/"))
		}
		if len(m.Group) > 1 {
			var cool bool
			mg := util.StringJoin(m.Group, "/")
			for _, x := range models {
				xg := util.StringJoin(x.Group, "/") + "/" + x.Package
				if mg == xg {
					cool = true
				}
			}
			if !cool {
				return errors.Errorf(msg, m.Name, util.StringJoin(m.Group, "/"))
			}
		}
	}
	return nil
}

func validateBasic(m *Model) error {
	if len(m.PKs()) == 0 {
		return errors.Errorf("model [%s] has no primary key", m.Name)
	}
	if m.Package == "vendor" {
		return errors.Errorf("model [%s] uses [vendor] package, which is reserved by Go", m.Name)
	}
	if m.IsSoftDelete() {
		if d := m.Columns.WithTag("deleted"); len(d) != 1 {
			return errors.Errorf("when set to soft delete, model [%s] must have one column tagged [deleted]", m.Name)
		}
	}
	if dupes := lo.FindDuplicates(m.Columns.Names()); len(dupes) > 0 {
		return errors.Errorf("model [%s] has duplicates columns [%s]", m.Name, util.StringJoin(dupes, ", "))
	}
	for _, col := range m.Columns {
		if lo.Contains(goKeywords, col.Name) {
			return errors.Errorf("model [%s] column [%s] uses reserved keyword", m.Name, col.Name)
		}
	}
	for _, rel := range m.Relations {
		for _, s := range rel.Src {
			if m.Columns.Get(s) == nil {
				return errors.Errorf("model [%s] relation [%s] references missing source column [%s]", m.Name, rel.Name, s)
			}
		}
	}
	return nil
}
