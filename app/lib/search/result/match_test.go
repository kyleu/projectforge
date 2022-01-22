// Content managed by Project Forge, see [projectforge.md] for details.
package result_test

import (
	"testing"

	"github.com/kyleu/projectforge/app/lib/search/result"
	"github.com/kyleu/projectforge/app/util"
)

var splitTests = []struct {
	q string
	t string
	r []string
}{
	{q: "foo", t: "there's a foo here", r: []string{"there's a ", "foo", " here"}},
	{q: "foo", t: "foo is what this is", r: []string{"foo", " is what this is"}},
	{q: "foo", t: "this is a foo", r: []string{"this is a ", "foo"}},
	{q: "foo", t: "there's a foo here and a foo there", r: []string{"there's a ", "foo", " here and a ", "foo", " there"}},
}

func TestEncryptDecrypt(t *testing.T) {
	t.Parallel()

	for _, tt := range splitTests {
		m := &result.Match{Key: "test", Value: tt.t}
		r := m.ValueSplit(tt.q)
		if len(tt.r) != len(r) {
			t.Errorf("%s :: %s", tt.t, util.ToJSONCompact(r))
		}
	}
}
