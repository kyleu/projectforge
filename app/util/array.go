// Content managed by Project Forge, see [projectforge.md] for details.
package util

import (
	"fmt"
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
