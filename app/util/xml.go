// Package util - Content managed by Project Forge, see [projectforge.md] for details.
package util

import (
	"bytes"
	"encoding/xml"
	"io"
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
	err := enc.Encode(x) //nolint:errchkxml // no chance of error
	if err != nil {
		return nil, err
	}
	return bts.Bytes(), nil
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

func FromXMLInterface(msg []byte) (any, error) {
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

func XMLToMap(i any) map[string]any {
	m := map[string]any{}
	_ = CycleXML(i, &m)
	return m
}
