package util

import (
	"fmt"
	"reflect"
)

func StringArrayMaxLength(a []string) int {
	ret := 0
	for _, x := range a {
		l := len(x)
		if l > ret {
			ret = l
		}
	}
	return ret
}

func StringArrayQuoted(a []string) []string {
	ret := make([]string, 0, len(a))
	for _, x := range a {
		ret = append(ret, fmt.Sprintf("%q", x))
	}
	return ret
}

func StringArrayFromInterfaces(a []any, maxLength int) []string {
	ret := make([]string, 0, len(a))
	for _, x := range a {
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
	}
	return ret
}

func ArrayRemoveDuplicates[T comparable](x []T) []T {
	m := make(map[T]struct{}, len(x))
	ret := make([]T, 0, len(x))
	for _, item := range x {
		if _, ok := m[item]; !ok {
			m[item] = struct{}{}
			ret = append(ret, item)
		}
	}
	return ret
}

func InterfaceArrayFrom[T any](x ...T) []any {
	ret := make([]any, len(x))
	for idx, item := range x {
		ret[idx] = item
	}
	return ret
}

func StringArrayOxfordComma(names []string, separator string) string {
	ret := ""
	for idx, name := range names {
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
	}
	return ret
}

func ArrayRemoveNil[T any](x []*T) []*T {
	ret := make([]*T, 0, len(x))
	for _, item := range x {
		if item != nil {
			ret = append(ret, item)
		}
	}
	return ret
}

func ArrayDefererence[T any](x []*T) []T {
	ret := make([]T, 0, len(x))
	for _, item := range x {
		if item != nil {
			ret = append(ret, *item)
		}
	}
	return ret
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
		ret := make([]any, 0, rfl.Len())
		for i := 0; i < rfl.Len(); i++ {
			e := rfl.Index(i)
			ret = append(ret, e.Interface())
		}
		return ret
	}
	return []any{dest}
}

// TakeFirstN returns the first N items from a slice
func TakeFirstN[V any](n int, items []V) []V {
	if n > len(items) {
		n = len(items)
	}
	return items[:n]
}

// TakeLastN returns the last N items from a slice
func TakeLastN[V any](n int, items []V) []V {
	if n > len(items) {
		return items
	}
	return items[len(items)-n:]
}
