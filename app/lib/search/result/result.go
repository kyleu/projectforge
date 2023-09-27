// Package result - Content managed by Project Forge, see [projectforge.md] for details.
package result

import (
	"cmp"
	"slices"
	"strings"
)

type Result struct {
	Type    string  `json:"type,omitempty"`
	ID      string  `json:"id,omitempty"`
	Title   string  `json:"title,omitempty"`
	Icon    string  `json:"icon,omitempty"`
	URL     string  `json:"url,omitempty"`
	Matches Matches `json:"matches,omitempty"`
	Data    any     `json:"data,omitempty"`
	HTML    string  `json:"-"`
}

func NewResult(t string, id string, url string, title string, icon string, diff any, data any, q string) *Result {
	return &Result{Type: t, ID: id, URL: url, Title: title, Icon: icon, Matches: MatchesFor("", diff, q), Data: data}
}

type Results []*Result

func (rs Results) Sort() Results {
	slices.SortFunc(rs, func(l *Result, r *Result) int {
		if l.Type == r.Type {
			return cmp.Compare(strings.ToLower(l.Title), strings.ToLower(r.Title))
		}
		return cmp.Compare(l.Type, r.Type)
	})
	return rs
}
