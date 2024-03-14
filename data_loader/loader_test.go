package data_loader

import (
	"pytest_durations_report/nodes"
	"reflect"
	"testing"
)

func TestParseStringToRecord1(t *testing.T) {
	line := "2.84s setup    my_project/package_a/test_values.py::test_value1"
	want := nodes.Record{
		FullPath: "my_project/package_a/test_values.py/test_value1",
		ItemName: "setup",
		Value:    2.84,
	}

	result := ParseStringToRecord(line)

	if !reflect.DeepEqual(result, want) {
		t.Errorf("record %v not equal wanted %v", result, want)
	}
}

func TestParseStringToRecord2(t *testing.T) {
	line := "1.5m call   my_project/package_a/test_values.py::test_value2[parametrize_values]"
	want := nodes.Record{
		FullPath: "my_project/package_a/test_values.py/test_value2/parametrize_values",
		ItemName: "call",
		Value:    90,
	}

	result := ParseStringToRecord(line)

	if !reflect.DeepEqual(result, want) {
		t.Errorf("record %v not equal wanted %v", result, want)
	}
}

func TestParseStringToRecord3(t *testing.T) {
	line := "1.5m call   my_project/package_a/test_values.py::test_value2[parametrize_values/some_value]"
	want := nodes.Record{
		FullPath: "my_project/package_a/test_values.py/test_value2/parametrize_values_some_value",
		ItemName: "call",
		Value:    90,
	}

	result := ParseStringToRecord(line)

	if !reflect.DeepEqual(result, want) {
		t.Errorf("record %v not equal wanted %v", result, want)
	}
}

func TestParseStringToRecord4(t *testing.T) {
	line := "1.5m call   my_project/package_a/test_values.py::test_value2[parametrize_values and some_value]"
	want := nodes.Record{
		FullPath: "my_project/package_a/test_values.py/test_value2/parametrize_values and some_value",
		ItemName: "call",
		Value:    90,
	}

	result := ParseStringToRecord(line)

	if !reflect.DeepEqual(result, want) {
		t.Errorf("record %v not equal wanted %v", result, want)
	}
}
