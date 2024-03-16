package data_loader

import (
	"bufio"
	"fmt"
	"os"
	"pytest_durations_report/nodes"
	"strings"
)

func splitLikePython(s string, sep string) []string {
	result := make([]string, 0)
	items := strings.Split(s, sep)
	for _, value := range items {
		if value != "" {
			result = append(result, value)
		}
	}
	return result
}

func ParseStringToRecord(s string) nodes.Record {
	var duration float64
	items := splitLikePython(s, " ")
	if len(items) != 3 {
		items = []string{items[0], items[1], strings.Join(items[2:], " ")}
	}

	// get duration
	lastChar := items[0][len(items[0])-1:]
	firstChars := items[0][:len(items[0])-1]
	_, error := fmt.Sscanf(firstChars, "%f", &duration)
	if error != nil {
		panic(fmt.Errorf("error while convert duration string %s on line %s", firstChars, s))
	}

	switch lastChar {
	case "s":
		duration = duration * 1
	case "m":
		duration = duration * 60
	case "h":
		duration = duration * 60 * 60
	case "d":
		duration = duration * 60 * 60 * 24
	default:
		panic(fmt.Errorf("unknown postfix \"%s\" in duration, line %s", lastChar, s))
	}

	// get full path. Transform string "a/b::c[d]" to "a/b/c/d"
	var first_path, last_path = "", ""
	if index := strings.Index(items[2], "["); index != -1 {
		last_path = strings.ReplaceAll(items[2][index+1:len(items[2])-1], string(os.PathSeparator), "_")
		first_path = items[2][:index]
	} else {
		first_path = items[2]
	}

	path_items := strings.Split(first_path, string(os.PathSeparator))
	last_path_items := strings.Split(path_items[len(path_items)-1], "::")
	path_items[len(path_items)-1] = last_path_items[0]
	path_items = append(path_items, last_path_items[1])
	if last_path != "" {
		path_items = append(path_items, last_path)
	}
	full_path := strings.Join(path_items, string(os.PathSeparator))

	return nodes.Record{FullPath: full_path, ItemName: items[1], Value: duration}

}

func LoadFromStdout() nodes.TreeNode {
	scanner := bufio.NewScanner(os.Stdin)
	return Load(scanner)
}

func LoadFromFile(filename string) nodes.TreeNode {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)

	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	return Load(scanner)
}

func Load(scanner *bufio.Scanner) nodes.TreeNode {
	var data_was_start bool
	nodes_obj := nodes.NewTreeNode("nodes", nil)

	for scanner.Scan() {
		line := scanner.Text()

		if !data_was_start {
			if strings.Contains(line, "slowest durations") && strings.HasPrefix(line, "=") && strings.HasSuffix(line, "=") {
				data_was_start = true
			}
			continue
		}

		// end data block, stop reading
		if data_was_start && (line == "" || (strings.HasPrefix(line, "=") && strings.HasSuffix(line, "="))) {
			break
		}

		record := ParseStringToRecord(line)
		nodes.AddRecordInnodes(&nodes_obj, record)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return nodes_obj
}
