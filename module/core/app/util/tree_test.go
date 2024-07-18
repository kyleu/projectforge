package util

import (
	"reflect"
	"testing"
)

func TestNodeGet(t *testing.T) {
	root := &Node[int]{
		Name: "root",
		Children: Nodes[int]{
			{Name: "child1", Children: Nodes[int]{{Name: "grandchild"}}},
			{Name: "child2"},
		},
	}

	tests := []struct {
		path     []string
		expected string
	}{
		{[]string{}, "root"},
		{[]string{"child1"}, "child1"},
		{[]string{"child1", "grandchild"}, "grandchild"},
		{[]string{"child2"}, "child2"},
		{[]string{"nonexistent"}, ""},
	}

	for _, tt := range tests {
		result := root.Get(tt.path...)
		if result == nil && tt.expected != "" {
			t.Errorf("Expected %s, got nil for path %v", tt.expected, tt.path)
		} else if result != nil && result.Name != tt.expected {
			t.Errorf("Expected %s, got %s for path %v", tt.expected, result.Name, tt.path)
		}
	}
}

func TestNodeFlatten(t *testing.T) {
	root := &Node[int]{
		Name: "root",
		Children: Nodes[int]{
			{Name: "child1", Children: Nodes[int]{{Name: "grandchild"}}},
			{Name: "child2"},
		},
	}

	expected := []string{
		"root/child1/grandchild",
		"root/child1",
		"root/child2",
		"root",
	}

	result := root.Flatten("")
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestNodesSort(t *testing.T) {
	nodes := Nodes[int]{
		{Name: "c"},
		{Name: "a"},
		{Name: "b"},
	}

	sorted := nodes.Sort()
	expected := []string{"a", "b", "c"}

	for i, node := range sorted {
		if node.Name != expected[i] {
			t.Errorf("Expected %s at position %d, got %s", expected[i], i, node.Name)
		}
	}
}

func TestNodesMerge(t *testing.T) {
	n1 := Nodes[int]{
		{Name: "a", Tags: []string{"tag1"}},
		{Name: "b", Children: Nodes[int]{{Name: "b1"}}},
	}
	n2 := Nodes[int]{
		{Name: "a", Tags: []string{"tag2"}},
		{Name: "c"},
	}

	merged := n1.Merge(n2)

	if len(merged) != 3 {
		t.Errorf("Expected 3 nodes, got %d", len(merged))
	}

	aNode := merged.Get("a")
	if aNode == nil || !reflect.DeepEqual(aNode.Tags, []string{"tag1", "tag2"}) {
		t.Errorf("Incorrect merge for node 'a'")
	}

	if merged.Get("b").Children.Get("b1") == nil {
		t.Errorf("Child node 'b1' not preserved in merge")
	}

	if merged.Get("c") == nil {
		t.Errorf("Node 'c' not added in merge")
	}
}

func TestTreeFlatten(t *testing.T) {
	tree := Tree[int]{
		Nodes: Nodes[int]{
			{Name: "root", Children: Nodes[int]{
				{Name: "child"},
			}},
		},
	}

	flattened := tree.Flatten()
	expected := []string{"root/child", "root"}

	if !reflect.DeepEqual(flattened, expected) {
		t.Errorf("Expected %v, got %v", expected, flattened)
	}
}

func TestTreeMerge(t *testing.T) {
	t1 := &Tree[int]{
		Nodes:  Nodes[int]{{Name: "a"}},
		Config: ValueMap{"key1": "value1"},
	}
	t2 := &Tree[int]{
		Nodes:  Nodes[int]{{Name: "b"}},
		Config: ValueMap{"key2": "value2"},
	}

	merged := t1.Merge(t2)

	if len(merged.Nodes) != 2 {
		t.Errorf("Expected 2 nodes, got %d", len(merged.Nodes))
	}

	if merged.Config["key1"] != "value1" || merged.Config["key2"] != "value2" {
		t.Errorf("Incorrect config merge")
	}
}
