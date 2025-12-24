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
	return len(lo.MaxBy(a, func(x string, mx string) bool {
		return len(x) > len(mx)
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
	ret := NewStringSliceWithSize(Choose(len(a) > maxLength, maxLength, len(a)))
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

func ArrayCopy[S ~[]T, T any](x S) S {
	if x == nil {
		return nil
	}
	return append(S{}, x...)
}

func ArrayToAnyArray[T any](x []T) []any {
	ret := make([]any, 0, len(x))
	for _, v := range x {
		ret = append(ret, v)
	}
	return ret
}

func ArrayRemoveDuplicates[T comparable](x []T) []T {
	return lo.Uniq(x)
}

func ArrayFindDuplicates[T comparable](items []T) []T {
	seen := make(map[T]int)
	var duplicates []T
	for _, item := range items {
		seen[item]++
	}
	for item, count := range seen {
		if count > 1 {
			duplicates = append(duplicates, item)
		}
	}
	return duplicates
}

func ArraySorted[T cmp.Ordered](x []T) []T {
	slices.Sort(x)
	return x
}

func ArrayLimit[S ~[]T, T any](x S, limit int) (S, int) {
	if limit == 0 || limit > len(x) {
		limit = len(x)
	}
	return x[:limit], len(x) - limit
}

func StringArrayOxfordComma(names []string, separator string) string {
	var ret string
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

func ArraySplit[S ~[]T, T any](xs S, fn func(T) bool) (S, S) {
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

func ArrayRemoveEmpty[T comparable](x []T) []T {
	var check T
	return lo.Reject(x, func(el T, _ int) bool {
		return el == check
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

func ArrayFromAny[T any](dest any) ([]T, error) {
	defer func() { _ = recover() }()
	rfl := reflect.ValueOf(dest)
	if rfl.Kind() == reflect.Ptr {
		rfl = rfl.Elem()
	}
	if k := rfl.Kind(); k == reflect.Array || k == reflect.Slice {
		l := rfl.Len()
		ret := make([]T, 0, l)
		for i := range l {
			x := rfl.Index(i).Interface()
			t, err := Cast[T](x)
			if err != nil {
				return nil, err
			}
			ret = append(ret, t)
		}
		return ret, nil
	}
	if t, err := Cast[T](dest); err == nil {
		return []T{t}, nil
	}
	return nil, errors.Errorf("unable to convert [%T] to an array", dest)
}

func ArrayFromAnyOK[T any](dest any) []T {
	ret, _ := ArrayFromAny[T](dest)
	return ret
}

func ArrayTest(dest any) bool {
	defer func() { _ = recover() }()
	k := reflect.ValueOf(dest).Kind()
	return k == reflect.Array || k == reflect.Slice
}

func ArrayFlatten[S ~[]T, T any](arrs ...S) S {
	return lo.Flatten(arrs)
}

func ArrayFirstN[S ~[]T, T any](items S, n int) S {
	if n > len(items) {
		return items
	}
	return items[:n]
}

func ArrayLastN[S ~[]T, T any](items S, n int) S {
	if n > len(items) {
		return items
	}
	return items[len(items)-n:]
}

func ArrayReplaceOrAdd[S ~[]T, T any](a S, fn func(el T) bool, x T) S {
	ret := make([]T, 0, len(a)+1)
	var hit bool
	for _, el := range a {
		if !hit && fn(el) {
			hit = true
			ret = append(ret, x)
		} else {
			ret = append(ret, el)
		}
	}
	if !hit {
		ret = append(ret, x)
	}
	return ret
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
