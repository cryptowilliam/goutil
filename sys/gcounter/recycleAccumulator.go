package gcounter

// A wrapper for recycle number accumulator, like network package serial number.

import (
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"sync"
)

type RecycleAccumulator struct {
	l   sync.RWMutex
	n   int64
	min int64
	max int64
}

// NewAccumulator returns a new pointer to an accumulator
func NewAccumulator(min, max int64) (*RecycleAccumulator, error) {
	if max < min {
		return nil, gerrors.Errorf("min(%d) is bigger than max(%d)", min, max)
	}
	return &RecycleAccumulator{n: 0, min: min, max: max}, nil
}

// Incr incremenets the accumulator by 1
func (i *RecycleAccumulator) Incr() int64 {
	return i.IncrN(1)
}

// IncrN incremenets the accumulator by N
func (i *RecycleAccumulator) IncrN(N int64) int64 {
	i.l.Lock()
	defer i.l.Unlock()

	if i.min == i.max {
		return i.n
	}

	if (i.n+N) > i.max || (i.n+N < i.min) {
		i.n = i.min + (i.n + N - i.max - 1)
	} else {
		i.n += N
	}
	return i.n
}

func (i *RecycleAccumulator) Get() int64 {
	i.l.RLock()
	defer i.l.RUnlock()

	return i.n
}

func (i *RecycleAccumulator) Set(val int64) int64 {
	if val < i.min || val > i.max {
		return i.Get()
	}

	i.l.Lock()
	defer i.l.Unlock()
	i.n = val
	return i.n
}

// Flush resets the accumulator to min
func (i *RecycleAccumulator) Reset() int64 {
	i.l.Lock()
	defer i.l.Unlock()

	i.n = i.min
	return i.n
}
