package model

import (
	"slices"

	"github.com/samber/lo"

	"{{{ .Package }}}/app/util"
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
	Provided    bool     `json:"provided,omitempty"`
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

func (g *Group) Proper() string {
	return util.StringToCamel(g.Key)
}

func (g *Group) String() string {
	if g.Title == "" {
		return g.Key
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

func (g Groups) Flatten() Groups {
	ret := slices.Clone(g)
	for _, x := range g {
		ret = append(ret, x.Children.Flatten()...)
	}
	return ret
}

func (g Groups) Strings(prefix string) []string {
	return lo.Map(g, func(x *Group, _ int) string {
		return prefix + x.Proper()
	})
}
