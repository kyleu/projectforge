package util

import (
	"cmp"
	"fmt"
	"reflect"
	"slices"

	"github.com/pkg/errors"
	"github.com/samber/lo"
)

func StringArrayMaxLength(a []string) int {
	return len(lo.MaxBy(a, func(x string, max string) bool {
		return len(x) > len(max)
	}))
}

func ArrayToStringArray[T any](a []T) []string {
	return lo.Map(a, func(x T, _ int) string {
		return fmt.Sprint(x)
	})
}

func StringArrayQuoted(a []string) []string {
	return lo.Map(a, func(x string, _ int) string {
		return fmt.Sprintf("%q", x)
	})
}

func StringArrayFromAny(a []any, maxLength int) []string {
	ret := NewStringSlice(make([]string, 0, len(a)))
	lo.ForEach(a, func(x any, _ int) {
		var v string
		switch t := x.(type) {
		case string:
			v = t
		case []byte:
			v = string(t)
		default:
			v = fmt.Sprint(x)
		}
		if maxLength > 0 && len(v) > maxLength {
			v = v[:maxLength] + "... (truncated)"
		}
		ret.Push(v)
	})
	return ret.Slice
}

func ArrayRemoveDuplicates[T comparable](x []T) []T {
	return lo.Uniq(x)
}

func ArraySorted[T cmp.Ordered](x []T) []T {
	slices.Sort(x)
	return x
}

func StringArrayOxfordComma(names []string, separator string) string {
	ret := ""
	lo.ForEach(names, func(name string, idx int) {
		if idx > 0 {
			if idx == (len(names) - 1) {
				if idx > 1 {
					ret += ","
				}
				ret += " " + separator + " "
			} else {
				ret += ", "
			}
		}
		ret += name
	})
	return ret
}

func ArrayTransform[T any, U any](x []T, f func(T) U) []U {
	return lo.Map(x, func(i T, _ int) U {
		return f(i)
	})
}

func ArraySplit[T any](xs []T, fn func(T) bool) ([]T, []T) {
	var t, f []T
	lo.ForEach(xs, func(x T, _ int) {
		if fn(x) {
			t = append(t, x)
		} else {
			f = append(f, x)
		}
	})
	return t, f
}

func ArrayRemoveNil[T any](x []*T) []*T {
	return lo.Reject(x, func(el *T, _ int) bool {
		return el == nil
	})
}

func ArrayDereference[T any](x []*T) []T {
	return lo.Map(x, func(el *T, _ int) T {
		return lo.FromPtr(el)
	})
}

func LengthAny(dest any) int {
	defer func() { _ = recover() }()
	rfl := reflect.ValueOf(dest)
	if rfl.Kind() == reflect.Ptr {
		rfl = rfl.Elem()
	}
	return rfl.Len()
}

func ArrayFromAny[T any](dest any) []T {
	defer func() { _ = recover() }()
	rfl := reflect.ValueOf(dest)
	if rfl.Kind() == reflect.Ptr {
		rfl = rfl.Elem()
	}
	if k := rfl.Kind(); k == reflect.Array || k == reflect.Slice {
		return lo.Times(rfl.Len(), func(i int) T {
			x := rfl.Index(i).Interface()
			if t, ok := x.(T); ok {
				return t
			}
			var t T
			return t
		})
	}
	if t, ok := dest.(T); ok {
		return []T{t}
	}
	return nil
}

func ArrayTest(dest any) bool {
	defer func() { _ = recover() }()
	k := reflect.ValueOf(dest).Kind()
	return k == reflect.Array || k == reflect.Slice
}

func ArrayFlatten[T any](arrs ...[]T) []T {
	return lo.Flatten(arrs)
}

func ArrayFirstN[V any](items []V, n int) []V {
	if n > len(items) {
		return items
	}
	return items[:n]
}

func ArrayLastN[V any](items []V, n int) []V {
	if n > len(items) {
		return items
	}
	return items[len(items)-n:]
}

func MapError[T any, U any](xa []T, f func(el T, idx int) (U, error)) ([]U, error) {
	ret := make([]U, 0, len(xa))
	for i, x := range xa {
		candidate, err := f(x, i)
		if err != nil {
			return nil, errors.Wrapf(err, "error processing element [%d]", i)
		}
		ret = append(ret, candidate)
	}
	return ret, nil
}
