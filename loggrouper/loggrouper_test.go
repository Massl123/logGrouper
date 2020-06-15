package loggrouper

import (
	"io"
	"regexp"
	"strings"
	"testing"
	"time"
)

// Test if regexp in Profiles are ok
func TestProfilesSyntax(t *testing.T) {
	for _, profile := range Profiles {
		_, err := regexp.Compile(profile.LogFormat)
		if err != nil {
			t.Errorf("Profile %q LogFormat not compiling: %s\n", profile.Name, err)
		}
	}
}

// Test whole parsing and setup
func TestApacheLogParsing(t *testing.T) {

	// Check that time slot and string grouping is correct
	// Dynamicly from testdata
	do := func(profileName, interval, logLine, timeGroup, logGroup string) {
		t.Run(profileName, func(t *testing.T) {
			t.Parallel()
			t.Log("Testing Profile: ", profileName)
			t.Log("Testing LogLine: ", logLine)
			// Setup analyzer with profile and intervall
			profile := Profiles[profileName]
			analyzer := NewLogGrouper(profile.LogFormat, profile.TimeFormat, interval)
			analyzer.Analyze([]io.Reader{strings.NewReader(logLine)})

			// Adjust timeGroup to local timezone
			tg, _ := time.Parse("2006-01-02 15:04:05 -0700 MST", timeGroup)
			timeGroup = tg.Format("2006-01-02 15:04:05 -0700 MST")

			// Check if unparsed lines is zero
			if len(analyzer.UnmatchedLines) != 0 {
				t.Errorf("Unmatched lines is not zero! Unmatched lines:\n%s\n", strings.Join(analyzer.UnmatchedLines, "\n"))
			}

			t.Log("Found these TimeSlots:")
			for k := range analyzer.TimeSlots {
				t.Logf("%q\n", k)
			}
			// Check if timeGrouping worked correctly
			if timeSlot, ok := analyzer.TimeSlots[timeGroup]; !ok {
				t.Errorf("LogLine did not match to right TimeSlot, expected slot %q\n", timeGroup)
			} else {
				t.Log("Found these LogGroups:")
				for k := range timeSlot.LogGroups {
					t.Logf("%q\n", k)
				}
				// timeSlot exists, check group
				if _, ok := timeSlot.LogGroups[logGroup]; !ok {
					t.Errorf("LogLine did not match to right LogGroup, expected group %q\n", logGroup)
				}
			}
		})
	}

	// As all logLines timezones are converted to local timezone the given group time may change
	// For the tests to work in different timezones add this string to every timestamp string
	tz := time.Now().Format("-0700 MST")

	// GET and single slash
	do("apacheAccessLog-full", "15m", `127.0.0.1 - frank [10/Oct/2000:13:55:36 -0700] "GET /apache_pb.gif HTTP/1.0" 200 2326 "http://www.example.com/start.html" "Mozilla/4.08 [en] (Win98; I ;Nav)"`, "2000-10-10 22:45:00 "+tz, "GET /apache_pb.gif")
	do("apacheAccessLog-first-group", "15m", `127.0.0.1 - frank [10/Oct/2000:13:55:36 -0700] "GET /apache_pb.gif HTTP/1.0" 200 2326 "http://www.example.com/start.html" "Mozilla/4.08 [en] (Win98; I ;Nav)"`, "2000-10-10 22:45:00 "+tz, "GET /apache_pb.gif ")

	// POST and single slash
	do("apacheAccessLog-full", "15m", `127.0.0.1 - frank [10/Oct/2000:13:55:36 -0700] "POST /apache_pb.gif HTTP/2.0" 200 2326 "http://www.example.com/start.html" "Mozilla/4.08 [en] (Win98; I ;Nav)"`, "2000-10-10 22:45:00 "+tz, "POST /apache_pb.gif")
	do("apacheAccessLog-first-group", "15m", `127.0.0.1 - frank [10/Oct/2000:13:55:36 -0700] "POST /apache_pb.gif HTTP/2.0" 200 2326 "http://www.example.com/start.html" "Mozilla/4.08 [en] (Win98; I ;Nav)"`, "2000-10-10 22:45:00 "+tz, "POST /apache_pb.gif ")

	// GET and two slashes
	do("apacheAccessLog-full", "15m", `127.0.0.1 - frank [10/Oct/2000:13:55:36 -0700] "GET /test123/apache_pb.gif HTTP/1.0" 200 2326 "http://www.example.com/start.html" "Mozilla/4.08 [en] (Win98; I ;Nav)"`, "2000-10-10 22:45:00 "+tz, "GET /test123/apache_pb.gif")
	do("apacheAccessLog-first-group", "15m", `127.0.0.1 - frank [10/Oct/2000:13:55:36 -0700] "GET /test123/apache_pb.gif HTTP/1.0" 200 2326 "http://www.example.com/start.html" "Mozilla/4.08 [en] (Win98; I ;Nav)"`, "2000-10-10 22:45:00 "+tz, "GET /test123/")

	// IPv6
	do("apacheAccessLog-full", "15m", `fe80:880:197d:0:250:fcfb:fe23:3279 - frank [10/Oct/2000:13:55:36 -0700] "GET /apache_pb.gif HTTP/1.0" 200 2326 "http://www.example.com/start.html" "Mozilla/4.08 [en] (Win98; I ;Nav)"`, "2000-10-10 22:45:00 "+tz, "GET /apache_pb.gif")
	do("apacheAccessLog-first-group", "15m", `fe80:880:197d:0:250:fcfb:fe23:3279 - frank [10/Oct/2000:13:55:36 -0700] "GET /apache_pb.gif HTTP/1.0" 200 2326 "http://www.example.com/start.html" "Mozilla/4.08 [en] (Win98; I ;Nav)"`, "2000-10-10 22:45:00 "+tz, "GET /apache_pb.gif ")
}
