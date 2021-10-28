package nlp

import (
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/encoding/gchart"
	"github.com/cryptowilliam/goutil/sys/gtime"
	"github.com/wcharczuk/go-chart/drawing"
	"sync"
	"time"
)

// 把时许轴上大量的Post的情绪指数进行分析

type SentimentReportItem struct {
	Time  time.Time
	Score float64
}

type SentimentReport struct {
	Items []SentimentReportItem
}

type Sentiment struct {
	Polarity     float64 `json:"Polarity"`
	Subjectivity float64 `json:"Subjectivity"`
}

type SentimentWithTime struct {
	Time time.Time
	Sentiment
}

type SentimentSet struct {
	Items            []SentimentWithTime
	itemsRwLock      sync.RWMutex
	dataKeptInterval time.Duration
}

func NewSet(dataKeptInterval time.Duration) *SentimentSet {
	return &SentimentSet{dataKeptInterval: dataKeptInterval}
}

func (ss *SentimentSet) IsEmpty() bool {
	ss.itemsRwLock.RLock()
	defer ss.itemsRwLock.RUnlock()

	return len(ss.Items) == 0
}

func (ss *SentimentSet) cleanExpiredData() {
	beginTime := gtime.Sub(time.Now(), ss.dataKeptInterval)
	beginIndex := 0

	ss.itemsRwLock.Lock()
	defer ss.itemsRwLock.Unlock()

	for i, v := range ss.Items {
		if v.Time.Before(beginTime) {
			beginIndex = i
		} else {
			break
		}
	}
	ss.Items = ss.Items[beginIndex:]
}

func (ss *SentimentSet) Append(tm time.Time, polarity, subjectivity float64) {
	item := SentimentWithTime{}
	item.Time = tm
	item.Polarity = polarity
	item.Subjectivity = subjectivity

	ss.itemsRwLock.Lock()
	defer ss.itemsRwLock.Unlock()
	ss.Items = append(ss.Items, item)
}

func (ss *SentimentSet) ParseReport(interval time.Duration) (*SentimentReport, error) {
	if ss.IsEmpty() {
		return nil, gerrors.Errorf("Empty SentimentSet items.")
	}
	ss.cleanExpiredData()

	ss.itemsRwLock.RLock()
	defer ss.itemsRwLock.RUnlock()

	resultItem := SentimentReportItem{}
	result := SentimentReport{}
	lastPack := []SentimentWithTime{}
	lastPackTime := ss.Items[0].Time

	for i, v := range ss.Items {
		// 当前Item不超过interval，作为一个pack的item进行缓存
		if v.Time.Sub(lastPackTime) < interval && i < len(ss.Items)-1 {
			lastPack = append(lastPack, v)
			continue
		}

		// 当前Item超过interval，需要把这段时间范围内的Pack进行平均计算，并作为一个ReportItem存档
		totalScore := float64(0)
		resultItem.Time = lastPackTime
		for _, it := range lastPack {
			totalScore += it.Polarity
		}
		resultItem.Score = totalScore / float64(len(lastPack))
		result.Items = append(result.Items, resultItem)
		// Reset.
		lastPack = nil
		lastPack = append(lastPack, v)
		lastPackTime = v.Time
	}

	return &result, nil
}

func (sp *SentimentReport) DrawImage() ([]byte, error) {
	line := gchart.TimeSeriesLine{}
	line.LineColor = drawing.ColorBlack
	for _, v := range sp.Items {
		dot := gchart.TimeSeriesDot{}
		dot.Time = v.Time
		dot.YValue = v.Score
		line.Dots = append(line.Dots, dot)
	}
	return gchart.DrawTimeSeriesLines("Date", "Index", 600, 400, []gchart.TimeSeriesLine{line})
}
