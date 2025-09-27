package util

import (
	"cmp"
	"fmt"
	"slices"
	"strings"

	"github.com/samber/lo"
)

type Diff struct {
	Path string `json:"path"`
	Old  string `json:"o,omitzero"`
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
	return StringJoin(lo.Map(d, func(x *Diff, _ int) string {
		return x.String()
	}), "; ")
}

func (d Diffs) Sorted() Diffs {
	ret := ArrayCopy(d)
	slices.SortFunc(ret, func(l *Diff, r *Diff) int {
		return cmp.Compare(l.Path, r.Path)
	})
	return ret
}

func (d Diffs) Without(paths ...string) Diffs {
	return lo.Filter(d, func(x *Diff, _ int) bool {
		return !lo.Contains(paths, x.Path)
	})
}

func (d Diffs) StringVerbose() string {
	return StringJoin(lo.Map(d, func(x *Diff, _ int) string {
		return x.StringVerbose()
	}), "; ")
}

type DiffsSet map[string]Diffs

func (d DiffsSet) Keys() []string {
	return MapKeysSorted(d)
}

func DiffObjects(l any, r any, path ...string) Diffs {
	return DiffObjectsIgnoring(l, r, nil, path...)
}

func DiffObjectsIgnoring(l any, r any, ignored []string, path ...string) Diffs {
	if len(path) > 0 && lo.Contains(ignored, path[len(path)-1]) {
		return nil
	}
	if l == nil && r == nil {
		return nil
	}
	if l == nil {
		return Diffs{NewDiff(StringJoin(path, "."), "", fmt.Sprint(r))}
	}
	if r == nil {
		return Diffs{NewDiff(StringJoin(path, "."), fmt.Sprint(l), "")}
	}
	if lt, rt := fmt.Sprintf("%T", l), fmt.Sprintf("%T", r); lt != rt {
		return Diffs{NewDiff(StringJoin(path, "."), ToJSONCompact(l), ToJSONCompact(r))}
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
		rm := CastOK[Diffs](r)
		lo.ForEach(t, func(v *Diff, idx int) {
			rv := rm[idx]
			ret = append(ret, DiffObjectsIgnoring(v, rv, ignored, ArrayCopy(path)...)...)
		})
	case int64:
		i := CastOK[int64](r)
		if t != i {
			ret = append(ret, NewDiff(StringJoin(path, "."), fmt.Sprint(t), fmt.Sprint(i)))
		}
	case int:
		i := CastOK[int](r)
		if t != i {
			ret = append(ret, NewDiff(StringJoin(path, "."), fmt.Sprint(t), fmt.Sprint(i)))
		}
	case float64:
		f := CastOK[float64](r)
		if t != f {
			ret = append(ret, NewDiff(StringJoin(path, "."), fmt.Sprint(t), fmt.Sprint(f)))
		}
	case string:
		s := CastOK[string](r)
		if t != s {
			ret = append(ret, NewDiff(StringJoin(path, "."), t, s))
		}
	default:
		lj, rj := ToJSONCompact(l), ToJSONCompact(r)
		if !recursed && (strings.HasPrefix(lj, "{") || strings.HasPrefix(lj, "[")) {
			lx, _ := FromJSONAny([]byte(lj))
			rx, _ := FromJSONAny([]byte(rj))
			ret = append(ret, diffType(lx, rx, ignored, true, path...)...)
		} else if lj != rj {
			ret = append(ret, NewDiff(StringJoin(path, "."), lj, rj))
		}
	}
	return ret
}

func diffArrays(l []any, r any, ignored []string, path ...string) Diffs {
	var ret Diffs
	rm := CastOK[[]any](r)
	lo.ForEach(l, func(v any, idx int) {
		if len(rm) > idx {
			rv := rm[idx]
			ret = append(ret, DiffObjectsIgnoring(v, rv, ignored, append(ArrayCopy(path), fmt.Sprint(idx))...)...)
		} else {
			ret = append(ret, DiffObjectsIgnoring(v, nil, ignored, append(ArrayCopy(path), fmt.Sprint(idx))...)...)
		}
	})
	if len(rm) > len(l) {
		for i := len(l); i < len(rm); i++ {
			ret = append(ret, DiffObjectsIgnoring(nil, rm[i], ignored, append(ArrayCopy(path), fmt.Sprint(i))...)...)
		}
	}
	return ret
}

func diffMaps(l map[string]any, r any, ignored []string, path ...string) Diffs {
	var ret Diffs
	rm, err := Cast[map[string]any](r)
	if err != nil {
		rm = CastOK[ValueMap](r)
	}
	for k, v := range l {
		if lo.Contains(ignored, k) {
			continue
		}
		rv := rm[k]
		ret = append(ret, DiffObjectsIgnoring(v, rv, ignored, append(ArrayCopy(path), k)...)...)
	}
	for k, v := range rm {
		if lo.Contains(ignored, k) {
			continue
		}
		if _, exists := l[k]; !exists {
			ret = append(ret, DiffObjectsIgnoring(nil, v, ignored, append(ArrayCopy(path), k)...)...)
		}
	}
	return ret
}

func diffIntMaps(l map[string]int, r any, ignored []string, path ...string) Diffs {
	var ret Diffs
	rm := CastOK[map[string]int](r)
	for k, v := range l {
		if lo.Contains(ignored, k) {
			continue
		}
		rv := rm[k]
		ret = append(ret, DiffObjectsIgnoring(v, rv, ignored, append(ArrayCopy(path), k)...)...)
	}
	for k, v := range rm {
		if lo.Contains(ignored, k) {
			continue
		}
		if _, exists := l[k]; !exists {
			ret = append(ret, DiffObjectsIgnoring(nil, v, ignored, append(ArrayCopy(path), k)...)...)
		}
	}
	return ret
}
