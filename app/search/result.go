package search

import (
	"sort"
	"strings"

	"github.com/kyleu/projectforge/app/util"
)

type Match struct {
	Key   string `json:"k"`
	Value string `json:"v"`
}

func (m *Match) ValueSplit(q string) []string {
	ql := strings.ToLower(q)
	vl := strings.ToLower(m.Value)
	idx := 0
	var ret []string
	for idx > -1 {
		nIdx := strings.Index(vl[idx:], ql)
		if nIdx == -1 {
			ret = append(ret, m.Value[idx:])
			break
		} else {
			if idx == 0 {
				ret = append(ret, "")
			}
			nIdx += idx
			ret = append(ret, m.Value[idx:nIdx])
		}
		idx = nIdx + len(q)
	}
	return ret
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
