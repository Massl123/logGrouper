package loggrouper

import (
	"sync"
	"time"
)

// newTimeslot creates a initialized timeSlot object
func newTimeslot(startTime time.Time) *timeSlot {
	ts := &timeSlot{}
	ts.StartTime = startTime
	ts.LogGroups = make(map[string]*logGroup)

	return ts
}

// timeSlot is used for grouping on time
type timeSlot struct {
	sync.Mutex
	// Timeslot datetime
	StartTime time.Time
	Count     uint

	LogGroups map[string]*logGroup
}

func (ts *timeSlot) addLineToLogGroup(logGroupName string) {
	ts.Lock()
	if group, ok := ts.LogGroups[logGroupName]; ok {
		ts.Unlock()
		group.Lock()
		group.Count++
		group.Unlock()
	} else {
		// Group doesn't exists, create new one
		group := &logGroup{}
		group.Name = logGroupName
		group.Count++
		ts.LogGroups[logGroupName] = group
		ts.Unlock()
	}
}

// logGroup is the group in a timeSlot group
type logGroup struct {
	sync.Mutex
	Name  string
	Count uint
}
