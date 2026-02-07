//go:build test_all || !func_test
// +build test_all !func_test

package util_test

import (
	"reflect"
	"strings"
	"testing"

	"{{{ .Package }}}/app/util"
)

const mapParseStr = "str"

func TestGetPathAndSetPath(t *testing.T) {
	t.Parallel()

	m := util.ValueMap{
		"files": util.ValueMap{
			"readme.txt": util.ValueMap{"size": 12},
		},
		"arr": []any{util.ValueMap{"k": "v"}, 42},
	}

	val, err := m.GetPath("files.\"readme.txt\".size", false)
	if err != nil || val != 12 {
		t.Fatalf("unexpected GetPath result: %v %v", val, err)
	}

	val, err = m.GetPath("arr.1", false)
	if err != nil || val != 42 {
		t.Fatalf("unexpected array path result: %v %v", val, err)
	}

	val, err = m.GetPath("missing", true)
	if err != nil || val != nil {
		t.Fatalf("expected nil for missing path, got %v %v", val, err)
	}

	if _, err = m.GetPath("missing", false); err == nil {
		t.Fatalf("expected error for missing path")
	}

	m2 := util.ValueMap{}
	if err := m2.SetPath("a.b.c", 7); err != nil {
		t.Fatalf("unexpected SetPath error: %v", err)
	}
	val, err = m2.GetPath("a.b.c", false)
	if err != nil || val != 7 {
		t.Fatalf("unexpected SetPath result: %v %v", val, err)
	}
}

func TestMapGetHelpers(t *testing.T) {
	t.Parallel()

	m := util.ValueMap{
		"s":    mapParseStr,
		"i":    42,
		"f":    1.5,
		"b":    true,
		"arr":  []any{"a", "b"},
		"json": []byte(`{"a":1}`),
	}

	if v, err := m.GetRequired("s"); err != nil || v != mapParseStr {
		t.Fatalf("unexpected GetRequired result: %v %v", v, err)
	}

	if _, err := m.GetRequired("missing"); err == nil || !strings.Contains(err.Error(), "candidates") {
		t.Fatalf("expected error for missing required key")
	}

	if v, err := m.GetString("s", false); err != nil || v != mapParseStr {
		t.Fatalf("unexpected GetString result: %v %v", v, err)
	}

	if v, err := m.GetInt("i", false); err != nil || v != 42 {
		t.Fatalf("unexpected GetInt result: %v %v", v, err)
	}

	if v, err := m.GetFloat("f", false); err != nil || v != 1.5 {
		t.Fatalf("unexpected GetFloat result: %v %v", v, err)
	}

	if v, err := m.GetBool("b", false); err != nil || v != true {
		t.Fatalf("unexpected GetBool result: %v %v", v, err)
	}

	if v, err := m.GetArray("arr", false); err != nil || len(v) != 2 {
		t.Fatalf("unexpected GetArray result: %v %v", v, err)
	}

	var out map[string]any
	if err := m.GetType("json", &out); err != nil || out["a"] != float64(1) {
		t.Fatalf("unexpected GetType result: %v %v", out, err)
	}

	if err := m.GetType("s", &out); err == nil || !strings.Contains(err.Error(), "expected binary json") {
		t.Fatalf("expected GetType error for non-bytes")
	}

	if v := util.MapGetOrElse(map[string]int{"a": 1}, "b", 2); v != 2 {
		t.Fatalf("unexpected MapGetOrElse result: %v", v)
	}

	m3 := util.ValueMap{"x": "1"}
	parsed, err := m3.ParseInt("x", false, false)
	if err != nil || parsed != 1 {
		t.Fatalf("unexpected ParseInt result: %v %v", parsed, err)
	}

	if _, err = m3.ParseInt("missing", false, false); err == nil {
		t.Fatalf("expected error for missing ParseInt")
	}

	opt := m3.GetStringOpt("missing")
	if opt != "" {
		t.Fatalf("expected empty opt string, got %q", opt)
	}

	if !reflect.DeepEqual(m3.GetStringArrayOr("missing", []string{"a"}, true), []string{"a"}) {
		t.Fatalf("unexpected GetStringArrayOr result")
	}
}
