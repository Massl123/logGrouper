package main

import (
	"github.com/massl123/logGrouper/loggrouper"
)

/*
	TODO
	- possiblity to align interval to e.g. 00:00
*/

func main() {
	// https://golang.org/pkg/regexp/syntax/
	analyzer := loggrouper.NewLogAnalyzer([]string{"demo-access_log"}, `^.+ .+ .+ \[(?P<timestamp>.+)\] ".+ (?P<group>/.*?[/?\ ]).*" .+$`, "2/Jan/2006:15:04:05 -0700", "12h")
	analyzer.Analyze()
	analyzer.Print(true)
}
