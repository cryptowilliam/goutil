package gtime

import (
	"github.com/bcampbell/fuzzytime"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/container/gstring"
	"github.com/tkuchiki/parsetime"
	"strings"
	"time"
)

func fixDateTime(s string) string {
	var tmp string

	ymdChs := []string{
		"年",
		"月",
		"日",
	}
	hmsChs := []string{
		"时",
		"分",
		"秒",
	}
	weeksChs := []string{
		"星期一",
		"星期二",
		"星期三",
		"星期四",
		"星期五",
		"星期六",
		"星期天",
		"星期日",
	}
	monshort := []string{
		"jan",
		"feb",
		"mar",
		"apr",
		"may",
		"june",
		"july",
		"aug",
		"sept",
		"oct",
		"nov",
		"dec",
	}
	monlong := []string{
		"january",
		"february",
		"march",
		"april",
		"may",
		"june",
		"july",
		"agust",
		"september",
		"october",
		"november",
		"december",
	}

	s = strings.ToLower(s)
	s = strings.Replace(s, "北京时间", "bjt", -1) // "北京时间" -> "bjt"
	for _, ymd := range ymdChs {              // "2016年 12月 20日" or "2016年12月20日" -> "2016-12-20"
		if ymd == "日" {
			s = strings.Replace(s, ymd, " ", -1)
		} else {
			s = strings.Replace(s, ymd+" ", "-", -1)
			s = strings.Replace(s, ymd, "-", -1)
		}
	}
	for _, hms := range hmsChs { // "14时 12分 20秒" or "14时12分20秒" -> "14:12:20"
		if hms == "秒" {
			s = strings.Replace(s, hms+" ", " ", -1)
		} else {
			s = strings.Replace(s, hms+" ", ":", -1)
			s = strings.Replace(s, hms, ":", -1)
		}
	}
	for _, weekday := range weeksChs { // "周一," or "周一" -> ""
		s = strings.Replace(s, weekday+",", "", -1)
		s = strings.Replace(s, weekday, "", -1)
	}
	s = strings.Replace(s, "上午", "", -1) // "上午9:30" -> "9:30 am"
	if strings.Contains(s, "下午") {
		s = strings.Replace(s, "下午", "", -1)
		s += " pm"
	}
	s = strings.Replace(s, "下午", "pm", -1)                    // "下午9:30" -> "9:30 pm"
	tmp, err := gstring.ReplaceWithTags(s, "(", ")", " ", -1) // "12月20日(二)13:06" -> "12月20日 13:06"
	if err == nil {
		s = tmp
	}
	// 把英文月份缩写前后补上空格, 但是后面带","和"."的不用补空格, 可以识别
	for i, mon := range monshort {
		s = strings.Replace(s, mon+".", " "+mon+" ", -1)
		s = strings.Replace(s, mon+",", " "+mon+" ", -1)
		// "20dec2016" -> "20 dec 2016"
		if strings.Index(s, mon) >= 0 && strings.Index(s, monlong[i]) < 0 {
			s = strings.Replace(s, mon, " "+mon+" ", 1)
		}
	}
	s = strings.Replace(s, "a.m.", "am", 1)
	s = strings.Replace(s, "p.m.", "pm", 1)
	if strings.Count(s, ".") == 2 {
		s = strings.Replace(s, ".", "-", -1)
	}

	return s
}

// 3rd choice package: github.com/olebedev/when
// Parse human readable date time string to machine-oriented time - unix timestamp
func ParseDatetimeStringFuzz(datetimeString string) (time.Time, error) {
	if len(datetimeString) == 0 {
		return time.Time{}, gerrors.New("empty time string")
	}
	if gstring.CountDigit(datetimeString) < 2 {
		return time.Time{}, gerrors.New("invalid time string")
	}

	datetimeString = fixDateTime(datetimeString)

	// this package can recognize nanoseconds
	parser, _ := parsetime.NewParseTime()
	tm, err := parser.Parse(datetimeString)
	if err == nil {
		return tm, nil
	}

	// another choice
	fuzzytime.ExtendYear(2016)
	dt, _, err := fuzzytime.Extract(datetimeString)
	if err == nil {
		return time.Date(dt.Year(), time.Month(dt.Month()), dt.Day(), dt.Hour(), dt.Minute(), dt.Second(), 0, time.FixedZone("UTC", dt.TZOffset())), nil
	}

	return time.Time{}, gerrors.New("can't parse this time string")
}

// strict == true : only can parse ISO standard date: "YYYY-MM-DD"
// strict == false : fuzzy parse date: "YYYY-MM-DD" or "YYYY.M.D"
func ParseDateString(s string, strict bool) (Date, error) {
	if len(s) == 0 {
		return ZeroDate, gerrors.New("empty input date string")
	}

	if strict {
		formats := []string{"2006-01-02"}
		for _, v := range formats {
			tm, err := time.ParseInLocation(v, s, time.UTC)
			if err == nil {
				return TimeToDate(tm, time.UTC), nil
			}
		}
	} else {
		tm, err := ParseDatetimeStringFuzz(s)
		if err == nil {
			return TimeToDate(tm, time.UTC), nil
		}
	}

	return ZeroDate, gerrors.Errorf("Invalid date string '%s'", s)
}

// strict == true: standard date range supported only: "YYYY-MM-DD/YYYY-MM-DD"
// Date Range standard reference：http://www.ukoln.ac.uk/metadata/dcmi/date-dccd-odrf/
// strict == false: fuzzy parse in web page parser, these formats supported:
// 2018-01-02 - 2018-01-03
// 2018-01-02 ~ 2018-01-03
// 2017.1.2-2017.1.7
func ParseDateRangeString(s string, strict bool) (DateRange, error) {
	// Check
	ErrDefault := gerrors.Errorf("Invalid date range string '%s'", s)
	if len(s) == 0 {
		return ZeroDateRange, nil
	}
	if len(s) < 5 {
		return DateRange{}, ErrDefault
	}

	// Parse
	var splits []string
	if strict {
		splits = []string{"/"}
	} else {
		splits = []string{" / ", " ~ ", " - ", "/", "~", "-"}
	}
	var ss []string
	for _, v := range splits {
		ss = strings.Split(s, v)
		if len(ss) == 2 {
			break
		}
	}
	if len(ss) != 2 {
		return DateRange{}, ErrDefault
	}
	begin, err := ParseDateString(ss[0], strict)
	if err != nil {
		return DateRange{}, err
	}
	if len(ss[0]) > 0 && begin.IsZero() {
		return DateRange{}, ErrDefault
	}
	end, err := ParseDateString(ss[1], strict)
	if err != nil {
		return DateRange{}, err
	}
	if len(ss[1]) > 0 && end.IsZero() {
		return DateRange{}, ErrDefault
	}
	return DateRange{Begin: begin, End: end}, nil
}

func ParseTimeStringStrict(s string) (time.Time, error) {
	formats := []string{
		LayoutRFC3339Nano,        // Formatted by MarshalJSON, like '2021-01-03T17:09:06.686235+08:00' and '2020-12-23T12:00:00+08:00'
		LayoutRFC3339GoExtension, // Formatted by time.String(), like '2020-12-04 18:00:00 +0800 CST'
		//time.RFC3339,
		//time.RFC1123Z,
		//time.RFC1123,
		//time.UnixDate,
	}
	res := time.Time{}
	err := error(nil)
	for _, format := range formats {
		res, err = time.Parse(format, s)
		if err == nil {
			return res, nil
		}
	}
	return time.Time{}, err
}
