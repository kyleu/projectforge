package util

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestToJSON(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		want  string
	}{
		{"Simple", map[string]int{"a": 1, "b": 2}, "{\n  \"a\": 1,\n  \"b\": 2\n}"},
		{"Complex", struct {
			A int
			B string
		}{1, "test"}, "{\n  \"A\": 1,\n  \"B\": \"test\"\n}"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToJSON(tt.input); got != tt.want {
				t.Errorf("ToJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToJSONCompact(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		want  string
	}{
		{"Simple", map[string]int{"a": 1, "b": 2}, "{\"a\":1,\"b\":2}"},
		{"Complex", struct {
			A int
			B string
		}{1, "test"}, "{\"A\":1,\"B\":\"test\"}"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToJSONCompact(tt.input); got != tt.want {
				t.Errorf("ToJSONCompact() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFromJSON(t *testing.T) {
	type testStruct struct {
		A int
		B string
	}

	tests := []struct {
		name  string
		input json.RawMessage
		want  testStruct
	}{
		{"Valid", json.RawMessage(`{"A":1,"B":"test"}`), testStruct{1, "test"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got testStruct
			if err := FromJSON(tt.input, &got); err != nil {
				t.Errorf("FromJSON() error = %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFromJSONString(t *testing.T) {
	tests := []struct {
		name  string
		input json.RawMessage
		want  string
	}{
		{"Valid", json.RawMessage(`"test"`), "test"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FromJSONString(tt.input)
			if err != nil {
				t.Errorf("FromJSONString() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("FromJSONString() = %v, want %v", got, tt.want)
			}
		})
	}
}
