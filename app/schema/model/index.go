package model

import (
	"strings"

	"github.com/kyleu/projectforge/app/schema/field"
)

type Index struct {
	Key      string          `json:"key"`
	Fields   []string        `json:"type"`
	Unique   bool            `json:"unique,omitempty"`
	Primary  bool            `json:"primary,omitempty"`
	Metadata *field.Metadata `json:"metadata,omitempty"`
}

func (i Index) String() string {
	return i.Key + "(" + strings.Join(i.Fields, ", ") + ")"
}

type Indexes []*Index

func (s Indexes) Get(key string) *Index {
	for _, x := range s {
		if x.Key == key {
			return x
		}
	}
	return nil
}
