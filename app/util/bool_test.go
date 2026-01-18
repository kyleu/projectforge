//go:build test_all || !func_test
// +build test_all !func_test

package util_test

import (
	"testing"

	"projectforge.dev/projectforge/app/util"
)

const boolTestHello = "hello"

func TestChoose(t *testing.T) {
	t.Parallel()

	t.Run("returns ifTrue when condition is true", func(t *testing.T) {
		t.Parallel()
		result := util.Choose(true, "yes", "no")
		if result != "yes" {
			t.Errorf("expected 'yes', got '%s'", result)
		}
	})

	t.Run("returns ifFalse when condition is false", func(t *testing.T) {
		t.Parallel()
		result := util.Choose(false, "yes", "no")
		if result != "no" {
			t.Errorf("expected 'no', got '%s'", result)
		}
	})

	t.Run("works with integers", func(t *testing.T) {
		t.Parallel()
		result := util.Choose(true, 42, 0)
		if result != 42 {
			t.Errorf("expected 42, got %d", result)
		}
	})

	t.Run("works with slices", func(t *testing.T) {
		t.Parallel()
		trueSlice := []int{1, 2, 3}
		falseSlice := []int{4, 5, 6}
		result := util.Choose(false, trueSlice, falseSlice)
		if len(result) != 3 || result[0] != 4 {
			t.Errorf("expected falseSlice, got %v", result)
		}
	})
}

func TestOrDefault(t *testing.T) {
	t.Parallel()

	t.Run("returns value when not zero", func(t *testing.T) {
		t.Parallel()
		result := util.OrDefault(boolTestHello, "default")
		if result != boolTestHello {
			t.Errorf("expected '%s', got '%s'", boolTestHello, result)
		}
	})

	t.Run("returns default when value is zero string", func(t *testing.T) {
		t.Parallel()
		result := util.OrDefault("", "default")
		if result != "default" {
			t.Errorf("expected 'default', got '%s'", result)
		}
	})

	t.Run("returns value when int is not zero", func(t *testing.T) {
		t.Parallel()
		result := util.OrDefault(42, 100)
		if result != 42 {
			t.Errorf("expected 42, got %d", result)
		}
	})

	t.Run("returns default when int is zero", func(t *testing.T) {
		t.Parallel()
		result := util.OrDefault(0, 100)
		if result != 100 {
			t.Errorf("expected 100, got %d", result)
		}
	})

	t.Run("works with pointers", func(t *testing.T) {
		t.Parallel()
		var nilPtr *int
		defaultVal := 42
		result := util.OrDefault(nilPtr, &defaultVal)
		if result != &defaultVal {
			t.Errorf("expected default pointer, got %v", result)
		}
	})
}
