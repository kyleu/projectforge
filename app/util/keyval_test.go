//go:build test_all || !func_test
// +build test_all !func_test

package util_test

import (
	"testing"

	"projectforge.dev/projectforge/app/util"
)

func TestKeyVal_String(t *testing.T) {
	t.Parallel()

	t.Run("formats correctly", func(t *testing.T) {
		t.Parallel()
		kv := util.KeyVal[string]{Key: "name", Val: "value"}
		result := kv.String()
		if result != "name: value" {
			t.Errorf("expected 'name: value', got '%s'", result)
		}
	})

	t.Run("formats int value", func(t *testing.T) {
		t.Parallel()
		kv := util.KeyVal[int]{Key: "count", Val: 42}
		result := kv.String()
		if result != "count: 42" {
			t.Errorf("expected 'count: 42', got '%s'", result)
		}
	})
}

func TestKeyVals_ToMap(t *testing.T) {
	t.Parallel()

	t.Run("converts to map", func(t *testing.T) {
		t.Parallel()
		kvs := util.KeyVals[string]{
			{Key: "a", Val: "1"},
			{Key: "b", Val: "2"},
		}
		m := kvs.ToMap()
		if m["a"] != "1" || m["b"] != "2" {
			t.Errorf("unexpected map: %v", m)
		}
	})
}

func TestKeyVals_String(t *testing.T) {
	t.Parallel()

	t.Run("joins with commas", func(t *testing.T) {
		t.Parallel()
		kvs := util.KeyVals[string]{
			{Key: "a", Val: "1"},
			{Key: "b", Val: "2"},
		}
		result := kvs.String()
		if result != "a: 1, b: 2" {
			t.Errorf("expected 'a: 1, b: 2', got '%s'", result)
		}
	})
}

func TestKeyVals_Values(t *testing.T) {
	t.Parallel()

	t.Run("returns values", func(t *testing.T) {
		t.Parallel()
		kvs := util.KeyVals[int]{
			{Key: "a", Val: 1},
			{Key: "b", Val: 2},
		}
		vals := kvs.Values()
		if len(vals) != 2 || vals[0] != 1 || vals[1] != 2 {
			t.Errorf("unexpected values: %v", vals)
		}
	})
}

func TestKeyTypeDesc_String(t *testing.T) {
	t.Parallel()

	t.Run("formats correctly", func(t *testing.T) {
		t.Parallel()
		ktd := util.KeyTypeDesc{Key: "name", Type: "string"}
		result := ktd.String()
		if result != "name (string)" {
			t.Errorf("expected 'name (string)', got '%s'", result)
		}
	})
}

func TestKeyTypeDesc_Array(t *testing.T) {
	t.Parallel()

	t.Run("returns array with replaced key", func(t *testing.T) {
		t.Parallel()
		ktd := util.KeyTypeDesc{Key: "{key}.field", Type: "string", Description: "Field for {key}"}
		result := ktd.Array("user")
		if result[0] != "`user.field`" || result[1] != "string" || result[2] != "Field for user" {
			t.Errorf("unexpected array: %v", result)
		}
	})
}

func TestKeyTypeDesc_Matches(t *testing.T) {
	t.Parallel()

	t.Run("matches same key", func(t *testing.T) {
		t.Parallel()
		ktd1 := &util.KeyTypeDesc{Key: "name"}
		ktd2 := &util.KeyTypeDesc{Key: "name"}
		if !ktd1.Matches(ktd2) {
			t.Error("expected match")
		}
	})

	t.Run("does not match different key", func(t *testing.T) {
		t.Parallel()
		ktd1 := &util.KeyTypeDesc{Key: "name"}
		ktd2 := &util.KeyTypeDesc{Key: "other"}
		if ktd1.Matches(ktd2) {
			t.Error("expected no match")
		}
	})
}

func TestKeyTypeDescs_Strings(t *testing.T) {
	t.Parallel()

	t.Run("returns strings", func(t *testing.T) {
		t.Parallel()
		ktds := util.KeyTypeDescs{
			{Key: "a", Type: "int"},
			{Key: "b", Type: "string"},
		}
		strs := ktds.Strings()
		if len(strs) != 2 || strs[0] != "a (int)" || strs[1] != "b (string)" {
			t.Errorf("unexpected strings: %v", strs)
		}
	})
}

func TestKeyTypeDescs_Sort(t *testing.T) {
	t.Parallel()

	t.Run("sorts by key", func(t *testing.T) {
		t.Parallel()
		ktds := util.KeyTypeDescs{
			{Key: "zebra", Type: "string"},
			{Key: "apple", Type: "string"},
			{Key: "Banana", Type: "string"},
		}
		sorted := ktds.Sort()
		if sorted[0].Key != "apple" || sorted[1].Key != "Banana" || sorted[2].Key != "zebra" {
			t.Errorf("unexpected sort order: %v", sorted)
		}
	})
}

func TestKeyTypeDescs_Array(t *testing.T) {
	t.Parallel()

	t.Run("returns sorted array", func(t *testing.T) {
		t.Parallel()
		ktds := util.KeyTypeDescs{
			{Key: "z", Type: "int", Description: "Z"},
			{Key: "a", Type: "string", Description: "A"},
		}
		result := ktds.Array("test")
		if len(result) != 2 || result[0][0] != "`a`" || result[1][0] != "`z`" {
			t.Errorf("unexpected result: %v", result)
		}
	})
}
