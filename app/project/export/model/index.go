package model

type Index struct {
	Name string `json:"name"`
	Decl string `json:"decl"`
}

type Indexes []*Index
