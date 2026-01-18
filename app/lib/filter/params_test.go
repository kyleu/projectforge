package filter_test

import (
	"net/url"
	"testing"

	"projectforge.dev/projectforge/app/lib/filter"
	"projectforge.dev/projectforge/app/lib/log"
)

func TestParamsWithDefaultOrdering(t *testing.T) {
	t.Parallel()
	defaultOrdering := &filter.Ordering{Column: "created_at", Asc: false}

	tests := []struct {
		name             string
		key              string
		params           *filter.Params
		defaultOrderings []*filter.Ordering
		expectOrderings  int
	}{
		{
			name:             "nil params gets default orderings",
			key:              "test",
			params:           nil,
			defaultOrderings: []*filter.Ordering{defaultOrdering},
			expectOrderings:  1,
		},
		{
			name:             "empty orderings gets defaults",
			key:              "test",
			params:           &filter.Params{Key: "test"},
			defaultOrderings: []*filter.Ordering{defaultOrdering},
			expectOrderings:  1,
		},
		{
			name:             "existing orderings preserved",
			key:              "test",
			params:           &filter.Params{Key: "test", Orderings: filter.Orderings{{Column: "name", Asc: true}}},
			defaultOrderings: []*filter.Ordering{defaultOrdering},
			expectOrderings:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := filter.ParamsWithDefaultOrdering(tt.key, tt.params, tt.defaultOrderings...)
			if len(result.Orderings) != tt.expectOrderings {
				t.Errorf("ParamsWithDefaultOrdering() orderings length = %v, want %v", len(result.Orderings), tt.expectOrderings)
			}
		})
	}
}

func TestParams_Sanitize(t *testing.T) {
	t.Parallel()
	defaultOrdering := &filter.Ordering{Column: "id", Asc: true}

	tests := []struct {
		name           string
		params         *filter.Params
		key            string
		expectedLimit  int
		expectedOffset int
	}{
		{
			name:           "nil params returns new params with defaults",
			params:         nil,
			key:            "test",
			expectedLimit:  0,
			expectedOffset: 0,
		},
		{
			name:           "zero limit gets default page size",
			params:         &filter.Params{Key: "test", Limit: 0},
			key:            "test",
			expectedLimit:  filter.PageSize,
			expectedOffset: 0,
		},
		{
			name:           "limit exceeding max gets default page size",
			params:         &filter.Params{Key: "test", Limit: filter.MaxRows + 1},
			key:            "test",
			expectedLimit:  filter.PageSize,
			expectedOffset: 0,
		},
		{
			name:           "negative offset becomes zero",
			params:         &filter.Params{Key: "test", Limit: 50, Offset: -10},
			key:            "test",
			expectedLimit:  50,
			expectedOffset: 0,
		},
		{
			name:           "valid params unchanged",
			params:         &filter.Params{Key: "test", Limit: 50, Offset: 100},
			key:            "test",
			expectedLimit:  50,
			expectedOffset: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := tt.params.Sanitize(tt.key, defaultOrdering)
			if result.Limit != tt.expectedLimit {
				t.Errorf("Params.Sanitize() limit = %v, want %v", result.Limit, tt.expectedLimit)
			}
			if result.Offset != tt.expectedOffset {
				t.Errorf("Params.Sanitize() offset = %v, want %v", result.Offset, tt.expectedOffset)
			}
		})
	}
}

func TestParams_WithKey(t *testing.T) {
	t.Parallel()
	p := &filter.Params{Key: "original"}
	result := p.WithKey("new_key")
	if result.Key != "new_key" {
		t.Errorf("Params.WithKey() = %v, want %v", result.Key, "new_key")
	}
	if result != p {
		t.Error("Params.WithKey() should return same pointer")
	}
}

func TestParams_WithLimit(t *testing.T) {
	t.Parallel()
	p := &filter.Params{Limit: 10}
	result := p.WithLimit(50)
	if result.Limit != 50 {
		t.Errorf("Params.WithLimit() = %v, want %v", result.Limit, 50)
	}
	if result != p {
		t.Error("Params.WithLimit() should return same pointer")
	}
}

func TestParams_WithOffset(t *testing.T) {
	t.Parallel()
	p := &filter.Params{Offset: 0}
	result := p.WithOffset(100)
	if result.Offset != 100 {
		t.Errorf("Params.WithOffset() = %v, want %v", result.Offset, 100)
	}
	if result != p {
		t.Error("Params.WithOffset() should return same pointer")
	}
}

func TestParams_CloneOrdering(t *testing.T) {
	t.Parallel()
	originalOrdering := &filter.Ordering{Column: "name", Asc: true}
	newOrdering := &filter.Ordering{Column: "age", Asc: false}

	tests := []struct {
		name        string
		params      *filter.Params
		newOrdering *filter.Ordering
		expectNil   bool
	}{
		{
			name:        "nil params returns nil",
			params:      nil,
			newOrdering: newOrdering,
			expectNil:   true,
		},
		{
			name:        "clones with new ordering",
			params:      &filter.Params{Key: "test", Orderings: filter.Orderings{originalOrdering}, Limit: 10, Offset: 5},
			newOrdering: newOrdering,
			expectNil:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := tt.params.CloneOrdering(tt.newOrdering)
			if tt.expectNil {
				if result != nil {
					t.Errorf("Params.CloneOrdering() = %v, want nil", result)
				}
				return
			}
			if result == nil {
				t.Fatal("Params.CloneOrdering() = nil, want non-nil")
			}
			if result == tt.params {
				t.Error("Params.CloneOrdering() should return a new instance")
			}
			if len(result.Orderings) != 1 || result.Orderings[0] != tt.newOrdering {
				t.Errorf("Params.CloneOrdering() orderings = %v, want %v", result.Orderings, filter.Orderings{tt.newOrdering})
			}
		})
	}
}

func TestParams_CloneLimit(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		params    *filter.Params
		newLimit  int
		expectNil bool
	}{
		{
			name:      "nil params returns nil",
			params:    nil,
			newLimit:  50,
			expectNil: true,
		},
		{
			name:      "clones with new limit",
			params:    &filter.Params{Key: "test", Limit: 10, Offset: 5},
			newLimit:  50,
			expectNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := tt.params.CloneLimit(tt.newLimit)
			if tt.expectNil {
				if result != nil {
					t.Errorf("Params.CloneLimit() = %v, want nil", result)
				}
				return
			}
			if result == nil {
				t.Fatal("Params.CloneLimit() = nil, want non-nil")
			}
			if result == tt.params {
				t.Error("Params.CloneLimit() should return a new instance")
			}
			if result.Limit != tt.newLimit {
				t.Errorf("Params.CloneLimit() limit = %v, want %v", result.Limit, tt.newLimit)
			}
		})
	}
}

func TestParams_HasNextPage(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		params   *filter.Params
		count    int
		expected bool
	}{
		{name: "nil params", params: nil, count: 100, expected: false},
		{name: "zero limit", params: &filter.Params{Limit: 0}, count: 100, expected: false},
		{name: "count less than offset plus limit", params: &filter.Params{Limit: 10, Offset: 0}, count: 5, expected: false},
		{name: "count equals offset plus limit", params: &filter.Params{Limit: 10, Offset: 0}, count: 10, expected: true},
		{name: "count greater than offset plus limit", params: &filter.Params{Limit: 10, Offset: 0}, count: 15, expected: true},
		{name: "with offset", params: &filter.Params{Limit: 10, Offset: 20}, count: 30, expected: true},
		{name: "with offset not enough", params: &filter.Params{Limit: 10, Offset: 20}, count: 25, expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := tt.params.HasNextPage(tt.count)
			if result != tt.expected {
				t.Errorf("Params.HasNextPage() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestParams_NextPage(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		params         *filter.Params
		expectedOffset int
	}{
		{
			name:           "from zero offset with limit",
			params:         &filter.Params{Key: "test", Limit: 10, Offset: 0},
			expectedOffset: 10,
		},
		{
			name:           "from non-zero offset",
			params:         &filter.Params{Key: "test", Limit: 10, Offset: 20},
			expectedOffset: 30,
		},
		{
			name:           "zero limit uses page size",
			params:         &filter.Params{Key: "test", Limit: 0, Offset: 0},
			expectedOffset: filter.PageSize,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := tt.params.NextPage()
			if result.Offset != tt.expectedOffset {
				t.Errorf("Params.NextPage() offset = %v, want %v", result.Offset, tt.expectedOffset)
			}
			if result == tt.params {
				t.Error("Params.NextPage() should return a new instance")
			}
		})
	}
}

func TestParams_HasPreviousPage(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		params   *filter.Params
		expected bool
	}{
		{name: "nil params", params: nil, expected: false},
		{name: "zero offset", params: &filter.Params{Offset: 0}, expected: false},
		{name: "positive offset", params: &filter.Params{Offset: 10}, expected: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := tt.params.HasPreviousPage()
			if result != tt.expected {
				t.Errorf("Params.HasPreviousPage() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestParams_PreviousPage(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		params         *filter.Params
		expectedOffset int
	}{
		{
			name:           "from offset greater than limit",
			params:         &filter.Params{Key: "test", Limit: 10, Offset: 20},
			expectedOffset: 10,
		},
		{
			name:           "from offset less than limit",
			params:         &filter.Params{Key: "test", Limit: 10, Offset: 5},
			expectedOffset: 0,
		},
		{
			name:           "from zero offset",
			params:         &filter.Params{Key: "test", Limit: 10, Offset: 0},
			expectedOffset: 0,
		},
		{
			name:           "zero limit uses page size",
			params:         &filter.Params{Key: "test", Limit: 0, Offset: 150},
			expectedOffset: 50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := tt.params.PreviousPage()
			if result.Offset != tt.expectedOffset {
				t.Errorf("Params.PreviousPage() offset = %v, want %v", result.Offset, tt.expectedOffset)
			}
			if result == tt.params {
				t.Error("Params.PreviousPage() should return a new instance")
			}
		})
	}
}

func TestParams_GetOrdering(t *testing.T) {
	t.Parallel()
	o1 := &filter.Ordering{Column: "name", Asc: true}
	o2 := &filter.Ordering{Column: "age", Asc: false}

	p := &filter.Params{Orderings: filter.Orderings{o1, o2}}

	tests := []struct {
		name     string
		col      string
		expected *filter.Ordering
	}{
		{name: "find existing", col: "name", expected: o1},
		{name: "find second", col: "age", expected: o2},
		{name: "not found", col: "nonexistent", expected: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := p.GetOrdering(tt.col)
			if result != tt.expected {
				t.Errorf("Params.GetOrdering() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestParams_OrderByString(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		params   *filter.Params
		expected string
	}{
		{
			name:     "empty orderings",
			params:   &filter.Params{Key: "test"},
			expected: "",
		},
		{
			name:     "single ascending",
			params:   &filter.Params{Key: "test", Orderings: filter.Orderings{{Column: "name", Asc: true}}},
			expected: `"name"`,
		},
		{
			name:     "single descending",
			params:   &filter.Params{Key: "test", Orderings: filter.Orderings{{Column: "name", Asc: false}}},
			expected: `"name" desc`,
		},
		{
			name: "multiple orderings",
			params: &filter.Params{Key: "test", Orderings: filter.Orderings{
				{Column: "name", Asc: true},
				{Column: "age", Asc: false},
			}},
			expected: `"name", "age" desc`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := tt.params.OrderByString()
			if result != tt.expected {
				t.Errorf("Params.OrderByString() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestParams_Filtered(t *testing.T) {
	t.Parallel()
	logger, _ := log.CreateTestLogger()

	tests := []struct {
		name             string
		params           *filter.Params
		key              string
		available        []string
		expectedOrderLen int
	}{
		{
			name:             "wildcard allows all",
			params:           &filter.Params{Key: "test", Orderings: filter.Orderings{{Column: "any", Asc: true}}},
			key:              "test",
			available:        []string{"*"},
			expectedOrderLen: 1,
		},
		{
			name:             "filters to allowed columns",
			params:           &filter.Params{Key: "test", Orderings: filter.Orderings{{Column: "name", Asc: true}, {Column: "invalid", Asc: false}}},
			key:              "test",
			available:        []string{"name", "age"},
			expectedOrderLen: 1,
		},
		{
			name:             "empty orderings unchanged",
			params:           &filter.Params{Key: "test"},
			key:              "test",
			available:        []string{"name"},
			expectedOrderLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := tt.params.Filtered(tt.key, tt.available, logger)
			if len(result.Orderings) != tt.expectedOrderLen {
				t.Errorf("Params.Filtered() orderings length = %v, want %v", len(result.Orderings), tt.expectedOrderLen)
			}
		})
	}
}

func TestParams_IsDefault(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		params   *filter.Params
		expected bool
	}{
		{name: "all defaults", params: &filter.Params{}, expected: true},
		{name: "has offset", params: &filter.Params{Offset: 10}, expected: false},
		{name: "has limit", params: &filter.Params{Limit: 10}, expected: false},
		{name: "has orderings", params: &filter.Params{Orderings: filter.Orderings{{Column: "name"}}}, expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := tt.params.IsDefault()
			if result != tt.expected {
				t.Errorf("Params.IsDefault() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestParams_String(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		params   *filter.Params
		contains []string
	}{
		{
			name:     "basic params",
			params:   &filter.Params{Key: "test"},
			contains: []string{"test"},
		},
		{
			name:     "with offset and limit",
			params:   &filter.Params{Key: "test", Offset: 10, Limit: 20},
			contains: []string{"test", "10/", "20"},
		},
		{
			name:     "with orderings",
			params:   &filter.Params{Key: "test", Orderings: filter.Orderings{{Column: "name", Asc: true}}},
			contains: []string{"test", "name"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := tt.params.String()
			for _, substr := range tt.contains {
				if !containsString(result, substr) {
					t.Errorf("Params.String() = %v, should contain %v", result, substr)
				}
			}
		})
	}
}

func TestParams_ToQueryString(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		params   *filter.Params
		url      *url.URL
		contains []string
		excludes []string
	}{
		{
			name:     "nil params",
			params:   nil,
			url:      mustParseURL("http://example.com"),
			contains: []string{},
		},
		{
			name:     "nil url",
			params:   &filter.Params{Key: "test"},
			url:      nil,
			contains: []string{},
		},
		{
			name:     "with orderings",
			params:   &filter.Params{Key: "test", Orderings: filter.Orderings{{Column: "name", Asc: true}}},
			url:      mustParseURL("http://example.com"),
			contains: []string{"test.o=name"},
		},
		{
			name:     "with descending ordering",
			params:   &filter.Params{Key: "test", Orderings: filter.Orderings{{Column: "name", Asc: false}}},
			url:      mustParseURL("http://example.com"),
			contains: []string{"test.o=name.d"},
		},
		{
			name:     "with limit",
			params:   &filter.Params{Key: "test", Limit: 50},
			url:      mustParseURL("http://example.com"),
			contains: []string{"test.l=50"},
		},
		{
			name:     "with offset",
			params:   &filter.Params{Key: "test", Offset: 100},
			url:      mustParseURL("http://example.com"),
			contains: []string{"test.x=100"},
		},
		{
			name:     "default limit 1000 excluded",
			params:   &filter.Params{Key: "test", Limit: 1000},
			url:      mustParseURL("http://example.com"),
			excludes: []string{"test.l="},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := tt.params.ToQueryString(tt.url)
			for _, substr := range tt.contains {
				if !containsString(result, substr) {
					t.Errorf("Params.ToQueryString() = %v, should contain %v", result, substr)
				}
			}
			for _, substr := range tt.excludes {
				if containsString(result, substr) {
					t.Errorf("Params.ToQueryString() = %v, should not contain %v", result, substr)
				}
			}
		})
	}
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || substr == "" ||
		(s != "" && substr != "" && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func mustParseURL(rawURL string) *url.URL {
	u, err := url.Parse(rawURL)
	if err != nil {
		panic(err)
	}
	return u
}
