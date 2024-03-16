package data_loader

import (
	"fmt"
	"pytest_durations_report/nodes"
	"reflect"
	"strings"
	"testing"
)

func TestParseStringToRecordSuccess(t *testing.T) {
	cases := []struct {
		line string
		want nodes.Record
	}{
		{
			"2.84s setup    my_project/package_a/test_values.py::test_value1",
			nodes.Record{
				FullPath: "my_project/package_a/test_values.py/test_value1",
				ItemName: "setup",
				Value:    2.84,
			},
		}, {
			"1.5m call   my_project/package_a/test_values.py::test_value2[parametrize_values]",
			nodes.Record{
				FullPath: "my_project/package_a/test_values.py/test_value2/parametrize_values",
				ItemName: "call",
				Value:    90,
			},
		}, {
			"1.5m call   my_project/package_a/test_values.py::test_value2[parametrize_values/some_value]",
			nodes.Record{
				FullPath: "my_project/package_a/test_values.py/test_value2/parametrize_values_some_value",
				ItemName: "call",
				Value:    90,
			},
		}, {
			"1.5m call   my_project/package_a/test_values.py::test_value2[parametrize_values and some_value]",
			nodes.Record{
				FullPath: "my_project/package_a/test_values.py/test_value2/parametrize_values and some_value",
				ItemName: "call",
				Value:    90,
			},
		}, {
			"2h call   my_project/package_a/test_values.py::test_value2[parametrize_values and some_value]",
			nodes.Record{
				FullPath: "my_project/package_a/test_values.py/test_value2/parametrize_values and some_value",
				ItemName: "call",
				Value:    7200,
			},
		}, {
			"1d call   my_project/package_a/test_values.py::test_value2[parametrize_values and some_value]",
			nodes.Record{
				FullPath: "my_project/package_a/test_values.py/test_value2/parametrize_values and some_value",
				ItemName: "call",
				Value:    86400,
			},
		},
	}
	for _, c := range cases {
		result := ParseStringToRecord(c.line)

		if !reflect.DeepEqual(result, c.want) {
			t.Errorf("record %v not equal wanted %v", result, c.want)
		}
	}
}

// Runner function for run test which must panic, do partial check of error message
func runner(line, prefix string, t *testing.T) {
	defer func(prefix string) {
		if r := recover(); r != nil {
			error_s := fmt.Sprint(r)
			if !strings.HasPrefix(error_s, prefix) {
				t.Errorf("record must failed with error prefix, got %q, wanted %q", error_s, prefix)
			}
		} else {
			t.Errorf("record must failed with error prefix, but was no panic")

		}
	}(prefix)
	ParseStringToRecord(line)
}

func TestParseStringToRecordFailed(t *testing.T) {
	cases := []struct {
		line, prefix string
	}{
		{
			"2.84a setup    my_project/package_a/test_values.py::test_value1",
			"unknown postfix",
		}, {
			"qm call   my_project/package_a/test_values.py::test_value2[parametrize_values]",
			"error while convert duration string",
		},
	}
	for _, c := range cases {
		runner(c.line, c.prefix, t)
	}
}
