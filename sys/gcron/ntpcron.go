package gcron

import (
	"github.com/cryptowilliam/goutil/sys/gtime"
	"time"
)

type NtpCron struct {
	nc          *gtime.NtpClock
	triggerTime time.Time
}

func NewNtpCronONLINE(ntpClock *gtime.NtpClock, triggerTime time.Time) (*NtpCron, error) {
	nc, err := gtime.GetNtpClockONLINE()
	if err != nil {
		return nil, err
	}
	return &NtpCron{nc: nc, triggerTime: triggerTime}, nil
}

func NewNtpCronWithClock(ntpClock *gtime.NtpClock, triggerTime time.Time) (*NtpCron, error) {
	return &NtpCron{nc: ntpClock, triggerTime: triggerTime}, nil
}

func (sc *NtpCron) Wait() {
	for {
		sleepMillis := 100 // default sleep 5000 milliseconds for each loop
		dur := sc.triggerTime.Sub(sc.nc.Now())
		durMillis := gtime.NsecToMillis(dur.Nanoseconds())
		if durMillis < 40 {
			break
		}
		time.Sleep(gtime.MillisToDuration(int64(sleepMillis)))
	}
}
