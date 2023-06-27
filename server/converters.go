package server

import (
	"strconv"
	"time"
)

func convertToFloat(num string) float64 {
	f, _ := strconv.ParseFloat(num, 64)
	return f
}

func convertToTime(date string) time.Time {
	t, _ := time.Parse("2006-01-02", date)
	return t
}
