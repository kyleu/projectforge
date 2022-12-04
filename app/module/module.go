package module

import (
	"golang.org/x/exp/slices"

	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/util"
)

type Module struct {
	Key         string                `json:"-"`
	Name        string                `json:"name,omitempty"`
	Icon        string                `json:"icon,omitempty"`
	Description string                `json:"description,omitempty"`
	AuthorName  string                `json:"authorName,omitempty"`
	AuthorEmail string                `json:"authorEmail,omitempty"`
	License     string                `json:"license,omitempty"`
	Sourcecode  string                `json:"sourcecode,omitempty"`
	ConfigVars  util.KeyTypeDescs     `json:"configVars,omitempty"`
	PortOffsets map[string]int        `json:"portOffsets,omitempty"`
	Requires    []string              `json:"requires,omitempty"`
	Priority    int                   `json:"priority,omitempty"`
	Files       filesystem.FileLoader `json:"-"`
	URL         string                `json:"-"`
	UsageMD     string                `json:"-"`
}

func (m *Module) Title() string {
	if m.Name == "" {
		return m.Key
	}
	return m.Name
}

func (m *Module) IconSafe() string {
	if _, ok := util.SVGLibrary[m.Icon]; !ok {
		return "compass"
	}
	return m.Icon
}

func (m *Module) WebPath() string {
	return "/m/" + m.Key
}

func (m *Module) DocPath() string {
	return "doc/module/" + m.Key + ".md"
}

type Modules []*Module

func (m Modules) Get(key string) *Module {
	for _, item := range m {
		if item.Key == key {
			return item
		}
	}
	return nil
}

func (m Modules) Keys() []string {
	ret := make([]string, 0, len(m))
	for _, x := range m {
		ret = append(ret, x.Key)
	}
	return ret
}

func (m Modules) Sort() Modules {
	slices.SortFunc(m, func(l *Module, r *Module) bool {
		return l.Priority < r.Priority
	})
	return m
}
