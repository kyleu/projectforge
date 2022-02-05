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
		return append(ret, NewDiff(strings.Join(path, "."), "", fmt.Sprint(r)))
	}
	if r == nil {
		return append(ret, NewDiff(strings.Join(path, "."), fmt.Sprint(l), ""))
	}
	if lt, rt := fmt.Sprintf("%T", l), fmt.Sprintf("%T", r); lt != rt {
		return append(ret, NewDiff(strings.Join(path, "."), ToJSONCompact(l), ToJSONCompact(r)))
	}

	switch t := l.(type) {
	case ValueMap:
		ret = append(ret, diffMaps(t, r, path...)...)
	case map[string]interface{}:
		ret = append(ret, diffMaps(t, r, path...)...)
	case map[string]int:
		ret = append(ret, diffIntMaps(t, r, path...)...)
	case []interface{}:
		ret = append(ret, diffArrays(t, r, path...)...)
	case Diffs:
		rm, _ := r.(Diffs)
		for idx, v := range t {
			rv := rm[idx]
			ret = append(ret, DiffObjects(v, rv, append([]string{}, path...)...)...)
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
		if lj, rj := ToJSONCompact(l), ToJSONCompact(r); lj != rj {
			ret = append(ret, NewDiff(strings.Join(path, "."), lj, rj))
		}
	}

	return ret
}

func diffArrays(l []interface{}, r interface{}, path ...string) Diffs {
	var ret Diffs
	rm, _ := r.([]interface{})
	for idx, v := range l {
		if len(rm) > idx {
			rv := rm[idx]
			ret = append(ret, DiffObjects(v, rv, append(append([]string{}, path...), fmt.Sprint(idx))...)...)
		} else {
			ret = append(ret, DiffObjects(v, nil, append(append([]string{}, path...), fmt.Sprint(idx))...)...)
		}
	}
	if len(rm) > len(l) {
		for i := len(l); i < len(rm); i++ {
			ret = append(ret, DiffObjects(nil, rm[i], append(append([]string{}, path...), fmt.Sprint(i))...)...)
		}
	}
	return ret
}

func diffMaps(l map[string]interface{}, r interface{}, path ...string) Diffs {
	var ret Diffs
	rm, ok := r.(map[string]interface{})
	if !ok {
		rm, _ = r.(ValueMap)
	}
	for k, v := range l {
		rv := rm[k]
		ret = append(ret, DiffObjects(v, rv, append(append([]string{}, path...), k)...)...)
	}
	for k, v := range rm {
		if _, exists := l[k]; !exists {
			ret = append(ret, DiffObjects(nil, v, append(append([]string{}, path...), k)...)...)
		}
	}
	return ret
}

func diffIntMaps(l map[string]int, r interface{}, path ...string) Diffs {
	var ret Diffs
	rm, _ := r.(map[string]int)
	for k, v := range l {
		rv := rm[k]
		ret = append(ret, DiffObjects(v, rv, append(append([]string{}, path...), k)...)...)
	}
	for k, v := range rm {
		if _, exists := l[k]; !exists {
			ret = append(ret, DiffObjects(nil, v, append(append([]string{}, path...), k)...)...)
		}
	}
	return ret
}
