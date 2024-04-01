package main

import (
	"flag"
	"fmt"
	"os"
	"pytest_durations_report/data_loader"
	"pytest_durations_report/nodes"
	"pytest_durations_report/report"
)

const VERSION = "1.1"

func main() {
	var version = flag.Bool("version", false, "print version")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "\tcat FILENAME | %s\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "\t%s FILENAME\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()
	if *version {
		fmt.Printf("Version: %s\n", VERSION)
		os.Exit(0)
	}

	filename := flag.Arg(0)

	var tree_node nodes.TreeNode

	if filename != "" {
		tree_node = data_loader.LoadFromFile(filename)
	} else {
		tree_node = data_loader.LoadFromStdout()
	}

	nodes.CalcChildsValues(&tree_node)
	report.CreateHtmlReport(&tree_node)
}
