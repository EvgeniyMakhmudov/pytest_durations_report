Pytest durations report
=======================

This utility is designed to process pytest output with the `durations` option to produce an HTML report.
The report helps you understand the hierarchy of testing time distribution across Python packages.
There are also statistics and visibility of the testing time of the setup, execution, and cleaning stages.

Build and run
=============

```
go build
./pytest_durations_report FILENAME
```

or

```
go run . FILENAME
```

Usage
=====

```
Usage of ./pytest_durations_report:
        cat FILENAME | ./pytest_durations_report
        ./pytest_durations_report FILENAME
```

Example
=======

Run pytest on tagert project
```
pytest --durations 0 path/to/project > project_pytest.log
```

Make report
```
./pytest_durations_report project_pytest.log
```

Report
======

The result of executing the utility is one HTML file. Each node represents a level of the Python project hierarchy.
For example using report of testing [YADM project](https://github.com/ticketscloud/yadm), see screenshot.
The "test_embedded_2.py" node of has a test setup phase of 0.09 seconds, a test execution phase of 0.05 seconds, and a cleanup phase of 0.02 seconds. The total execution time for the package tests is 0.160 seconds. The color bar signifies the distribution of these times: red, green, black in respectively.
