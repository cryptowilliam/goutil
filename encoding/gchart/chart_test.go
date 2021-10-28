package gchart

import (
	"github.com/cryptowilliam/goutil/sys/gfs"
	"github.com/wcharczuk/go-chart"
	"time"
)

func sample() {
	line := TimeSeriesLine{}
	line.Dots = []TimeSeriesDot{{Time: time.Now().AddDate(0, 0, -9), YValue: 1},
		{Time: time.Now().AddDate(0, 0, -8), YValue: 2},
		{Time: time.Now().AddDate(0, 0, -7), YValue: 3},
		{Time: time.Now().AddDate(0, 0, -6), YValue: 4},
		{Time: time.Now().AddDate(0, 0, -5), YValue: 5},
		{Time: time.Now().AddDate(0, 0, -4), YValue: 6},
		{Time: time.Now().AddDate(0, 0, -3), YValue: 7},
		{Time: time.Now().AddDate(0, 0, -2), YValue: 8},
		{Time: time.Now().AddDate(0, 0, -1), YValue: 7},
		{Time: time.Now().AddDate(0, 0, -0), YValue: 6},
		{Time: time.Now().AddDate(0, 0, 1), YValue: 5}}
	line.LineColor = chart.ColorAlternateGreen
	buf, _ := DrawTimeSeriesLines("x", "y", 1000, 600, []TimeSeriesLine{line})
	gfs.BytesToFile(buf, "sample.png")
}
