package module

import (
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
	Files       filesystem.FileLoader `json:"-"`
}

func (m *Module) Title() string {
	if m.Name == "" {
		return m.Key
	}
	return m.Name
}

type Modules []*Module

func (i Modules) Get(key string) *Module {
	for _, item := range i {
		if item.Key == key {
			return item
		}
	}
	return nil
}

func (i Modules) Keys() []string {
	var ret []string
	for _, x := range i {
		ret = append(ret, x.Key)
	}
	return ret
}
