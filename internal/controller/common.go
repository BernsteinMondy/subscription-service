package controller

import (
	"fmt"
	"time"
)

const timeFormat = "01-2006"

func parseStartAndEndDate(startDateStr, endDateStr string) (time.Time, time.Time, error) {
	startDate, err := time.Parse(timeFormat, startDateStr)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("start date parse failed: %w", err)
	}

	endDate, err := time.Parse(timeFormat, endDateStr)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("end date parse failed: %w", err)
	}

	if startDate.After(endDate) {
		return time.Time{}, time.Time{}, fmt.Errorf("start date is after end date")
	}

	if endDate.Before(startDate) {
		return time.Time{}, time.Time{}, fmt.Errorf("end date is before start date")
	}

	return startDate, endDate, nil
}
