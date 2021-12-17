package gtime

import (
	"time"
	_ "unsafe" // required to use //go:linkname
)

/*func Uptime() (uint64, error) {
	if runtime.GOOS == "windows" {
		return 0, gerrors.New(runtime.GOOS + " not supported for now")
	} else {
		return host.Uptime()
	}
}*/

//go:noescape
//go:linkname nanotime runtime.nanotime
func nanotime() int64

func Uptime() time.Duration {
	return NsecToDuration(nanotime())
}

// 不可以用tm.Seconds() 和 tm.Nanoseconds()作为判断标准，
// 因为这两个都只是钟表刻度盘上的零头而已
// 1970-01-01 00:00:00 +0000 UTC
func IsEpochBeginning(tm time.Time) bool {
	return tm.Unix() == 0 && tm.UnixNano() == 0
}

// 0001-01-01 00:00:00 +0000 UTC
func IsZero(tm time.Time) bool {
	return tm.IsZero()
}

/*
func SetSystemTimeROOT(t time.Time) error {
	if runtime.GOOS == "windows" {
		_, err := exec.Command("CMD", "/C", "DATE", t.Format("2006-01-02")).Output()
		if err != nil {
			return err
		}
		_, err = exec.Command("CMD", "/C", "TIME", t.Format("15:04:05")).Output()
		if err != nil {
			return err
		}
		return nil
	} else {
		var tv syscall.Timeval
		tv.Sec = t.Unix()
		tv.Usec = 0
		if err := syscall.Settimeofday(&tv); err != nil {
			isAdmin, err2 := xuser.IsRunAsAdmin()
			if err2 == nil && !isAdmin {
				return gerrors.Errorf(err.Error() + ", modifying system time requires administrator privileges")
			} else {
				return err
			}
		}
		return nil
	}
}*/

// returns current time in milliseconds
func CurrentUnixMillis() uint32 { return uint32(time.Now().UnixNano() / int64(time.Millisecond)) }
