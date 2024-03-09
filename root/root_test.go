package root

import (
	"math"
	"testing"
)

func TestAddRecordInRoot(t *testing.T) {
	root := NewLeaf("ROOT", nil)

	records := []Record{
		{"l1/l2/test1", "setup", 0.1},
		{"l1/l2/test1", "call", 0.2},
		{"l1/test2", "call", 0.7},
		{"l1/test2", "teardown", 0.3},
	}
	for _, record := range records {
		AddRecordInRoot(&root, record)
	}

	if _f := root.Childs["l1"].Childs["l2"].Childs["test1"].TimeTotal; math.Round(_f*1000)/1000 != 0.3 {
		t.Errorf("test1 total = %x, want %x\n", _f, 0.3)
	}

	if _f := root.Childs["l1"].Childs["test2"].TimeTotal; _f != 1 {
		t.Errorf("test2 total = %f, want %f\n", _f, 1.0)
	}

	root.Childs["l1"].TimeTotal = 12.0
	if _f := root.Childs["l1"].Childs["l2"].Parent.TimeTotal; _f != 12 {
		t.Errorf("parent total = %f, want %f\n", _f, 12.0)
	}
}

func TestCalcChildsValues(t *testing.T) {
	root := NewLeaf("ROOT", nil)

	records := []Record{
		{"l1/l2/test1", "setup", 1.1},
		{"l1/l2/test1", "call", 0.2},
		{"l1/l3/test2", "call", 10.1},
		{"l1/l3/test2", "teardown", 10.2},
		{"l1/l3/test3", "call", 110.1},
	}
	for _, record := range records {
		AddRecordInRoot(&root, record)
	}

	CalcChildsValues(&root)

	if _f := root.Childs["l1"].Childs["l2"].TimeTotal; math.Round(_f*1000)/1000 != 1.3 {
		t.Errorf("test1 total = %x, want %x\n", _f, 1.3)
	}

	if _f := root.Childs["l1"].Childs["l3"].TimeTotal; math.Round(_f*1000)/1000 != 130.4 {
		t.Errorf("test1 total = %x, want %x\n", _f, 130.4)
	}

	if _f := root.Childs["l1"].TimeTotal; math.Round(_f*1000)/1000 != 131.7 {
		t.Errorf("test1 total = %x, want %x\n", _f, 131.7)
	}
}
