package module

import (
	"fmt"
)

type Module struct {
	Key         string `json:"-"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	AuthorName  string `json:"authorName,omitempty"`
	AuthorEmail string `json:"authorEmail,omitempty"`
	License     string `json:"license,omitempty"`
	Sourcecode  string `json:"sourcecode,omitempty"`
}

func (m *Module) Title() string {
	if m.Name == "" {
		return m.Key
	}
	return m.Name
}

func (m *Module) Path() string {
	return fmt.Sprintf("module/%s", m.Key)
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
