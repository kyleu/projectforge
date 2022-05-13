package diff_test

import (
	"testing"

	"projectforge.dev/projectforge/app/diff"
)

func TestDiffs(t *testing.T) {
	t.Parallel()
	for _, x := range diff.AllExamples {
		diffs := x.Calc()
		t.Logf("%s: %d edits", x.File, len(diffs.Edits))
	}
}
