package diff_test

import (
	"fmt"
	"testing"

	"github.com/kyleu/projectforge/app/diff"
	"github.com/kyleu/projectforge/app/util"
)

func TestErr(t *testing.T) {
	for _, x := range diff.AllExamples {
		diffs := x.Calc()
		println(fmt.Sprintf("%s: %s", x.File, util.ToJSON(diffs)))
	}
}
