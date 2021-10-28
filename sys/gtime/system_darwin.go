package gtime

import (
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/sys/gusers"
	"syscall"
	"time"
)

func SetSystemTimeROOT(t time.Time) error {
	var tv syscall.Timeval
	tv.Sec = t.Unix()
	tv.Usec = 0
	if err := syscall.Settimeofday(&tv); err != nil {
		isAdmin, err2 := gusers.IsRunAsAdmin()
		if err2 == nil && !isAdmin {
			return gerrors.Errorf(err.Error() + ", modifying system time requires administrator privileges")
		} else {
			return err
		}
	}
	return nil
}
