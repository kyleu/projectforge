//go:build test_all || !func_test
// +build test_all !func_test

package util_test

import (
	"testing"

	"projectforge.dev/projectforge/app/util"
)

const testHello = "hello"

func TestCast(t *testing.T) {
	t.Parallel()

	t.Run("casts string successfully", func(t *testing.T) {
		t.Parallel()
		var v any = testHello
		result, err := util.Cast[string](v)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result != testHello {
			t.Errorf("expected '%s', got '%s'", testHello, result)
		}
	})

	t.Run("casts int successfully", func(t *testing.T) {
		t.Parallel()
		var v any = 42
		result, err := util.Cast[int](v)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result != 42 {
			t.Errorf("expected 42, got %d", result)
		}
	})

	t.Run("returns error for invalid cast", func(t *testing.T) {
		t.Parallel()
		var v any = testHello
		_, err := util.Cast[int](v)
		if err == nil {
			t.Error("expected error for invalid cast")
		}
	})

	t.Run("casts struct successfully", func(t *testing.T) {
		t.Parallel()
		type MyStruct struct {
			Name string
		}
		var v any = MyStruct{Name: "test"}
		result, err := util.Cast[MyStruct](v)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Name != "test" {
			t.Errorf("expected 'test', got '%s'", result.Name)
		}
	})
}

func TestCastOK(t *testing.T) {
	t.Parallel()

	t.Run("returns value on successful cast", func(t *testing.T) {
		t.Parallel()
		var v any = testHello
		result := util.CastOK[string](v)
		if result != testHello {
			t.Errorf("expected '%s', got '%s'", testHello, result)
		}
	})

	t.Run("returns zero value on failed cast", func(t *testing.T) {
		t.Parallel()
		var v any = testHello
		result := util.CastOK[int](v)
		if result != 0 {
			t.Errorf("expected 0, got %d", result)
		}
	})

	t.Run("returns zero value for nil", func(t *testing.T) {
		t.Parallel()
		var v any
		result := util.CastOK[string](v)
		if result != "" {
			t.Errorf("expected empty string, got '%s'", result)
		}
	})
}
