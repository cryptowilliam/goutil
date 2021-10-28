package ghtml

import (
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/sys/gtime"
	"time"
)

/*
 XAxisSync is used to sync kline & dots by x axis
*/

type (
	inDot struct {
		tm  time.Time
		val interface{}
	}

	inLine struct {
		name string
		dots []inDot
	}

	XAxisSync struct {
		// input data source
		in []inLine

		// sorted output cache
		outTime  []time.Time
		outNames []string
		outVals  [][]interface{}
	}
)

func NewXAxisSync() *XAxisSync {
	return new(XAxisSync)
}

func (s *XAxisSync) inHasName(name string) bool {
	for i := range s.in {
		if s.in[i].name == name {
			return true
		}
	}
	return false
}

func (s *XAxisSync) AddFloats(name string, times []time.Time, vals []float64) error {
	if s.inHasName(name) {
		return gerrors.Errorf("name %s already exists", name)
	}
	if len(times) != len(vals) {
		return gerrors.Errorf("time len %d != values len %d", len(times), len(vals))
	}

	newIn := inLine{name: name}
	for i := 0; i < len(times); i++ {
		pVal := new(float64)
		*pVal = vals[i]
		newIn.dots = append(newIn.dots, inDot{times[i], pVal})
	}
	s.in = append(s.in, newIn)
	return nil
}

func (s *XAxisSync) Sync() {
	s.outTime = nil
	s.outNames = nil
	s.outVals = nil

	// 根据缓存输入数据，对所有时间点进行去重和排序
	tmap := gtime.NewTimeMap(nil)
	for i := range s.in {
		for _, dot := range s.in[i].dots {
			tmap.Add(dot.tm)
		}
	}
	s.outTime = gtime.SortTimes(tmap.Export(nil))

	// build fixed size nil two-dimensional array
	nName := len(s.in)
	nTime := len(s.outTime)
	s.outVals = [][]interface{}{}
	for i := 0; i < nName; i++ {
		var nilLine []interface{}
		for j := 0; j < nTime; j++ {
			nilLine = append(nilLine, nil)
		}
		s.outVals = append(s.outVals, nilLine)
	}

	// set names & values
	// 根据操作之后的时间点建立数组，然后按坐标向数组填充数据
	timePosMap := make(map[int64]int)
	for i, v := range s.outTime {
		timePosMap[v.UnixNano()] = i
	}
	for nameIdx := range s.in {
		s.outNames = append(s.outNames, s.in[nameIdx].name)
		for dotIdx := range s.in[nameIdx].dots {
			tm := s.in[nameIdx].dots[dotIdx].tm
			val := s.in[nameIdx].dots[dotIdx].val
			s.outVals[nameIdx][timePosMap[tm.UnixNano()]] = val
		}
	}
}

func (s *XAxisSync) GetNames() []string {
	r := []string{}
	for i := range s.in {
		r = append(r, s.in[i].name)
	}
	return r
}

func (s *XAxisSync) GetTimes() []time.Time {
	return s.outTime
}

// returns []*float / []*KDot
func (s *XAxisSync) GetValues(name string) []interface{} {
	for i, v := range s.outNames {
		if v == name {
			return s.outVals[i]
		}
	}
	return nil
}

func (s *XAxisSync) GetFloatValuesPanic(name string) []*float64 {
	itfs := s.GetValues(name)

	var r []*float64
	for i := range itfs {
		if itfs[i] == nil {
			r = append(r, nil)
		} else {
			r = append(r, itfs[i].(*float64))
		}
	}
	return r
}
