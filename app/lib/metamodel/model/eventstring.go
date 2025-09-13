package model

import (
	"fmt"
	"strings"

	"projectforge.dev/projectforge/app/util"
)

func (e *Event) Camel() string {
	return util.StringToCamel(e.Name, e.acronyms...)
}

func (e *Event) CamelLower() string {
	return strings.ToLower(e.Camel())
}

func (e *Event) CamelPlural() string {
	return util.StringToPlural(e.Camel())
}

func (e *Event) Title() string {
	if e.TitleOverride == "" {
		return util.StringToTitle(e.Name, e.acronyms...)
	}
	return e.TitleOverride
}

func (e *Event) Proper() string {
	if e.ProperOverride == "" {
		return util.StringToProper(e.Name, e.acronyms...)
	}
	return util.StringToProper(e.ProperOverride, e.acronyms...)
}

func (e *Event) ProperWithGroup() string {
	if len(e.Group) > 0 {
		return util.StringToProper(e.Group[len(e.Group)-1], e.acronyms...) + util.StringToProper(e.Package, e.acronyms...)
	}
	return e.Proper()
}

func (e *Event) TitleLower() string {
	return strings.ToLower(e.Title())
}

func (e *Event) TitlePlural() string {
	return util.StringToPlural(e.Title())
}

func (e *Event) TitlePluralLower() string {
	return util.StringToPlural(e.TitleLower())
}

func (e *Event) Plural() string {
	if e.PluralOverride != "" {
		return e.PluralOverride
	}
	ret := util.StringToPlural(e.Name)
	if ret == e.Name {
		return ret + tSet
	}
	return ret
}

func (e *Event) PluralLower() string {
	return strings.ToLower(e.Plural())
}

func (e *Event) ProperPlural() string {
	ret := util.StringToPlural(e.Proper())
	if ret == e.Proper() {
		return ret + tSet
	}
	return ret
}

func (e *Event) FirstLetter() string {
	return strings.ToLower(e.Name[0:1])
}

func (e *Event) Pointer() string {
	return fmt.Sprintf("*%s.%s", e.Package, e.Proper())
}

func (e *Event) IconSafe() string {
	if _, ok := util.SVGLibrary[e.Icon]; ok {
		return e.Icon
	}
	return defaultIcon
}

func (e *Event) ClassRef() string {
	return e.Package + "." + e.Proper()
}

func (e *Event) LastGroup(prefix string, dflt string) string {
	if len(e.Group) == 0 {
		if dflt != "" {
			return dflt
		}
		return strings.ToLower(e.Package)
	}
	return strings.ToLower(prefix + e.Group[len(e.Group)-1])
}
