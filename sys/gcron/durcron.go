package gcron

import (
	"github.com/cryptowilliam/goutil/sys/gtime"
	"sync"
	"time"
)

// Equal duration cron, it's concurrent-safe.
type DurCron struct {
	// User given origin time for cron calc.
	origin time.Time

	// User given duration for cron.
	duration time.Duration

	// Last return trigger time value.
	lastReturn time.Time

	// Mutex for DurCron, avoid repeat return same trigger time in the case of concurrency.
	mu sync.Mutex

	// true: return true in first call Check, false: returns true only duration is reached
	returnTrueFirstCheck bool

	// Flag about whether it is the first time to call Check
	fisrtCheck bool
}

// NewDurCron creates new DurCron object.
// The origin is calc origin time, if it's nil, time.Now() used as origin.
// The d is cron interval.
func NewDurCron(origin *time.Time, returnTrueFirstCheck bool, d time.Duration) *DurCron {
	if origin == nil {
		now := gtime.RoundEarlier(time.Now(), d)
		origin = &now
	}
	c := DurCron{
		origin:               *origin,
		duration:             d,
		lastReturn:           *origin,
		returnTrueFirstCheck: returnTrueFirstCheck,
		fisrtCheck:           true,
	}
	return &c
}

// CheckNowUnblock returns true when now is a trigger point, otherwise false.
// Backlog trigger point history will be ignored, it just return latest trigger point.
func (c *DurCron) CheckNowUnblock() (triggerCount int, trigger bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.fisrtCheck {
		c.fisrtCheck = false
		if c.returnTrueFirstCheck {
			return 1, true
		}
	}

	// Calc float64 times and ceil it to int
	count := int(time.Now().Sub(c.lastReturn) / c.duration)
	//fmt.Println(count, time.Now(), c.lastReturn, time.Now().Sub(c.lastReturn), c.duration)

	if count <= 0 {
		return 0, false
	} else {
		c.lastReturn = c.lastReturn.Add(gtime.MulDuration(int64(count), c.duration))
		return count, true
	}
}
