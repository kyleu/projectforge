package model

import (
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/util"
)

func (m *Model) Camel() string {
	return util.StringToLowerCamel(m.Name)
}

func (m *Model) CamelPlural() string {
	return util.StringToPlural(m.Camel())
}

func (m *Model) Title() string {
	if m.TitleOverride == "" {
		return util.StringToTitle(m.Name)
	}
	return m.TitleOverride
}

func (m *Model) Proper() string {
	if m.ProperOverride == "" {
		return util.StringToCamel(m.Name)
	}
	return util.StringToCamel(m.ProperOverride)
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

func (m *Model) Table() string {
	if m.TableOverride == "" {
		return m.Name
	}
	return m.TableOverride
}

func (m *Model) Route() string {
	if m.RouteOverride == "" {
		return m.PackageWithGroup("")
	}
	return m.RouteOverride
}

func (m *Model) IconSafe() string {
	if _, ok := util.SVGLibrary[m.Icon]; ok {
		return m.Icon
	}
	return defaultIcon
}

func (m *Model) URLPath(prefix string) string {
	url := "\"/" + m.Route() + "\""
	lo.ForEach(m.PKs(), func(pk *Column, _ int) {
		url += "+\"/\"+" + pk.ToGoString(prefix)
	})
	return url
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
