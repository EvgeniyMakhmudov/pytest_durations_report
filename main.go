package main

import (
	"flag"
	"pytest_durations_report/data_loader"
	"pytest_durations_report/nodes"
	"pytest_durations_report/report"
)

func main() {
	var tree_node nodes.TreeNode

	flag.Parse()
	filename := flag.Arg(0)

	if filename != "" {
		tree_node = data_loader.LoadFromFile(filename)
	} else {
		tree_node = data_loader.LoadFromStdout()
	}

	nodes.CalcChildsValues(&tree_node)
	report.CreateHtmlReport(&tree_node)
}
