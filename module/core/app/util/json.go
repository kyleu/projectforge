package util

import (
	"bytes"
	"encoding/json"
	"io"
)

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
	return json.Unmarshal(msg, tgt)
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
	dec := json.NewDecoder(bytes.NewReader(msg))
	dec.DisallowUnknownFields()
	return dec.Decode(tgt)
}

func CycleJSON(src interface{}, tgt interface{}) error {
	return FromJSON(ToJSONBytes(src, true), tgt)
}
