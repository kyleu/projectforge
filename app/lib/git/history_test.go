package git_test

import (
	"strings"
	"testing"
	"time"

	"projectforge.dev/projectforge/app/lib/git"
)

const (
	historyFieldDelim = "\u00bb\u00a6\u00ab"
	historyLineDelim  = "\u00bb\u00a6\u00a6\u00a6\u00ab"
)

func TestParseResultsDelimited(t *testing.T) {
	t.Parallel()
	dateStr := "Mon May 15 14:30:45 2023 +0000"
	entry1 := strings.Join([]string{"sha1", "Alice", "alice@example.com", dateStr, "First commit"}, historyFieldDelim)
	entry2 := strings.Join([]string{"sha2", "Bob", "bob@example.com", dateStr, "Second commit"}, historyFieldDelim)
	output := entry1 + historyLineDelim + entry2

	res, err := git.ParseResultsDelimited(output)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(res))
	}
	if res[0].SHA != "sha1" {
		t.Errorf("unexpected sha for entry 0: %s", res[0].SHA)
	}
	if res[0].AuthorName != "Alice" || res[0].AuthorEmail != "alice@example.com" {
		t.Errorf("unexpected author for entry 0: %s <%s>", res[0].AuthorName, res[0].AuthorEmail)
	}
	if res[0].Message != "First commit" {
		t.Errorf("unexpected message for entry 0: %s", res[0].Message)
	}
	expectedTime := time.Date(2023, 5, 15, 14, 30, 45, 0, time.UTC)
	if !res[0].Occurred.Equal(expectedTime) {
		t.Errorf("unexpected time for entry 0: %s", res[0].Occurred)
	}
}

func TestParseResultsDelimitedInvalidParts(t *testing.T) {
	t.Parallel()
	output := strings.Join([]string{"sha1", "Alice", "alice@example.com", "First commit"}, historyFieldDelim)
	_, err := git.ParseResultsDelimited(output)
	if err == nil {
		t.Fatal("expected error for invalid parts, got nil")
	}
}

func TestParseResultsDelimitedInvalidTime(t *testing.T) {
	t.Parallel()
	output := strings.Join([]string{"sha1", "Alice", "alice@example.com", "not-a-time", "First commit"}, historyFieldDelim)
	_, err := git.ParseResultsDelimited(output)
	if err == nil {
		t.Fatal("expected error for invalid time, got nil")
	}
}

func TestHistoryEntriesGetAndAuthors(t *testing.T) {
	t.Parallel()
	when := time.Date(2023, 5, 15, 14, 30, 45, 0, time.UTC)
	entries := git.HistoryEntries{
		&git.HistoryEntry{SHA: "a", AuthorName: "Alice", AuthorEmail: "alice@example.com", Occurred: when},
		&git.HistoryEntry{SHA: "b", AuthorName: "Bob", AuthorEmail: "bob@example.com", Occurred: when},
		&git.HistoryEntry{SHA: "c", AuthorName: "Alice", AuthorEmail: "alice@example.com", Occurred: when},
	}
	if got := entries.Get("b"); got == nil || got.SHA != "b" {
		t.Errorf("unexpected entry for sha b: %v", got)
	}
	if got := entries.Get("missing"); got != nil {
		t.Errorf("expected nil entry for missing sha, got %v", got)
	}

	authors := entries.Authors()
	if len(authors) != 2 {
		t.Fatalf("expected 2 authors, got %d", len(authors))
	}
	byEmail := map[string]*git.HistoryAuthor{}
	for _, author := range authors {
		byEmail[author.Key] = author
	}
	if got := byEmail["alice@example.com"]; got == nil || got.Count != 2 || got.Name != "Alice" {
		t.Errorf("unexpected author summary for Alice: %v", got)
	}
	if got := byEmail["bob@example.com"]; got == nil || got.Count != 1 || got.Name != "Bob" {
		t.Errorf("unexpected author summary for Bob: %v", got)
	}
}

func TestHistoryFilesStrings(t *testing.T) {
	t.Parallel()
	files := git.HistoryFiles{
		&git.HistoryFile{Status: "M", File: "a.txt"},
		&git.HistoryFile{Status: "A", File: "b.txt"},
	}
	res := files.Strings()
	if len(res) != 2 {
		t.Fatalf("expected 2 strings, got %d", len(res))
	}
	if res[0] != "a.txt" || res[1] != "b.txt" {
		t.Errorf("unexpected file list: %v", res)
	}
}
