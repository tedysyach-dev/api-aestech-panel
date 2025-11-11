package utils

import (
	"fmt"
	"strings"
	"time"
)

// bentuk rangeDate "2025-01-01|2025-01-03"
func GetStartEndDate(rangeDate string) (int64, int64, error) {
	const layout = "2006-01-02"

	parts := strings.Split(rangeDate, "|")

	var startDateStr, endDateStr string
	switch len(parts) {
	case 1:
		startDateStr = parts[0]
		endDateStr = parts[0]
	case 2:
		startDateStr = parts[0]
		endDateStr = parts[1]
	default:
		return 0, 0, fmt.Errorf("invalid rangeDate format, expected 'YYYY-MM-DD' or 'YYYY-MM-DD|YYYY-MM-DD'")
	}

	// parse start & end
	start, err := time.ParseInLocation(layout, strings.TrimSpace(startDateStr), time.Local)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid start date: %v", err)
	}

	end, err := time.ParseInLocation(layout, strings.TrimSpace(endDateStr), time.Local)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid end date: %v", err)
	}

	// set jam ke awal & akhir hari
	start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.Local)
	end = time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 999_000_000, time.Local)

	return start.UnixMilli(), end.UnixMilli(), nil
}
