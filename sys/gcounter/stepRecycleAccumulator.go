package gcounter

// A wrapper for step & recycle number accumulator
// 多次递增只计做一次，且循环计数.

import (
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"sync"
)

type StepRecycleAccumulator struct {
	l       sync.RWMutex
	min     int64
	max     int64
	curr    int64
	stepLen int64

	tmpStepLeftLastTime int64
}

// NewStepAdd returns a new pointer to an StepRecycleAccumulator
func NewStepRecycleAccumulator(min, max, stepLen int64) (*StepRecycleAccumulator, error) {
	if max < min {
		return nil, gerrors.Errorf("min(%d) is bigger than max(%d)", min, max)
	}
	if stepLen <= 0 {
		return nil, gerrors.Errorf("Invalid stepLen(%d)", stepLen)
	}
	return &StepRecycleAccumulator{min: min, max: max, curr: min, stepLen: stepLen}, nil
}

// Incr incremenets the StepRecycleAccumulator by 1
func (s *StepRecycleAccumulator) Incr() int64 {
	return s.IncrN(1)
}

// IncrN incremenets the StepRecycleAccumulator by N
func (s *StepRecycleAccumulator) IncrN(N int64) int64 {
	s.l.Lock()
	defer s.l.Unlock()

	if s.min == s.max {
		return s.curr
	}

	step := (s.tmpStepLeftLastTime + N) / s.stepLen
	s.tmpStepLeftLastTime = (s.tmpStepLeftLastTime + N) % s.stepLen

	if (s.curr+step) > s.max || (s.curr+step < s.min) {
		s.curr = s.min + (s.curr + step - s.max - 1)
	} else {
		s.curr += step
	}

	return s.curr
}

func (s *StepRecycleAccumulator) Get() int64 {
	s.l.RLock()
	defer s.l.RUnlock()

	return s.curr
}
