package git_test

import (
	"reflect"
	"testing"
	"time"

	"projectforge.dev/projectforge/app/lib/git"
	"projectforge.dev/projectforge/app/util"
)

func TestResultDataHelpers(t *testing.T) {
	t.Parallel()
	data := util.ValueMap{
		"branch":        "main",
		"dirty":         []string{"a.go", "b.go"},
		"commitsAhead":  2,
		"commitsBehind": "3",
		"status":        "ignored",
		"extra":         "keep",
	}
	res := git.NewResult("proj", "clean", data)

	if got := res.DataString("branch"); got != "main" {
		t.Errorf("expected branch 'main', got %q", got)
	}
	if got := res.DataInt("commitsAhead"); got != 2 {
		t.Errorf("expected commitsAhead 2, got %d", got)
	}
	if got := res.DataStringArray("dirty"); len(got) != 2 || got[0] != "a.go" || got[1] != "b.go" {
		t.Errorf("unexpected dirty list: %v", got)
	}

	cleaned := res.CleanData()
	if _, ok := cleaned["branch"]; ok {
		t.Error("expected branch to be removed from CleanData")
	}
	if _, ok := cleaned["dirty"]; ok {
		t.Error("expected dirty to be removed from CleanData")
	}
	if _, ok := cleaned["status"]; ok {
		t.Error("expected status to be removed from CleanData")
	}
	if got := cleaned.GetStringOpt("extra"); got != "keep" {
		t.Errorf("expected extra to remain, got %q", got)
	}

	if got := res.String(); got != "clean[main] (2 ahead, 3 behind)" {
		t.Errorf("unexpected String(): %s", got)
	}
}

func TestResultHistory(t *testing.T) {
	t.Parallel()
	when := time.Date(2023, 5, 15, 14, 30, 45, 0, time.UTC)
	history := &git.HistoryResult{
		Args: &git.HistoryArgs{Limit: 1},
		Entries: git.HistoryEntries{
			&git.HistoryEntry{SHA: "abc", AuthorName: "Alice", AuthorEmail: "alice@example.com", Message: "msg", Occurred: when},
		},
	}
	res := git.NewResult("proj", util.KeyOK, util.ValueMap{"history": history})
	got := res.History()
	if got == nil {
		t.Fatal("expected history, got nil")
	}
	if got.Args == nil || got.Args.Limit != 1 {
		t.Errorf("unexpected history args: %v", got.Args)
	}
	if len(got.Entries) != 1 || got.Entries[0].SHA != "abc" {
		t.Errorf("unexpected history entries: %v", got.Entries)
	}
}

func TestResultsGetAndToCSV(t *testing.T) {
	t.Parallel()
	res1 := git.NewResult("p1", "ok", util.ValueMap{"branch": "main"})
	res2 := git.NewResult("p2", "bad", util.ValueMap{})
	results := git.Results{res1, res2}

	if got := results.Get("p2"); got != res2 {
		t.Errorf("expected to find p2, got %v", got)
	}
	if got := results.Get("missing"); got != nil {
		t.Errorf("expected nil for missing result, got %v", got)
	}

	fields, rows := results.ToCSV()
	if !reflect.DeepEqual(fields, git.ResultFields) {
		t.Errorf("unexpected fields: %v", fields)
	}
	if len(rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(rows))
	}
	if !reflect.DeepEqual(rows[0], res1.Strings()) {
		t.Errorf("unexpected row 0: %v", rows[0])
	}
	if !reflect.DeepEqual(rows[1], res2.Strings()) {
		t.Errorf("unexpected row 1: %v", rows[1])
	}
}
