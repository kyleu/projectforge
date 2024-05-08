package model

import (
	"fmt"

	"github.com/samber/lo"

	"{{{ .Package }}}/app/util"
)

var IndexFieldDescs = util.FieldDescs{
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

func (i Indexes) Get(name string) *Index {
	return lo.FindOrElse(i, nil, func(x *Index) bool {
		return x.Name == name
	})
}
