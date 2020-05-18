package loggrouper

import (
	"io"
	"strings"
	"testing"
	"time"
)

// Common Testing stuff
// Dont change this without adjusting testdata
const Interval string = "15m"

type TestData struct {
	Logline         string
	TimestampString string
	GroupString     string
}

// TODO: Add wayyyyy more log formats (OPTIONS, HEAD, Timezones, ...)
var ApacheTestData = []TestData{
	{`127.0.0.1 - frank [10/Oct/2000:13:55:36 -0700] "GET /apache_pb.gif HTTP/1.0" 200 2326 "http://www.example.com/start.html" "Mozilla/4.08 [en] (Win98; I ;Nav)"`, "2000-10-10 22:45:00 " + time.Now().Local().Format("-0700 MST"), "/apache_pb.gif"},
	{`127.0.0.1 - frank [10/Oct/2000:13:55:36 -0700] "GET /apache_pb.gif HTTP/1.0" 200 2326`, "2000-10-10 22:45:00 " + time.Now().Local().Format("-0700 MST"), "/apache_pb.gif"},
}

func getReaderFromData(data []TestData) io.Reader {
	var logLines string
	for _, test := range data {
		logLines += test.Logline + "\n"
	}
	return strings.NewReader(logLines)
}

// Test if LogGrouper setup is correct
func TestNewLogGrouper(t *testing.T) {
	// TODO
}

// Test whole parsing and setup
func TestApacheLogParsing(t *testing.T) {
	testData := ApacheTestData

	// Setup logGrouper with apache profile
	profile := Profiles["apacheAccessLog"]
	analyzer := NewLogGrouper(profile.LogFormat, profile.TimeFormat, Interval)

	analyzer.Analyze([]io.Reader{getReaderFromData(testData)})

	// Check if unparsed lines is zero
	if len(analyzer.UnmatchedLines) != 0 {
		t.Errorf("Unmatched lines is not zero!\n%s", strings.Join(analyzer.UnmatchedLines, "\n"))
	}

	// Check that time slot grouping is correct
	// Dynamicly from testdata
}
