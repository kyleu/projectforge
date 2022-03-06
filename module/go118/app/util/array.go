package util

import (
	"fmt"
)

func ArrayIndexOf[T comparable](a []T, v T) int {
	for idx, x := range a {
		if x == v {
			return idx
		}
	}
	return -1
}

func InterfaceArray[T any](s []T) []interface{} {
	ret := make([]interface{}, len(s))
	for idx, item := range s {
		ret[idx] = item
	}
	return ret
}

func StringArrayContains(a []string, str string) bool {
	return ArrayIndexOf[string](a, str) >= 0
}

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

func StringArrayCopy(a []string) []string {
	ret := make([]string, 0, len(a))
	return append(ret, a...)
}

func StringArrayQuoted(a []string) []string {
	ret := make([]string, 0, len(a))
	for _, x := range a {
		ret = append(ret, fmt.Sprintf("%q", x))
	}
	return ret
}

func StringArrayFromInterfaces(a []interface{}, maxLength int) []string {
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
