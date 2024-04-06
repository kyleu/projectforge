// Package util - Content managed by Project Forge, see [projectforge.md] for details.
package util

import (
	"fmt"
	"sort"
	"strings"
)

type StringSlice struct {
	Slice []string
}

func NewStringSlice(a []string) *StringSlice {
	return &StringSlice{Slice: a}
}

func (s *StringSlice) Empty() bool {
	return len(s.Slice) == 0
}

func (s *StringSlice) Push(strs ...string) {
	s.Slice = append(s.Slice, strs...)
}

func (s *StringSlice) Pushf(msg string, args ...any) {
	s.Push(fmt.Sprintf(msg, args...))
}

func (s *StringSlice) Join(x string) string {
	return strings.Join(s.Slice, x)
}

func (s *StringSlice) Sort() {
	sort.Strings(s.Slice)
}
