package model

import (
	"projectforge.dev/projectforge/app/util"
)

var LinkFieldDescs = util.FieldDescs{
	{Key: "title", Title: "Title", Description: "The title of the link"},
	{Key: "icon", Title: "Icon", Description: "The icon of this link"},
	{Key: "url", Title: "URL", Description: "The href of this link"},
}

type Link struct {
	Title string `json:"title,omitempty"`
	URL   string `json:"url"`
	Icon  string `json:"icon,omitempty"`
}

type Links []*Link
