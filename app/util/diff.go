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
	sb := make([]string, 0, len(d))
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

	if lt, rt := fmt.Sprintf("%T", l), fmt.Sprintf("%T", r); lt != rt {
		lj := ToJSON(l) //nolint
		rj := ToJSON(r) //nolint
		ret = append(ret, NewDiff(strings.Join(path, "."), lj, rj))
	}

	switch t := l.(type) {
	case ValueMap:
		rm, _ := r.(ValueMap)
		for k, v := range t {
			rv := rm[k]
			ret = append(ret, DiffObjects(v, rv, append(append([]string{}, path...), k)...)...)
		}
	case map[string]interface{}:
		rm, _ := r.(map[string]interface{})
		for k, v := range t {
			rv := rm[k]
			ret = append(ret, DiffObjects(v, rv, append(append([]string{}, path...), k)...)...)
		}
	case map[string]int:
		rm, _ := r.(map[string]int)
		for k, v := range t {
			rv := rm[k]
			ret = append(ret, DiffObjects(v, rv, append(append([]string{}, path...), k)...)...)
		}
	case int:
		i, _ := r.(int)
		if t != i {
			ret = append(ret, NewDiff(strings.Join(path, "."), fmt.Sprint(t), fmt.Sprint(i)))
		}
	case string:
		s, _ := r.(string)
		if t != s {
			ret = append(ret, NewDiff(strings.Join(path, "."), t, s))
		}
	default:
		if lj, rj := ToJSON(l), ToJSON(r); lj != rj {
			ret = append(ret, NewDiff(strings.Join(path, "."), lj, rj))
		}
	}

	return ret
}
