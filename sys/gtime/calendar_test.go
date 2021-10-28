package gtime

import (
	"github.com/cryptowilliam/goutil/basic/gtest"
	"testing"
	"time"
)

func TestCountWorkDays(t *testing.T) {
	begin := time.Date(2020, 1, 1, 0, 0, 0, 0, TimeZoneAsiaShanghai)
	end := time.Date(2020, 1, 17, 0, 0, 0, 0, TimeZoneAsiaShanghai)
	if CountWorkDays(begin, end) != 13 {
		gtest.PrintlnExit(t, "days should be 13")
	}

	begin = time.Date(2020, 1, 1, 0, 0, 0, 0, TimeZoneAsiaShanghai)
	end = time.Date(2020, 1, 17, 0, 1, 0, 0, TimeZoneAsiaShanghai)
	if CountWorkDays(begin, end) != 13 {
		gtest.PrintlnExit(t, "days should be 13")
	}
}
