package file

import (
	"sort"
	"strings"
)

type Replacement struct {
	K string `json:"k"`
	V string `json:"v"`
}

type Replacements []*Replacement

func (r Replacements) Sort() {
	sort.Slice(r, func(i, j int) bool {
		return len(r[i].V) > len(r[j].V)
	})
}

func (r Replacements) ToReplacer(prefix string, suffix string) *strings.Replacer {
	args := make([]string, 0, len(r)*2)
	for _, x := range r {
		args = append(args, prefix+x.K+suffix, x.V)
	}
	return strings.NewReplacer(args...)
}

type Changeset struct {
	Replacements Replacements `json:"replacements,omitempty"`
}
