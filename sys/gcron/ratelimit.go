package gcron

// Frequency limiter designed to limit web api access frequency.
// For example, **.com http api could be visited 3 times in one second,
// you should create a RateLimiter which duration is (time.second / 3),
// then, MarkAndWaitBlock() before every http request.

import (
	"github.com/cryptowilliam/goutil/sys/gtime"
	"sync"
	"time"
)

type RateLimiter struct {
	d      time.Duration
	c      gtime.Clock
	last   time.Time
	lastMu sync.Mutex
}

func NewRateLimiter(d time.Duration) *RateLimiter {
	return NewRateLimiterEx(d, nil)
}

func NewRateLimiterEx(d time.Duration, c gtime.Clock) *RateLimiter {
	/*if d < time.Millisecond {
		return nil, gerrors.Errorf("Unacceptable too small duration %s for frequency limiter", d.String())
	}*/
	if c == nil {
		c = gtime.GetSysClock()
	}
	return &RateLimiter{d: d, c: c}
}

func (r *RateLimiter) Get() time.Duration {
	return r.d
}

func (r *RateLimiter) MarkAndWaitUnblock() bool {
	if r.d <= 0 {
		return true
	}

	r.lastMu.Lock()
	defer r.lastMu.Unlock()
	if r.c.Now().Sub(r.last) >= r.d {
		r.last = r.c.Now()
		return true
	} else {
		return false
	}
}

func (r *RateLimiter) MarkAndWaitBlock() {
	if r.d <= 0 {
		return
	}

	for {
		time.Sleep(time.Millisecond * 2)
		r.lastMu.Lock()
		if r.c.Now().Sub(r.last) >= r.d {
			r.last = r.c.Now()
			r.lastMu.Unlock()
			return
		} else {
			r.lastMu.Unlock()
			continue
		}
	}
}
