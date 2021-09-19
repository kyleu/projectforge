package diff_test

import (
	"fmt"
	"testing"

	"github.com/kyleu/projectforge/app/diff"
	"github.com/kyleu/projectforge/app/util"
)

func TestErr(t *testing.T) {
	t.Parallel()
	for _, x := range diff.AllExamples {
		diffs := x.Calc()
		t.Log(fmt.Sprintf("%s: %s", x.File, util.ToJSON(diffs)))
	}
}