package model

import (
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/util"
)

const defaultIcon = "star"

type Group struct {
	Key         string   `json:"key"`
	Title       string   `json:"title,omitempty"`
	Description string   `json:"description,omitempty"`
	Route       string   `json:"route,omitempty"`
	Icon        string   `json:"icon,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	Children    Groups   `json:"children,omitempty"`
}

func (g *Group) HasTag(t string) bool {
	return lo.Contains(g.Tags, t)
}

func (g *Group) IconSafe() string {
	if g.Icon == "" {
		return defaultIcon
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
