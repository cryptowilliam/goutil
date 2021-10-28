package gsysinfo

import (
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/shirou/gopsutil/host"
	"runtime"
	"time"
)

// TODO: test under Windows/linux

// NOTICE: accurate to second only.
func UpDuration() (time.Duration, error) {
	durationSeconds := uint64(0)
	err := error(nil)

	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		durationSeconds, err = host.Uptime()
	} else if runtime.GOOS == "windows" {
		durationSeconds, err = host.BootTime()
	} else {
		return 0, gerrors.Errorf("current system(%s) not implemented", runtime.GOOS)
	}
	if err != nil {
		return 0, err
	}
	return time.Duration(durationSeconds) * time.Second, nil
}

func UpTime() (time.Time, error) {
	unixSeconds := uint64(0)
	err := error(nil)

	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		unixSeconds, err = host.BootTime()
	} else if runtime.GOOS == "windows" {
		unixSeconds, err = host.Uptime()
	} else {
		return time.Time{}, gerrors.Errorf("current system(%s) not implemented", runtime.GOOS)
	}
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(int64(unixSeconds), 0), nil // epoch seconds to time
}
