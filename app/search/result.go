package search

import (
	"sort"

	"github.com/kyleu/projectforge/app/util"
)

type Match struct {
	Key   string `json:"k"`
	Value string `json:"v"`
}

type Matches []*Match

func MatchesFrom(a []string) Matches {
	ret := make(Matches, 0, len(a))
	for _, x := range a {
		k, v := util.SplitString(x, ':', true)
		ret = append(ret, &Match{Key: k, Value: v})
	}
	return ret
}

type Result struct {
	ID      string      `json:"id,omitempty"`
	Type    string      `json:"type,omitempty"`
	Title   string      `json:"title,omitempty"`
	Icon    string      `json:"icon,omitempty"`
	URL     string      `json:"url,omitempty"`
	Matches []*Match    `json:"matches,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	HTML    string      `json:"-"`
}

type Results []*Result

func (rs Results) Sort() {
	sort.Slice(rs, func(i, j int) bool {
		l, r := rs[i], rs[j]
		if l.Type == r.Type {
			return l.Title < r.Title
		}
		return l.Type < r.Type
	})
}
