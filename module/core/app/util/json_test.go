package util_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"{{{ .Package }}}/app/util"
)

func TestToJSON(t *testing.T) {
	t.Parallel()
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
		x := tt
		t.Run(x.name, func(t *testing.T) {
			t.Parallel()
			if got := util.ToJSON(x.input); got != x.want {
				t.Errorf("ToJSON() = %v, want %v", got, x.want)
			}
		})
	}
}

func TestToJSONCompact(t *testing.T) {
	t.Parallel()
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
		x := tt
		t.Run(x.name, func(t *testing.T) {
			t.Parallel()
			if got := util.ToJSONCompact(x.input); got != x.want {
				t.Errorf("ToJSONCompact() = %v, want %v", got, x.want)
			}
		})
	}
}

func TestFromJSON(t *testing.T) {
	t.Parallel()
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
		x := tt
		t.Run(x.name, func(t *testing.T) {
			t.Parallel()
			var got testStruct
			if err := util.FromJSON(x.input, &got); err != nil {
				t.Errorf("FromJSON() error = %v", err)
			}
			if !reflect.DeepEqual(got, x.want) {
				t.Errorf("FromJSON() = %v, want %v", got, x.want)
			}
		})
	}
}

func TestFromJSONString(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		input json.RawMessage
		want  string
	}{
		{"Valid", json.RawMessage(`"test"`), "test"},
	}

	for _, tt := range tests {
		x := tt
		t.Run(x.name, func(t *testing.T) {
			t.Parallel()
			got, err := util.FromJSONString(x.input)
			if err != nil {
				t.Errorf("FromJSONString() error = %v", err)
			}
			if got != x.want {
				t.Errorf("FromJSONString() = %v, want %v", got, x.want)
			}
		})
	}
}
