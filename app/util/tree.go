package util

import (
	"cmp"
	"path"
	"slices"

	"github.com/samber/lo"
)

type Node[T any] struct {
	Name     string   `json:"name"`
	Children Nodes[T] `json:"children,omitempty"`
	Tags     []string `json:"tags,omitempty"`
}

func (n *Node[T]) Get(pth ...string) *Node[T] {
	if len(pth) == 0 {
		return n
	}
	return n.Children.Get(pth...)
}

func (n *Node[T]) Flatten(curr string) []string {
	x := path.Join(curr, n.Name)
	ret := n.Children.Flatten(x)
	ret = append(ret, x)
	return ret
}

type Nodes[T any] []*Node[T]

func (n Nodes[T]) Flatten(curr string) []string {
	ret := NewStringSlice(make([]string, 0, len(n)))
	for _, node := range n {
		ret.Push(node.Flatten(curr)...)
	}
	return ret.Slice
}

func (n Nodes[T]) Sort() Nodes[T] {
	slices.SortFunc(n, func(l *Node[T], r *Node[T]) int {
		return cmp.Compare(l.Name, r.Name)
	})
	return n
}

func (n Nodes[T]) Get(pth ...string) *Node[T] {
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

func (n Nodes[T]) Merge(x Nodes[T]) Nodes[T] {
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

type Tree[T any] struct {
	Nodes  Nodes[T] `json:"nodes,omitempty"`
	Config ValueMap `json:"config,omitempty"`
}

func (t Tree[T]) Flatten() []string {
	return t.Nodes.Flatten("")
}

func (t Tree[T]) Merge(x *Tree[T]) *Tree[T] {
	return &Tree[T]{Nodes: t.Nodes.Merge(x.Nodes), Config: t.Config.Merge(x.Config)}
}
