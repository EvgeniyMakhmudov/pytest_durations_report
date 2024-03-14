package report

import (
	"pytest_durations_report/root"
	"strings"
	"testing"
)

func TestMakeColorBar(t *testing.T) {
	cases := []struct {
		a, b, c float64
		want    string
	}{
		{1, 2, 1, "background-image: linear-gradient(to right, red 0% 25% , lime 25% 75%, black 75% );"},
		{4, 5, 1, "background-image: linear-gradient(to right, red 0% 40% , lime 40% 90%, black 90% );"},
	}

	for _, c := range cases {
		result := makeColorBar(c.a, c.b, c.c)
		if result != c.want {
			t.Errorf("Failed result %q, want %q", result, c.want)
		}
	}
}

func TestFormatFloatToContentString(t *testing.T) {
	cases := []struct {
		a    float64
		want string
	}{
		{1, "&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;1"},
		{123, "&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;123"},
		{1.2345, "&nbsp;&nbsp;&nbsp;1.234"},
	}

	for _, c := range cases {
		result := formatFloatToContentString(c.a)
		if result != c.want {
			t.Errorf("Failed result %q, want %q", result, c.want)
		}
	}
}

func TestFillNBSP(t *testing.T) {
	cases := []struct {
		s         string
		width     int
		direction string
		want      string
	}{
		{"qwe", 4, "left", "&nbsp;qwe"},
		{"qwe", 4, "right", "qwe&nbsp;"},
		{"qwe", 1, "left", "qwe"},
		{"qwe", 1, "right", "qwe"},
	}

	for _, c := range cases {
		result := fillNBSP(c.s, c.width, c.direction)
		if result != c.want {
			t.Errorf("Failed result %q, want %q", result, c.want)
		}
	}
}

func TestBuildLines(t *testing.T) {
	top := root.NewLeaf("ROOT", nil)
	top.TimeTotal = 10.5
	top.TimeSetup = 1.0
	top.TimeCall = 9.5

	child := root.NewLeaf("Child", &top)
	top.TimeTotal = 10.5
	top.TimeSetup = 1.0
	top.TimeCall = 9.5
	top.Childs["Child"] = &child

	lines := buildLines(&top, 4)
	if lines[0] != "<li>" {
		t.Errorf("Failed check first string: %q, want %q", lines[0], "<li>")
	}
	if !strings.Contains(lines[1], "drop dropM") {
		t.Errorf("Failed check second string: %q, must contain %q", lines[1], "drop dropM")
	}
	if lines[len(lines)-1] != "</li>" {
		t.Errorf("Failed check last string: %q, want %q", lines[len(lines)-1], "</li>")
	}

}
