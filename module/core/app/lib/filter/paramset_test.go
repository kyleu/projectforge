package filter_test

import (
	"strings"
	"testing"

	"{{{ .Package }}}/app/lib/filter"
	"{{{ .Package }}}/app/lib/log"
)

func TestParamSet_Get(t *testing.T) {
	t.Parallel()
	logger, _ := log.CreateTestLogger()

	tests := []struct {
		name        string
		set         filter.ParamSet
		key         string
		allowed     []string
		expectKey   string
		expectEmpty bool
	}{
		{
			name:        "key not in set returns new params",
			set:         filter.ParamSet{},
			key:         "test",
			allowed:     nil,
			expectKey:   "test",
			expectEmpty: true,
		},
		{
			name: "key exists returns filtered params",
			set: filter.ParamSet{
				"test": &filter.Params{Key: "test", Limit: 50, Orderings: filter.Orderings{{Column: "name", Asc: true}}},
			},
			key:         "test",
			allowed:     []string{"name"},
			expectKey:   "test",
			expectEmpty: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := tt.set.Get(tt.key, tt.allowed, logger)
			if result.Key != tt.expectKey {
				t.Errorf("ParamSet.Get() key = %v, want %v", result.Key, tt.expectKey)
			}
		})
	}
}

func TestParamSet_Sanitized(t *testing.T) {
	t.Parallel()
	logger, _ := log.CreateTestLogger()
	defaultOrdering := &filter.Ordering{Column: "id", Asc: true}

	tests := []struct {
		name             string
		set              filter.ParamSet
		key              string
		expectKey        string
		expectHasDefault bool
	}{
		{
			name:             "key not in set returns sanitized with defaults",
			set:              filter.ParamSet{},
			key:              "test",
			expectKey:        "test",
			expectHasDefault: true,
		},
		{
			name: "key exists returns sanitized params",
			set: filter.ParamSet{
				"test": &filter.Params{Key: "test", Limit: 50},
			},
			key:              "test",
			expectKey:        "test",
			expectHasDefault: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := tt.set.Sanitized(tt.key, logger, defaultOrdering)
			if result.Key != tt.expectKey {
				t.Errorf("ParamSet.Sanitized() key = %v, want %v", result.Key, tt.expectKey)
			}
			if tt.expectHasDefault && len(result.Orderings) == 0 {
				t.Error("ParamSet.Sanitized() should have default orderings")
			}
		})
	}
}

func TestParamSet_Specifies(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		set      filter.ParamSet
		key      string
		expected bool
	}{
		{
			name:     "key not in set",
			set:      filter.ParamSet{},
			key:      "test",
			expected: false,
		},
		{
			name: "key exists with defaults",
			set: filter.ParamSet{
				"test": &filter.Params{Key: "test"},
			},
			key:      "test",
			expected: false,
		},
		{
			name: "key exists with non-default limit",
			set: filter.ParamSet{
				"test": &filter.Params{Key: "test", Limit: 50},
			},
			key:      "test",
			expected: true,
		},
		{
			name: "key exists with non-default offset",
			set: filter.ParamSet{
				"test": &filter.Params{Key: "test", Offset: 10},
			},
			key:      "test",
			expected: true,
		},
		{
			name: "key exists with orderings",
			set: filter.ParamSet{
				"test": &filter.Params{Key: "test", Orderings: filter.Orderings{{Column: "name"}}},
			},
			key:      "test",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := tt.set.Specifies(tt.key)
			if result != tt.expected {
				t.Errorf("ParamSet.Specifies() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestParamSet_String(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		set      filter.ParamSet
		contains []string
	}{
		{
			name:     "empty set",
			set:      filter.ParamSet{},
			contains: []string{},
		},
		{
			name: "single entry",
			set: filter.ParamSet{
				"test": &filter.Params{Key: "test", Limit: 50},
			},
			contains: []string{"test"},
		},
		{
			name: "multiple entries",
			set: filter.ParamSet{
				"one": &filter.Params{Key: "one"},
				"two": &filter.Params{Key: "two"},
			},
			contains: []string{"one", "two"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := tt.set.String()
			for _, substr := range tt.contains {
				if !strings.Contains(result, substr) {
					t.Errorf("ParamSet.String() = %v, should contain %v", result, substr)
				}
			}
		})
	}
}
