package report

import (
	"fmt"
	"os"
	"pytest_durations_report/root"
	"slices"
	"strings"
	"text/template"
	"time"
)

func CreateHtmlReport(leaf *root.Leaf) {
	dt := time.Now()
	data := make(map[string]string)

	filename_output := fmt.Sprintf("pytest_durations_html_report_%4d%02d%02d_%02d%02d%02d.html", dt.Year(), dt.Month(), dt.Day(), dt.Hour(), dt.Minute(), dt.Second())
	file, error := os.Create(filename_output)
	if error != nil {
		panic(error)
	}

	// TODO: get title from first leaf title, based by cmd of launch pytest
	data["title"] = fmt.Sprintf("Pytest durations report %s.html", dt.Format("2000-01-01 12:00:00"))
	data["tree_body"] = strings.Join(buildLines(leaf), "\n")

	tmpl, error := template.New("template.html").ParseFiles("report/template.html")
	if error != nil {
		panic(error)
	}

	tmpl.Execute(file, data)
}

func buildLines(leaf *root.Leaf) []string {
	result := make([]string, 0)

	tmp := `
	<span>%s</span>
	<span>%.3f</span>
	<span>%.3f</span>
	<span>%.3f</span>
	<span style="font-weight: bold">%.3f</span>
	`
	content := fmt.Sprintf(tmp, leaf.Title, leaf.TimeSetup, leaf.TimeCall, leaf.TimeTearDown, leaf.TimeTotal)

	if len(leaf.Childs) != 0 {
		leafs := make([]*root.Leaf, 0, len(leaf.Childs))
		for _, v := range leaf.Childs {
			leafs = append(leafs, v)
		}
		slices.SortFunc(leafs, root.LeafSortFunc)
		result = append(result, "<li>", "<div class=\"drop dropM\">-</div>", content, "<ul>")

		for _, leaf := range leafs {
			lines := buildLines((leaf))
			result = append(result, lines...)
		}

		result = append(result, "</ul>", "</li>")

	} else {
		result = append(result, "<li>", content, "</li>")
	}

	return result
}
