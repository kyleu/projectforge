package diff_test

import (
	"fmt"
	"testing"

	"github.com/kyleu/projectforge/app/diff"
)

func TestDiffs(t *testing.T) {
	t.Parallel()
	for _, x := range diff.AllExamples {
		diffs := x.Calc()
		t.Log(fmt.Sprintf("%s: %d edits", x.File, len(diffs.Edits)))
	}
}
