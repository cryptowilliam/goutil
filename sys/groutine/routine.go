package groutine

import (
	"sync"
	"time"
)

type Routine struct {
	wg sync.WaitGroup

	timeout   *time.Duration
	timeoutMu sync.RWMutex

	exitMsg chan bool

	isOver   bool
	isOverMu sync.RWMutex
}

type RoutineFunc func(args ...interface{})

func RunOnce(function RoutineFunc, args ...interface{}) *Routine {
	r := Routine{exitMsg: nil}
	r.wg.Add(1)

	go func(args ...interface{}) {
		defer r.wg.Done()
		defer func() {
			r.isOverMu.Lock()
			r.isOver = true
			r.isOverMu.Unlock()
		}()
		// One time call.
		function(args...)
	}(args...)

	return &r
}

// Create new routine.
func RunLoop(function RoutineFunc, loopInterval time.Duration, timeout *time.Duration, args ...interface{}) *Routine {
	r := Routine{}
	r.exitMsg = make(chan bool, 1)
	if timeout != nil && *timeout > 0 {
		r.timeout = timeout
	}
	r.wg.Add(1)

	go func(args ...interface{}) {
		defer r.wg.Done()
		defer func() {
			r.isOverMu.Lock()
			r.isOver = true
			r.isOverMu.Unlock()
		}()
		begin := time.Now()

		// Loop call.
		for {
			// Check exit message.
			select {
			case <-r.exitMsg:
				return
			default:
			}

			// Timeout in for loop.
			if r.timeout != nil && time.Now().Sub(begin) > *r.timeout {
				return
			}

			// Do job.
			function(args...)

			// Limit the loop interval.
			if loopInterval.Seconds() > 0 {
				time.Sleep(loopInterval)
			}
		}

	}(args...)

	return &r
}

// Is alive or not.
func (r *Routine) IsAlive() bool {
	r.isOverMu.RLock()
	defer r.isOverMu.RUnlock()

	return r.isOver
}

// Wait until go routine end.
func (r *Routine) Wait() {
	r.wg.Wait()
}

func (r *Routine) Close() {
	// SendEth exit message.
	if r.exitMsg != nil {
		r.exitMsg <- true
	}
}
