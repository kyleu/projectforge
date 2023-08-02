// Content managed by Project Forge, see [projectforge.md] for details.
package filesystem_test

import (
	"testing"

	"golang.org/x/exp/slices"

	"projectforge.dev/projectforge/app/lib/filesystem"
)

var (
	testNodeChildren = filesystem.Nodes{{Name: "bar", Dir: true, Children: filesystem.Nodes{{Name: "bark"}, {Name: "baz"}}}}
	testNode         = &filesystem.Node{Name: "foo", Dir: true, Children: testNodeChildren}
	testTree         = &filesystem.Tree{Nodes: filesystem.Nodes{testNode}}
	testTreeString   = []string{"foo/bar/bark", "foo/bar/baz"}
)

func TestTree_Flatten(t *testing.T) {
	s := testTree.Flatten()
	if !slices.Equal(s, testTreeString) {
		t.Errorf("invalid tree string [%s]", s)
	}
}
