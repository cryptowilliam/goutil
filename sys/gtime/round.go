package gtime

import (
	"github.com/andreyvit/timerounding"
	"time"
)

// RoundEarlier("2017-01-07 09:35:00 +0000 UTC", 20 * time.Minute) -> "2017-01-07 09:20:00 +0000 UTC"
func RoundEarlier(t time.Time, d time.Duration) time.Time {
	return timerounding.Round(t, d)
}

// RoundLater("2017-01-07 09:35:00 +0000 UTC", 20 * time.Minute) -> "2017-01-07 09:40:00 +0000 UTC"
func RoundLater(t time.Time, d time.Duration) time.Time {
	return timerounding.Round(t, d).Add(d)
}
