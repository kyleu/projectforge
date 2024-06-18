package filesystem

import (
	"cmp"
	"path"
	"slices"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/util"
)

type Node struct {
	Name     string   `json:"name"`
	Dir      bool     `json:"dir,omitempty"`
	Size     int      `json:"size,omitempty"`
	Children Nodes    `json:"children,omitempty"`
	Tags     []string `json:"tags,omitempty"`
}

func (n *Node) Get(pth ...string) *Node {
	if len(pth) == 0 {
		return n
	}
	return n.Children.Get(pth...)
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
	ret := util.NewStringSlice(make([]string, 0, len(n)))
	for _, node := range n {
		ret.Push(node.Flatten(curr)...)
	}
	return ret.Slice
}

func (n Nodes) Sort() Nodes {
	slices.SortFunc(n, func(l *Node, r *Node) int {
		return cmp.Compare(l.Name, r.Name)
	})
	return n
}

func (n Nodes) Get(pth ...string) *Node {
	if len(pth) == 0 {
		return nil
	}
	for _, x := range n {
		if x.Name == pth[0] {
			return x.Get(pth[1:]...)
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
}

func (t Tree) Flatten() []string {
	return t.Nodes.Flatten("")
}

func (t Tree) Merge(x *Tree) *Tree {
	return &Tree{Nodes: t.Nodes.Merge(x.Nodes), Config: t.Config.Merge(x.Config)}
}
