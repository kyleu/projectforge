package util_test

import (
	"testing"

	"projectforge.dev/projectforge/app/util"
)

func TestParseArrayHelpers(t *testing.T) {
	t.Parallel()

	arr, err := util.ParseArray("[1,2]", "p", false, false)
	if err != nil || len(arr) != 2 {
		t.Fatalf("unexpected ParseArray result: %v %v", arr, err)
	}

	arr, err = util.ParseArray("a,b", "p", false, true)
	if err != nil || len(arr) != 2 {
		t.Fatalf("unexpected coerce ParseArray result: %v %v", arr, err)
	}

	if _, err = util.ParseArray([]any{}, "p", false, false); err == nil {
		t.Fatalf("expected error for empty array")
	}

	strs, err := util.ParseArrayString([]any{"a", 2}, "p", false)
	if err != nil || len(strs) != 2 || strs[1] != "2" {
		t.Fatalf("unexpected ParseArrayString result: %v %v", strs, err)
	}

	ints, err := util.ParseArrayInt([]any{"1", 2}, "p", false)
	if err != nil || len(ints) != 2 || ints[1] != 2 {
		t.Fatalf("unexpected ParseArrayInt result: %v %v", ints, err)
	}

	floats, err := util.ParseArrayFloat([]any{"1.5", 2}, "p", false)
	if err != nil || len(floats) != 2 || floats[0] != 1.5 {
		t.Fatalf("unexpected ParseArrayFloat result: %v %v", floats, err)
	}

	maps, err := util.ParseArrayMap([]any{map[string]any{"k": "v"}}, "p", false)
	if err != nil || len(maps) != 1 || maps[0]["k"] != "v" {
		t.Fatalf("unexpected ParseArrayMap result: %v %v", maps, err)
	}
}
