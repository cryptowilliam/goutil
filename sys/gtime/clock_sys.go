package gtime

import (
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"time"
)

var _globalSysClock_ *SysClock

type SysClock struct {
}

func GetSysClock() *SysClock {
	return _globalSysClock_
}

func (sc *SysClock) Name() string {
	return "system"
}

func (sc *SysClock) Now() time.Time {
	return time.Now()
}

func (sc *SysClock) Sleep(d time.Duration) {
	time.Sleep(d)
}

func (sc *SysClock) Set(tm time.Time) error {
	return gerrors.New("system clock doesn't support Set interface")
}
