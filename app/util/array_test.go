package util_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/util"
)

func TestStringArrayMaxLength(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		input    []string
		expected int
	}{
		{name: "empty array", input: []string{}, expected: 0},
		{name: "single element", input: []string{"hello"}, expected: 5},
		{name: "multiple elements with same length", input: []string{"abc", "def", "ghi"}, expected: 3},
		{name: "multiple elements with different lengths", input: []string{"a", "abc", "abcde"}, expected: 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := util.StringArrayMaxLength(tt.input)
			if result != tt.expected {
				t.Errorf("StringArrayMaxLength() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestArrayToStringArray(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		input    []any
		expected []string
	}{
		{name: "empty array", input: []any{}, expected: []string{}},
		{name: "integers", input: []any{1, 2, 3}, expected: []string{"1", "2", "3"}},
		{name: "mixed types", input: []any{1, "hello", true, 3.14}, expected: []string{"1", "hello", "true", "3.14"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			anySlice := util.ArrayCopy(tt.input)
			result := util.ArrayToStringArray(anySlice)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ArrayToStringArray() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestStringArrayQuoted(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{name: "empty array", input: []string{}, expected: []string{}},
		{name: "simple strings", input: []string{"hello", "world"}, expected: []string{`"hello"`, `"world"`}},
		{name: "strings with quotes", input: []string{`a"b`, `c'd`}, expected: []string{`"a\"b"`, `"c'd"`}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := util.StringArrayQuoted(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("StringArrayQuoted() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestStringArrayFromAny(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		input     []any
		maxLength int
		expected  []string
	}{
		{name: "empty array", input: []any{}, maxLength: 0, expected: []string{}},
		{name: "mixed types no truncation", input: []any{1, "hello", true, []byte("world")}, maxLength: 0, expected: []string{"1", "hello", "true", "world"}},
		{name: "with truncation", input: []any{"abcdefghij", "12345"}, maxLength: 5, expected: []string{"abcde... (truncated)", "12345"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := util.StringArrayFromAny(tt.input, tt.maxLength)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("StringArrayFromAny() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestArrayCopy(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{name: "nil array", input: nil, expected: nil},
		{name: "empty array", input: []int{}, expected: []int{}},
		{name: "integers", input: []int{1, 2, 3}, expected: []int{1, 2, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := util.ArrayCopy(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ArrayCopy() = %v, want %v", result, tt.expected)
			}

			if len(tt.input) > 0 {
				original := make([]int, len(tt.input))
				copy(original, tt.input)
				tt.input[0] = 999
				if len(result) > 0 && reflect.DeepEqual(result, tt.input) {
					t.Errorf("ArrayCopy() did not create a true copy")
				}
				copy(tt.input, original)
			}
		})
	}
}

func TestArrayRemoveDuplicates(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{name: "empty array", input: []int{}, expected: []int{}},
		{name: "no duplicates", input: []int{1, 2, 3}, expected: []int{1, 2, 3}},
		{name: "with duplicates", input: []int{1, 2, 2, 3, 1, 4}, expected: []int{1, 2, 3, 4}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := util.ArrayRemoveDuplicates(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("ArrayRemoveDuplicates() length = %v, want %v", len(result), len(tt.expected))
			}

			resultMap := make(map[int]bool)
			for _, v := range result {
				resultMap[v] = true
			}

			expectedMap := make(map[int]bool)
			for _, v := range tt.expected {
				expectedMap[v] = true
			}

			for k := range expectedMap {
				if !resultMap[k] {
					t.Errorf("ArrayRemoveDuplicates() missing expected element %v", k)
				}
			}

			for k := range resultMap {
				if !expectedMap[k] {
					t.Errorf("ArrayRemoveDuplicates() has unexpected element %v", k)
				}
			}
		})
	}
}

func TestArraySorted(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{name: "empty array", input: []int{}, expected: []int{}},
		{name: "already sorted", input: []int{1, 2, 3}, expected: []int{1, 2, 3}},
		{name: "unsorted", input: []int{3, 1, 2}, expected: []int{1, 2, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := util.ArraySorted(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ArraySorted() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestArrayLimit(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		input         []int
		limit         int
		expectedArray []int
		expectedRest  int
	}{
		{name: "empty array", input: []int{}, limit: 5, expectedArray: []int{}, expectedRest: 0},
		{name: "limit larger than array", input: []int{1, 2, 3}, limit: 5, expectedArray: []int{1, 2, 3}, expectedRest: 0},
		{name: "limit smaller than array", input: []int{1, 2, 3, 4, 5}, limit: 3, expectedArray: []int{1, 2, 3}, expectedRest: 2},
		{name: "limit zero", input: []int{1, 2, 3}, limit: 0, expectedArray: []int{1, 2, 3}, expectedRest: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result, rest := util.ArrayLimit(tt.input, tt.limit)
			if !reflect.DeepEqual(result, tt.expectedArray) {
				t.Errorf("ArrayLimit() array = %v, want %v", result, tt.expectedArray)
			}
			if rest != tt.expectedRest {
				t.Errorf("ArrayLimit() rest = %v, want %v", rest, tt.expectedRest)
			}
		})
	}
}

func TestStringArrayOxfordComma(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		input     []string
		separator string
		expected  string
	}{
		{name: "empty array", input: []string{}, separator: "and", expected: ""},
		{name: "single element", input: []string{"apple"}, separator: "and", expected: "apple"},
		{name: "two elements", input: []string{"apple", "banana"}, separator: "and", expected: "apple and banana"},
		{name: "three elements", input: []string{"apple", "banana", "cherry"}, separator: "and", expected: "apple, banana, and cherry"},
		{name: "four elements", input: []string{"apple", "banana", "cherry", "date"}, separator: "or", expected: "apple, banana, cherry, or date"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := util.StringArrayOxfordComma(tt.input, tt.separator)
			if result != tt.expected {
				t.Errorf("StringArrayOxfordComma() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestArrayTransform(t *testing.T) {
	tx := func(i int) string {
		return "x" + string(rune(i+'0'))
	}
	t.Parallel()
	tests := []struct {
		name     string
		input    []int
		fn       func(int) string
		expected []string
	}{
		{name: "empty array", input: []int{}, fn: tx, expected: []string{}},
		{name: "transform integers to strings", input: []int{1, 2, 3}, fn: tx, expected: []string{"x1", "x2", "x3"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := util.ArrayTransform(tt.input, tt.fn)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ArrayTransform() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestArraySplit(t *testing.T) {
	fn := func(i int) bool {
		return i%2 == 0
	}
	t.Parallel()
	tests := []struct {
		name          string
		input         []int
		fn            func(int) bool
		expectedTrue  []int
		expectedFalse []int
	}{
		{name: "empty array", input: []int{}, fn: fn, expectedTrue: []int{}, expectedFalse: []int{}},
		{name: "split even/odd", input: []int{1, 2, 3, 4, 5}, fn: fn, expectedTrue: []int{2, 4}, expectedFalse: []int{1, 3, 5}},
		{name: "all true", input: []int{2, 4, 6}, fn: fn, expectedTrue: []int{2, 4, 6}, expectedFalse: []int{}},
		{name: "all false", input: []int{1, 3, 5}, fn: fn, expectedTrue: []int{}, expectedFalse: []int{1, 3, 5}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			trueResult, falseResult := util.ArraySplit(tt.input, tt.fn)

			if len(trueResult) != len(tt.expectedTrue) {
				t.Errorf("ArraySplit() true result length = %v, want %v", len(trueResult), len(tt.expectedTrue))
			} else {
				for i, v := range trueResult {
					if i < len(tt.expectedTrue) && v != tt.expectedTrue[i] {
						t.Errorf("ArraySplit() true result[%d] = %v, want %v", i, v, tt.expectedTrue[i])
					}
				}
			}

			if len(falseResult) != len(tt.expectedFalse) {
				t.Errorf("ArraySplit() false result length = %v, want %v", len(falseResult), len(tt.expectedFalse))
			} else {
				for i, v := range falseResult {
					if i < len(tt.expectedFalse) && v != tt.expectedFalse[i] {
						t.Errorf("ArraySplit() false result[%d] = %v, want %v", i, v, tt.expectedFalse[i])
					}
				}
			}
		})
	}
}

func TestArrayRemoveNil(t *testing.T) {
	t.Parallel()
	i1, i2, i3 := 1, 2, 3
	tests := []struct {
		name     string
		input    []*int
		expected []*int
	}{
		{name: "empty array", input: []*int{}, expected: []*int{}},
		{name: "no nil elements", input: []*int{&i1, &i2, &i3}, expected: []*int{&i1, &i2, &i3}},
		{name: "some nil elements", input: []*int{&i1, nil, &i2, nil, &i3}, expected: []*int{&i1, &i2, &i3}},
		{name: "all nil elements", input: []*int{nil, nil, nil}, expected: []*int{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := util.ArrayRemoveNil(tt.input)

			if len(result) != len(tt.expected) {
				t.Errorf("ArrayRemoveNil() length = %v, want %v", len(result), len(tt.expected))
				return
			}

			for i, v := range result {
				switch {
				case v == nil:
					t.Errorf("ArrayRemoveNil() contains nil at index %d", i)
				case tt.expected[i] == nil:
					t.Errorf("ArrayRemoveNil() expected nil at index %d, got %v", i, *v)
				case *v != *tt.expected[i]:
					t.Errorf("ArrayRemoveNil() at index %d = %v, want %v", i, *v, *tt.expected[i])
				}
			}
		})
	}
}

func TestArrayRemoveEmpty(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{name: "empty array", input: []string{}, expected: []string{}},
		{name: "no empty elements", input: []string{"a", "b", "c"}, expected: []string{"a", "b", "c"}},
		{name: "some empty elements", input: []string{"a", "", "b", "", "c"}, expected: []string{"a", "b", "c"}},
		{name: "all empty elements", input: []string{"", "", ""}, expected: []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := util.ArrayRemoveEmpty(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ArrayRemoveEmpty() = %v, want %v", result, tt.expected)
			}
		})
	}

	intTests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{name: "empty array", input: []int{}, expected: []int{}},
		{name: "no zero elements", input: []int{1, 2, 3}, expected: []int{1, 2, 3}},
		{name: "some zero elements", input: []int{1, 0, 2, 0, 3}, expected: []int{1, 2, 3}},
		{name: "all zero elements", input: []int{0, 0, 0}, expected: []int{}},
	}

	for _, tt := range intTests {
		t.Run(tt.name+" (int)", func(t *testing.T) {
			t.Parallel()
			result := util.ArrayRemoveEmpty(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ArrayRemoveEmpty() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestArrayDereference(t *testing.T) {
	t.Parallel()
	i1, i2, i3 := 1, 2, 3
	tests := []struct {
		name     string
		input    []*int
		expected []int
	}{
		{name: "empty array", input: []*int{}, expected: []int{}},
		{name: "simple array", input: []*int{&i1, &i2, &i3}, expected: []int{1, 2, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := util.ArrayDereference(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ArrayDereference() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestLengthAny(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		input    any
		expected int
	}{
		{name: "empty slice", input: []int{}, expected: 0},
		{name: "int slice", input: []int{1, 2, 3}, expected: 3},
		{name: "string slice", input: []string{"a", "b", "c"}, expected: 3},
		{name: "pointer to slice", input: &[]int{1, 2, 3, 4}, expected: 4},
		{name: "array", input: [3]int{1, 2, 3}, expected: 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := util.LengthAny(tt.input)
			if result != tt.expected {
				t.Errorf("LengthAny() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestArrayFromAny(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		input    any
		expected []int
	}{
		{name: "empty slice", input: []int{}, expected: []int{}},
		{name: "int slice", input: []int{1, 2, 3}, expected: []int{1, 2, 3}},
		{name: "pointer to slice", input: &[]int{1, 2, 3, 4}, expected: []int{1, 2, 3, 4}},
		{name: "array", input: [3]int{1, 2, 3}, expected: []int{1, 2, 3}},
		{name: "single value", input: 42, expected: []int{42}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result, err := util.ArrayFromAny[int](tt.input)
			if err != nil {
				t.Errorf("ArrayFromAny() = %v, want %v", result, tt.expected)
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ArrayFromAny() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestArrayTest(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		input    any
		expected bool
	}{
		{name: "empty slice", input: []int{}, expected: true},
		{name: "non-empty slice", input: []int{1, 2, 3}, expected: true},
		{name: "array", input: [3]int{1, 2, 3}, expected: true},
		{name: "string", input: "not an array", expected: false},
		{name: "integer", input: 42, expected: false},
		{name: "struct", input: struct{ name string }{"test"}, expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := util.ArrayTest(tt.input)
			if result != tt.expected {
				t.Errorf("ArrayTest() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestArrayFlatten(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		input    [][]int
		expected []int
	}{
		{name: "empty arrays", input: [][]int{}, expected: []int{}},
		{name: "single empty array", input: [][]int{{}}, expected: []int{}},
		{name: "multiple arrays", input: [][]int{{1, 2}, {3, 4}, {5, 6}}, expected: []int{1, 2, 3, 4, 5, 6}},
		{name: "mixed empty and non-empty arrays", input: [][]int{{}, {1, 2}, {}, {3, 4}}, expected: []int{1, 2, 3, 4}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := util.ArrayFlatten(tt.input...)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ArrayFlatten() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestArrayFirstN(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		input    []int
		n        int
		expected []int
	}{
		{name: "empty array", input: []int{}, n: 3, expected: []int{}},
		{name: "n larger than array", input: []int{1, 2, 3}, n: 5, expected: []int{1, 2, 3}},
		{name: "n smaller than array", input: []int{1, 2, 3, 4, 5}, n: 3, expected: []int{1, 2, 3}},
		{name: "n equal to array length", input: []int{1, 2, 3}, n: 3, expected: []int{1, 2, 3}},
		{name: "n is zero", input: []int{1, 2, 3}, n: 0, expected: []int{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := util.ArrayFirstN(tt.input, tt.n)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ArrayFirstN() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestArrayLastN(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		input    []int
		n        int
		expected []int
	}{
		{name: "empty array", input: []int{}, n: 3, expected: []int{}},
		{name: "n larger than array", input: []int{1, 2, 3}, n: 5, expected: []int{1, 2, 3}},
		{name: "n smaller than array", input: []int{1, 2, 3, 4, 5}, n: 3, expected: []int{3, 4, 5}},
		{name: "n equal to array length", input: []int{1, 2, 3}, n: 3, expected: []int{1, 2, 3}},
		{name: "n is zero", input: []int{1, 2, 3}, n: 0, expected: []int{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := util.ArrayLastN(tt.input, tt.n)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ArrayLastN() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestArrayReplaceOrAdd(t *testing.T) {
	chk := func(i int) bool {
		return i == 2
	}
	t.Parallel()
	tests := []struct {
		name     string
		input    []int
		fn       func(int) bool
		replace  int
		expected []int
	}{
		{name: "empty array", input: []int{}, fn: chk, replace: 42, expected: []int{42}},
		{name: "element found", input: []int{1, 2, 3}, fn: chk, replace: 42, expected: []int{1, 42, 3}},
		{name: "element not found", input: []int{1, 3, 5}, fn: chk, replace: 42, expected: []int{1, 3, 5, 42}},
		{name: "multiple matches (should replace first)", input: []int{1, 2, 3, 2, 4}, fn: chk, replace: 42, expected: []int{1, 42, 3, 2, 4}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := util.ArrayReplaceOrAdd(tt.input, tt.fn, tt.replace)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ArrayReplaceOrAdd() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestMapError(t *testing.T) {
	chk := func(i, idx int) (string, error) {
		return fmt.Sprintf("x%d", i), nil
	}
	chk3 := func(i, idx int) (string, error) {
		if i == 3 {
			return "", errors.Errorf("error at value %d", i)
		}
		return fmt.Sprintf("x%d", i), nil
	}
	t.Parallel()
	tests := []struct {
		name        string
		input       []int
		fn          func(int, int) (string, error)
		expected    []string
		expectError bool
	}{
		{name: "empty array", input: []int{}, fn: chk, expected: []string{}},
		{name: "no errors", input: []int{1, 2, 3}, fn: chk, expected: []string{"x1", "x2", "x3"}},
		{name: "with error", input: []int{1, 2, 3, 4, 5}, fn: chk3, expected: nil, expectError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result, err := util.MapError(tt.input, tt.fn)

			if tt.expectError {
				if err == nil {
					t.Errorf("MapError() expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("MapError() unexpected error: %v", err)
				}
				if !reflect.DeepEqual(result, tt.expected) {
					t.Errorf("MapError() = %v, want %v", result, tt.expected)
				}
			}
		})
	}
}
