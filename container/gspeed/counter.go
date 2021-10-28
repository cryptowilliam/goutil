package gspeed

import (
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"time"
)

type statSlice struct {
	bytes uint64
	tm    time.Time
}

type SpeedCounter struct {
	queue     []statSlice
	dur       time.Duration
	beginTime time.Time // 开始计时时间戳
}

// dur: 多长时间之内的数据作为参与统计速度的有效数据。此时长越长，计算得出的速度变化越均匀
func NewCounter(dur time.Duration) *SpeedCounter {
	var counter SpeedCounter
	counter.dur = dur
	counter.queue = make([]statSlice, 0)
	return &counter
}

// 清空超时的统计切片
func (c *SpeedCounter) clearOutTimeSlice() {
	for len(c.queue) > 0 {
		element := c.queue[0]
		if time.Now().Sub(element.tm) > c.dur {
			c.queue = c.queue[1:] // Dicounterard top element
		} else {
			break
		}
	}
}

// 这个接口什么时候调用？
func (c *SpeedCounter) BeginCount() {
	c.beginTime = time.Now()
}

func (c *SpeedCounter) Add(byteSize uint64) {
	c.clearOutTimeSlice()
	var element statSlice
	element.tm = time.Now()
	element.bytes = byteSize
	c.queue = append(c.queue, element)
}

func _nsecToMillis(nsec int64) int64 {
	return nsec / 1000000
}

func (c *SpeedCounter) Get() (*Speed, error) {
	// Check beginTime member
	if c.beginTime.Unix() <= 0 || time.Now().Before(c.beginTime) {
		return nil, gerrors.New("You need call BeginCount() before Get()")
	}

	var s Speed

	// Clear timeout slices
	c.clearOutTimeSlice()
	if len(c.queue) == 0 {
		s = 0
		return &s, nil
	}

	// Calculate total download bytes
	var totalBytes uint64 = 0
	for _, element := range c.queue {
		totalBytes += element.bytes
	}

	// Calculate total duration
	var totalDur time.Duration
	totalDur = c.queue[len(c.queue)-1].tm.Sub(c.beginTime)
	if totalDur > c.dur {
		totalDur = c.dur
	}

	s = Speed(float64(totalBytes) * 8 * float64(1000) / float64(_nsecToMillis(totalDur.Nanoseconds())))
	return &s, nil
}

func (c *SpeedCounter) Reset() {
	c.queue = nil
}
