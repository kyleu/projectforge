// Package util - Content managed by Project Forge, see [projectforge.md] for details.
package util

import (
	"fmt"
	"strings"

	"github.com/samber/lo"
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

func (d Diff) StringVerbose() string {
	return fmt.Sprintf("%s (%q != %q)", d.Path, d.Old, d.New)
}

type Diffs []*Diff

func (d Diffs) String() string {
	return strings.Join(lo.Map(d, func(x *Diff, _ int) string {
		return x.String()
	}), "; ")
}

func (d Diffs) StringVerbose() string {
	return strings.Join(lo.Map(d, func(x *Diff, _ int) string {
		return x.StringVerbose()
	}), "; ")
}

type DiffsSet map[string]Diffs

func DiffObjects(l any, r any, path ...string) Diffs {
	return DiffObjectsIgnoring(l, r, nil, path...)
}

func DiffObjectsIgnoring(l any, r any, ignored []string, path ...string) Diffs {
	if len(path) > 0 && lo.Contains(ignored, path[len(path)-1]) {
		return nil
	}
	if l == nil {
		return Diffs{NewDiff(strings.Join(path, "."), "", fmt.Sprint(r))}
	}
	if r == nil {
		return Diffs{NewDiff(strings.Join(path, "."), fmt.Sprint(l), "")}
	}
	if lt, rt := fmt.Sprintf("%T", l), fmt.Sprintf("%T", r); lt != rt {
		return Diffs{NewDiff(strings.Join(path, "."), ToJSONCompact(l), ToJSONCompact(r))}
	}
	return diffType(l, r, ignored, false, path...)
}

func diffType(l any, r any, ignored []string, recursed bool, path ...string) Diffs {
	var ret Diffs
	switch t := l.(type) {
	case ValueMap:
		ret = append(ret, diffMaps(t, r, ignored, path...)...)
	case map[string]any:
		ret = append(ret, diffMaps(t, r, ignored, path...)...)
	case map[string]int:
		ret = append(ret, diffIntMaps(t, r, ignored, path...)...)
	case []any:
		ret = append(ret, diffArrays(t, r, ignored, path...)...)
	case Diffs:
		rm, _ := r.(Diffs)
		lo.ForEach(t, func(v *Diff, idx int) {
			rv := rm[idx]
			ret = append(ret, DiffObjectsIgnoring(v, rv, ignored, append([]string{}, path...)...)...)
		})
	case int64:
		i, _ := r.(int64)
		if t != i {
			ret = append(ret, NewDiff(strings.Join(path, "."), fmt.Sprint(t), fmt.Sprint(i)))
		}
	case int:
		i, _ := r.(int)
		if t != i {
			ret = append(ret, NewDiff(strings.Join(path, "."), fmt.Sprint(t), fmt.Sprint(i)))
		}
	case float64:
		i, _ := r.(float64)
		if t != i {
			ret = append(ret, NewDiff(strings.Join(path, "."), fmt.Sprint(t), fmt.Sprint(i)))
		}
	case string:
		s, _ := r.(string)
		if t != s {
			ret = append(ret, NewDiff(strings.Join(path, "."), t, s))
		}
	default:
		lj, rj := ToJSONCompact(l), ToJSONCompact(r)
		if !recursed && (strings.HasPrefix(lj, "{") || strings.HasPrefix(lj, "[")) {
			lx, _ := FromJSONAny([]byte(lj))
			rx, _ := FromJSONAny([]byte(rj))
			ret = append(ret, diffType(lx, rx, ignored, true, path...)...)
		} else if lj != rj {
			ret = append(ret, NewDiff(strings.Join(path, "."), lj, rj))
		}
	}
	return ret
}

func diffArrays(l []any, r any, ignored []string, path ...string) Diffs {
	var ret Diffs
	rm, _ := r.([]any)
	lo.ForEach(l, func(v any, idx int) {
		if len(rm) > idx {
			rv := rm[idx]
			ret = append(ret, DiffObjectsIgnoring(v, rv, ignored, append(append([]string{}, path...), fmt.Sprint(idx))...)...)
		} else {
			ret = append(ret, DiffObjectsIgnoring(v, nil, ignored, append(append([]string{}, path...), fmt.Sprint(idx))...)...)
		}
	})
	if len(rm) > len(l) {
		for i := len(l); i < len(rm); i++ {
			ret = append(ret, DiffObjectsIgnoring(nil, rm[i], ignored, append(append([]string{}, path...), fmt.Sprint(i))...)...)
		}
	}
	return ret
}

func diffMaps(l map[string]any, r any, ignored []string, path ...string) Diffs {
	var ret Diffs
	rm, ok := r.(map[string]any)
	if !ok {
		rm, _ = r.(ValueMap)
	}
	for k, v := range l {
		if lo.Contains(ignored, k) {
			continue
		}
		rv := rm[k]
		ret = append(ret, DiffObjectsIgnoring(v, rv, ignored, append(append([]string{}, path...), k)...)...)
	}
	for k, v := range rm {
		if lo.Contains(ignored, k) {
			continue
		}
		if _, exists := l[k]; !exists {
			ret = append(ret, DiffObjectsIgnoring(nil, v, ignored, append(append([]string{}, path...), k)...)...)
		}
	}
	return ret
}

func diffIntMaps(l map[string]int, r any, ignored []string, path ...string) Diffs {
	var ret Diffs
	rm, _ := r.(map[string]int)
	for k, v := range l {
		if lo.Contains(ignored, k) {
			continue
		}
		rv := rm[k]
		ret = append(ret, DiffObjectsIgnoring(v, rv, ignored, append(append([]string{}, path...), k)...)...)
	}
	for k, v := range rm {
		if lo.Contains(ignored, k) {
			continue
		}
		if _, exists := l[k]; !exists {
			ret = append(ret, DiffObjectsIgnoring(nil, v, ignored, append(append([]string{}, path...), k)...)...)
		}
	}
	return ret
}
