package result

import (
	"strings"

	"golang.org/x/exp/slices"
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

func NewResult(t string, id string, url string, title string, icon string, obj any, q string) *Result {
	return &Result{Type: t, ID: id, URL: url, Title: title, Icon: icon, Matches: MatchesFor("", obj, q), Data: obj}
}

type Results []*Result

func (rs Results) Sort() {
	slices.SortFunc(rs, func(l *Result, r *Result) bool {
		if l.Type == r.Type {
			return strings.ToLower(l.Title) < strings.ToLower(r.Title)
		}
		return l.Type < r.Type
	})
}
