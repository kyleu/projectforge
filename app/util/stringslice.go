package util

import (
	"fmt"
	"slices"
	"sort"
)

type StringSlice struct {
	Slice []string
}

func NewStringSlice(a ...string) *StringSlice {
	return &StringSlice{Slice: a}
}

func NewStringSliceWithSize(size int) *StringSlice {
	return &StringSlice{Slice: make([]string, 0, size)}
}

func (s *StringSlice) Empty() bool {
	return len(s.Slice) == 0
}

func (s *StringSlice) Length() int {
	return len(s.Slice)
}

func (s *StringSlice) SliceSafe() []string {
	if s == nil || len(s.Slice) == 0 {
		return nil
	}
	return s.Slice
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

func (s *StringSlice) PushUnique(strs ...string) {
	for _, x := range strs {
		if !slices.Contains(s.Slice, x) {
			s.Push(x)
		}
	}
}

func (s *StringSlice) Pushf(msg string, args ...any) {
	s.Push(fmt.Sprintf(msg, args...))
}

func (s *StringSlice) PushfUnlessNil(msg string, args ...any) {
	if s != nil {
		s.Pushf(msg, args...)
	}
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

func (s *StringSlice) UnmarshalJSON(data []byte) error {
	slc, err := FromJSONObj[[]string](data)
	if err != nil {
		return err
	}
	s.Slice = slc
	return nil
}

func (s *StringSlice) MarshalJSON() ([]byte, error) {
	return ToJSONBytes(s.Slice, true), nil
}
