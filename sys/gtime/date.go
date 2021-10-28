package gtime

import (
	"fmt"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/container/gstring"
	"math"
	"time"
)

// There is no time.Location info in Date,
// if there is, Location will be used to compare two date, or convert Date from current Location to another,
// but you can't convert Date's time.Location without hours even minutes,
// so, Location is not included.

// Why not define Date as time.Time and keep most info for it?
// If did that, there will be problem when NewDate(),Cmp(),In(loc *Location),JSON Marshal/Unmarshal.
// If you need time.Location, use time.Time but not Date.

type Date int32

func Today(loc *time.Location) Date {
	return TimeToDate(time.Now(), loc)
}

func Yesterday(loc *time.Location) Date {
	return TimeToDate(Sub(time.Now(), time.Hour*24), loc)
}

func TimeToDate(tm time.Time, loc *time.Location) Date {
	return Date(tm.In(loc).Year()*10000 + int(tm.In(loc).Month())*100 + tm.In(loc).Day())
}

func MinDate(a, b Date) Date {
	if a.IntYYYYMMDD() < b.IntYYYYMMDD() {
		return a
	}
	return b
}

func MaxDate(a, b Date) Date {
	if a.IntYYYYMMDD() > b.IntYYYYMMDD() {
		return a
	}
	return b
}

// check if date is valid
// invalid date example: 2018-2-30
func DateValid(year, month, day int) bool {
	if month <= 0 || month >= 13 {
		return false
	}
	if day <= 0 || day >= 32 {
		return false
	}
	tm := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	if tm.Year() != year || tm.Month() != time.Month(month) || tm.Day() != day {
		return false
	}
	return true
}

func NewDate(year, month, day int) (Date, error) {
	if !DateValid(year, month, day) {
		return ZeroDate, gerrors.Errorf("Invalid date input %d-%d-%d", year, month, day)
	}
	return Date((year * 10000) + (month * 100) + day), nil
}

func NewDatePanic(year, month, day int) Date {
	if !DateValid(year, month, day) {
		panic(gerrors.Errorf("Invalid date input %d-%d-%d", year, month, day))
	}
	return Date((year * 10000) + (month * 100) + day)
}

func (d Date) Year() int {
	return int(d) / 10000
}

func (d Date) Month() time.Month {
	return time.Month(int(math.Abs(float64(d))) / 100 % 100)
}

func (d Date) Day() int {
	return int(math.Abs(float64(d))) % 100
}

func (d Date) Equal(cmp Date) bool {
	return int(d) == int(cmp)
}

func (d Date) Before(cmp Date) bool {
	return int(d) < int(cmp)
}

func (d Date) BeforeEqual(cmp Date) bool {
	return int(d) <= int(cmp)
}

func (d Date) After(cmp Date) bool {
	return int(d) > int(cmp)
}

func (d Date) AfterEqual(cmp Date) bool {
	return int(d) >= int(cmp)
}

func (d Date) IsZero() bool {
	return d.Equal(ZeroDate)
}

// days from unix epoch date
func (d Date) UnixDays() int {
	return int(d.ToTime(0, 0, 0, 0, time.UTC).Sub(EpochBeginDate.ToTime(0, 0, 0, 0, time.UTC)).Hours() / 24)
}

func (d Date) Sub(cmp Date) int {
	return d.UnixDays() - cmp.UnixDays()
}

func (d Date) AddDays(days int) Date {
	return TimeToDate(d.ToTime(0, 0, 0, 0, time.UTC).Add(Day*time.Duration(days)), time.UTC)
}

func (d Date) PreviousDay() Date {
	return d.AddDays(-1)
}

func (d Date) NextDay() Date {
	return d.AddDays(1)
}

func (d Date) String() string {
	return d.StringYYYY_MM_DD()
}

func (d Date) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", d.String())), nil
}

func (d *Date) UnmarshalJSON(b []byte) error {
	s := string(b)
	if len(s) <= 1 {
		return gerrors.Errorf("Invalid json date '%s'", s)
	}
	if s[0] != '"' || s[len(s)-1] != '"' {
		return gerrors.Errorf("Invalid json date '%s'", s)
	}
	s = gstring.RemoveHead(s, 1)
	s = gstring.RemoveTail(s, 1)
	dt, err := ParseDateString(s, true)
	if err != nil {
		*d = ZeroDate
		return err
	}
	*d = dt
	return nil
}

// yyyymmdd
// Notice:
// 如果你写成了fmt.Sprintf("%04d%02d%02d", d.Year, d.YearMonth, d.Day)，编译也能通过
// 但是返回结果却是很大很大的数字，因为它们代表函数地址
func (d Date) StringYYYYMMDD() string {
	return fmt.Sprintf("%04d%02d%02d", d.Year(), d.Month(), d.Day())
}

func (d Date) IntYYYYMMDD() int {
	return int(d)
}

// yyyy-mm-dd
func (d Date) StringYYYY_MM_DD() string {
	return fmt.Sprintf("%04d-%02d-%02d", d.Year(), d.Month(), d.Day())
}

func (d Date) ToTime(hour, minute, sec, nsec int, loc *time.Location) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), hour, minute, sec, nsec, loc)
}

func (d Date) ToTimeLocation(loc *time.Location) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, loc)
}

func (d Date) ToTimeUTC() time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC)
}

type DateRange struct {
	Begin Date
	End   Date
}

func (dr DateRange) String() string {
	if dr.Begin.IsZero() && dr.End.IsZero() {
		return ""
	}
	return dr.Begin.String() + "/" + dr.End.String()
}

func (dr DateRange) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", dr.String())), nil
}

func (dr *DateRange) UnmarshalJSON(b []byte) error {
	// Init
	dr.Begin = ZeroDate
	dr.End = ZeroDate

	// Remove '"'
	s := string(b)
	ErrDefault := gerrors.Errorf("invalid json date range '%s'", s)
	if len(s) <= 1 {
		return ErrDefault
	}
	if s[0] != '"' || s[len(s)-1] != '"' {
		return ErrDefault
	}
	s = gstring.RemoveHead(s, 1)
	s = gstring.RemoveTail(s, 1)

	// Parse
	res, err := ParseDateRangeString(s, true)
	if err != nil {
		return ErrDefault
	}
	*dr = res
	return nil
}

func (dr DateRange) IsZero() bool {
	return dr.Begin.IsZero() && dr.End.IsZero()
}
