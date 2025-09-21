package menu

import (
	"strings"

	"projectforge.dev/projectforge/app/util"
)

var Separator = &Item{}

type Item struct {
	Key         string `json:"key"`
	Title       string `json:"title,omitzero"`
	Description string `json:"description,omitzero"`
	Badge       string `json:"badge,omitzero"`
	Icon        string `json:"icon,omitzero"`
	Route       string `json:"route,omitzero"`
	Hidden      bool   `json:"hidden"`
	Warning     string `json:"warning,omitzero"`
	Children    Items  `json:"children,omitempty"`
}

func ItemFromString(bc string, dflt string) *Item {
	icon := util.OrDefault(dflt, "file")
	if iconIdx := strings.Index(bc, "**"); iconIdx > -1 {
		icon = bc[iconIdx+2:]
		bc = bc[:iconIdx]
	}
	var bcLink string
	if bci := strings.Index(bc, "||"); bci > -1 {
		bcLink = bc[bci+2:]
		bc = bc[:bci]
	}
	if bc == "" {
		bc = dflt
	}
	return &Item{Key: bc, Title: bc, Icon: icon, Route: bcLink}
}

func (i *Item) AddChild(child *Item) {
	i.Children = append(i.Children, child)
}

func (i *Item) Desc() string {
	if i.Description != "" {
		return i.Description
	}
	return i.Title
}

func (i *Item) Clone() *Item {
	return &Item{
		Key: i.Key, Title: i.Title, Description: i.Description, Badge: i.Badge, Icon: i.Icon,
		Route: i.Route, Hidden: i.Hidden, Warning: i.Warning, Children: i.Children,
	}
}
