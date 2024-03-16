// Package util - Content managed by Project Forge, see [projectforge.md] for details.
package util

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"reflect"
	"strings"
)

func ToXML(x any) (string, error) {
	ret, err := ToXMLBytes(x, true)
	if err != nil {
		return "", err
	}
	return string(ret), nil
}

func ToXMLCompact(x any) (string, error) {
	ret, err := ToXMLBytes(x, false)
	if err != nil {
		return "", err
	}
	return string(ret), nil
}

func ToXMLBytes(x any, indent bool) ([]byte, error) {
	bts := &bytes.Buffer{}
	enc := xml.NewEncoder(bts)
	if indent {
		enc.Indent("", "  ")
	}
	if m, ok := x.(map[string]any); ok {
		x = ValueMap(m)
	}
	err := enc.Encode(x) //nolint:errchkxml // no chance of error
	if err != nil {
		return nil, err
	}
	b := bts.Bytes()
	if rf := reflect.ValueOf(x); rf.Kind() == reflect.Slice {
		n := fmt.Sprintf("%T", x)
		n = strings.NewReplacer("[", "", "]", "", "{", "", "}", "", " ", "").Replace(n)
		if idx := strings.LastIndex(n, "."); idx > -1 {
			n = n[idx+1:]
		}
		if n == "interface" || n == "any" {
			n = "array"
		}
		if indent {
			b = append([]byte(fmt.Sprintf("<%s>\n", n)), b...)
			b = append(b, []byte(fmt.Sprintf("\n</%s>", n))...)
		} else {
			b = append([]byte(fmt.Sprintf("<%s>", n)), b...)
			b = append(b, []byte(fmt.Sprintf("</%s>", n))...)
		}
	}
	return b, nil
}

func FromXML(msg []byte, tgt any) error {
	return xml.Unmarshal(msg, tgt)
}

func FromXMLString(msg []byte) (string, error) {
	var tgt string
	err := xml.Unmarshal(msg, &tgt)
	return tgt, err
}

func FromXMLMap(msg []byte) (ValueMap, error) {
	var tgt ValueMap
	err := xml.Unmarshal(msg, &tgt)
	return tgt, err
}

func FromXMLAny(msg []byte) (any, error) {
	var tgt any
	err := FromXML(msg, &tgt)
	return tgt, err
}

func FromXMLReader(r io.Reader, tgt any) error {
	return xml.NewDecoder(r).Decode(tgt)
}

func FromXMLStrict(msg []byte, tgt any) error {
	dec := xml.NewDecoder(bytes.NewReader(msg))
	dec.Strict = true
	return dec.Decode(tgt)
}

func CycleXML(src any, tgt any) error {
	b, _ := xml.Marshal(src)
	return FromXML(b, tgt)
}
