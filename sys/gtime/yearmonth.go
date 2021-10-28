package gtime

import (
	"fmt"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"strconv"
	"strings"
	"time"
)

//isLearYear 判断年份是否是润年
func IsLeapYear(year int) bool {
	return year%400 == 0 || year%4 == 0 && year%100 != 0
}

//getMonthDays 获取月份对应的月份天数
func GetMonthDays(year, month int) int {
	days := 30
	switch {
	case month == 1 || month == 3 || month == 5 || month == 7 || month == 8 || month == 10 || month == 12:
		days = 31
	case month == 2:
		if IsLeapYear(year) {
			days = 29
		} else {
			days = 28
		}
	default:
		days = 30
	}
	return days
}

func GetMonthDuration(year, month int) time.Duration {
	return time.Duration(GetMonthDays(year, month)) * Day
}

func GetYearDuration(year int) time.Duration {
	if IsLeapYear(year) {
		return time.Duration(366) * Day
	} else {
		return time.Duration(365) * Day
	}
}

func ThisMonthString() string {
	return time.Now().Format("2006-01")
}

type YearMonth int32

// support formats: 201601, 2016-01
func ParseYearMonthString(s string) (YearMonth, error) {
	deferr := gerrors.Errorf("invalid year month string %s", s)

	if ymint, err := strconv.ParseInt(s, 10, 32); err == nil {
		return ParseYearMonthInt(int(ymint))
	} else {
		ss := strings.Split(s, "-")
		if len(ss) != 2 {
			return ZeroYearMonth, deferr
		}
		intyear, err := strconv.ParseInt(ss[0], 10, 32)
		if err != nil {
			return ZeroYearMonth, deferr
		}
		intmonth, err := strconv.ParseInt(ss[1], 10, 32)
		if err != nil || intmonth < 1 || intmonth > 12 {
			return ZeroYearMonth, deferr
		}
		return NewYearMonth(int(intyear), int(intmonth))
	}
}

func ParseYearMonthInt(ym int) (YearMonth, error) {
	year := ym / 100
	month := ym % 100
	if month < 0 { // when ym == -198702, month == -02
		month = -month
	}
	return NewYearMonth(year, month)
}

// NOTICE: 不同的时区可能是不同的月份
func TimeToMonth(tm time.Time) YearMonth {
	if tm.Year() >= 0 {
		return YearMonth((tm.Year() * 100) + int(tm.Month()))
	} else { // -1986, 02
		return YearMonth((tm.Year() * 100) - int(tm.Month()))
	}
}

func NewYearMonth(year, month int) (YearMonth, error) {
	if month < 1 || month > 12 {
		return ZeroYearMonth, gerrors.Errorf("Invalid month input %d-%d", year, month)
	}
	if year >= 0 {
		return YearMonth((year * 100) + month), nil
	} else { // -1986, 02
		return YearMonth((year * 100) - month), nil
	}
}

func (ym *YearMonth) Year() int {
	return int(*ym) / 100
}

func (ym *YearMonth) Month() time.Month {
	tmpym := *ym
	if tmpym < 0 {
		tmpym = -tmpym
	}
	return time.Month(tmpym % 100)
}

func (ym *YearMonth) Equal(cmp YearMonth) bool {
	return *ym == cmp
}

func (ym *YearMonth) IsZero() bool {
	return ym.Equal(ZeroYearMonth)
}

func (ym *YearMonth) Int() int {
	return int(*ym)
}

func (ym *YearMonth) ToTime(day, hour, minute, sec, nsec int, tz time.Location) time.Time {
	return time.Date(ym.Year(), ym.Month(), day, hour, minute, sec, nsec, &tz)
}

func (ym *YearMonth) ToTimeDefault() time.Time {
	return time.Date(ym.Year(), ym.Month(), 1, 0, 0, 0, 0, time.UTC)
}

func (ym YearMonth) String() string {
	return ym.StringYYYY_MM()
}

func (ym YearMonth) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%d", int(ym))), nil
}

func (ym *YearMonth) UnmarshalJSON(b []byte) error {
	s := string(b)
	if len(s) <= 1 {
		return gerrors.Errorf("Invalid json date '%s'", s)
	}
	intym, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		*ym = ZeroYearMonth
		return err
	}
	_ym, err := ParseYearMonthInt(int(intym))
	if err != nil {
		*ym = ZeroYearMonth
		return err
	}
	*ym = _ym
	return nil
}

// yyyymm
func (ym *YearMonth) StringYYYYMM() string {
	return fmt.Sprintf("%04d%02d", ym.Year(), int(ym.Month()))
}

// yyyy-mm
func (ym *YearMonth) StringYYYY_MM() string {
	return fmt.Sprintf("%04d-%02d", ym.Year(), int(ym.Month()))
}
