//go:build test_all || !func_test
// +build test_all !func_test

package util_test

import (
	"testing"

	"projectforge.dev/projectforge/app/util"
)

func TestRandomID(t *testing.T) {
	t.Parallel()

	t.Run("returns 16 character string", func(t *testing.T) {
		t.Parallel()
		id := util.RandomID()
		if len(id) != 16 {
			t.Errorf("expected length 16, got %d", len(id))
		}
	})

	t.Run("generates unique IDs", func(t *testing.T) {
		t.Parallel()
		id1 := util.RandomID()
		id2 := util.RandomID()
		if id1 == id2 {
			t.Error("expected different IDs")
		}
	})
}

func TestRandomString(t *testing.T) {
	t.Parallel()

	cases := []int{0, 1, 5, 10, 100}
	for _, length := range cases {
		l := length
		t.Run("", func(t *testing.T) {
			t.Parallel()
			s := util.RandomString(l)
			if len(s) != l {
				t.Errorf("expected length %d, got %d", l, len(s))
			}
		})
	}
}

func TestRandomInt(t *testing.T) {
	t.Parallel()

	t.Run("returns value in range", func(t *testing.T) {
		t.Parallel()
		for i := 0; i < 100; i++ {
			n := util.RandomInt(10)
			if n < 0 || n >= 10 {
				t.Errorf("expected 0 <= n < 10, got %d", n)
			}
		}
	})
}

func TestRandomInt16(t *testing.T) {
	t.Parallel()

	t.Run("returns value in range", func(t *testing.T) {
		t.Parallel()
		for i := 0; i < 100; i++ {
			n := util.RandomInt16(100)
			if n < 0 || n >= 100 {
				t.Errorf("expected 0 <= n < 100, got %d", n)
			}
		}
	})
}

func TestRandomInt32(t *testing.T) {
	t.Parallel()

	t.Run("returns value in range", func(t *testing.T) {
		t.Parallel()
		for i := 0; i < 100; i++ {
			n := util.RandomInt32(1000)
			if n < 0 || n >= 1000 {
				t.Errorf("expected 0 <= n < 1000, got %d", n)
			}
		}
	})
}

func TestRandomInt64(t *testing.T) {
	t.Parallel()

	t.Run("returns value in range", func(t *testing.T) {
		t.Parallel()
		for i := 0; i < 100; i++ {
			n := util.RandomInt64(10000)
			if n < 0 || n >= 10000 {
				t.Errorf("expected 0 <= n < 10000, got %d", n)
			}
		}
	})
}

func TestRandomFloat(t *testing.T) {
	t.Parallel()

	t.Run("returns value in range", func(t *testing.T) {
		t.Parallel()
		for i := 0; i < 100; i++ {
			f := util.RandomFloat(10)
			if f < 0 || f >= 10 {
				t.Errorf("expected 0 <= f < 10, got %f", f)
			}
		}
	})
}

func TestRandomBool(t *testing.T) {
	t.Parallel()

	t.Run("returns boolean", func(t *testing.T) {
		t.Parallel()
		trueCount := 0
		falseCount := 0
		for i := 0; i < 100; i++ {
			if util.RandomBool() {
				trueCount++
			} else {
				falseCount++
			}
		}
		if trueCount == 0 || falseCount == 0 {
			t.Error("expected mix of true and false values")
		}
	})
}

func TestRandomValueMap(t *testing.T) {
	t.Parallel()

	t.Run("returns map with specified keys", func(t *testing.T) {
		t.Parallel()
		m := util.RandomValueMap(5)
		if len(m) != 5 {
			t.Errorf("expected 5 keys, got %d", len(m))
		}
	})
}

func TestRandomBytes(t *testing.T) {
	t.Parallel()

	t.Run("returns bytes of specified size", func(t *testing.T) {
		t.Parallel()
		b := util.RandomBytes(32)
		if len(b) != 32 {
			t.Errorf("expected 32 bytes, got %d", len(b))
		}
	})

	t.Run("generates unique bytes", func(t *testing.T) {
		t.Parallel()
		b1 := util.RandomBytes(16)
		b2 := util.RandomBytes(16)
		same := true
		for i := range b1 {
			if b1[i] != b2[i] {
				same = false
				break
			}
		}
		if same {
			t.Error("expected different byte slices")
		}
	})
}

func TestRandomDate(t *testing.T) {
	t.Parallel()

	t.Run("returns future date", func(t *testing.T) {
		t.Parallel()
		d := util.RandomDate()
		now := util.TimeCurrent()
		if d.Before(now) {
			t.Error("expected date in the future")
		}
	})
}

func TestRandomURL(t *testing.T) {
	t.Parallel()

	t.Run("returns valid URL", func(t *testing.T) {
		t.Parallel()
		u := util.RandomURL()
		if u == nil {
			t.Fatal("expected non-nil URL")
		}
		if u.Scheme != "https" {
			t.Errorf("expected https scheme, got %s", u.Scheme)
		}
		if u.Host == "" {
			t.Error("expected non-empty host")
		}
	})
}

func TestRandomElement(t *testing.T) {
	t.Parallel()

	t.Run("returns element from slice", func(t *testing.T) {
		t.Parallel()
		slice := []string{"a", "b", "c"}
		elem := util.RandomElement(slice)
		found := false
		for _, s := range slice {
			if s == elem {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("element %s not in slice", elem)
		}
	})

	t.Run("returns zero value for empty slice", func(t *testing.T) {
		t.Parallel()
		var slice []string
		elem := util.RandomElement(slice)
		if elem != "" {
			t.Errorf("expected empty string, got %s", elem)
		}
	})
}

func TestRandomElements(t *testing.T) {
	t.Parallel()

	t.Run("returns specified number of elements", func(t *testing.T) {
		t.Parallel()
		slice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		elems := util.RandomElements(slice, 3)
		if len(elems) != 3 {
			t.Errorf("expected 3 elements, got %d", len(elems))
		}
	})

	t.Run("handles index larger than slice", func(t *testing.T) {
		t.Parallel()
		slice := []int{1, 2, 3}
		elems := util.RandomElements(slice, 10)
		if len(elems) != 3 {
			t.Errorf("expected 3 elements, got %d", len(elems))
		}
	})
}

func TestRandomDiffs(t *testing.T) {
	t.Parallel()

	t.Run("returns specified number of diffs", func(t *testing.T) {
		t.Parallel()
		diffs := util.RandomDiffs(5)
		if len(diffs) != 5 {
			t.Errorf("expected 5 diffs, got %d", len(diffs))
		}
	})

	t.Run("diffs have non-empty fields", func(t *testing.T) {
		t.Parallel()
		diffs := util.RandomDiffs(1)
		if diffs[0].Path == "" || diffs[0].Old == "" || diffs[0].New == "" {
			t.Error("expected non-empty diff fields")
		}
	})
}
