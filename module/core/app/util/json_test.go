package util_test

import (
	"encoding/json/jsontext"
	"reflect"
	"testing"

	"{{{ .Package }}}/app/util"
)

func TestToJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input any
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
		input any
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
		input []byte
		want  testStruct
	}{
		{"Valid", []byte(`{"A":1,"B":"test"}`), testStruct{1, "test"}},
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
		input []byte
		want  string
	}{
		{"Valid", []byte(`"test"`), "test"},
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

type jsonThing struct {
	A int `json:"a"`
}

func TestFromJSONAnyAndStrict(t *testing.T) {
	t.Parallel()

	msg := []byte(`"{\"a\":1}"`)
	res, err := util.FromJSONAny(msg)
	if err != nil {
		t.Fatalf("unexpected FromJSONAny error: %v", err)
	}
	m, ok := res.(map[string]any)
	if !ok || m["a"] != float64(1) {
		t.Fatalf("unexpected FromJSONAny result: %v", res)
	}

	okRes := util.FromJSONAnyOK([]byte(`{"a":2}`))
	if okMap, ok := okRes.(map[string]any); !ok || okMap["a"] != float64(2) {
		t.Fatalf("unexpected FromJSONAnyOK result: %v", okRes)
	}

	var tgt jsonThing
	if err := util.FromJSONStrict([]byte(`{"a":3}`), &tgt); err != nil || tgt.A != 3 {
		t.Fatalf("unexpected FromJSONStrict result: %v %v", tgt, err)
	}
	if err := util.FromJSONStrict([]byte(`{"a":3,"b":4}`), &tgt); err == nil {
		t.Fatalf("expected FromJSONStrict error for unknown field")
	}
}

func TestCycleJSONHelpers(t *testing.T) {
	t.Parallel()

	src := jsonThing{A: 5}
	var tgt jsonThing
	if err := util.CycleJSON(src, &tgt); err != nil || tgt.A != 5 {
		t.Fatalf("unexpected CycleJSON result: %v %v", tgt, err)
	}

	m := util.JSONToMap(src)
	if m["a"] != float64(5) {
		t.Fatalf("unexpected JSONToMap result: %v", m)
	}
}

func TestFromJSONBytesArrayAndMap(t *testing.T) {
	t.Parallel()

	vals := []jsontext.Value{jsontext.Value(`{"a":1}`), jsontext.Value(`{"a":2}`)}
	arr, err := util.FromJSONBytesArray[jsonThing](vals...)
	if err != nil || len(arr) != 2 || arr[1].A != 2 {
		t.Fatalf("unexpected FromJSONBytesArray result: %v %v", arr, err)
	}

	mp, err := util.FromJSONBytesMap[jsonThing](map[string]jsontext.Value{"x": jsontext.Value(`{"a":3}`)})
	if err != nil || mp["x"].A != 3 {
		t.Fatalf("unexpected FromJSONBytesMap result: %v %v", mp, err)
	}
}
