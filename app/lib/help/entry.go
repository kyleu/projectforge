// Package help - Content managed by Project Forge, see [projectforge.md] for details.
package help

import (
	"strings"

	"projectforge.dev/projectforge/app/util"
)

type Entry struct {
	Key      string `json:"key"`
	Title    string `json:"title,omitempty"`
	Markdown string `json:"markdown,omitempty"`
	HTML     string `json:"html,omitempty"`
}

type Entries []*Entry

func (e Entries) Get(key string) *Entry {
	key = strings.TrimSuffix(key, util.ExtMarkdown)
	for _, x := range e {
		if x.Key == key {
			return x
		}
	}
	return nil
}
