package main

import (
	"fmt"
	"time"
)

// Report keeps up with logs and timestamps
type Report struct {
	messages   []string
	timeStamps []time.Time
}

// NewReport creates a report to write to.
func NewReport() *Report {
	r := Report{
		messages:   []string{},
		timeStamps: []time.Time{},
	}
	return &r
}

// Log an entry to the report for later
func (report *Report) Log(msg string) {
	report.timeStamps = append(report.timeStamps, time.Now())
	report.messages = append(report.messages, msg)
}

// GetTime returns our current time formatted
func GetTime() string {
	return time.Now().Format("Mon Jan _2 15:04:05 UTC-01:00 2006")
}

func (report Report) Generate() []string {
	final := make([]string, len(report.messages))
	for i := len(report.messages) - 1; i >= 0; i-- {
		final[len(report.messages)-1-i] = fmt.Sprintf("%s: %s", report.timeStamps[i].Format("Mon Jan _2 15:04:05 UTC-01:00 2006"), report.messages[i])
	}
	return final
}
