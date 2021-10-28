package gtime

import (
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"time"
)

type Clock interface {
	Name() string
	Now() time.Time
	Sleep(d time.Duration)
	Set(tm time.Time) error
}

func NewClock(clockName string) (Clock, error) {
	switch clockName {
	case "system":
		return GetSysClock(), nil
	case "ntp":
		return GetNtpClockONLINE()
	case "mock":
		return NewMockClock(time.Now(), time.UTC), nil
	}
	return nil, gerrors.Errorf("unsupported clock name %s", clockName)
}
