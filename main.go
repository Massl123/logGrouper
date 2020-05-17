package main

import (
	"fmt"
	"os"

	"github.com/massl123/logGrouper/loggrouper"
	flag "github.com/spf13/pflag"
)

/*
	TODO
	- possiblity to align interval to e.g. 00:00
		- problem is truncate, which alignes to UTC time
		- other timezones shift this alignment then
	- limit group output count

*/

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Usage: logGrouper [options] [file1 file2 file...]\n")
		fmt.Fprintf(os.Stderr, "Use file name \"-\" or give no file name to read from stdin.\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Group lines in log by time and occurance.\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Copyright (c) 2020 Marcel Freundl <github.com/Massl123>\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	fVerbose := flag.BoolP("verbose", "v", false, "Verbose output (show unparsed lines)")
	fGroupLimit := flag.IntP("limit", "l", 20, "Limit output per timeslot.")
	fInterval := flag.StringP("interval", "i", "15m", "Interval to group by. Format like 15m (Units supported: ns, us (or µs), ms, s, m, h).")
	fLogFormat := flag.StringP("format", "f", `^.*\[(?P<timestamp>.+)\].*"(?P<group>.+ [/\*].*?[/?\ ]).*".*$`, "LogFormat regexp in GoLang Format. Match group \"timestamp\" and \"group\" have to exist. See https://golang.org/pkg/regexp/syntax/")
	fTimeFormat := flag.StringP("timeFormat", "t", "2/Jan/2006:15:04:05 -0700", "Time format for \"timestamp\" match group. Given in GoLang format, see https://golang.org/pkg/time/#Parse")

	flag.Parse()

	fArgs := flag.Args()

	// Read from stdin if no file is given
	if len(fArgs) == 0 {
		fArgs = []string{"-"}
	}

	// https://golang.org/pkg/regexp/syntax/
	analyzer := loggrouper.NewLogAnalyzer(fArgs, *fLogFormat, *fTimeFormat, *fInterval)
	analyzer.Analyze()
	analyzer.Print(*fGroupLimit, *fVerbose)
}
