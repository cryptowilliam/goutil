package gtime

import (
	"time"
)

type TimeMap struct {
	vals map[int64]bool
}

func NewTimeMap(times []time.Time) *TimeMap {
	r := new(TimeMap)
	r.vals = make(map[int64]bool)

	for _, v := range times {
		r.Add(v)
	}

	return r
}

func (m *TimeMap) Add(tm time.Time) {
	un := tm.UnixNano()
	m.vals[un] = true
}

func (m *TimeMap) Exist(tm time.Time) bool {
	ok, _ := m.vals[tm.UnixNano()]
	return ok
}

func (m *TimeMap) Export(loc *time.Location) []time.Time {
	if loc == nil {
		loc = time.UTC
	}

	var r []time.Time
	for v := range m.vals {
		r = append(r, EpochNsecToTime(v).In(loc))
	}
	return r
}
