package gtime

// A clock who is independent with system clock and sync with NTP server.
// Change system clock need ROOT, but NtpClock doesn't.

import (
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"math"
	"sync"
	"time"
)

// github.com/hlandau/degoutils/clock

type NtpClock struct {
	diff     time.Duration
	diffRwmu sync.RWMutex
}

var _globalNtpClock_ *NtpClock

func GetNtpClockONLINE() (*NtpClock, error) {
	ntpTime, err := GetNetTimeInLocalONLINE()
	if err != nil {
		return nil, err
	}

	nc := NtpClock{
		diff: ntpTime.Sub(time.Now()),
	}
	return &nc, nil
}

func GetNtpSyncClockONLINE() (*NtpClock, error) {
	ntpTime, err := GetNetTimeInLocalONLINE()
	if err != nil {
		return nil, err
	}

	nc := NtpClock{
		diff: ntpTime.Sub(time.Now()),
	}
	go nc.sync()
	return &nc, nil
}

func (nc *NtpClock) sync() {
	for {
		time.Sleep(time.Second * 5)
		ntpTime, err := GetNetTimeInLocalONLINE()
		if err != nil {
			continue
		}
		nc.diffRwmu.Lock()
		nc.diff = ntpTime.Sub(time.Now())
		nc.diffRwmu.Unlock()
	}
}

func (nc *NtpClock) Name() string {
	return "ntp"
}

func (nc *NtpClock) Now() time.Time {
	nc.diffRwmu.RLock()
	diff := nc.diff
	nc.diffRwmu.RUnlock()
	return time.Now().Add(diff)
}

func (nc *NtpClock) Sleep(d time.Duration) {
	time.Sleep(d)
}

func (nc *NtpClock) IsMock() bool {
	return false
}

func (nc *NtpClock) Set(tm time.Time) error {
	return gerrors.New("ntp clock doesn't support Set interface")
}

func VerifySystemTimeWithNtpOL(intervalAllowed time.Duration) error {
	ntpTime, err := GetNetTimeInLocalONLINE()
	if err != nil {
		return err
	}
	now := time.Now()
	diff := ntpTime.Sub(now)

	if math.Abs(float64(diff.Nanoseconds())) > math.Abs(float64(intervalAllowed.Nanoseconds())) {
		return gerrors.Errorf("system time(%s) different from ntp time(%s) more than %s", now.String(), ntpTime.String(), intervalAllowed.String())
	}
	return nil
}
