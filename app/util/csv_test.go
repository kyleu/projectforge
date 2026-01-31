//go:build test_all || !func_test
// +build test_all !func_test

package util_test

import (
	"strings"
	"testing"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/util"
)

type testCSV struct{}

func (t testCSV) ToCSV() ([]string, [][]string) {
	return []string{"h1", "h2"}, [][]string{{"a", "b"}}
}

func TestToCSV(t *testing.T) {
	cases := []struct {
		name    string
		input   any
		expectC []string
		expectR [][]string
		err     bool
	}{
		{name: "interface", input: testCSV{}, expectC: []string{"h1", "h2"}, expectR: [][]string{{"a", "b"}}},
		{name: "string csv", input: "c1,c2\n1,2\n3,4\n", expectC: []string{"c1", "c2"}, expectR: [][]string{{"1", "2"}, {"3", "4"}}},
		{name: "string invalid", input: "a,\"b", expectC: []string{"a,\"b"}, expectR: [][]string{}},
		{name: "error", input: errors.New("boom"), expectC: []string{"Error", "Message"}},
		{name: "detail", input: &util.ErrorDetail{Type: "t", Message: "m"}, expectC: []string{"Error", "Message"}, expectR: [][]string{{"t", "m"}}},
		{name: "unknown", input: 42, err: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cols, rows, err := util.ToCSV(tc.input)
			if tc.err {
				if err == nil {
					t.Fatalf("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if strings.Join(cols, ",") != strings.Join(tc.expectC, ",") {
				t.Fatalf("unexpected cols: %v", cols)
			}
			if tc.name == "error" {
				if len(rows) != 1 || len(rows[0]) != 2 || rows[0][1] != "boom" || !strings.Contains(rows[0][0], "errors.") {
					t.Fatalf("unexpected error row: %v", rows)
				}
				return
			}
			if len(rows) != len(tc.expectR) {
				t.Fatalf("unexpected rows length: %v", rows)
			}
			for i := range rows {
				if strings.Join(rows[i], ",") != strings.Join(tc.expectR[i], ",") {
					t.Fatalf("unexpected row %d: %v", i, rows[i])
				}
			}
		})
	}
}

func TestToCSVBytes(t *testing.T) {
	b, err := util.ToCSVBytes(testCSV{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(string(b), "h1") || !strings.Contains(string(b), "a") {
		t.Fatalf("unexpected csv bytes: %s", string(b))
	}
}
