// Content managed by Project Forge, see [projectforge.md] for details.
package util

import (
	"bytes"
	"encoding/json"
	"io"

	jsoniter "github.com/json-iterator/go"
)

var jsoniterParser = jsoniter.Config{EscapeHTML: false, SortMapKeys: true, ValidateJsonRawMessage: true}.Froze()

func ToJSON(x any) string {
	return string(ToJSONBytes(x, true))
}

func ToJSONCompact(x any) string {
	return string(ToJSONBytes(x, false))
}

func ToJSONBytes(x any, indent bool) []byte {
	if indent {
		bts := &bytes.Buffer{}
		enc := json.NewEncoder(bts)
		if indent {
			enc.SetIndent("", "  ")
		}
		enc.SetEscapeHTML(false)
		_ = enc.Encode(x) //nolint:errchkjson // no chance of error
		return bts.Bytes()
	}
	b, _ := json.Marshal(x) //nolint:errchkjson // no chance of error
	return b
}

func FromJSON(msg json.RawMessage, tgt any) error {
	return jsoniterParser.Unmarshal(msg, tgt)
}

func FromJSONString(msg json.RawMessage) (string, error) {
	var tgt string
	err := jsoniterParser.Unmarshal(msg, &tgt)
	return tgt, err
}

func FromJSONMap(msg json.RawMessage) (ValueMap, error) {
	var tgt ValueMap
	err := jsoniterParser.Unmarshal(msg, &tgt)
	return tgt, err
}

func FromJSONInterface(msg json.RawMessage) (any, error) {
	var tgt any
	err := FromJSON(msg, &tgt)
	return tgt, err
}

func FromJSONReader(r io.Reader, tgt any) error {
	return json.NewDecoder(r).Decode(tgt)
}

func FromJSONStrict(msg json.RawMessage, tgt any) error {
	dec := jsoniterParser.NewDecoder(bytes.NewReader(msg))
	dec.DisallowUnknownFields()
	return dec.Decode(tgt)
}

func CycleJSON(src any, tgt any) error {
	b, _ := jsoniterParser.Marshal(src)
	return FromJSON(b, tgt)
}

func JSONToMap(i any) map[string]any {
	m := map[string]any{}
	_ = CycleJSON(i, &m)
	return m
}
