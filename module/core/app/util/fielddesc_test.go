//go:build test_all || !func_test
// +build test_all !func_test

package util_test

import (
	"testing"

	"{{{ .Package }}}/app/util"
)

func TestFieldDescParseAndCollect(t *testing.T) {
	d := &util.FieldDesc{Key: "flag", Type: "bool"}
	v, err := d.Parse("true")
	if err != nil || v != true {
		t.Fatalf("unexpected bool parse: %v %v", v, err)
	}

	d = &util.FieldDesc{Key: "num", Type: "int"}
	v, err = d.Parse("42")
	if err != nil || v != 42 {
		t.Fatalf("unexpected int parse: %v %v", v, err)
	}

	d = &util.FieldDesc{Key: "flt", Type: "float"}
	v, err = d.Parse("3.5")
	if err != nil || v != 3.5 {
		t.Fatalf("unexpected float parse: %v %v", v, err)
	}

	d = &util.FieldDesc{Key: "str", Type: "string"}
	v, err = d.Parse("hello")
	if err != nil || v != "hello" {
		t.Fatalf("unexpected string parse: %v %v", v, err)
	}

	d = &util.FieldDesc{Key: "arr", Type: "[]string"}
	v, err = d.Parse("a,b")
	arr, _ := v.([]string)
	if err != nil || len(arr) != 2 {
		t.Fatalf("unexpected array parse: %v %v", v, err)
	}

	if (&util.FieldDesc{Key: "x"}).TitleSafe() != "x" {
		t.Fatalf("expected TitleSafe to return key when title empty")
	}

	d = &util.FieldDesc{Key: "def", Default: "true"}
	if !d.DefaultBool() {
		t.Fatalf("expected DefaultBool true")
	}
	if (&util.FieldDesc{Key: "def", Default: "12"}).DefaultInt() != 12 {
		t.Fatalf("expected DefaultInt 12")
	}
	if (&util.FieldDesc{Key: "def", Default: "2.5"}).DefaultFloat() != 2.5 {
		t.Fatalf("expected DefaultFloat 2.5")
	}

	args := util.FieldDescs{
		{Key: "a", Default: "1"},
		{Key: "b"},
	}
	res := util.FieldDescsCollectMap(util.ValueMap{"b": "x"}, args)
	if !res.HasMissing() || len(res.Missing) != 1 || res.Missing[0] != "a" {
		t.Fatalf("expected missing field a, got %v", res.Missing)
	}
	if res.Values["a"] != "1" || res.Values["b"] != "x" {
		t.Fatalf("unexpected values: %v", res.Values)
	}
}
