package filter_test

import (
	"testing"

	"projectforge.dev/projectforge/app/lib/filter"
)

func TestOrdering_String(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		ordering filter.Ordering
		expected string
	}{
		{name: "ascending", ordering: filter.Ordering{Column: "name", Asc: true}, expected: "name"},
		{name: "descending", ordering: filter.Ordering{Column: "name", Asc: false}, expected: "name:desc"},
		{name: "empty column ascending", ordering: filter.Ordering{Column: "", Asc: true}, expected: ""},
		{name: "empty column descending", ordering: filter.Ordering{Column: "", Asc: false}, expected: ":desc"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := tt.ordering.String()
			if result != tt.expected {
				t.Errorf("Ordering.String() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestOrderings_Get(t *testing.T) {
	t.Parallel()
	o1 := &filter.Ordering{Column: "name", Asc: true}
	o2 := &filter.Ordering{Column: "age", Asc: false}
	o3 := &filter.Ordering{Column: "created_at", Asc: true}

	orderings := filter.Orderings{o1, o2, o3}

	tests := []struct {
		name     string
		key      string
		expected *filter.Ordering
	}{
		{name: "find existing column", key: "name", expected: o1},
		{name: "find another existing column", key: "age", expected: o2},
		{name: "find third column", key: "created_at", expected: o3},
		{name: "column not found", key: "nonexistent", expected: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := orderings.Get(tt.key)
			if result != tt.expected {
				t.Errorf("Orderings.Get() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestOrderings_Get_EmptySlice(t *testing.T) {
	t.Parallel()
	orderings := filter.Orderings{}
	result := orderings.Get("any")
	if result != nil {
		t.Errorf("Orderings.Get() on empty slice = %v, want nil", result)
	}
}
