//go:build test_all || !func_test
// +build test_all !func_test

package search_test

import (
	"testing"

	"projectforge.dev/projectforge/app/lib/search"
)

func TestParams_String(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name     string
		q        string
		expected string
	}{
		{name: "simple query", q: "foo", expected: "foo"},
		{name: "empty query", q: "", expected: ""},
		{name: "query with spaces", q: "foo bar baz", expected: "foo bar baz"},
	}
	for _, tc := range cases {
		c := tc
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			p := &search.Params{Q: c.q}
			if res := p.String(); res != c.expected {
				t.Errorf("String() = %q, expected %q", res, c.expected)
			}
		})
	}
}

func TestParams_Parts(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name     string
		q        string
		expected []string
	}{
		{name: "single word", q: "foo", expected: []string{"foo"}},
		{name: "multiple words", q: "foo bar baz", expected: []string{"foo", "bar", "baz"}},
		{name: "extra spaces", q: "  foo   bar  ", expected: []string{"foo", "bar"}},
		{name: "empty query", q: "", expected: []string{}},
		{name: "with keyed params", q: "foo type:bar baz", expected: []string{"foo", "type:bar", "baz"}},
	}
	for _, tc := range cases {
		c := tc
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			p := &search.Params{Q: c.q}
			res := p.Parts()
			if len(res) != len(c.expected) {
				t.Errorf("Parts() returned %d items, expected %d: %v", len(res), len(c.expected), res)
				return
			}
			for i, v := range res {
				if v != c.expected[i] {
					t.Errorf("Parts()[%d] = %q, expected %q", i, v, c.expected[i])
				}
			}
		})
	}
}

func TestParams_General(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name     string
		q        string
		expected []string
	}{
		{name: "no keyed params", q: "foo bar baz", expected: []string{"foo", "bar", "baz"}},
		{name: "with keyed params", q: "foo type:bar baz", expected: []string{"foo", "baz"}},
		{name: "only keyed params", q: "type:bar status:active", expected: []string{}},
		{name: "empty query", q: "", expected: []string{}},
		{name: "mixed with multiple colons", q: "foo url:http://example.com bar", expected: []string{"foo", "bar"}},
	}
	for _, tc := range cases {
		c := tc
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			p := &search.Params{Q: c.q}
			res := p.General()
			if len(res) != len(c.expected) {
				t.Errorf("General() returned %d items, expected %d: %v", len(res), len(c.expected), res)
				return
			}
			for i, v := range res {
				if v != c.expected[i] {
					t.Errorf("General()[%d] = %q, expected %q", i, v, c.expected[i])
				}
			}
		})
	}
}

func TestParams_Keyed(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name     string
		q        string
		expected map[string]string
	}{
		{name: "no keyed params", q: "foo bar baz", expected: map[string]string{}},
		{name: "single keyed param", q: "foo type:bar baz", expected: map[string]string{"type": "bar"}},
		{name: "multiple keyed params", q: "type:bar status:active", expected: map[string]string{"type": "bar", "status": "active"}},
		{name: "empty query", q: "", expected: map[string]string{}},
		{name: "keyed with url value", q: "url:http://example.com", expected: map[string]string{"url": "http://example.com"}},
	}
	for _, tc := range cases {
		c := tc
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			p := &search.Params{Q: c.q}
			res := p.Keyed()
			if len(res) != len(c.expected) {
				t.Errorf("Keyed() returned %d items, expected %d: %v", len(res), len(c.expected), res)
				return
			}
			for k, v := range c.expected {
				if res[k] != v {
					t.Errorf("Keyed()[%q] = %q, expected %q", k, res[k], v)
				}
			}
		})
	}
}
