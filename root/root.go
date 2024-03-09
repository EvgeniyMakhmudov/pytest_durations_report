package root

import (
	"fmt"
	"strings"
)

func NewLeaf(title string, parent *Leaf) Leaf {
	leaf := Leaf{}
	leaf.Parent = parent
	leaf.Title = title
	leaf.Childs = make(map[string]*Leaf)
	return leaf
}

type Leaf struct {
	Title        string
	FullPath     string
	Parent       *Leaf
	Childs       map[string]*Leaf
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

func AddRecordInRoot(root *Leaf, record Record) {
	var parent, child *Leaf
	var ok bool
	parent = root
	path_items := strings.Split(record.FullPath, "/")

	for _, path_item := range path_items {
		child, ok = parent.Childs[path_item]
		if !ok {
			leaf := NewLeaf(path_item, parent)
			child = &leaf
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

func CalcChildsValues(root *Leaf) values {
	var setup, call, teardown float64
	for _, child := range root.Childs {
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
	root.TimeSetup = setup
	root.TimeCall = call
	root.TimeTearDown = teardown
	root.TimeTotal = setup + call + teardown
	return values{setup, call, teardown}
}

func LeafSortFunc(l1, l2 *Leaf) int {
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
