package help

import (
	"strings"

	"projectforge.dev/projectforge/app/util"
)

type Entry struct {
	Key      string `json:"key"`
	Title    string `json:"title,omitzero"`
	Markdown string `json:"markdown,omitzero"`
	HTML     string `json:"html,omitzero"`
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
