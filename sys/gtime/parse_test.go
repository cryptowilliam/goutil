package gtime

import (
	"fmt"
	"github.com/cryptowilliam/goutil/basic/gtest"
	"testing"
)

func TestParseDatetimeStringFuzz(t *testing.T) {
	timeStrArr := []string{
		"12-23 14:01:12",
		"2010-12-23 14:01:12",
		"星期四, 12/01/2016 - 12:19",
		"星期四, 12-01-2016 - 12:19",
		"12月20日(二)13:06",
		"2016年 12月 20日 11:12",
		"2016年 12月 20日 星期二 17:22 BJT",
		"2016年12月20日 星期二 03:30 AM",
		"12月 19日 (星期一) 18:59",
		"2016年12月25日",
		"2016年12月20日 17:43",
		"北京时间2016年12月20日",
		"2016年12月20日",
		"2016年12月20日 16:26",
		"2016年12月19日 下午 8:39",
		"2016年12月20日 13:33",
		"17:19 2016年12月20日",
		"(HKT): 1220 00:00",

		"2018-11-25 13:21:37.400",
		"2016-12-19 09:36 PM",
		"2016-12-20",
		"2016-12-20 10:29",
		"2016-12-20 4:57 PM",
		"2016-12-19 04:36:48",
		"2016-12-20 10:29",
		"2016.12.20 15:57",
		"December 14, 2016",
		"2016/12/14",
		"1220 12:04",
		"20-12-2016",
		"20-12-2016",
		"20.12.2016",
		"2016-12-20 17:32",
		"2016-12-20 17:24",
		"2016-12-20",
		"2016/12/20 17:27:00",
		"2016.12.20 / 15:44",
		"2016-12-20 01:43:22",
		"December 20, 2016, 2:32 am",
		"2016-12-20|11:27",
		"4:20 PM Tuesday Dec 20, 2016",
		"2016-12-11 15:01",
		"19 December 2016 13:24",
		"DECEMBER 19, 2016",
		"DEC 16 2016, 9:25 PM ET",
		"20 Dec, 2016 09:50",
		"2016-12-16",
		"20 DEC 2016 - 7:56PM",
		"December 17, 2016",
		"Dec. 20, 2016 4:41 a.m. ET",
		"19/12/2016",
		"2:43 AM ET",
		"2016-12-19T03:00:00PST",
		"December 19, 2016, 3:00 AM",
		"4:34 pm, December 20, 2016",
		//"December 20, 2016 2:16 AM ET", // 是不是时间录入错误?
		"December 20, 2016 2:16 AM ET",
		"20.12.2016",
		"2016/12/21 11:20",
		"2016-12-21 09:46:06 KST",
		"December 21, 2016 6:35 am JST",
		"DEC. 16, 2016",
		"December 19, 2016, 2:50 PM",
		"Dec. 20, 2016 at 7:42 PM",
		"1400 GMT (2200 HKT) December 20, 2016",
		"21 Dec 2016",
		"December 20, 2016",
		"Tuesday 20 December 2016 23.30 GMT",
		"2016-12-20",
		"Dec. 20, 2016 9:32 PM EST",
		"03:43 21.12.2016",
		"20DEC2016",
		"20 DEC 2016",
		"Tue Dec 20, 2016 | 9:29pm EST",
		"00:09, UK,Wednesday 21 December 2016",
		"Dec 20, 2016, 3:59 PM ET",
		"12/16/2016 - 08:00am",
		"11:37, 20 DEC 2016",
		"19 DECEMBER 2016 • 8:49PM",
		"18:07 GMT, 20 December 2016",
		"12/21/16 AT 1:12 AM",
		"December 20, 2016",
		"DECEMBER 21, 201610:27AM",
		"12/20/2016 03:52 pm ET",
		"Tuesday, Dec. 20, 2016 9:47PM EST",
		"7:29 p.m. EST December 20, 2016",
		"December 20, 2016 | 11:49am",
		"Wed Apr 16 17:32:51 NZST 2014",
		"2010-02-01T13:14:43Z", // an iso 8601 form
		"March 10th, 1999",
		"2:51pm",
		"no date or time info here",
	}

	for _, timeStr := range timeStrArr {
		tm, err := ParseDatetimeStringFuzz(timeStr)
		if err != nil {
			fmt.Println("time string \"" + timeStr + "\" parse error, " + err.Error())
		} else {
			fmt.Println(timeStr)
			fmt.Println(tm.Format("2006-01-02 03:04:05 PM"))
			fmt.Println("")
		}
	}
}

func TestEpochMillisToTime(t *testing.T) {
	tm := EpochMillisToTime(1417536000000)
	if tm.Format("2006-01-02 15:04:05") != "2014-12-03 00:00:00" {
		t.Fail()
	}
}

func TestParseDateString(t *testing.T) {

	// "0001-01-01"

	tmstrs := []string{"", "1901-01-01"}
	for _, v := range tmstrs {
		dt, err := ParseDateString(v, true)
		if err != nil {
			t.Error(err)
			continue
		}
		if dt.String() != v {
			t.Errorf("ParseDateStringFuzz(\"%s\") error, String() = \"%s\"", v, dt.String())
		}
	}
}

func TestParseDateRangeString(t *testing.T) {
	tmstrs := []string{"1901-01-01/1901-01-02", "1901-01-01/", "/1901-01-01"}
	for _, v := range tmstrs {
		dr, err := ParseDateRangeString(v, true)
		if err != nil {
			t.Error(err)
			continue
		}
		if dr.String() != v {
			t.Errorf("ParseDateStringFuzz(\"%s\") error, String() = \"%s\"", v, dr.String())
		}
	}
}

func TestParseTimeStringStrict(t *testing.T) {
	_, err := ParseTimeStringStrict("2020-12-23T12:00:00+08:00")
	gtest.Assert(t, err)
	_, err = ParseTimeStringStrict("2021-01-03T17:09:06.686235+08:00")
	gtest.Assert(t, err)
	_, err = ParseTimeStringStrict("2020-12-04 18:00:00 +0800 CST")
	gtest.Assert(t, err)
}