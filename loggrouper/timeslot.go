package loggrouper

import (
	"sync"
	"time"
)

// NewTimeslot ...
func NewTimeslot(startTime time.Time) *Timeslot {
	ts := &Timeslot{}
	ts.StartTime = startTime
	ts.LogGroups = make(map[string]*LogGroup)

	return ts
}

// Timeslot ...
type Timeslot struct {
	sync.Mutex
	// Timeslot datetime
	StartTime time.Time
	Count     uint

	LogGroups map[string]*LogGroup
}

func (ts *Timeslot) addLineToLogGroup(logGroupName string) {
	ts.Lock()
	if group, ok := ts.LogGroups[logGroupName]; ok {
		ts.Unlock()
		group.Lock()
		group.Count++
		group.Unlock()
	} else {
		// Group doesn't exists, create new one
		group := &LogGroup{}
		group.Name = logGroupName
		group.Count++
		ts.LogGroups[logGroupName] = group
		ts.Unlock()
	}
}

// LogGroup ...
type LogGroup struct {
	sync.Mutex
	Name  string
	Count uint
}
