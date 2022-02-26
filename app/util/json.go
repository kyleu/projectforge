// Content managed by Project Forge, see [projectforge.md] for details.
package util

import (
	"bytes"
	"encoding/json"
	"io"

	jsoniter "github.com/json-iterator/go"
)

var jsoniterParser = jsoniter.ConfigCompatibleWithStandardLibrary

func ToJSON(x interface{}) string {
	return string(ToJSONBytes(x, true))
}

func ToJSONCompact(x interface{}) string {
	return string(ToJSONBytes(x, false))
}

func ToJSONBytes(x interface{}, indent bool) []byte {
	if indent {
		b, _ := json.MarshalIndent(x, "", "  ")
		return b
	}
	b, _ := json.Marshal(x)
	return b
}

func FromJSON(msg json.RawMessage, tgt interface{}) error {
	return jsoniterParser.Unmarshal(msg, tgt)
}

func FromJSONInterface(msg json.RawMessage) (interface{}, error) {
	var tgt interface{}
	err := FromJSON(msg, &tgt)
	return tgt, err
}

func FromJSONReader(r io.Reader, tgt interface{}) error {
	return json.NewDecoder(r).Decode(tgt)
}

func FromJSONStrict(msg json.RawMessage, tgt interface{}) error {
	dec := jsoniterParser.NewDecoder(bytes.NewReader(msg))
	dec.DisallowUnknownFields()
	return dec.Decode(tgt)
}

func CycleJSON(src interface{}, tgt interface{}) error {
	b, _ := jsoniterParser.Marshal(src)
	return FromJSON(b, tgt)
}

func JSONToMap(i interface{}) map[string]interface{} {
	var m map[string]interface{}
	_ = CycleJSON(i, &m)
	return m
}
