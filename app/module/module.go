package module

import (
	"sort"

	"github.com/gomarkdown/markdown"
	"github.com/kyleu/projectforge/app/filesystem"
	"github.com/kyleu/projectforge/app/util"
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
	PortOffsets map[string]int        `json:"portOffsets,omitempty"`
	Priority    int                   `json:"priority,omitempty"`
	Files       filesystem.FileLoader `json:"-"`
	URL         string                `json:"-"`
	UsageMD     string                `json:"-"`
	usageHTML   string                `json:"-"`
}

func (m *Module) Title() string {
	if m.Name == "" {
		return m.Key
	}
	return m.Name
}

func (m *Module) IconSafe() string {
	_, ok := util.SVGLibrary[m.Icon]
	if !ok {
		return "compass"
	}
	return m.Icon
}

func (m *Module) UsageHTML() string {
	if m.usageHTML == "" {
		m.usageHTML = string(markdown.ToHTML([]byte(m.UsageMD), nil, nil))
	}
	return m.usageHTML
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
	sort.Slice(m, func(i int, j int) bool {
		return m[i].Priority < m[j].Priority
	})
	return m
}
