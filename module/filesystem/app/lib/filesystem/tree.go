package filesystem

import (
	"path"

	"github.com/samber/lo"
	"golang.org/x/exp/slices"

	"{{{ .Package }}}/app/util"
)

type Node struct {
	Name     string   `json:"name"`
	Dir      bool     `json:"dir,omitempty"`
	Size     int      `json:"size,omitempty"`
	Children Nodes    `json:"children,omitempty"`
	Tags     []string `json:"tags,omitempty"`
}

func (n *Node) Get(path ...string) *Node {
	if len(path) == 0 {
		return n
	}
	return n.Children.Get(path...)
}

func (n *Node) Flatten(curr string) []string {
	x := path.Join(curr, n.Name)
	ret := n.Children.Flatten(x)
	if !n.Dir {
		ret = append(ret, x)
	}
	return ret
}

type Nodes []*Node

func (n Nodes) Flatten(curr string) []string {
	ret := make([]string, 0, len(n))
	for _, node := range n {
		ret = append(ret, node.Flatten(curr)...)
	}
	return ret
}

func (n Nodes) Sort() Nodes {
	slices.SortFunc(n, func(l *Node, r *Node) bool {
		return l.Name < r.Name
	})
	return n
}

func (n Nodes) Get(path ...string) *Node {
	if len(path) == 0 {
		return nil
	}
	for _, x := range n {
		if x.Name == path[0] {
			return x.Get(path[1:]...)
		}
	}
	return nil
}

func (n Nodes) Merge(x Nodes) Nodes {
	ret := slices.Clone(n)
	for _, xn := range x {
		if curr := ret.Get(xn.Name); curr != nil {
			curr.Tags = lo.Uniq(append(slices.Clone(curr.Tags), xn.Tags...))
			if len(curr.Children) == 0 && len(xn.Children) > 0 {
				curr.Children = xn.Children
			} else {
				curr.Children = curr.Children.Merge(xn.Children)
			}
		} else {
			ret = append(ret, xn)
		}
	}
	return ret
}

type Tree struct {
	Nodes  Nodes         `json:"nodes,omitempty"`
	Config util.ValueMap `json:"config,omitempty"`
	keys   []string
}

func (t Tree) Flatten() []string {
	return t.Nodes.Flatten("")
}

func (t Tree) Merge(x *Tree) *Tree {
	return &Tree{Nodes: t.Nodes.Merge(x.Nodes), Config: t.Config.Merge(x.Config)}
}

func (f *FileSystem) listNodes(pth string, ign []string, logger util.Logger, tags ...string) (Nodes, error) {
	files := f.ListFiles(pth, ign, logger)
	nodes := make(Nodes, 0, len(files))
	for _, de := range files {
		x := path.Join(pth, de.Name())
		if de.IsDir() {
			kids, err := f.listNodes(x, ign, logger)
			if err != nil {
				return nil, err
			}
			nodes = append(nodes, &Node{Name: de.Name(), Dir: true, Children: kids, Tags: tags})
		} else {
			nodes = append(nodes, &Node{Name: de.Name(), Tags: tags, Size: f.Size(x)})
		}
	}
	return nodes.Sort(), nil
}

func (f *FileSystem) ListTree(cfg util.ValueMap, pth string, ign []string, logger util.Logger, tags ...string) (*Tree, error) {
	nodes, err := f.listNodes(pth, ign, logger)
	if err != nil {
		return nil, err
	}
	return &Tree{Config: cfg, Nodes: nodes}, nil
}
