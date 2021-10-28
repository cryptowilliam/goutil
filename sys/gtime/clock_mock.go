package gtime

import (
	mc "github.com/benbjohnson/clock"
	"time"
)

type (
	MockClock struct {
		in       *mc.Mock
		location *time.Location
	}
)

func NewMockClock(begin time.Time, location *time.Location) *MockClock {
	r := &MockClock{}
	r.in = mc.NewMock()
	r.in.Set(begin)
	if location == nil {
		r.location = time.Now().Location()
	} else {
		r.location = location
	}
	return r
}

func (m *MockClock) Set(t time.Time) error {
	m.in.Set(t)
	return nil
}

func (m *MockClock) SetLocation(location *time.Location) {
	if location != nil {
		m.location = location
	}
}

func (m *MockClock) MockAdd(d time.Duration) {
	m.in.Add(d)
}

func (m *MockClock) Name() string {
	return "mock"
}

func (m *MockClock) Now() time.Time {
	return m.in.Now().In(m.location)
}

func (m *MockClock) Sleep(d time.Duration) {
	m.in.Sleep(d)
}
