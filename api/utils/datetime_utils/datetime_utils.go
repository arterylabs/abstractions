package datetime_utils

import (
	"time"

	timeago "github.com/caarlos0/timea.go"
)

func GetTimeAgo(timeValue time.Time) string {
	return timeago.Of(timeValue)
}

func BeginningOfMonth(date time.Time) time.Time {
	return date.AddDate(0, 0, -date.Day()+1)
}

func EndOfMonth(date time.Time) time.Time {
	return date.AddDate(0, 1, -date.Day())
}
