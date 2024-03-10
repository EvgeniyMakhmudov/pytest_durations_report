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

const CELLWIDTH = 8

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
	data["tree_body"] = strings.Join(buildLines(leaf, len(leaf.Title)), "\n")

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

func makeColorBar(a, b, c float64) string {
	total := a + b + c
	na := int((a * 100) / total)
	nb := na + int((b*100)/total)
	tmpl := "background-image: linear-gradient(to right, red 0%% %d%% , green %d%% %d%%, blue %d%% );"
	//background-image: linear-gradient(to right, red 0% 25% , green 25% 75%, blue 75% );
	return fmt.Sprintf(tmpl, na, na, nb, nb)
}

func makeContent(leaf *root.Leaf, max_title_length int) string {

	tmp := `<span>%s</span>
	<span style="background-color:#FFF5EE">%s</span>
	<span style="background-color:#F0FFFF">%s</span>
	<span style="background-color:#FFF5EE">%s</span>
	<span style="background-color:#000;color:#FFF;font-weight: bold">%s</span>
	<div style="display:inline-block;height:5px;width:150px;%s"></div>
	`
	content := fmt.Sprintf(
		tmp,
		fillNBSP(leaf.Title, max_title_length, "right"),
		formatFloatToContentString(leaf.TimeSetup),
		formatFloatToContentString(leaf.TimeCall),
		formatFloatToContentString(leaf.TimeTearDown),
		formatFloatToContentString(leaf.TimeTotal),
		makeColorBar(leaf.TimeSetup, leaf.TimeCall, leaf.TimeTearDown),
	)
	return content
}

func buildLines(leaf *root.Leaf, max_title_length int) []string {
	result := make([]string, 0)
	content := makeContent(leaf, max_title_length)

	if len(leaf.Childs) != 0 {
		leafs := make([]*root.Leaf, 0, len(leaf.Childs))
		max_title_length := 0
		for _, v := range leaf.Childs {
			leafs = append(leafs, v)
			if len(v.Title) > max_title_length {
				max_title_length = len(v.Title)
			}
		}
		slices.SortFunc(leafs, root.LeafSortFunc)
		result = append(result, "<li>", "<div class=\"drop dropM\">-</div>", content, "<ul>")

		for _, leaf := range leafs {
			lines := buildLines(leaf, max_title_length)
			result = append(result, lines...)
		}

		result = append(result, "</ul>", "</li>")

	} else {
		result = append(result, "<li>", content, "</li>")
	}

	return result
}
