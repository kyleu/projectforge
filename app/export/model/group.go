package model

import (
	"golang.org/x/exp/slices"
	"projectforge.dev/projectforge/app/util"
)

type Group struct {
	Key         string   `json:"key"`
	Title       string   `json:"title,omitempty"`
	Description string   `json:"description,omitempty"`
	Icon        string   `json:"icon,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	Children    Groups   `json:"children,omitempty"`
}

func (g *Group) HasTag(t string) bool {
	return slices.Contains(g.Tags, t)
}

func (g *Group) IconSafe() string {
	if g.Icon == "" {
		return "star"
	}
	return g.Icon
}

func (g *Group) TitleSafe() string {
	if g.Title == "" {
		return util.StringToTitle(g.Key)
	}
	return g.Title
}

type Groups []*Group

func (g Groups) Get(keys ...string) *Group {
	if g == nil || len(keys) == 0 {
		return nil
	}
	for _, x := range g {
		if x.Key == keys[0] {
			if len(keys) == 1 {
				return x
			}
			return x.Children.Get(keys[1:]...)
		}
	}
	return nil
}
