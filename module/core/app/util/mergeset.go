package util

import (
	"cmp"
	"encoding/xml"
	"fmt"
	"maps"
	"strings"

	"github.com/samber/lo"
)

type Mergeable[K cmp.Ordered, T any] interface {
	Key() K
	Merge(other T) (T, error)
}

type MergeSet[K cmp.Ordered, T Mergeable[K, T]] struct {
	Map map[K]T
}

func (m *MergeSet[K, T]) Merge(other *MergeSet[K, T]) (*MergeSet[K, T], error) {
	ret := NewMergeSet[K, T]()
	for _, v := range m.Map {
		if err := ret.Set(v); err != nil {
			return nil, err
		}
	}
	for _, v := range other.Map {
		if err := ret.Set(v); err != nil {
			return nil, err
		}
	}
	return ret, nil
}

func NewMergeSet[K cmp.Ordered, T Mergeable[K, T]](capacity ...int) *MergeSet[K, T] {
	return &MergeSet[K, T]{Map: make(map[K]T, lo.Sum(capacity))}
}

func (m *MergeSet[K, T]) Set(v T) error {
	k := v.Key()
	curr, ok := m.Map[k]
	if ok {
		var err error
		v, err = v.Merge(curr)
		if err != nil {
			return err
		}
	}
	m.Map[k] = v
	return nil
}

func (m *MergeSet[K, T]) Contains(x T) bool {
	_, ok := m.Map[x.Key()]
	return ok
}

func (m *MergeSet[K, T]) Remove(x T) {
	delete(m.Map, x.Key())
}

func (m *MergeSet[K, T]) Length() int {
	return len(m.Map)
}

func (m *MergeSet[K, T]) Entries() []T {
	return lo.Map(MapKeysSorted(m.Map), func(k K, _ int) T {
		return m.Map[k]
	})
}

func (m *MergeSet[K, T]) Clone() *MergeSet[K, T] {
	if m == nil {
		return nil
	}
	return &MergeSet[K, T]{Map: maps.Clone(m.Map)}
}

func (m *MergeSet[K, T]) MarshalYAML() (any, error) {
	return MapKeysSorted(m.Map), nil
}

func (m MergeSet[K, T]) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if idx := strings.LastIndex(start.Name.Local, "["); idx > -1 {
		start.Name.Local = start.Name.Local[:idx]
	}

	if err := e.EncodeToken(start); err != nil {
		return err
	}

	for _, v := range m.Entries() {
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

func (m *MergeSet[K, T]) UnmarshalJSON(b []byte) error {
	var ret []T
	if err := FromJSON(b, &ret); err != nil {
		return err
	}
	for _, x := range ret {
		if err := m.Set(x); err != nil {
			return err
		}
	}
	return nil
}

func (m *MergeSet[K, T]) MarshalJSON() ([]byte, error) {
	return ToJSONBytes(m.Entries(), true), nil
}
