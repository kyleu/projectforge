// Content managed by Project Forge, see [projectforge.md] for details.
package util

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"sort"
)

type OrderedMap struct {
	Lexical bool
	Order   []string
	Map     map[string]interface{}
}

func NewOrderedMap(lexical bool, capacity int) *OrderedMap {
	return &OrderedMap{Lexical: lexical, Order: make([]string, 0, capacity), Map: make(map[string]interface{}, capacity)}
}

func (o *OrderedMap) Append(k string, v interface{}) {
	o.Order = append(o.Order, k)
	o.Map[k] = v
	if o.Lexical {
		sort.Strings(o.Order)
	}
}

func (o *OrderedMap) Get(k string) (interface{}, bool) {
	ret, ok := o.Map[k]
	return ret, ok
}

func (o OrderedMap) MarshalYAML() (interface{}, error) {
	return o.Map, nil
}

func (o OrderedMap) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	err := e.EncodeToken(start)
	if err != nil {
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

	err = e.EncodeToken(xml.EndElement{Name: start.Name})
	if err != nil {
		return err
	}

	return e.Flush()
}

func (o *OrderedMap) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &o.Map)
	if err != nil {
		return err
	}

	index := make(map[string]int)
	for key := range o.Map {
		o.Order = append(o.Order, key)
		esc, _ := json.Marshal(key) // Escape the key
		index[key] = bytes.Index(b, esc)
	}

	if o.Lexical {
		sort.Slice(o.Order, func(i, j int) bool { return index[o.Order[i]] < index[o.Order[j]] })
	}
	return nil
}

func (o OrderedMap) MarshalJSON() ([]byte, error) {
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
