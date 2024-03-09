package data_loader

import (
	"bufio"
	"fmt"
	"os"
	"pytest_durations_report/root"
	"strings"
)

func Load(stream bufio.Reader) {

}

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

func ParseStringToRecord(s string) root.Record {
	var duration float64
	items := splitLikePython(s, " ")
	if len(items) != 3 {
		panic(fmt.Errorf("line of data have %d items,not 3. Line: %s", len(items), s))
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
	path_items := strings.Split(items[2], string(os.PathSeparator))
	last_path := path_items[len(path_items)-1]

	last_path_items := strings.Split(last_path, "::")
	path_items[len(path_items)-1] = last_path_items[0]

	if index := strings.Index(last_path_items[1], "["); index != -1 {
		path_items = append(
			path_items,
			last_path_items[1][:index],
			last_path_items[1][index+1:len(last_path_items[1])-1],
		)
	} else {
		path_items = append(path_items, last_path_items[1])
	}
	full_path := strings.Join(path_items, string(os.PathSeparator))

	return root.Record{FullPath: full_path, ItemName: items[1], Value: duration}

}

func LoadFromStdout() root.Leaf {
	root_obj := root.NewLeaf("ROOT", nil)
	// TODO: implement this
	return root_obj
}

func LoadFromFile(filename string) root.Leaf {
	var data_was_start bool
	root_obj := root.NewLeaf("ROOT", nil)

	file, err := os.Open(filename)
	if err != nil {
		panic(err)

	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if !data_was_start {
			if strings.Contains(line, "slowest durations") && strings.HasPrefix(line, "=") && strings.HasSuffix(line, "=") {
				data_was_start = true
			}
			continue
		}

		// end data block, stop reading
		if data_was_start && line == "" {
			break
		}

		record := ParseStringToRecord(line)
		root.AddRecordInRoot(&root_obj, record)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return root_obj
}
