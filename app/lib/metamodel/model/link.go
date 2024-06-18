package model

import (
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/util"
)

var LinkFieldDescs = util.FieldDescs{
	{Key: "title", Title: "Title", Description: "The title of the link"},
	{Key: "url", Title: "URL", Description: "The href of this link"},
	{Key: "icon", Title: "Icon", Description: "The icon of this link"},
	{Key: "dangerous", Title: "Dangerous", Description: "If set, this link will require confirmation", Type: "bool"},
	{Key: "tags", Title: "Tags", Description: "The tags that apply to this link", Type: "[]string"},
}

type Link struct {
	Title     string   `json:"title,omitempty"`
	URL       string   `json:"url"`
	Icon      string   `json:"icon,omitempty"`
	Dangerous bool     `json:"dangerous,omitempty"`
	Tags      []string `json:"tags,omitempty"`
}

type Links []*Link

func (l Links) WithTags(includeEmpty bool, tags ...string) Links {
	return lo.Filter(l, func(x *Link, _ int) bool {
		return (includeEmpty && len(x.Tags) == 0) || (len(lo.Intersect(tags, x.Tags)) > 0)
	})
}
