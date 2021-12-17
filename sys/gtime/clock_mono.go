package gtime

/**
The monotonic clock is a process start-based clock that is not affected by the system clock,
which may be toggled forward or backward, while the monotonic clock is not.
This clock is used to calculate the exact time difference within the process,
rather than providing the functionality that a wall clock like the year, month and day should provide.
*/

import (
	"time"
)

type (
	MonoClock struct {
		loc *time.Location
		baselineTime time.Time
	}
)

func NewMonoClock() *MonoClock {
	return &MonoClock{loc: time.UTC, baselineTime: EpochBeginTime}
}

func (c *MonoClock) Name() string {
	return "mono"
}

// Now returns the current time from a monotonic clock.
//
// The time returned is based on some arbitrary platform-specific point in the
// past. The time returned is guaranteed to increase monotonically without
// notable jumps, unlike time.Now() from the Go standard library, which may
// jump forward or backward significantly due to system time changes or leap
// seconds.
//
// It's implemented using runtime.nanotime(), which uses CLOCK_MONOTONIC on
// Linux. Note that unlike CLOCK_MONOTONIC_RAW, CLOCK_MONOTONIC is affected
// by time changes. However, time changes never cause clock jumps; instead,
// clock frequency is adjusted slowly.
func (c *MonoClock) Now() time.Time {
	return c.baselineTime.In(c.loc).Add(Uptime())
}

func (c *MonoClock) Sleep(d time.Duration) {
	panic("can't sleep monotonic clock")
}

func (c *MonoClock) Set(tm time.Time) error {
	panic("can't set time for monotonic clock")
}