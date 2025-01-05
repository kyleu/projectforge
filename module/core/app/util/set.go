package util

import (
	"cmp"
	"encoding/xml"
	"fmt"
	"maps"
	"strings"

	"github.com/samber/lo"
)

type Set[T cmp.Ordered] struct {
	Map map[T]struct{}
}

func NewSet[T cmp.Ordered](capacity ...int) *Set[T] {
	return &Set[T]{Map: make(map[T]struct{}, lo.Sum(capacity))}
}

func (s *Set[T]) Set(v T) {
	s.Map[v] = struct{}{}
}

func (s *Set[T]) Contains(x T) bool {
	_, ok := s.Map[x]
	return ok
}

func (s *Set[T]) Remove(x T) {
	delete(s.Map, x)
}

func (s *Set[T]) Entries() []T {
	return MapKeysSorted(s.Map)
}

func (s *Set[T]) Clone() *Set[T] {
	if s == nil {
		return nil
	}
	return &Set[T]{Map: maps.Clone(s.Map)}
}

func (s Set[T]) MarshalYAML() (any, error) {
	return MapKeysSorted(s.Map), nil
}

func (s Set[T]) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if idx := strings.LastIndex(start.Name.Local, "["); idx > -1 {
		start.Name.Local = start.Name.Local[:idx]
	}

	if err := e.EncodeToken(start); err != nil {
		return err
	}

	for _, v := range s.Entries() {
		n := xml.Name{Local: fmt.Sprintf("%T", v)}
		t := xml.StartElement{Name: n}
		e := e.EncodeElement(v, t)
		if e != nil {
			return e
		}
	}

	if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
		return err
	}

	return e.Flush()
}

func (s *Set[T]) UnmarshalJSON(b []byte) error {
	var ret []T
	if err := FromJSON(b, &ret); err != nil {
		return err
	}
	for _, x := range ret {
		s.Set(x)
	}
	return nil
}

func (s Set[T]) MarshalJSON() ([]byte, error) {
	return ToJSONBytes(s.Entries(), true), nil
}
