package filesystem_test

import (
	"testing"

	"golang.org/x/exp/slices"

	"{{{ .Package }}}/app/lib/filesystem"
)

var (
	testNodeChildren = filesystem.Nodes{{Name: "bar", Dir: true, Children: filesystem.Nodes{{Name: "bark"}, {Name: "baz"}}}}
	testNode         = &filesystem.Node{Name: "foo", Dir: true, Children: testNodeChildren}
	testTree         = &filesystem.Tree{Nodes: filesystem.Nodes{testNode}}
	testTreeString   = []string{"foo/bar/bark", "foo/bar/baz"}
)

func TestTree_Flatten(t *testing.T) {
	t.Parallel()
	s := testTree.Flatten()
	if !slices.Equal(s, testTreeString) {
		t.Errorf("invalid tree string [%s]", s)
	}
}
