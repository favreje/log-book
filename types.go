package main

import (
	"errors"
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
		return 0, ErrStartTimeNotSet
	}
	if l.endTime.IsZero() {
		return 0, ErrEndTimeNotSet
	}
	if l.endTime.Before(l.startTime) {
		return 0, ErrEndTimeBeforeStartTime
	}
	return l.endTime.Sub(l.startTime), nil
}

type InputState struct {
	dateEntered      bool
	startTimeEntered bool
	endTimeEntered   bool
	baseDate         time.Time
	statusMsg        string
}

type Boundary string

const (
	Start Boundary = "Start"
	End   Boundary = "End"
)

// Define sentinel errors
var (
	ErrStartTimeNotSet        = errors.New("Start Time is not set")
	ErrEndTimeNotSet          = errors.New("End Time is not set")
	ErrEndTimeBeforeStartTime = errors.New("End Time is before Start Time")
)
