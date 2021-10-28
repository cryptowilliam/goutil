package gtime

import (
	"os/exec"
	"time"
)

// 下述实现只能在windows下编译通过，也可以通过syscall实现纯go的版本

// Get retrieves the current system time, either via syscall.Gettimeofday on Linux/Darwin
// or via kernel32.GetSystemTime() on Windows. It then parses the result into a
// standard time.Time struct.
/*
func GetSystemTime() (*time.Time, error) {

	// Gets the system time from the kernel32 API
	st := w32.GetSystemTime()

	// Convert the SYSTEMTIME to time.Time
	t := time.Date(
		int(st.Year),
		time.Month(st.Month),
		int(st.Day),
		int(st.Hour),
		int(st.Minute),
		int(st.Second),
		0, time.UTC)

	return &t, nil

}

// Set sets the current system time, either via syscall.Settimeofday on Linux/Darwin
// or via kernel32.SetSystemtime() on Windows.
func SetSystemTime(input time.Time) error {

	st := &w32.SYSTEMTIME{
		Year:         uint16(input.Year()),
		YearMonth:    uint16(input.Month()),
		DayOfWeek:    uint16(input.Weekday()),
		Day:          uint16(input.Day()),
		Hour:         uint16(input.Hour()),
		Minute:       uint16(input.Minute()),
		Second:       uint16(input.Second()),
		Milliseconds: 0,
	}

	if success := w32.SetSystemTime(st); success != true {
		return gerrors.New("unable to set system time")
	}

	return nil

}*/

func SetSystemTimeROOT(t time.Time) error {
	_, err := exec.Command("CMD", "/C", "DATE", t.Format("2006-01-02")).Output()
	if err != nil {
		return err
	}
	_, err = exec.Command("CMD", "/C", "TIME", t.Format("15:04:05")).Output()
	if err != nil {
		return err
	}
	return nil
}
