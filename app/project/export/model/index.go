package model

import (
	"fmt"

	"projectforge.dev/projectforge/app/util"
)

var IndexModelDescs = util.ModelDescs{
	{Key: "name", Title: "Name", Description: "The name of the index"},
	{Key: "decl", Title: "Declaration", Description: "The declaration of the index"},
}

type Index struct {
	Name   string `json:"name"`
	Decl   string `json:"decl"`
	Unique bool   `json:"unique,omitempty"`
}

func (i *Index) SQL() string {
	t := ""
	if i.Unique {
		t = " unique"
	}
	return fmt.Sprintf("create%s index if not exists %q on %s;", t, i.Name, i.Decl)
}

type Indexes []*Index
