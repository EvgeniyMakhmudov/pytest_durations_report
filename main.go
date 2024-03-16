package main

import (
	"flag"
	"fmt"
	"os"
	"pytest_durations_report/data_loader"
	"pytest_durations_report/nodes"
	"pytest_durations_report/report"
)

func main() {
	var tree_node nodes.TreeNode

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "\tcat FILENAME | %s\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "\t%s FILENAME\n", os.Args[0])
		flag.PrintDefaults()
	}

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
