package module

import (
	"sort"

	"github.com/kyleu/projectforge/app/filesystem"
)

type Module struct {
	Key         string                `json:"-"`
	Name        string                `json:"name,omitempty"`
	Description string                `json:"description,omitempty"`
	AuthorName  string                `json:"authorName,omitempty"`
	AuthorEmail string                `json:"authorEmail,omitempty"`
	License     string                `json:"license,omitempty"`
	Sourcecode  string                `json:"sourcecode,omitempty"`
	Priority    int                   `json:"priority,omitempty"`
	Files       filesystem.FileLoader `json:"-"`
}

func (m *Module) Title() string {
	if m.Name == "" {
		return m.Key
	}
	return m.Name
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
	sort.Slice(m, func(i, j int) bool {
		return m[i].Priority < m[j].Priority
	})
	return m
}
