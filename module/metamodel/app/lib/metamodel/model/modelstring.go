package model

import (
	"fmt"
	"strings"

	"{{{ .Package }}}/app/util"
)

const tSet = "Set"

func (m *Model) Camel() string {
	return util.StringToLowerCamel(m.Name, m.acronyms...)
}

func (m *Model) CamelLower() string {
	return strings.ToLower(m.Camel())
}

func (m *Model) CamelPlural() string {
	return util.StringToPlural(m.Camel())
}

func (m *Model) Title() string {
	if m.TitleOverride == "" {
		return util.StringToTitle(m.Name, m.acronyms...)
	}
	return m.TitleOverride
}

func (m *Model) Proper() string {
	if m.ProperOverride == "" {
		return util.StringToCamel(m.Name, m.acronyms...)
	}
	return util.StringToCamel(m.ProperOverride, m.acronyms...)
}

func (m *Model) ProperWithGroup() string {
	if len(m.Group) > 0 {
		return util.StringToCamel(m.Group[len(m.Group)-1], m.acronyms...) + util.StringToCamel(m.Package, m.acronyms...)
	}
	return m.Proper()
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
	if m.PluralOverride != "" {
		return m.PluralOverride
	}
	ret := util.StringToPlural(m.Name)
	if ret == m.Name {
		return ret + tSet
	}
	return ret
}

func (m *Model) PluralLower() string {
	return strings.ToLower(m.Plural())
}

func (m *Model) ProperPlural() string {
	ret := util.StringToPlural(m.Proper())
	if ret == m.Proper() {
		return ret + tSet
	}
	return ret
}

func (m *Model) FirstLetter() string {
	return strings.ToLower(m.Name[0:1])
}

func (m *Model) Table() string {
	return util.OrDefault(m.TableOverride, m.Name)
}

func (m *Model) TableNS() string {
	ret := m.Table()
	if m.Schema != "" {
		ret = fmt.Sprintf("%q.%q", m.Schema, m.Table())
	}
	return ret
}

func (m *Model) Pointer() string {
	return fmt.Sprintf("*%s.%s", m.Package, m.Proper())
}

func (m *Model) Route() string {
	if m.RouteOverride == "" {
		return m.PackageWithGroup("")
	}
	return m.RouteOverride
}

func (m *Model) CSRoute() string {
	if m.RouteOverride == "" {
		return "/" + m.CamelLower()
	}
	return m.RouteOverride
}

func (m *Model) IconSafe() string {
	if _, ok := util.SVGLibrary[m.Icon]; ok {
		return m.Icon
	}
	return defaultIcon
}

func (m *Model) ClassRef() string {
	return m.Package + "." + m.Proper()
}

func (m *Model) LastGroup(prefix string, dflt string) string {
	if len(m.Group) == 0 {
		if dflt != "" {
			return dflt
		}
		return strings.ToLower(m.Package)
	}
	return strings.ToLower(prefix + m.Group[len(m.Group)-1])
}
