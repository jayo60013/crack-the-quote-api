package main

import (
	"fmt"
	"time"
)

func FormatDateString() string {
	now := time.Now()
	day := now.Day()
	dateFormat := fmt.Sprintf("January 2%s, 2006", getDaySuffix(day))

	return time.Now().Format(dateFormat)
}

func getDaySuffix(day int) string {
	switch day % 10 {
	case 1:
		return "st"
	case 2:
		return "nd"
	case 3:
		return "rd"
	default:
		return "th"
	}
}
