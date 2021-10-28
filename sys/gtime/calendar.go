package gtime

import (
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/container/gnum"
	"time"
)

/** holidays: github.com/rickar/cal/
 */

func IsWorkDay(tm time.Time) bool {
	return tm.Weekday() != time.Saturday && tm.Weekday() != time.Sunday
}

func ElapsedTimeInDay(tm time.Time) time.Duration {
	dayBegin := time.Date(tm.Year(), tm.Month(), tm.Day(), 0, 0, 0, 0, tm.Location())
	return tm.Sub(dayBegin)
}

func LeftTimeInDay(tm time.Time) time.Duration {
	dayBegin := time.Date(tm.Year(), tm.Month(), tm.Day(), 0, 0, 0, 0, tm.Location())
	nextDayBegin := dayBegin.Add(Day)
	return nextDayBegin.Sub(tm)
}

// get how many days across begin time and end time
func CountDays(begin, end time.Time) int {
	firstDayDuration := LeftTimeInDay(begin)
	lastDayDuration := ElapsedTimeInDay(end)
	return int((end.Sub(begin)-firstDayDuration-lastDayDuration)/Day) + 2
}

/**
    Count how many work days from a specified Weekday in this week

    Weekday-Index 	Weekday 	Work-Days-From	Work-Days-To
	0 				Sunday 		5				0
	1 				Monday 		5				1
	2 				Tuesday 	4				2
	3			 	Wednesday 	3				3
	4 				Thursday 	2				4
	5 				Friday 		1				5
	6 				Saturday 	0				5
*/
// 'from' is included
func CountWorkDaysThisWeekFrom(from time.Weekday) int {
	return gnum.MinInt(6-int(from), 5)
}

// 'to' is included
func CountWorkDaysThisWeekTo(to time.Weekday) int {
	return gnum.MinInt(int(to), 5)
}

// this is not accurate, holidays(christmas, thanks giving day, spring festival...) are marked as workday
func CountWorkDays(begin, end time.Time) int {
	// [1-7,7,7...7,1-7]
	// 把第一周和最后一周拎出来计算，中间很方便计算了
	min := MinTime(begin, end)
	max := MaxTime(begin, end)
	firstWeekday := min.Weekday()            // 第一个时间属于周几
	lastWeekday := max.Weekday()             // 最后一个时间属于周几
	daysOfFirstWeek := int(firstWeekday) + 1 // 第一周的天数
	daysOfLastWeek := int(lastWeekday) + 1   // 最后一周的天数
	daysOfAll := CountDays(min, max)
	daysOfMiddle := 0
	if daysOfAll > daysOfFirstWeek+daysOfLastWeek {
		daysOfMiddle = daysOfAll - daysOfFirstWeek - daysOfLastWeek // 除去第一周和最后一周，中间有多少天
	} else {
		daysOfMiddle = 0
	}
	//fmt.Println(daysOfAll, daysOfFirstWeek, daysOfLastWeek)
	if daysOfMiddle%7 != 0 {
		panic(gerrors.New("daysOfMiddle %d error", daysOfMiddle))
	}
	return CountWorkDaysThisWeekFrom(firstWeekday) + ((daysOfMiddle / 7) * 5) + CountWorkDaysThisWeekTo(lastWeekday)
}
