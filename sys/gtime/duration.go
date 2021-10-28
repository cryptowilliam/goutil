package gtime

import (
	"fmt"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/container/gnum"
	"github.com/cryptowilliam/goutil/container/gstring"
	"github.com/hako/durafmt"
	"math"
	"strconv"
	"strings"
	"time"
)

const (
	Day     = 24 * time.Hour
	Week    = 7 * Day
	Month30 = 30 * Day
	Year365 = 365 * Day
)

// 计算两个时间直接的间隔天数，不足一天的算一天称为biggerDays，不足一天不计算称为smallerDays
func DaysBetween(since, to time.Time) (exactDays float64, biggerDays, smallerDays int) {
	days := to.Sub(since).Hours() / float64(24)
	return days, int(math.Ceil(days)), int(math.Floor(days))
}

func StringHumanReadable(duration time.Duration) string {
	return durafmt.Parse(duration).String()
}

type HumanDuration time.Duration

func ToHumanDuration(d time.Duration) HumanDuration {
	return HumanDuration(d)
}

func ParseHumanDuration(s string) (*HumanDuration, error) {
	for {
		s = strings.Replace(s, "  ", " ", -1)
		if !strings.Contains(s, "  ") {
			break
		}
	}
	s = strings.ToLower(s)
	ss := strings.Split(s, " ")
	if len(ss) == 0 {
		return nil, gerrors.Errorf("invalid duration (%s)", s)
	}
	if gnum.IsOddInt64(int64(len(ss))) {
		return nil, gerrors.Errorf("invalid duration (%s)", s)
	}

	// parse
	type valUnit struct {
		val  int64
		unit time.Duration
	}
	var vus []valUnit
	err := error(nil)
	for i := 1; i < len(ss); i += 2 {
		item := valUnit{}
		item.val, err = strconv.ParseInt(ss[i-1], 10, 64)
		if err != nil {
			return nil, err
		}
		switch ss[i] {
		case "year":
			item.unit = Year365
		case "month":
			item.unit = Month30
		case "week":
			item.unit = Week
		case "day":
			item.unit = Day
		case "hour":
			item.unit = time.Hour
		case "minute":
			item.unit = time.Minute
		case "second":
			item.unit = time.Second
		case "millisecond":
			item.unit = time.Millisecond
		case "microsecond":
			item.unit = time.Microsecond
		case "years":
			item.unit = Year365
		case "months":
			item.unit = Month30
		case "weeks":
			item.unit = Week
		case "days":
			item.unit = Day
		case "hours":
			item.unit = time.Hour
		case "minutes":
			item.unit = time.Minute
		case "seconds":
			item.unit = time.Second
		case "milliseconds":
			item.unit = time.Millisecond
		case "microseconds":
			item.unit = time.Microsecond
		default:
			err = gerrors.Errorf("invalid duration (%s)", s)
		}

		if err != nil {
			return nil, err
		}
		vus = append(vus, item)
	}

	dura := time.Duration(0)
	for _, v := range vus {
		dura += v.unit * time.Duration(v.val)
	}
	r := HumanDuration(dura)
	return &r, nil
}

func (d HumanDuration) ToDuration() time.Duration {
	return time.Duration(d)
}

func (d HumanDuration) String() string {
	return StringHumanReadable(d.ToDuration())
}

func (d HumanDuration) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", d.String())), nil
}

func (d *HumanDuration) UnmarshalJSON(b []byte) error {
	s := string(b)
	if len(s) <= 1 {
		return gerrors.Errorf("Invalid json HumanDuration '%s'", s)
	}
	if s[0] != '"' || s[len(s)-1] != '"' {
		return gerrors.Errorf("Invalid json HumanDuration '%s'", s)
	}
	s = gstring.RemoveHead(s, 1)
	s = gstring.RemoveTail(s, 1)
	dura, err := ParseHumanDuration(s)
	if err != nil {
		*d = HumanDuration(time.Duration(0))
		return err
	}
	*d = HumanDuration(*dura)
	return nil
}
