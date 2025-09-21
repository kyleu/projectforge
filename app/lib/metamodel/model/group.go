package model

import (
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/util"
)

const defaultIcon = "star"

type Group struct {
	Key         string   `json:"key"`
	Title       string   `json:"title,omitzero"`
	Description string   `json:"description,omitzero"`
	Route       string   `json:"route,omitzero"`
	Icon        string   `json:"icon,omitzero"`
	Tags        []string `json:"tags,omitempty"`
	Children    Groups   `json:"children,omitempty"`
	Provided    bool     `json:"provided,omitzero"`
}

func (g *Group) HasTag(t string) bool {
	return lo.Contains(g.Tags, t)
}

func (g *Group) IconSafe() string {
	return util.OrDefault(g.Icon, defaultIcon)
}

func (g *Group) TitleSafe() string {
	if g.Title == "" {
		return util.StringToTitle(g.Key)
	}
	return g.Title
}

func (g *Group) Proper() string {
	return util.StringToProper(g.Key)
}

func (g *Group) String() string {
	return util.OrDefault(g.Title, g.Key)
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
	ret := util.ArrayCopy(g)
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
