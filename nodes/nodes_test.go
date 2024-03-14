package nodes

import (
	"math"
	"slices"
	"testing"
)

func TestAddRecordInnodes(t *testing.T) {
	nodes := NewTreeNode("nodes", nil)

	records := []Record{
		{"l1/l2/test1", "setup", 0.1},
		{"l1/l2/test1", "call", 0.2},
		{"l1/test2", "call", 0.7},
		{"l1/test2", "teardown", 0.3},
	}
	for _, record := range records {
		AddRecordInnodes(&nodes, record)
	}

	if _f := nodes.Childs["l1"].Childs["l2"].Childs["test1"].TimeTotal; math.Round(_f*1000)/1000 != 0.3 {
		t.Errorf("test1 total = %x, want %x\n", _f, 0.3)
	}

	if _f := nodes.Childs["l1"].Childs["test2"].TimeTotal; _f != 1 {
		t.Errorf("test2 total = %f, want %f\n", _f, 1.0)
	}

	nodes.Childs["l1"].TimeTotal = 12.0
	if _f := nodes.Childs["l1"].Childs["l2"].Parent.TimeTotal; _f != 12 {
		t.Errorf("parent total = %f, want %f\n", _f, 12.0)
	}
}

func TestCalcChildsValues(t *testing.T) {
	nodes := NewTreeNode("nodes", nil)

	records := []Record{
		{"l1/l2/test1", "setup", 1.1},
		{"l1/l2/test1", "call", 0.2},
		{"l1/l3/test2", "call", 10.1},
		{"l1/l3/test2", "teardown", 10.2},
		{"l1/l3/test3", "call", 110.1},
	}
	for _, record := range records {
		AddRecordInnodes(&nodes, record)
	}

	CalcChildsValues(&nodes)

	if _f := nodes.Childs["l1"].Childs["l2"].TimeTotal; math.Round(_f*1000)/1000 != 1.3 {
		t.Errorf("test1 total = %x, want %x\n", _f, 1.3)
	}

	if _f := nodes.Childs["l1"].Childs["l3"].TimeTotal; math.Round(_f*1000)/1000 != 130.4 {
		t.Errorf("test1 total = %x, want %x\n", _f, 130.4)
	}

	if _f := nodes.Childs["l1"].TimeTotal; math.Round(_f*1000)/1000 != 131.7 {
		t.Errorf("test1 total = %x, want %x\n", _f, 131.7)
	}
}

func TestTreeNodeSortFunc(t *testing.T) {
	l3 := NewTreeNode("3", nil)
	l1 := NewTreeNode("1", nil)
	l4 := NewTreeNode("4", nil)
	l2 := NewTreeNode("2", nil)

	tree_nodes := []*TreeNode{&l3, &l1, &l4, &l2}
	tree_nodes[0].TimeTotal = 3
	tree_nodes[1].TimeTotal = 1
	tree_nodes[2].TimeTotal = 4
	tree_nodes[3].TimeTotal = 2

	slices.SortFunc(tree_nodes, TreeNodeSortFunc)
	order := make([]string, len(tree_nodes))
	for index, tree_node := range tree_nodes {
		order[index] = tree_node.Title
	}

	want := []string{"1", "2", "3", "4"}
	if !slices.Equal(order, want) {
		t.Errorf("order = %q, want %q\n", order, want)
	}
}

func TestTreeNodeSortReverseFunc(t *testing.T) {
	l3 := NewTreeNode("3", nil)
	l1 := NewTreeNode("1", nil)
	l4 := NewTreeNode("4", nil)
	l2 := NewTreeNode("2", nil)

	tree_nodes := []*TreeNode{&l3, &l1, &l4, &l2}
	tree_nodes[0].TimeTotal = 3
	tree_nodes[1].TimeTotal = 1
	tree_nodes[2].TimeTotal = 4
	tree_nodes[3].TimeTotal = 2

	slices.SortFunc(tree_nodes, TreeNodeSortReverseFunc)
	order := make([]string, len(tree_nodes))
	for index, tree_node := range tree_nodes {
		order[index] = tree_node.Title
	}

	want := []string{"4", "3", "2", "1"}
	if !slices.Equal(order, want) {
		t.Errorf("order = %q, want %q\n", order, want)
	}
}
