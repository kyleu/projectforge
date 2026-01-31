package util

import (
	"strings"

	"github.com/samber/lo"
)

func JoinLines(ss []string, delim string, maxLen int) []string {
	if maxLen == 0 {
		return []string{StringJoin(ss, delim)}
	}
	ret := &StringSlice{}
	var curr string
	lo.ForEach(ss, func(s string, _ int) {
		if curr != "" && (len(curr)+len(delim)+len(s)) > maxLen {
			ret.Push(curr)
			curr = ""
		}
		if curr != "" {
			curr += delim
		}
		curr += s
	})
	if curr != "" {
		ret.Push(curr)
	}
	return ret.Slice
}

func JoinLinesFull(ss []string, delim string, maxLen int, prefix string, indent string, suffix string) string {
	lines := JoinLines(ss, delim, maxLen)
	if len(lines) == 1 && len(lines[0]) < 100 {
		return prefix + strings.TrimSuffix(lines[0], ",") + suffix
	}
	ret := prefix
	lo.ForEach(lines, func(l string, _ int) {
		ret += indent + l
	})
	ret += suffix
	return ret
}
