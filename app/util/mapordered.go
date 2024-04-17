// Package util - Content managed by Project Forge, see [projectforge.md] for details.
package util

import (
	"bytes"
	"cmp"
	"encoding/json"
	"encoding/xml"
	"slices"

	"github.com/samber/lo"
)

type OrderedPair[V any] struct {
	K string `json:"k"`
	V V      `json:"v"`
}

type OrderedMap[V any] struct {
	Lexical bool
	Order   []string
	Map     map[string]V
}

func NewOrderedMap[V any](lexical bool, capacity int) *OrderedMap[V] {
	return &OrderedMap[V]{Lexical: lexical, Order: make([]string, 0, capacity), Map: make(map[string]V, capacity)}
}

func NewOMap[V any]() *OrderedMap[V] {
	return NewOrderedMap[V](false, 0)
}

func (o *OrderedMap[V]) Append(k string, v V) {
	o.Order = append(o.Order, k)
	o.Map[k] = v
	if o.Lexical {
		slices.Sort(o.Order)
	}
}

func (o *OrderedMap[V]) Set(k string, v V) {
	if _, ok := o.Map[k]; !ok {
		o.Order = append(o.Order, k)
	}
	o.Map[k] = v
	if o.Lexical {
		slices.Sort(o.Order)
	}
}

func (o *OrderedMap[V]) HasKey(k string) bool {
	_, ok := o.Map[k]
	return ok
}

func (o *OrderedMap[V]) IndexOf(k string) int {
	return slices.Index(o.Order, k)
}

func (o *OrderedMap[V]) Get(k string) (V, bool) {
	ret, ok := o.Map[k]
	return ret, ok
}

func (o *OrderedMap[V]) GetSimple(k string) V {
	return o.Map[k]
}

func (o *OrderedMap[V]) Pairs() []*OrderedPair[V] {
	return lo.Map(o.Order, func(k string, _ int) *OrderedPair[V] {
		return &OrderedPair[V]{K: k, V: o.GetSimple(k)}
	})
}

func (o *OrderedMap[V]) Remove(k string) {
	o.Order = lo.Filter(o.Order, func(x string, _ int) bool {
		return x != k
	})
	delete(o.Map, k)
}

func (o OrderedMap[V]) MarshalYAML() (any, error) {
	return o.Map, nil
}

func (o OrderedMap[V]) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if err := e.EncodeToken(start); err != nil {
		return err
	}

	for _, key := range o.Order {
		n := xml.Name{Local: key}
		t := xml.StartElement{Name: n}

		v := o.Map[key]
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

func (o *OrderedMap[V]) UnmarshalJSON(b []byte) error {
	if err := FromJSON(b, &o.Map); err != nil {
		return err
	}

	index := make(map[string]int)
	lo.ForEach(lo.Keys(o.Map), func(key string, _ int) {
		o.Order = append(o.Order, key)
		esc := ToJSONBytes(key, false) // escape the key
		index[key] = bytes.Index(b, esc)
	})

	if o.Lexical {
		slices.Sort(o.Order)
	} else {
		slices.SortFunc(o.Order, func(l string, r string) int {
			return cmp.Compare(index[l], index[r])
		})
	}
	return nil
}

func (o OrderedMap[V]) MarshalJSON() ([]byte, error) {
	var b []byte
	buf := bytes.NewBuffer(b)
	buf.WriteByte('{')
	l := len(o.Order)
	for i, key := range o.Order {
		km, err := json.Marshal(key)
		if err != nil {
			return nil, err
		}
		buf.Write(km)
		buf.WriteByte(':')
		vm, err := json.Marshal(o.Map[key])
		if err != nil {
			return nil, err
		}
		buf.Write(vm)
		if i != l-1 {
			buf.WriteByte(',')
		}
	}
	buf.WriteByte('}')
	return buf.Bytes(), nil
}
