package main

import (
	"testing"
	"time"
)

func TestParsePathString(t *testing.T) {
	testcases := []struct {
		input          string
		expProjectName string
		expFilename    string
		expDate        string
	}{
		// test 1 - 5
		{"postgres/weekly/iotsignals_2019-04-03.gz", "iotsignals", "iotsignals_2019-04-03.gz", "2019-04-03"},
		{"postgres/weekly/predictiveparking_2019-03-20.gz", "predictiveparking", "predictiveparking_2019-03-20.gz", "2019-03-20"},
		{"postgres/weekly/predictiveparking_2019-03-27.gz", "predictiveparking", "predictiveparking_2019-03-27.gz", "2019-03-27"},
		{"postgres/weekly/signals_2019-03-20.gz", "signals", "signals_2019-03-20.gz", "2019-03-20"},
		{"postgres/weekly/signals_2019-03-27.gz", "signals", "signals_2019-03-27.gz", "2019-03-27"},
		// test 6 - 10
		{"postgres/weekly/signals_2019-04-03.gz", "signals", "signals_2019-04-03.gz", "2019-04-03"},
		{"postgres/weekly/signals_export_2019-03-20.gz", "signals_export", "signals_export_2019-03-20.gz", "2019-03-20"},
		{"postgres/weekly/signals_export_2019-03-27.gz", "signals_export", "signals_export_2019-03-27.gz", "2019-03-27"},
		{"postgres/weekly/signals_export_2019-04-03.gz", "signals_export", "signals_export_2019-04-03.gz", "2019-04-03"},
		{"postgres/weekly/tellus_2019-03-20.gz", "tellus", "tellus_2019-03-20.gz", "2019-03-20"},
		// test 11 - 15
		{"postgres/weekly/tellus_2019-03-27.gz", "tellus", "tellus_2019-03-27.gz", "2019-03-27"},
		{"postgres/weekly/tellus_2019-04-03.gz", "tellus", "tellus_2019-04-03.gz", "2019-04-03"},
		{"postgres/weekly/timetell_2019-03-20.gz", "timetell", "timetell_2019-03-20.gz", "2019-03-20"},
		{"postgres/weekly/timetell_2019-03-27.gz", "timetell", "timetell_2019-03-27.gz", "2019-03-27"},
		{"postgres/weekly/timetell_2019-04-03.gz", "timetell", "timetell_2019-04-03.gz", "2019-04-03"},
		// test 16 - 20
		{"postgres/weekly/timetelldienstverlening_2019-03-20.gz", "timetelldienstverlening", "timetelldienstverlening_2019-03-20.gz", "2019-03-20"},
		{"postgres/weekly/timetelldienstverlening_2019-03-27.gz", "timetelldienstverlening", "timetelldienstverlening_2019-03-27.gz", "2019-03-27"},
		{"postgres/weekly/timetelldienstverlening_2019-04-03.gz", "timetelldienstverlening", "timetelldienstverlening_2019-04-03.gz", "2019-04-03"},
		{"postgres/weekly/various_small_datasets_2019-03-20.gz", "various_small_datasets", "various_small_datasets_2019-03-20.gz", "2019-03-20"},
		{"postgres/weekly/various_small_datasets_2019-03-27.gz", "various_small_datasets", "various_small_datasets_2019-03-27.gz", "2019-03-27"},
	}
	for testcaseNumber, testcase := range testcases {
		result := parseContainerName(testcase.input)
		expTime, _ := time.Parse("2006-01-02", testcase.expDate)
		if result.ProjectName != testcase.expProjectName {
			t.Error("testcase: ", testcaseNumber+1, "projectName result:", result.ProjectName, "!=", testcase.expProjectName, "fullItem:", result)
		}
		if result.Filename != testcase.expFilename {
			t.Error("testcase: ", testcaseNumber+1, "fileName result:", result.Filename, "!=", testcase.expFilename, result)
		}
		if !result.TimeStamp.Equal(expTime) {
			t.Error("testcase: ", testcaseNumber+1, "date result:", result.TimeStamp, "!=", expTime, result)
		}
	}
}

func TestBlackList(t *testing.T) {
	testcases := []struct {
		input    string
		project  string
		expected bool
	}{
		// happy path test 1 - 6
		{"foo", "foo", true},
		{"foo,bar", "foo", true},
		{"foo,bar", "bar", true},
		{"foo,bar,bla", "foo", true},
		{"foo,bar,bla", "bar", true},
		{"foo,bar,bla", "bla", true},
		// test 7 - 10
		{"", "", false},
		{"", "foo", false},
		{"foo", "", false},
		{"foo", "bar", false},
		{"foo,bar", "", false},
		{"foo,bar", "foobar", false},
		{"foo,bar", "foo,bar", false},
	}
	for testcaseNumber, testcase := range testcases {
		resultBlacklist := createProjectBlackList(testcase.input)

		if resultBlacklist.IsBlackListed(testcase.project) != testcase.expected {
			t.Error("testcase: ", testcaseNumber+1, "result:", resultBlacklist.IsBlackListed(testcase.project), "!=", testcase.expected, "item:", testcase.input, "blacklist", resultBlacklist)
		}
	}
}
