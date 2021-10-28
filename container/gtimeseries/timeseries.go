package gtimeseries

import "time"

// github.com/codesuki/go-time-series 数据结构、排序，基础工具
// github.com/timkaye11/goTS
// github.com/lytics/anomalyzer 时间序列概率异常检测
// github.com/jianc94538/stats 时间序列数据结构与时序统计

type (
	TimeSeries interface {
		Time(i int) time.Time
		High(i int) float64
		Open(i int) float64
		Close(i int) float64
		Low(i int) float64
		Volume(i int) float64
		Len() int
	}
)
