package util

import (
	"fmt"
)

func StringArrayContains(a []string, str string) bool {
	return StringArrayIndexOf(a, str) >= 0
}

func StringArrayIndexOf(a []string, str string) int {
	for idx, x := range a {
		if x == str {
			return idx
		}
	}
	return -1
}

func StringArrayCopy(a []string) []string {
	ret := make([]string, 0, len(a))
	return append(ret, a...)
}

func StringArrayQuoted(a []string) []string {
	ret := make([]string, 0, len(a))
	for _, x := range a {
		ret = append(ret, `"`+x+`"`)
	}
	return ret
}

func StringArrayFromInterfaces(a []interface{}) []string {
	ret := make([]string, 0, len(a))
	for _, x := range a {
		ret = append(ret, fmt.Sprint(x))
	}
	return ret
}
