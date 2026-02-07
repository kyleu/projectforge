//go:build test_all || !func_test
// +build test_all !func_test

package util_test

import (
	"testing"

	"{{{ .Package }}}/app/util"
)

func TestJoinLines(t *testing.T) {
	t.Parallel()

	lines := util.JoinLines([]string{"aa", "bb", "cc"}, ",", 4)
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %v", lines)
	}

	single := util.JoinLines([]string{"a", "b"}, ",", 0)
	if len(single) != 1 || single[0] != "a,b" {
		t.Fatalf("unexpected join result: %v", single)
	}

	full := util.JoinLinesFull([]string{"a", "b"}, ",", 0, "{", " ", "}")
	if full != "{a,b}" {
		t.Fatalf("unexpected JoinLinesFull result: %q", full)
	}
}
