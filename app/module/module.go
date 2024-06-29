package module

import (
	"cmp"
	"slices"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/util"
)

type Module struct {
	Key         string                `json:"-"`
	Name        string                `json:"name,omitempty"`
	Icon        string                `json:"icon,omitempty"`
	Description string                `json:"description,omitempty"`
	Hidden      bool                  `json:"hidden,omitempty"`
	AuthorName  string                `json:"authorName,omitempty"`
	AuthorEmail string                `json:"authorEmail,omitempty"`
	License     string                `json:"license,omitempty"`
	Sourcecode  string                `json:"sourcecode,omitempty"`
	ConfigVars  util.KeyTypeDescs     `json:"configVars,omitempty"`
	PortOffsets map[string]int        `json:"portOffsets,omitempty"`
	Dangerous   bool                  `json:"dangerous,omitempty"`
	Requires    []string              `json:"requires,omitempty"`
	Priority    int                   `json:"priority,omitempty"`
	Technology  []string              `json:"technology,omitempty"`
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
	if _, ok := util.SVGLibrary[m.Icon]; ok {
		return m.Icon
	}
	return "compass"
}

func (m *Module) WebPath() string {
	return "/m/" + m.Key
}

func (m *Module) FeaturesPath() string {
	return "/features/" + m.Key
}

func (m *Module) FeaturesFilePath() string {
	return m.FeaturesPath() + "/files"
}

func (m *Module) DocPath() string {
	return "doc/module/" + m.Key + util.ExtMarkdown
}

type Modules []*Module

func (m Modules) Get(key string) *Module {
	return lo.FindOrElse(m, nil, func(item *Module) bool {
		return item.Key == key
	})
}

func (m Modules) Keys() []string {
	return lo.Map(m, func(x *Module, _ int) string {
		return x.Key
	})
}

func (m Modules) Sort() Modules {
	slices.SortFunc(m, func(l *Module, r *Module) int {
		return cmp.Compare(l.Priority, r.Priority)
	})
	return m
}
