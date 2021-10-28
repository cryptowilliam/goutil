package gprogress

import (
	"sync"
	"time"
)

type timeRange struct {
	begin time.Time
	end   time.Time
	curr  time.Time
}

type Progress struct {
	mu     sync.RWMutex
	ranges map[string]timeRange
}

func (tr timeRange) total() time.Duration {
	return tr.end.Sub(tr.begin)
}

func (tr timeRange) current() time.Duration {
	d := tr.curr.Sub(tr.begin)
	if d > 0 {
		return d
	}
	return time.Duration(0)
}

func NewTimeRangeProgress() *Progress {
	return &Progress{ranges: map[string]timeRange{}}
}

func (p *Progress) Add(name string, begin, end time.Time) {
	p.mu.Lock()
	defer p.mu.Unlock()

	tr := timeRange{begin: begin, end: end}
	p.ranges[name] = tr
}

func (p *Progress) Set(name string, tm time.Time) {
	p.mu.Lock()
	defer p.mu.Unlock()

	item, ok := p.ranges[name]
	if !ok {
		return
	}
	item.curr = tm
	p.ranges[name] = item
}

// 0.1: 10%
func (p *Progress) Get() float64 {
	p.mu.RLock()
	defer p.mu.RUnlock()

	total := 0.0
	for _, item := range p.ranges {
		total += item.total().Minutes() // if not use minutes, total may be overflow
	}

	curr := 0.0
	for _, item := range p.ranges {
		curr += item.current().Minutes()
	}

	return curr / total
}
