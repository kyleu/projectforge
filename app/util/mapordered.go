package util

import (
	"bytes"
	"encoding/xml"
	"maps"
	"slices"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/samber/lo"
)

type OrderedPair[V any] struct {
	K string `json:"k"`
	V V      `json:"v"`
}

func NewOrderedPair[V any](k string, v V) *OrderedPair[V] {
	return &OrderedPair[V]{K: k, V: v}
}

type OrderedPairs[V any] []*OrderedPair[V]

type OrderedMap[V any] struct {
	Lexical bool
	Order   []string
	Map     map[string]V
}

func NewOrderedMap[V any](lexical bool, capacity int, pairs ...*OrderedPair[V]) *OrderedMap[V] {
	ret := &OrderedMap[V]{Lexical: lexical, Order: make([]string, 0, capacity), Map: make(map[string]V, capacity)}
	for _, p := range pairs {
		ret.Set(p.K, p.V)
	}
	return ret
}

func NewOMap[V any]() *OrderedMap[V] {
	return NewOrderedMap[V](false, 0)
}

func (o *OrderedMap[V]) Set(k string, v V) {
	if o.Map == nil {
		o.Map = map[string]V{}
	}
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
		return NewOrderedPair(k, o.GetSimple(k))
	})
}

func (o *OrderedMap[V]) Remove(k string) {
	o.Order = lo.Filter(o.Order, func(x string, _ int) bool {
		return x != k
	})
	delete(o.Map, k)
}

func (o *OrderedMap[V]) Clone() *OrderedMap[V] {
	if o == nil {
		return nil
	}
	return &OrderedMap[V]{Lexical: o.Lexical, Order: ArrayCopy(o.Order), Map: maps.Clone(o.Map)}
}

func (o *OrderedMap[V]) Clear() {
	o.Order = nil
	o.Map = map[string]V{}
}

func (o OrderedMap[V]) MarshalYAML() (any, error) {
	return o.Map, nil
}

func (o OrderedMap[V]) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if idx := strings.LastIndex(start.Name.Local, "["); idx > -1 {
		start.Name.Local = start.Name.Local[:idx]
	}

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

func (o *OrderedMap[V]) UnmarshalJSON(data []byte) error {
	o.Clear()
	err := jsonparser.ObjectEach(data, func(keyData []byte, valueData []byte, dataType jsonparser.ValueType, offset int) error {
		if dataType == jsonparser.String {
			valueData = data[offset-len(valueData)-2 : offset]
		}

		key, err := DecodeUTF8(keyData)
		if err != nil {
			return err
		}
		var value V
		if err := FromJSON(valueData, &value); err != nil {
			return err
		}
		o.Set(key, value)
		return nil
	})
	if err != nil {
		return err
	}
	if o.Lexical {
		slices.Sort(o.Order)
	}
	return nil
}

func (o OrderedMap[V]) MarshalJSON() ([]byte, error) {
	var b []byte
	buf := bytes.NewBuffer(b)
	buf.WriteByte('{')
	l := len(o.Order)
	for i, key := range o.Order {
		buf.Write(ToJSONBytes(key, false))
		buf.WriteByte(':')
		buf.Write(ToJSONBytes(o.Map[key], false))
		if i != l-1 {
			buf.WriteByte(',')
		}
	}
	buf.WriteByte('}')
	return buf.Bytes(), nil
}

type OrderedMaps[V any] []*OrderedMap[V]

type ToOrderedMap[T any] interface {
	ToOrderedMap() *OrderedMap[T]
}

type ToOrderedMaps[T any] interface {
	ToOrderedMaps() OrderedMaps[T]
}
