//go:build test_all || !func_test
// +build test_all !func_test

package util_test

import (
	"testing"

	"projectforge.dev/projectforge/app/util"
)

func TestMarkdownTable(t *testing.T) {
	header := []string{"A", "B"}
	rows := [][]string{{"1", "2"}, {"3", "4"}}
	md, err := util.MarkdownTable(header, rows, "\n")
	if err != nil || md == "" {
		t.Fatalf("unexpected markdown table: %v %q", err, md)
	}

	ph, pr := util.MarkdownTableParse(md)
	if len(ph) != 2 || len(pr) != 2 {
		t.Fatalf("unexpected parsed table: %v %v", ph, pr)
	}

	_, err = util.MarkdownTable(header, [][]string{{"1"}}, "\n")
	if err == nil {
		t.Fatalf("expected error for mismatched row length")
	}
}
