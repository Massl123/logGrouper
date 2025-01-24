package loggrouper

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"sync"
	"time"
)

// NewLogGrouper creates LogGrouper object for usage
// logformat is the regexp for matching the lines. It has to include match groups "timestamp" and "group"
// timeformat is the golang time format for time grouping
// interval is the interval to group by
func NewLogGrouper(logformat, timeformat, interval string) *LogGrouper {
	var err error
	analyzer := &LogGrouper{}
	analyzer.Loglines = make(chan string, 10000)
	analyzer.TimeSlots = make(map[string]*timeSlot)

	// Set up Logformat
	analyzer.LogFormat = regexp.MustCompile(logformat)
	// Ensure needed match names are set
	var neededSubexpNames = [...]string{"timestamp", "group"}
	for _, neededName := range neededSubexpNames {
		found := false
		for _, name := range analyzer.LogFormat.SubexpNames() {
			if name == neededName {
				found = true
			}
		}
		if !found {
			fmt.Fprintf(os.Stderr, "Did not find match group name: %s\nRegex has to look like this: %s\n", neededName, profileApacheAccessLogFull.LogFormat)
			os.Exit(1)
		}
	}

	// Set timeformat
	analyzer.TimeFormat = timeformat

	// Set interval
	analyzer.TimeInterval, err = time.ParseDuration(interval)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error with interval: %s\n", err)
		os.Exit(1)
	}

	return analyzer
}

// LogGrouper groups logs (or any string with timestamp and some text) by time and a common string
type LogGrouper struct {
	sync.Mutex
	// WaitGroup for all filereaders
	wgFileReaders sync.WaitGroup
	// WaitGroup for all LogLine Processors
	wgLineProcessors sync.WaitGroup

	// Contains all raw Loglines which are sorted into the right TimeSlots
	Loglines chan string
	// Regex for the format of the loglines
	LogFormat *regexp.Regexp
	// TimeFormat in GoFormat for "timestamp" loglines match group
	// See https://golang.org/pkg/time/#Parse
	TimeFormat string

	// Time Interval used for grouping to TimeSlots
	// See https://golang.org/pkg/time/#ParseDuration
	TimeInterval time.Duration
	// Contains all TimeSlots which are found
	TimeSlots map[string]*timeSlot

	UnmatchedLines []string
}

// Analyze is the generic Analyzer using io.Reader
func (analyzer *LogGrouper) Analyze(readers []io.Reader) {
	// Start readers for given io.Readers
	for _, reader := range readers {

		analyzer.wgFileReaders.Add(1)
		go analyzer.read(reader)
	}

	// Wait for all FileReaders to finish before closing channel
	analyzer.wgFileReaders.Wait()
	close(analyzer.Loglines)

	// Block until LineProcessors are finished
	analyzer.wgLineProcessors.Wait()
}

// AnalyzeFiles wraps Analyze for easily using files
func (analyzer *LogGrouper) AnalyzeFiles(logFilePaths []string) {
	var readers []io.Reader
	for _, logFilePath := range logFilePaths {
		var logFile *os.File
		var err error
		if logFilePath == "-" {
			logFile = os.Stdin
		} else {
			logFile, err = os.Open(logFilePath)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening file: %s\n", err)
			os.Exit(1)
		}
		defer logFile.Close()
		readers = append(readers, logFile)
	}

	analyzer.Analyze(readers)
}

// read creates the line processors and continues to read from io.Reader until it is empty
func (analyzer *LogGrouper) read(reader io.Reader) {
	// Start processors
	for i := 0; i < 10; i++ {
		analyzer.wgLineProcessors.Add(1)
		go analyzer.lineProcessor()
	}

	scnr := bufio.NewScanner(reader)
	for scnr.Scan() {
		//fmt.Println(scnr.Text())
		analyzer.Loglines <- scnr.Text()
	}
	analyzer.wgFileReaders.Done()
}

// lineProcessor reads lines from analyzer.Loglines and parses them
// it calls addLineToTimeSlot for grouping
func (analyzer *LogGrouper) lineProcessor() {
	for line := range analyzer.Loglines {
		// Parse logline to regexp match
		match := analyzer.LogFormat.FindStringSubmatch(line)
		result := make(map[string]string)

		// If no match skip this line and increment UnmatchedLines counter
		if match == nil {
			analyzer.Lock()
			analyzer.UnmatchedLines = append(analyzer.UnmatchedLines, line)
			analyzer.Unlock()
			continue
		}

		for i, name := range analyzer.LogFormat.SubexpNames() {
			// index 0 is the full match, ignore that and ensure that name is set
			if i != 0 && name != "" {
				result[name] = match[i]
			}
		}

		// Parse timestamp from regexp
		logTime, err := time.Parse(analyzer.TimeFormat, result["timestamp"])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error during timestamp parsing: %s\nwith line\n%s\n", err, line)
			os.Exit(1)
		}

		// Adjust time to local timezone so different timezones are merged
		logTime = logTime.Local()

		// align given Time to nearest lower Interval
		// e.g. interval 15min
		// timestamp 13:20 -> aligned to 13:15
		logTime = logTime.Truncate(analyzer.TimeInterval)

		// Add logline to TimeSlot
		analyzer.addLineToTimeSlot(logTime, result["group"])

	}
	analyzer.wgLineProcessors.Done()
}

// addLineToTimeSlot find the right TimeSlot for given line and creates TimeSlot if needed
func (analyzer *LogGrouper) addLineToTimeSlot(logTime time.Time, logGroupName string) {
	analyzer.Lock()
	if slot, ok := analyzer.TimeSlots[logTime.String()]; ok {
		// Slot exists, increment counter and match to LogGroup
		analyzer.Unlock()
		slot.Lock()
		slot.Count++
		slot.Unlock()
		slot.addLineToLogGroup(logGroupName)

	} else {
		// Slot doesn't exists, create new one and then add logline
		slot := newTimeslot(logTime)
		analyzer.TimeSlots[logTime.String()] = slot
		analyzer.Unlock()

		// Recursive call to process line
		analyzer.addLineToTimeSlot(logTime, logGroupName)

	}
}

// Print the result nicely formatted to stdout
func (analyzer *LogGrouper) Print(logGroupLimit int, withUnmatchedLines bool) {
	analyzer.Lock()
	defer analyzer.Unlock()

	// Sorting by string is not the most accurate but works for now
	var sortedTimeSlotNames []string
	for slotName := range analyzer.TimeSlots {
		sortedTimeSlotNames = append(sortedTimeSlotNames, slotName)
	}

	sort.Strings(sortedTimeSlotNames)

	var currentDay time.Time

	/*
		Output Layout
		<Date>
			<Time>        <Count>
			    <URI>     <Count>
	*/

	for _, slotName := range sortedTimeSlotNames {
		slot := analyzer.TimeSlots[slotName]
		slotDay := time.Date(slot.StartTime.Year(), slot.StartTime.Month(), slot.StartTime.Day(), 0, 0, 0, 0, slot.StartTime.Location())
		if currentDay != slotDay {
			currentDay = slotDay
			fmt.Printf("%s\n", currentDay.Format("2006-01-02 Mon"))
		}

		// Print time and count in this TimeSlot
		fmt.Printf("%15s - %-15s%87d\n", slot.StartTime.Format("15:04:05"), slot.StartTime.Add(analyzer.TimeInterval).Format("15:04:05"), slot.Count)

		// Sort LogGroups
		var sortedLogGroups []*logGroup
		for _, group := range slot.LogGroups {
			sortedLogGroups = append(sortedLogGroups, group)
		}

		sort.SliceStable(sortedLogGroups, func(i, j int) bool { return sortedLogGroups[i].Count > sortedLogGroups[j].Count })

		// Print LogGroups
		for i, group := range sortedLogGroups {
			if i > logGroupLimit-1 {
				break
			}
			fmt.Printf("%10s%-100s%10d\n", "", group.Name, group.Count)
		}
	}

	fmt.Printf("\n\n")
	fmt.Printf("Interval: %s, Output timezone: %s, Unmatched Lines: %d\n", analyzer.TimeInterval, currentDay.Format("-0700 MST"), len(analyzer.UnmatchedLines))

	if withUnmatchedLines {
		fmt.Println("Unmatched lines:")
		for _, line := range analyzer.UnmatchedLines {
			fmt.Printf("%s\n", line)
		}
	}
}
