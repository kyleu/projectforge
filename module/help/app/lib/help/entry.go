package help

import "strings"

type Entry struct {
	Key      string `json:"key"`
	Title    string `json:"title,omitempty"`
	Markdown string `json:"markdown,omitempty"`
	HTML     string `json:"html,omitempty"`
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
