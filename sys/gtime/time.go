package gtime

import (
	"time"
)

type (
	Time time.Time
)

func New(t time.Time) Time {
	return Time(t)
}

func Sub(t time.Time, d time.Duration) time.Time {
	return t.Add(0 - d)
}

func StringTimeZone(tm time.Time, tz time.Location) string {
	return tm.In(&tz).String()
}

func AfterEqual(a, b time.Time) bool {
	return a.After(b) || a.Equal(b)
}

func BeforeEqual(a, b time.Time) bool {
	return a.Before(b) || a.Equal(b)
}

func MinTime(a time.Time, b ...time.Time) time.Time {
	min := a
	for _, v := range b {
		if v.Before(min) {
			min = v
		}
	}
	return min
}

func MaxTime(a time.Time, b ...time.Time) time.Time {
	max := a
	for _, v := range b {
		if v.After(max) {
			max = v
		}
	}
	return max
}
