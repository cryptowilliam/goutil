package gtime

import "time"

var (
	EpochBeginTime           = EpochSecToTime(0)              // 1970-01-01 00:00:00 +0000 UTC
	EpochBeginDate           = Date(19700101)                 // 0001-01-01 00:00:00 +0000 UTC
	ZeroTime                 = time.Time{}                    // 0001-01-01 00:00:00 +0000 UTC
	ZeroDate                 = TimeToDate(ZeroTime, time.UTC) // 0001-01-01 00:00:00 +0000 UTC
	ZeroYearMonth  YearMonth = 0                              // 0000-00
	ZeroDateRange            = DateRange{Begin: ZeroDate, End: ZeroDate}
)
