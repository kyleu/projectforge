package util

import (
	"bytes"
	"encoding/json"
	"io"

	jsoniter "github.com/json-iterator/go"
)

var (
	jsoniterParser  = jsoniter.Config{EscapeHTML: false, SortMapKeys: true, ValidateJsonRawMessage: true}.Froze()
	trailingNewline = []byte{'\n'}
)

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
		jsonHandleError(enc.Encode(x))
		return bytes.TrimSuffix(bts.Bytes(), trailingNewline)
	}
	b, err := json.Marshal(x)
	jsonHandleError(err)
	return bytes.TrimSuffix(b, trailingNewline)
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

func FromJSONAny(msg json.RawMessage) (any, error) {
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

func FromJSONAnyOK(msg json.RawMessage) any {
	ret, _ := FromJSONAny(msg)
	return ret
}

func FromJSONObj[T any](msg json.RawMessage) (T, error) {
	var tgt T
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
	jsonHandleError(CycleJSON(i, &m))
	return m
}

func jsonHandleError(err error) {
	if err != nil && RootLogger != nil {
		RootLogger.Warnf("error encountered serializing JSON: %+v", err)
	}
}
