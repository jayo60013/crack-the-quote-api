package utils

import (
	"fmt"
	"time"
)

func FormatTodaysDate() string {
	now := time.Now()
	day := now.Day()
	dateFormat := fmt.Sprintf("January 2%s, 2006", getDaySuffix(day))

	return time.Now().Format(dateFormat)
}

func getDaySuffix(day int) string {
	switch day {
	case 1, 21, 31:
		return "st"
	case 2, 22:
		return "nd"
	case 3, 23:
		return "rd"
	default:
		return "th"
	}
}
