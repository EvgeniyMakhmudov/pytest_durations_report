package nodes

import (
	"fmt"
	"strings"
)

func NewTreeNode(title string, parent *TreeNode) TreeNode {
	tree_node := TreeNode{}
	tree_node.Parent = parent
	tree_node.Title = title
	tree_node.Childs = make(map[string]*TreeNode)
	return tree_node
}

type TreeNode struct {
	Title        string
	FullPath     string
	Parent       *TreeNode
	Childs       map[string]*TreeNode
	TimeTotal    float64
	TimeSetup    float64
	TimeCall     float64
	TimeTearDown float64
}

type Record struct {
	FullPath string
	ItemName string // one of setup, call, teardown
	Value    float64
}

func AddRecordInnodes(nodes *TreeNode, record Record) {
	var parent, child *TreeNode
	var ok bool
	parent = nodes
	path_items := strings.Split(record.FullPath, "/")

	for _, path_item := range path_items {
		child, ok = parent.Childs[path_item]
		if !ok {
			tree_node := NewTreeNode(path_item, parent)
			child = &tree_node
			parent.Childs[path_item] = child
		}
		parent = child
	}

	switch record.ItemName {
	case "call":
		child.TimeCall = record.Value
	case "setup":
		child.TimeSetup = record.Value
	case "teardown":
		child.TimeTearDown = record.Value
	default:
		panic(fmt.Errorf("unexpected record name %s, want one of setup,call,teardown", record.ItemName))
	}

	// calculate timetotal by current values
	child.TimeTotal = child.TimeCall + child.TimeSetup + child.TimeTearDown
}

type values struct {
	TimeSetup    float64
	TimeCall     float64
	TimeTearDown float64
}

func CalcChildsValues(nodes *TreeNode) values {
	var setup, call, teardown float64
	for _, child := range nodes.Childs {
		if len(child.Childs) == 0 {
			setup += child.TimeSetup
			call += child.TimeCall
			teardown += child.TimeTearDown
		} else {
			values := CalcChildsValues(child)
			setup += values.TimeSetup
			call += values.TimeCall
			teardown += values.TimeTearDown
		}
	}
	nodes.TimeSetup = setup
	nodes.TimeCall = call
	nodes.TimeTearDown = teardown
	nodes.TimeTotal = setup + call + teardown
	return values{setup, call, teardown}
}

func TreeNodeSortFunc(l1, l2 *TreeNode) int {
	if l1.TimeTotal > l2.TimeTotal {
		return 1
	} else {
		if l1.TimeTotal < l2.TimeTotal {
			return -1
		} else {
			return 0
		}
	}
}

func TreeNodeSortReverseFunc(l1, l2 *TreeNode) int {
	if l1.TimeTotal > l2.TimeTotal {
		return -1
	} else {
		if l1.TimeTotal < l2.TimeTotal {
			return 1
		} else {
			return 0
		}
	}
}
