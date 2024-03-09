package main

import (
	"flag"
	"pytest_durations_report/data_loader"
	"pytest_durations_report/report"
	"pytest_durations_report/root"
)

func main() {
	var leaf root.Leaf

	flag.Parse()
	filename := flag.Arg(0)

	if filename != "" {
		leaf = data_loader.LoadFromFile(filename)
	} else {
		leaf = data_loader.LoadFromStdout()
	}

	root.CalcChildsValues(&leaf)
	report.CreateHtmlReport(&leaf)
}
