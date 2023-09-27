// Package util - Content managed by Project Forge, see [projectforge.md] for details.
package util

import (
	"cmp"
	"fmt"
	"reflect"
	"slices"

	"github.com/samber/lo"
)

func StringArrayMaxLength(a []string) int {
	return len(lo.MaxBy(a, func(x string, max string) bool {
		return len(x) > len(max)
	}))
}

func StringArrayQuoted(a []string) []string {
	return lo.Map(a, func(x string, _ int) string {
		return fmt.Sprintf("%q", x)
	})
}

func StringArrayFromInterfaces(a []any, maxLength int) []string {
	ret := make([]string, 0, len(a))
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
		ret = append(ret, v)
	})
	return ret
}

func ArrayRemoveDuplicates[T comparable](x []T) []T {
	return lo.Uniq(x)
}

func ArraySorted[T cmp.Ordered](x []T) []T {
	slices.Sort(x)
	return x
}

func InterfaceArrayFrom[T any](x ...T) []any {
	return lo.Map(x, func(item T, idx int) any {
		return item
	})
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

func ArrayRemoveNil[T any](x []*T) []*T {
	return lo.Reject(x, func(item *T, _ int) bool {
		return item == nil
	})
}

func ArrayDereference[T any](x []*T) []T {
	return lo.Map(x, func(item *T, _ int) T {
		return lo.FromPtr(item)
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

func ArrayFromAny(dest any) []any {
	defer func() { _ = recover() }()
	rfl := reflect.ValueOf(dest)
	if rfl.Kind() == reflect.Ptr {
		rfl = rfl.Elem()
	}
	if k := rfl.Kind(); k == reflect.Array || k == reflect.Slice {
		return lo.Times(rfl.Len(), func(i int) any {
			return rfl.Index(i).Interface()
		})
	}
	return []any{dest}
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
