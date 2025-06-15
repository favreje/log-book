package main

import (
	"time"
)

func parseTimeFromDb(timeStr string) (time.Time, error) {
	format := "2006-01-02 15:04:05 -0700 MST"
	t, err := time.Parse(format, timeStr)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}
