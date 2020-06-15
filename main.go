package main

import (
	"fmt"
	"os"

	"github.com/massl123/logGrouper/loggrouper"
	flag "github.com/spf13/pflag"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Usage: logGrouper [options] file1 [file2 file...]\n")
		fmt.Fprintf(os.Stderr, "Use file name \"-\" to read from stdin.\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Group lines in log by time and occurance.\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Copyright (c) 2020 Marcel Freundl <github.com/Massl123>\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Profiles:\n")
		fmt.Fprintf(os.Stderr, "%-40s %-100s %-25s\n", "Name", "Log format", "Time format")
		for _, profile := range loggrouper.Profiles {
			fmt.Fprintf(os.Stderr, "%-40s %-100s %-25s\n", profile.Name, profile.LogFormat, profile.TimeFormat)
		}
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	fVerbose := flag.BoolP("verbose", "v", false, "Verbose output (show unparsed lines)")
	fGroupLimit := flag.IntP("limit", "l", 20, "Limit output per timeslot.")
	fInterval := flag.StringP("interval", "i", "15m", "Interval to group by. Format like 15m (Units supported: ns, us (or µs), ms, s, m, h).")
	fProfile := flag.StringP("profile", "p", "apacheAccessLog-full", "Profile to use. Profile loads LogFormat and TimeFormat. Use LogFormat and TimeFormat parameters to override.")
	fLogFormat := flag.StringP("format", "f", "", "LogFormat regexp in GoLang Format. Match group \"timestamp\" and \"group\" have to exist. See https://golang.org/pkg/regexp/syntax/.")
	fTimeFormat := flag.StringP("timeFormat", "t", "", "Time format for \"timestamp\" match group. Given in GoLang format, see https://golang.org/pkg/time/#Parse.")

	flag.Parse()

	fArgs := flag.Args()

	// Show error if no files are given
	if len(fArgs) == 0 {
		flag.Usage()
		os.Exit(0)
	}

	// Load profile
	profile := loggrouper.Profiles[*fProfile]

	// Override profile values with custom values if specified
	if *fLogFormat != "" {
		profile.LogFormat = *fLogFormat
	}
	if *fTimeFormat != "" {
		profile.LogFormat = *fTimeFormat
	}

	// https://golang.org/pkg/regexp/syntax/
	analyzer := loggrouper.NewLogGrouper(profile.LogFormat, profile.TimeFormat, *fInterval)
	analyzer.AnalyzeFiles(fArgs)
	analyzer.Print(*fGroupLimit, *fVerbose)
}
