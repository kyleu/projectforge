//go:build test_all || !func_test
// +build test_all !func_test

package util_test

import (
	"testing"

	"projectforge.dev/projectforge/app/util"
)

type mergeThing struct {
	ID    string `json:"id" xml:"id"`
	Count int    `json:"count" xml:"count"`
}

func (m mergeThing) Key() string { return m.ID }

func (m mergeThing) Merge(other mergeThing) (mergeThing, error) {
	m.Count += other.Count
	return m, nil
}

func TestMergeSet(t *testing.T) {
	ms := util.NewMergeSet[string, mergeThing]()
	if err := ms.Set(mergeThing{ID: "a", Count: 1}); err != nil {
		t.Fatalf("unexpected set error: %v", err)
	}
	if err := ms.Set(mergeThing{ID: "a", Count: 2}); err != nil {
		t.Fatalf("unexpected merge error: %v", err)
	}
	if ms.Map["a"].Count != 3 {
		t.Fatalf("unexpected merge count: %d", ms.Map["a"].Count)
	}

	_ = ms.Set(mergeThing{ID: "b", Count: 4})
	if !ms.Contains(mergeThing{ID: "a"}) {
		t.Fatalf("expected contains to be true")
	}
	ms.Remove(mergeThing{ID: "b"})
	if ms.Contains(mergeThing{ID: "b"}) {
		t.Fatalf("expected remove to delete entry")
	}

	entries := ms.Entries()
	if len(entries) != 1 || entries[0].ID != "a" {
		t.Fatalf("unexpected entries: %v", entries)
	}

	clone := ms.Clone()
	clone.Map["a"] = mergeThing{ID: "a", Count: 10}
	if ms.Map["a"].Count == 10 {
		t.Fatalf("expected clone to be independent")
	}

	other := util.NewMergeSet[string, mergeThing]()
	_ = other.Set(mergeThing{ID: "c", Count: 1})
	merged, err := ms.Merge(other)
	if err != nil || merged.Length() != 2 {
		t.Fatalf("unexpected merge result: %v %v", merged, err)
	}

	data, err := merged.MarshalJSON()
	if err != nil {
		t.Fatalf("unexpected marshal error: %v", err)
	}
	copy := util.NewMergeSet[string, mergeThing]()
	if err := copy.UnmarshalJSON(data); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}
	if copy.Length() != 2 {
		t.Fatalf("unexpected unmarshaled length: %d", copy.Length())
	}
}
