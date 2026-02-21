//go:build test_all || !func_test
// +build test_all !func_test

package util_test

import (
	"testing"

	"{{{ .Package }}}/app/util"
)

const helloWorld = "Hello World"

func TestNewStringSlice(t *testing.T) {
	t.Parallel()

	t.Run("creates empty slice", func(t *testing.T) {
		t.Parallel()
		s := util.NewStringSlice()
		if s == nil {
			t.Fatal("expected non-nil StringSlice")
		}
		if !s.Empty() {
			t.Error("expected empty slice")
		}
	})

	t.Run("creates slice with initial values", func(t *testing.T) {
		t.Parallel()
		s := util.NewStringSlice("a", "b", "c")
		if s.Length() != 3 {
			t.Errorf("expected length 3, got %d", s.Length())
		}
	})
}

func TestNewStringSliceWithSize(t *testing.T) {
	t.Parallel()

	t.Run("creates empty slice with capacity", func(t *testing.T) {
		t.Parallel()
		s := util.NewStringSliceWithSize(10)
		if s == nil {
			t.Fatal("expected non-nil StringSlice")
		}
		if !s.Empty() {
			t.Error("expected empty slice")
		}
		if cap(s.Slice) != 10 {
			t.Errorf("expected capacity 10, got %d", cap(s.Slice))
		}
	})
}

func TestStringSliceEmpty(t *testing.T) {
	t.Parallel()

	t.Run("returns true for empty slice", func(t *testing.T) {
		t.Parallel()
		s := util.NewStringSlice()
		if !s.Empty() {
			t.Error("expected true for empty slice")
		}
	})

	t.Run("returns false for non-empty slice", func(t *testing.T) {
		t.Parallel()
		s := util.NewStringSlice("a")
		if s.Empty() {
			t.Error("expected false for non-empty slice")
		}
	})
}

func TestStringSliceLength(t *testing.T) {
	t.Parallel()

	t.Run("returns correct length", func(t *testing.T) {
		t.Parallel()
		s := util.NewStringSlice("a", "b", "c")
		if s.Length() != 3 {
			t.Errorf("expected 3, got %d", s.Length())
		}
	})
}

func TestStringSliceSliceSafe(t *testing.T) {
	t.Parallel()

	t.Run("returns nil for nil receiver", func(t *testing.T) {
		t.Parallel()
		var s *util.StringSlice
		if s.SliceSafe() != nil {
			t.Error("expected nil")
		}
	})

	t.Run("returns nil for empty slice", func(t *testing.T) {
		t.Parallel()
		s := util.NewStringSlice()
		if s.SliceSafe() != nil {
			t.Error("expected nil for empty slice")
		}
	})

	t.Run("returns slice for non-empty", func(t *testing.T) {
		t.Parallel()
		s := util.NewStringSlice("a", "b")
		result := s.SliceSafe()
		if len(result) != 2 {
			t.Errorf("expected length 2, got %d", len(result))
		}
	})
}

func TestStringSliceTotalLength(t *testing.T) {
	t.Parallel()

	t.Run("returns sum of string lengths", func(t *testing.T) {
		t.Parallel()
		s := util.NewStringSlice("abc", "de", "f")
		if s.TotalLength() != 6 {
			t.Errorf("expected 6, got %d", s.TotalLength())
		}
	})

	t.Run("returns 0 for empty slice", func(t *testing.T) {
		t.Parallel()
		s := util.NewStringSlice()
		if s.TotalLength() != 0 {
			t.Errorf("expected 0, got %d", s.TotalLength())
		}
	})
}

func TestStringSlicePush(t *testing.T) {
	t.Parallel()

	t.Run("adds strings to slice", func(t *testing.T) {
		t.Parallel()
		s := util.NewStringSlice()
		s.Push("a", "b")
		if s.Length() != 2 {
			t.Errorf("expected 2, got %d", s.Length())
		}
	})
}

func TestStringSlicePushUnique(t *testing.T) {
	t.Parallel()

	t.Run("adds only unique strings", func(t *testing.T) {
		t.Parallel()
		s := util.NewStringSlice("a", "b")
		s.PushUnique("b", "c")
		if s.Length() != 3 {
			t.Errorf("expected 3, got %d", s.Length())
		}
	})
}

func TestStringSlicePushf(t *testing.T) {
	t.Parallel()

	t.Run("adds formatted string", func(t *testing.T) {
		t.Parallel()
		s := util.NewStringSlice()
		s.Pushf("Hello %s", "World")
		if s.Slice[0] != helloWorld {
			t.Errorf("expected '%s', got %s", helloWorld, s.Slice[0])
		}
	})
}

func TestStringSlicePushfUnlessNil(t *testing.T) {
	t.Parallel()

	t.Run("does nothing for nil receiver", func(t *testing.T) {
		t.Parallel()
		var s *util.StringSlice
		s.PushfUnlessNil("test %s", "value")
	})

	t.Run("adds formatted string for non-nil", func(t *testing.T) {
		t.Parallel()
		s := util.NewStringSlice()
		s.PushfUnlessNil("Hello %s", "World")
		if s.Slice[0] != helloWorld {
			t.Errorf("expected '%s', got %s", helloWorld, s.Slice[0])
		}
	})
}

func TestStringSliceString(t *testing.T) {
	t.Parallel()

	t.Run("joins with newlines", func(t *testing.T) {
		t.Parallel()
		s := util.NewStringSlice("a", "b", "c")
		result := s.String()
		if result != "a\nb\nc" {
			t.Errorf("expected 'a\\nb\\nc', got %s", result)
		}
	})
}

func TestStringSliceJoin(t *testing.T) {
	t.Parallel()

	t.Run("joins with custom separator", func(t *testing.T) {
		t.Parallel()
		s := util.NewStringSlice("a", "b", "c")
		result := s.Join("-")
		if result != "a-b-c" {
			t.Errorf("expected 'a-b-c', got %s", result)
		}
	})
}

func TestStringSliceJoinSimple(t *testing.T) {
	t.Parallel()

	t.Run("joins without separator", func(t *testing.T) {
		t.Parallel()
		s := util.NewStringSlice("a", "b", "c")
		result := s.JoinSimple()
		if result != "abc" {
			t.Errorf("expected 'abc', got %s", result)
		}
	})
}

func TestStringSliceJoinSpace(t *testing.T) {
	t.Parallel()

	t.Run("joins with spaces", func(t *testing.T) {
		t.Parallel()
		s := util.NewStringSlice("a", "b", "c")
		result := s.JoinSpace()
		if result != "a b c" {
			t.Errorf("expected 'a b c', got %s", result)
		}
	})
}

func TestStringSliceJoinCommas(t *testing.T) {
	t.Parallel()

	t.Run("joins with commas", func(t *testing.T) {
		t.Parallel()
		s := util.NewStringSlice("a", "b", "c")
		result := s.JoinCommas()
		if result != "a, b, c" {
			t.Errorf("expected 'a, b, c', got %s", result)
		}
	})
}

func TestStringSliceSort(t *testing.T) {
	t.Parallel()

	t.Run("sorts alphabetically", func(t *testing.T) {
		t.Parallel()
		s := util.NewStringSlice("c", "a", "b")
		s.Sort()
		if s.Slice[0] != "a" || s.Slice[1] != "b" || s.Slice[2] != "c" {
			t.Errorf("expected sorted slice, got %v", s.Slice)
		}
	})
}

func TestStringSliceJSON(t *testing.T) {
	t.Parallel()

	t.Run("marshals to JSON array", func(t *testing.T) {
		t.Parallel()
		s := util.NewStringSlice("a", "b", "c")
		data := util.ToJSON(s)
		expected := `["a","b","c"]`
		if data != expected {
			t.Errorf("expected %s, got %s", expected, data)
		}
	})

	t.Run("unmarshals from JSON array", func(t *testing.T) {
		t.Parallel()
		s := &util.StringSlice{}
		err := util.FromJSON([]byte(`["x","y","z"]`), s)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if s.Length() != 3 || s.Slice[0] != "x" {
			t.Errorf("unexpected result: %v", s.Slice)
		}
	})
}
