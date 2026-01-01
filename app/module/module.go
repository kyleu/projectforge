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
	Name        string                `json:"name,omitzero"`
	Icon        string                `json:"icon,omitzero"`
	Description string                `json:"description,omitzero"`
	Hidden      bool                  `json:"hidden,omitzero"`
	AuthorName  string                `json:"authorName,omitzero"`
	AuthorEmail string                `json:"authorEmail,omitzero"`
	License     string                `json:"license,omitzero"`
	Sourcecode  string                `json:"sourcecode,omitzero"`
	ConfigVars  util.KeyTypeDescs     `json:"configVars,omitzero"`
	PortOffsets map[string]int        `json:"portOffsets,omitzero"`
	Dangerous   bool                  `json:"dangerous,omitzero"`
	Requires    []string              `json:"requires,omitempty"`
	Priority    int                   `json:"priority,omitzero"`
	Technology  []string              `json:"technology,omitempty"`
	Files       filesystem.FileLoader `json:"-"`
	URL         string                `json:"-"`
	UsageMD     string                `json:"-"`
}

func (m *Module) Title() string {
	return util.OrDefault(m.Name, m.Key)
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

func (m Modules) Sorted() Modules {
	ret := util.ArrayCopy(m)
	slices.SortFunc(ret, func(l *Module, r *Module) int {
		if l.Key == "core" {
			return -1
		}
		if r.Key == "core" {
			return 1
		}
		return cmp.Compare(l.Key, r.Key)
	})
	return ret
}

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
