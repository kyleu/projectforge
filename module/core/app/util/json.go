package util

import (
	"bytes"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"io"
)

var (
	jsonDefaultOpts = []json.Options{jsontext.EscapeForHTML(false), json.Deterministic(true)}
	jsonIndentOpts  = []json.Options{jsontext.EscapeForHTML(false), json.Deterministic(true), jsontext.WithIndent("  ")}
	trailingNewline = []byte{'\n'}
)

func ToJSON(x any) string {
	return string(ToJSONBytes(x, true))
}

func ToJSONCompact(x any) string {
	return string(ToJSONBytes(x, false))
}

func ToJSONBytes(x any, indent bool) []byte {
	opts := Choose(indent, jsonIndentOpts, jsonDefaultOpts)
	b, err := json.Marshal(x, opts...)
	jsonHandleError(x, err)
	return bytes.TrimSuffix(b, trailingNewline)
}

func FromJSON(msg []byte, tgt any) error {
	return FromJSONReader(bytes.NewReader(msg), tgt)
}

func FromJSONString(msg []byte) (string, error) {
	var tgt string
	err := FromJSON(msg, &tgt)
	return tgt, err
}

func FromJSONMap(msg []byte) (ValueMap, error) {
	var tgt ValueMap
	err := FromJSON(msg, &tgt)
	return tgt, err
}

func FromJSONOrderedMap[V any](msg []byte) (*OrderedMap[V], error) {
	var tgt *OrderedMap[V]
	err := FromJSON(msg, &tgt)
	return tgt, err
}

func FromJSONAny(msg []byte) (any, error) {
	if bytes.HasPrefix(msg, []byte("\"{")) {
		if str, err := FromJSONString(msg); err == nil {
			var tgt any
			if err = FromJSON([]byte(str), &tgt); err == nil {
				return tgt, nil
			}
		}
	}
	var tgt any
	err := FromJSON(msg, &tgt)
	return tgt, err
}

func FromJSONAnyOK(msg []byte) any {
	ret, _ := FromJSONAny(msg)
	return ret
}

func FromJSONObj[T any](msg []byte) (T, error) {
	var tgt T
	err := FromJSON(msg, &tgt)
	return tgt, err
}

func FromJSONStrict(msg []byte, tgt any) error {
	return FromJSONReader(bytes.NewReader(msg), tgt, json.RejectUnknownMembers(true))
}

func FromJSONReader(r io.Reader, tgt any, opts ...json.Options) error {
	return json.UnmarshalRead(r, tgt, opts...)
}

func CycleJSON(src any, tgt any) error {
	b := ToJSONBytes(src, false)
	return FromJSON(b, tgt)
}

func JSONToMap(i any) map[string]any {
	m := map[string]any{}
	jsonHandleError(i, CycleJSON(i, &m))
	return m
}

func jsonHandleError(src any, err error) {
	if err != nil && RootLogger != nil {
		RootLogger.Warnf("error [%s] encountered serializing JSON for type [%T]", err.Error(), src)
	}
}
