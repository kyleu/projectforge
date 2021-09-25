package util

import (
	"fmt"
	"strings"
)

type Diff struct {
	Path string `json:"path"`
	Old  string `json:"o,omitempty"`
	New  string `json:"n"`
}

func NewDiff(p string, o string, n string) *Diff {
	return &Diff{Path: p, Old: o, New: n}
}

func (d Diff) String() string {
	return d.Path
}

type Diffs []*Diff

func (d Diffs) String() string {
	var sb []string
	for _, x := range d {
		sb = append(sb, x.String())
	}
	return strings.Join(sb, "; ")
}

func DiffObjects(l interface{}, r interface{}, path ...string) Diffs {
	var ret Diffs

	if l == nil {
		ret = append(ret, NewDiff(strings.Join(path, "."), "", fmt.Sprint(r)))
	}
	if r == nil {
		ret = append(ret, NewDiff(strings.Join(path, "."), fmt.Sprint(l), ""))
	}

	lt := fmt.Sprintf("%T", l)
	rt := fmt.Sprintf("%T", r)

	if lt != rt {
		lj := ToJSON(l) //nolint
		rj := ToJSON(r) //nolint
		ret = append(ret, NewDiff(strings.Join(path, "."), lj, rj))
	}

	switch t := l.(type) {
	case ValueMap:
		rm := r.(ValueMap)
		for k, v := range t {
			rv := rm[k]
			ret = append(ret, DiffObjects(v, rv, append(append([]string{}, path...), k)...)...)
		}
	case map[string]interface{}:
		rm := r.(map[string]interface{})
		for k, v := range t {
			rv := rm[k]
			ret = append(ret, DiffObjects(v, rv, append(append([]string{}, path...), k)...)...)
		}
	case map[string]int:
		rm := r.(map[string]int)
		for k, v := range t {
			rv := rm[k]
			ret = append(ret, DiffObjects(v, rv, append(append([]string{}, path...), k)...)...)
		}
	case int:
		if t != r.(int) {
			ret = append(ret, NewDiff(strings.Join(path, "."), fmt.Sprint(t), fmt.Sprint(r.(int))))
		}
	case string:
		if t != r.(string) {
			ret = append(ret, NewDiff(strings.Join(path, "."), t, r.(string)))
		}
	default:
		if lj, rj := ToJSON(l), ToJSON(r); lj != rj {
			ret = append(ret, NewDiff(strings.Join(path, "."), lj, rj))
		}
	}

	return ret
}
