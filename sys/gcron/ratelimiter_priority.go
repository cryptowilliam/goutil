package gcron

// Frequency limiter designed to limit web api access frequency.
// For example, **.com http api could be visited 3 times in one second,
// you should create a RateLimiter which duration is (time.second / 3),
// then, MarkAndWaitBlock() before every http request.

// priority: 0->max_uint (high->low)

import (
	"github.com/cryptowilliam/goutil/container/gnum"
	"sync"
	"time"
)

type RateLimiterPriority struct {
	d               time.Duration
	loopSleep       time.Duration
	last            time.Time
	waitingPriority []uint
	mu              sync.RWMutex
}

func NewRateLimiterPriority(d time.Duration) *RateLimiterPriority {
	loopSleep := time.Millisecond * 2
	if d > time.Hour {
		loopSleep = time.Second
	}

	return &RateLimiterPriority{d: d, loopSleep: loopSleep}
}

func (r *RateLimiterPriority) Duration() time.Duration {
	return r.d
}

// mark and wait block
func (r *RateLimiterPriority) TakeWait(priority uint) {
	if r.d <= 0 {
		return
	}

	r.mu.Lock()
	r.waitingPriority = append(r.waitingPriority, priority)
	r.mu.Unlock()

	for {
		time.Sleep(r.loopSleep)

		r.mu.RLock()
		currWaitingPriority := r.waitingPriority
		r.mu.RUnlock()

		if len(currWaitingPriority) > 0 && gnum.MinUintArray(currWaitingPriority) < priority {
			continue
		}

		r.mu.Lock()
		if time.Now().Sub(r.last) >= r.d {
			r.last = time.Now()
			r.waitingPriority = gnum.RemoveUint(r.waitingPriority, priority, 1)
			r.mu.Unlock()
			return
		} else {
			r.mu.Unlock()
			continue
		}
	}
}
