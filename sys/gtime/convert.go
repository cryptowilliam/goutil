package gtime

import (
	"github.com/hako/durafmt"
	"time"
)

func MillisToNsec(millis int64) int64 {
	return millis * 1000000
}

func MillisToMicros(millis int64) int64 {
	return millis * 1000
}

func NsecToMillis(nsec int64) int64 {
	return nsec / 1000000
}

func NsecToSec(nsec int64) int64 {
	return nsec / 1000000000
}

func MicrosToMillis(micros int64) int64 {
	return micros / 1000
}

func NsecToDuration(nsec int64) time.Duration {
	return time.Duration(nsec)
}

func MillisToDuration(millis int64) time.Duration {
	return time.Duration(MillisToNsec(millis))
}

func EpochSecToTime(sec int64) time.Time {
	return time.Unix(sec, 0)
}

func EpochMillisToTime(millis int64) time.Time {
	return time.Unix(0, MillisToNsec(millis))
}

func EpochNsecToTime(nsec int64) time.Time {
	return time.Unix(0, nsec)
}

// Convert int64 type duration to duration type
// MulDuration(3, time.Second) -> 3 seconds duration
func MulDuration(size int64, unit time.Duration) time.Duration {
	return unit * time.Duration(size)
}

func TimeToEpochSec(t time.Time) int64 {
	return t.Unix()
}

func TimeToEpochMillis(t time.Time) int64 {
	return NsecToMillis(t.UnixNano())
}

func TimeToEpochNsec(t time.Time) int64 {
	return t.UnixNano()
}

// tm.UnixNano()的逆运算
func UnixNanoToTime(un int64, loc *time.Location) time.Time {
	if loc == nil {
		loc = time.UTC
	}
	return time.Unix(0, un).In(loc)
}

// time.Duaration.String() = "354h22m3.24s"
// PrettyFormat() = "2 weeks 18 hours 22 minutes 3 seconds"
func PrettyFormat(d time.Duration) string {
	return durafmt.Parse(d).String()
}
