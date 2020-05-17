package loggrouper

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
)

//NewLogAnalyzer ...
func NewLogAnalyzer(files []string, logformat, timeformat, interval string) *LogAnalyzer {
	var err error
	analyzer := &LogAnalyzer{}
	analyzer.Logfiles = files
	analyzer.Loglines = make(chan string, 10000)
	analyzer.TimeSlots = make(map[string]*Timeslot)

	// Set up Logformat
	analyzer.Logformat = regexp.MustCompile(logformat)
	// Ensure needed match names are set
	var neededSubexpNames = [...]string{"timestamp", "group"}
	for _, neededName := range neededSubexpNames {
		found := false
		for _, name := range analyzer.Logformat.SubexpNames() {
			if name == neededName {
				found = true
			}
		}
		if !found {
			fmt.Fprintf(os.Stderr, "Did not find match group name: %s\n", neededName)
			os.Exit(1)
		}
	}

	// Set timeformat
	analyzer.Timeformat = timeformat

	// Set interval
	analyzer.TimeInterval, err = time.ParseDuration(interval)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error with interval: %s\n", err)
		os.Exit(1)
	}

	return analyzer
}

// LogAnalyzer ...
type LogAnalyzer struct {
	sync.Mutex
	// WaitGroup for all filereaders
	wgFileReaders sync.WaitGroup
	// WaitGroup for all LogLine Processors
	wgLineProcessors sync.WaitGroup

	// Path to logfile
	Logfiles []string
	// Contains all raw loglines which are sorted into the right TimeSlots
	Loglines chan string
	// Regex for the format of the loglines
	Logformat *regexp.Regexp
	// Timeformat in GoFormat for "timestamp" loglines match group
	// See https://golang.org/pkg/time/#Parse
	Timeformat string

	// Time Interval used for grouping to TimeSlots
	// See https://golang.org/pkg/time/#ParseDuration
	TimeInterval time.Duration
	// Contains all TimeSlots which are found
	TimeSlots map[string]*Timeslot

	UnmatchedLines []string
}

// Analyze ...
// Blocks until finished
func (analyzer *LogAnalyzer) Analyze() {

	for _, logfile := range analyzer.Logfiles {

		analyzer.wgFileReaders.Add(1)
		go analyzer.readFile(logfile)
	}

	// Wait for all FileReaders to finish before closing channel
	analyzer.wgFileReaders.Wait()
	close(analyzer.Loglines)

	// Block until LineProcessors are finished
	analyzer.wgLineProcessors.Wait()
}

func (analyzer *LogAnalyzer) readFile(logFileName string) {
	// Start processors
	for i := 0; i < 10; i++ {
		analyzer.wgLineProcessors.Add(1)
		go analyzer.lineProcessor()
	}

	// Read file and write to channel
	var logFile *os.File
	// Special case: "-" means opens stdin
	if logFileName == "-" {
		logFile = os.Stdin
	} else {
		var err error
		logFile, err = os.Open(logFileName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening file: %s\n", err)
			os.Exit(1)
		}

	}
	defer logFile.Close()
	scnr := bufio.NewScanner(logFile)
	for scnr.Scan() {
		//fmt.Println(scnr.Text())
		analyzer.Loglines <- scnr.Text()
	}
	analyzer.wgFileReaders.Done()
}

// lineProcessor reads lines from analyzer.Loglines and parses them
func (analyzer *LogAnalyzer) lineProcessor() {
	for line := range analyzer.Loglines {
		// Parse logline to regexp match
		match := analyzer.Logformat.FindStringSubmatch(line)
		result := make(map[string]string)

		// If no match skip this line and increment UnmatchedLines counter
		if match == nil {
			analyzer.Lock()
			analyzer.UnmatchedLines = append(analyzer.UnmatchedLines, line)
			analyzer.Unlock()
			continue
		}

		for i, name := range analyzer.Logformat.SubexpNames() {
			// index 0 is the full match, ignore that and ensure that name is set
			if i != 0 && name != "" {
				result[name] = match[i]
			}
		}

		// Parse timestamp from regexp
		logTime, err := time.Parse(analyzer.Timeformat, result["timestamp"])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error during timestamp parsing: %s\n", err)
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
func (analyzer *LogAnalyzer) addLineToTimeSlot(logTime time.Time, logGroupName string) {
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
		slot := NewTimeslot(logTime)
		analyzer.TimeSlots[logTime.String()] = slot
		analyzer.Unlock()

		// Recursive call to process line
		analyzer.addLineToTimeSlot(logTime, logGroupName)

	}
}

// Print prints the result to stdout
func (analyzer *LogAnalyzer) Print(logGroupLimit int, withUnmatchedLines bool) {
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
		fmt.Printf("%15s - %-15s%27d\n", slot.StartTime.Format("15:04:05"), slot.StartTime.Add(analyzer.TimeInterval).Format("15:04:05"), slot.Count)

		// Sort LogGroups
		var sortedLogGroups []*LogGroup
		for _, group := range slot.LogGroups {
			sortedLogGroups = append(sortedLogGroups, group)
		}

		sort.SliceStable(sortedLogGroups, func(i, j int) bool { return sortedLogGroups[i].Count > sortedLogGroups[j].Count })

		// Print LogGroups
		for i, group := range sortedLogGroups {
			if i > logGroupLimit-1 {
				break
			}
			fmt.Printf("%10s%-40s%10d\n", "", group.Name, group.Count)
		}
	}

	fmt.Printf("\n\n")
	fmt.Printf("Interval: %s, Output timezone: %s, Unmatched Lines: %d\n", analyzer.TimeInterval, currentDay.Format("-0700 MST"), len(analyzer.UnmatchedLines))
	fmt.Printf("Files: %s\n", strings.Join(analyzer.Logfiles, ", "))

	if withUnmatchedLines {
		fmt.Println("Unmatched lines:")
		for _, line := range analyzer.UnmatchedLines {
			fmt.Printf("%s\n", line)
		}
	}
}
