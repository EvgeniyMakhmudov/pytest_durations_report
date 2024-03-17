package report

import (
	"fmt"
	"os"
	"pytest_durations_report/nodes"
	"slices"
	"strings"
	"text/template"
	"time"
)

const CELLWIDTH = 8

func CreateHtmlReport(tree_node *nodes.TreeNode) {
	dt := time.Now()
	data := make(map[string]string)

	filename_output := fmt.Sprintf("pytest_durations_html_report_%4d%02d%02d_%02d%02d%02d.html", dt.Year(), dt.Month(), dt.Day(), dt.Hour(), dt.Minute(), dt.Second())
	file, error := os.Create(filename_output)
	if error != nil {
		panic(error)
	}

	// TODO: get title from first tree_node title, based by cmd of launch pytest
	data["title"] = fmt.Sprintf("Pytest durations report %4d%02d%02d_%02d%02d%02d.html", dt.Year(), dt.Month(), dt.Day(), dt.Hour(), dt.Minute(), dt.Second())
	data["tree_body"] = strings.Join(buildLines(tree_node, len(tree_node.Title)), "\n")

	tmpl, error := template.New("template.html").ParseFiles("report/template.html")
	if error != nil {
		panic(error)
	}

	tmpl.Execute(file, data)
}

func formatFloatToContentString(f float64) string {
	fs := fmt.Sprintf("%.3f", f)
	fs = strings.TrimSuffix(fs, ".000")

	return fillNBSP(fs, CELLWIDTH, "left")
}

// Populate input string by HTML's symbol &nbsp by target width and direction of filling
func fillNBSP(s string, width int, direction string) string {
	if l := (width - len(s)); l > 0 {
		if direction == "left" {
			return strings.Repeat("&nbsp;", l) + s
		} else {
			return s + strings.Repeat("&nbsp;", l)
		}
	} else {
		return s
	}
}

// Make CSS style for color bar element by values
func makeColorBar(a, b, c float64) string {
	total := a + b + c
	na := int((a * 100) / total)
	nb := na + int((b*100)/total)
	tmpl := "background-image: linear-gradient(to right, red 0%% %d%% , lime %d%% %d%%, black %d%% );"
	return fmt.Sprintf(tmpl, na, na, nb, nb)
}

func makeContent(tree_node *nodes.TreeNode, max_title_length int) string {
	tmp := `<span>%s</span>
	<span style="background-color:#FFF5EE">%s</span>
	<span style="background-color:#F0FFFF">%s</span>
	<span style="background-color:#FFF5EE">%s</span>
	<span style="background-color:#000;color:#FFF;font-weight: bold">%s</span>
	<div class="tooltip cmd_short" data-tooltip="%s" style="display:inline-block;height:5px;width:150px;%s"></div>
	`

	na := int((tree_node.TimeSetup * 100) / tree_node.TimeTotal)
	nb := int((tree_node.TimeCall * 100) / tree_node.TimeTotal)
	nc := 100 - na - nb
	tooltip := fmt.Sprintf("%d%% %d%% %d%%", na, nb, nc)

	content := fmt.Sprintf(
		tmp,
		fillNBSP(tree_node.Title, max_title_length, "right"),
		formatFloatToContentString(tree_node.TimeSetup),
		formatFloatToContentString(tree_node.TimeCall),
		formatFloatToContentString(tree_node.TimeTearDown),
		formatFloatToContentString(tree_node.TimeTotal),
		tooltip,
		makeColorBar(tree_node.TimeSetup, tree_node.TimeCall, tree_node.TimeTearDown),
	)
	return content
}

func buildLines(tree_node *nodes.TreeNode, max_title_length int) []string {
	result := make([]string, 0)
	content := makeContent(tree_node, max_title_length)

	if len(tree_node.Childs) != 0 {
		tree_nodes := make([]*nodes.TreeNode, 0, len(tree_node.Childs))
		max_title_length := 0
		for _, v := range tree_node.Childs {
			tree_nodes = append(tree_nodes, v)
			if len(v.Title) > max_title_length {
				max_title_length = len(v.Title)
			}
		}
		slices.SortFunc(tree_nodes, nodes.TreeNodeSortReverseFunc)
		result = append(result, "<li>", "<div class=\"drop dropM\">-</div>", content, "<ul>")

		for _, tree_node := range tree_nodes {
			lines := buildLines(tree_node, max_title_length)
			result = append(result, lines...)
		}

		result = append(result, "</ul>", "</li>")

	} else {
		result = append(result, "<li>", content, "</li>")
	}

	return result
}
