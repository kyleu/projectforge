//go:build test_all || !func_test
// +build test_all !func_test

package result_test

import (
	"testing"

	"{{{ .Package }}}/app/lib/search/result"
)

const (
	typeUser  = "user"
	typeAdmin = "admin"
	nameAlice = "Alice"
)

func TestNewResult(t *testing.T) {
	t.Parallel()

	r := result.NewResult(typeUser, "123", "/users/123", "John Doe", "user-icon", "john", nil, "john")
	if r.Type != typeUser {
		t.Errorf("Type = %q, expected %q", r.Type, typeUser)
	}
	if r.ID != "123" {
		t.Errorf("ID = %q, expected %q", r.ID, "123")
	}
	if r.URL != "/users/123" {
		t.Errorf("URL = %q, expected %q", r.URL, "/users/123")
	}
	if r.Title != "John Doe" {
		t.Errorf("Title = %q, expected %q", r.Title, "John Doe")
	}
	if r.Icon != "user-icon" {
		t.Errorf("Icon = %q, expected %q", r.Icon, "user-icon")
	}
}

func TestNewResult_WithMatchingDiff(t *testing.T) {
	t.Parallel()

	r := result.NewResult("item", "1", "/items/1", "Test Item", "star", "hello world", nil, "world")
	if len(r.Matches) != 1 {
		t.Errorf("expected 1 match, got %d", len(r.Matches))
	}
}

func TestNewResult_NoMatchingDiff(t *testing.T) {
	t.Parallel()

	r := result.NewResult("item", "1", "/items/1", "Test Item", "star", "hello", nil, "xyz")
	if len(r.Matches) != 0 {
		t.Errorf("expected 0 matches, got %d", len(r.Matches))
	}
}

func TestResults_Sort_ByType(t *testing.T) {
	t.Parallel()

	results := result.Results{
		{Type: typeUser, Title: nameAlice},
		{Type: "item", Title: "Book"},
		{Type: typeAdmin, Title: "Charlie"},
	}
	sorted := results.Sort()

	expected := []string{typeAdmin, "item", typeUser}
	for i, r := range sorted {
		if r.Type != expected[i] {
			t.Errorf("Sort()[%d].Type = %q, expected %q", i, r.Type, expected[i])
		}
	}
}

func TestResults_Sort_ByTitleWithinType(t *testing.T) {
	t.Parallel()

	results := result.Results{
		{Type: typeUser, Title: "Charlie"},
		{Type: typeUser, Title: nameAlice},
		{Type: typeUser, Title: "Bob"},
	}
	sorted := results.Sort()

	expected := []string{nameAlice, "Bob", "Charlie"}
	for i, r := range sorted {
		if r.Title != expected[i] {
			t.Errorf("Sort()[%d].Title = %q, expected %q", i, r.Title, expected[i])
		}
	}
}

func TestResults_Sort_CaseInsensitiveTitle(t *testing.T) {
	t.Parallel()

	results := result.Results{
		{Type: typeUser, Title: "charlie"},
		{Type: typeUser, Title: nameAlice},
		{Type: typeUser, Title: "BOB"},
	}
	sorted := results.Sort()

	expected := []string{nameAlice, "BOB", "charlie"}
	for i, r := range sorted {
		if r.Title != expected[i] {
			t.Errorf("Sort()[%d].Title = %q, expected %q", i, r.Title, expected[i])
		}
	}
}

func TestResults_Sort_Empty(t *testing.T) {
	t.Parallel()

	results := result.Results{}
	sorted := results.Sort()

	if len(sorted) != 0 {
		t.Errorf("expected empty results, got %d", len(sorted))
	}
}

func TestResults_Sort_Single(t *testing.T) {
	t.Parallel()

	results := result.Results{
		{Type: typeUser, Title: nameAlice},
	}
	sorted := results.Sort()

	if len(sorted) != 1 || sorted[0].Title != nameAlice {
		t.Errorf("unexpected result for single item sort")
	}
}

func TestResults_Sort_MixedTypeAndTitle(t *testing.T) {
	t.Parallel()

	results := result.Results{
		{Type: typeUser, Title: "Zach"},
		{Type: typeAdmin, Title: "Bob"},
		{Type: typeUser, Title: nameAlice},
		{Type: typeAdmin, Title: nameAlice},
	}
	sorted := results.Sort()

	// Should be sorted by type first, then by title within type
	if sorted[0].Type != typeAdmin || sorted[0].Title != nameAlice {
		t.Errorf("sorted[0] = {%s, %s}, expected {%s, %s}", sorted[0].Type, sorted[0].Title, typeAdmin, nameAlice)
	}
	if sorted[1].Type != typeAdmin || sorted[1].Title != "Bob" {
		t.Errorf("sorted[1] = {%s, %s}, expected {%s, Bob}", sorted[1].Type, sorted[1].Title, typeAdmin)
	}
	if sorted[2].Type != typeUser || sorted[2].Title != nameAlice {
		t.Errorf("sorted[2] = {%s, %s}, expected {%s, %s}", sorted[2].Type, sorted[2].Title, typeUser, nameAlice)
	}
	if sorted[3].Type != typeUser || sorted[3].Title != "Zach" {
		t.Errorf("sorted[3] = {%s, %s}, expected {%s, Zach}", sorted[3].Type, sorted[3].Title, typeUser)
	}
}
