package main

import (
	"fmt"
	"time"
)

type LogData struct {
	projectId   int
	startTime   time.Time
	endTime     time.Time
	duration    time.Duration
	category    string
	description string
}

func (l LogData) calculateDuration() (time.Duration, error) {
	if l.startTime.IsZero() {
		return 0, fmt.Errorf("Start Time is not set")
	}
	if l.endTime.IsZero() {
		return 0, fmt.Errorf("End Time is not set")
	}
	if l.endTime.Before(l.startTime) {
		return 0, fmt.Errorf("End Time is before Start Time")
	}
	return l.endTime.Sub(l.startTime), nil
}
