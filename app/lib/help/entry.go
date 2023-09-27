// Package help - Content managed by Project Forge, see [projectforge.md] for details.
package help

import "strings"

type Entry struct {
	Key      string `json:"key"`
	Markdown string `json:"markdown"`
	HTML     string `json:"html"`
}

type Entries []*Entry

func (e Entries) Get(key string) *Entry {
	key = strings.TrimSuffix(key, ".md")
	for _, x := range e {
		if x.Key == key {
			return x
		}
	}
	return nil
}
