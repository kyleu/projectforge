package model

type Group struct {
	Key         string `json:"key"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Children    Groups `json:"children,omitempty"`
}

type Groups []*Group
