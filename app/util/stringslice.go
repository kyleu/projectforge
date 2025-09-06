package util

import (
	"fmt"
	"sort"
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

func (s *StringSlice) TotalLength() int {
	var total int
	for _, str := range s.Slice {
		total += len(str)
	}
	return total
}

func (s *StringSlice) Push(strs ...string) {
	s.Slice = append(s.Slice, strs...)
}

func (s *StringSlice) Pushf(msg string, args ...any) {
	s.Push(fmt.Sprintf(msg, args...))
}

func (s *StringSlice) Join(x string) string {
	return StringJoin(s.Slice, x)
}

func (s *StringSlice) String() string {
	return s.Join("\n")
}

func (s *StringSlice) Sort() {
	sort.Strings(s.Slice)
}
