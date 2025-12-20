package model

import (
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/util"
)

type Event struct {
	Name           string        `json:"name"`
	Package        string        `json:"package"`
	Group          []string      `json:"group,omitempty"`
	Schema         string        `json:"schema,omitzero"`
	Description    string        `json:"description,omitzero"`
	Icon           string        `json:"icon,omitzero"`
	Tags           []string      `json:"tags,omitempty"`
	TitleOverride  string        `json:"title,omitzero"`
	PluralOverride string        `json:"plural,omitzero"`
	ProperOverride string        `json:"proper,omitzero"`
	Config         util.ValueMap `json:"config,omitzero"`
	Columns        Columns       `json:"columns"`
	Imports        Imports       `json:"imports,omitempty"`
	acronyms       []string
}

func (e *Event) HasTag(t string) bool {
	return lo.Contains(e.Tags, t)
}

func (e *Event) AddTag(t string) {
	if !e.HasTag(t) {
		e.Tags = util.ArraySorted(append(e.Tags, t))
	}
}

func (e *Event) RemoveTag(t string) {
	e.Tags = lo.Filter(e.Tags, func(x string, _ int) bool {
		return x != t
	})
}

func (e *Event) PackageName() string {
	return e.Package
}

func (e *Event) PackageWithGroup(prefix string) string {
	if x := e.Config.GetStringOpt("pkg-" + prefix); x != "" {
		return x
	}
	if len(e.Group) == 0 {
		return prefix + e.Package
	}
	x := lo.Map(e.Group, func(g string, _ int) string {
		return strings.ToLower(prefix + g)
	})
	x = append(x, prefix+e.Package)
	return util.StringJoin(x, "/")
}

func (e *Event) GroupAndPackage() []string {
	return append(util.ArrayCopy(e.Group), e.Package)
}

func (e *Event) ID() string {
	return util.StringPath(e.PackageWithGroup(""), e.Name)
}

func (e *Event) GroupLen() int {
	return len(e.Group)
}

func (e *Event) GroupString(prefix string, dflt string) string {
	if len(e.Group) == 0 {
		return dflt
	}
	x := lo.Map(e.Group, func(g string, _ int) string {
		return strings.ToLower(prefix + g)
	})
	return util.StringJoin(x, "/")
}

func (e *Event) SetAcronyms(acronyms ...string) {
	e.acronyms = acronyms
	for _, col := range e.Columns {
		col.SetAcronyms(acronyms...)
	}
}
