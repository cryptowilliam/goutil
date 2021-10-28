package ghtml

import (
	"encoding/json"
	"fmt"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/container/gnum"
	"github.com/cryptowilliam/goutil/sys/gtime"
	"time"
)

// TODO 可以用ChartTemplate取代吗

/*
ChartLines is a Web/Html/eCharts oriented data structure, for good looking when shown K & indicators.
*/

/*

Compare 2 （可以2条X轴，不支持蜡烛线，不支持各种指标）
Compare N (必须共享X轴，不支持蜡烛线，不支持各种指标)
Backtest （支持CandleStick，支持各种指标）

*/

type (
	ChartLines struct {
		line  *_clJsonLine
		lines *_clJsonLines
		oclh  *_clJsonOclh
	}

	ClLineMode  ChartLines
	ClLinesMode ChartLines
	ClOclhMode  ChartLines

	_clJsonLine struct {
		Times    []gtime.ElegantTime
		LineName string
		Line     []gnum.ElegantFloat
	}

	_clJsonLines struct {
		Times     []gtime.ElegantTime
		LineNames []string
		Lines     [][]*gnum.ElegantFloat

		_sync_x_axis_cache_ *XAxisSync
	}

	_clJsonOclh struct {
		Times     []gtime.ElegantTime
		OCLHName  string
		OCLHs     [][4]gnum.ElegantFloat
		LineNames []string
		Lines     [][]gnum.ElegantFloat
	}
)

func NewChartLines() *ChartLines {
	r := new(ChartLines)
	return r
}

func (cl *ChartLines) LineMode() *ClLineMode {
	if cl.line == nil {
		cl.line = new(_clJsonLine)
	}
	return (*ClLineMode)(cl)
}

func (cl *ChartLines) LinesMode() *ClLinesMode {
	if cl.lines == nil {
		cl.lines = new(_clJsonLines)
		cl.lines._sync_x_axis_cache_ = NewXAxisSync()
	}
	return (*ClLinesMode)(cl)
}

func (cl *ChartLines) OclhMode() *ClOclhMode {
	if cl.oclh == nil {
		cl.oclh = new(_clJsonOclh)
	}
	return (*ClOclhMode)(cl)
}

func (cl *ClLineMode) Set(lineName string, times []time.Time, values []float64) {
	cl.line = nil
	cl.lines = nil
	cl.oclh = nil

	line := new(_clJsonLine)
	line.Times = gtime.NewElegantTimeArray(times, gtime.LayoutRFC3339GoExtension)
	line.LineName = lineName
	line.Line = gnum.NewElegantFloatArray(values, -1)
	cl.line = line
}

func (cl *ClLinesMode) Append(lineName string, times []time.Time, values []float64) {
	cl.line = nil
	cl.oclh = nil

	// sync x axis and copy results to cl.lines
	_ = cl.lines._sync_x_axis_cache_.AddFloats(lineName, times, values)
	cl.lines._sync_x_axis_cache_.Sync()
	cl.lines.LineNames = cl.lines._sync_x_axis_cache_.GetNames()
	cl.lines.Times = gtime.NewElegantTimeArray(cl.lines._sync_x_axis_cache_.GetTimes(), "")
	cl.lines.Lines = nil
	for _, name := range cl.lines.LineNames {
		vals := cl.lines._sync_x_axis_cache_.GetFloatValuesPanic(name)
		cl.lines.Lines = append(cl.lines.Lines, gnum.NewElegantFloatPtrArray(vals, -1))
	}
}

func (cl *ClOclhMode) Init(OCLHName string, times []time.Time, OCLHs [][4]float64) {
	cl.line = nil
	cl.lines = nil
	cl.oclh = new(_clJsonOclh)

	// clean indicators
	//cl.oclh.LineNames = nil
	//cl.oclh.Lines = nil

	cl.oclh.Times = gtime.NewElegantTimeArray(times, gtime.LayoutRFC3339GoExtension)
	cl.oclh.OCLHName = OCLHName
	cl.oclh.OCLHs = nil
	for i := range OCLHs {
		_OCLH_ := [4]gnum.ElegantFloat{
			gnum.NewElegantFloat(OCLHs[i][0], -1),
			gnum.NewElegantFloat(OCLHs[i][1], -1),
			gnum.NewElegantFloat(OCLHs[i][2], -1),
			gnum.NewElegantFloat(OCLHs[i][3], -1),
		}
		cl.oclh.OCLHs = append(cl.oclh.OCLHs, _OCLH_)
	}
}

func (cl *ClOclhMode) Append(name string, vals []float64) error {
	cl.line = nil
	cl.lines = nil

	if cl.oclh == nil || len(cl.oclh.OCLHs) == 0 {
		return gerrors.Errorf("OCLH is null")
	}

	if len(cl.oclh.OCLHs) != len(vals) {
		return gerrors.Errorf("OCLH len is %d, Indicator %s len is %d", len(cl.oclh.OCLHs), name, len(vals))
	}

	cl.oclh.LineNames = append(cl.oclh.LineNames, name)
	cl.oclh.Lines = append(cl.oclh.Lines, gnum.NewElegantFloatArray(vals, -1))
	return nil
}

func (cl *ClOclhMode) GetLine(i int) []float64 {
	if cl.oclh == nil {
		return nil
	}
	if i > len(cl.oclh.Lines)-1 {
		return nil
	}
	return gnum.ElegantFloatArrayToFloatArray(cl.oclh.Lines[i])
}

func (cl *ChartLines) Times() []time.Time {
	var r []time.Time
	ets := cl.JSONTimes()
	for _, v := range ets {
		r = append(r, v.Raw())
	}
	return r
}

func (cl *ChartLines) JSONTimes() []gtime.ElegantTime {
	if cl.line != nil {
		return cl.line.Times
	}
	if cl.lines != nil {
		return cl.lines.Times
	}
	if cl.oclh != nil {
		return cl.oclh.Times
	}
	return nil
}

func (cl *ChartLines) SetTimeLayout(layout string) {
	if cl.line != nil {
		for i := range cl.line.Times {
			cl.line.Times[i].SetLayout(layout)
		}
	}
	if cl.lines != nil {
		for i := range cl.lines.Times {
			cl.lines.Times[i].SetLayout(layout)
		}
	}
	if cl.oclh != nil {
		for i := range cl.oclh.Times {
			cl.oclh.Times[i].SetLayout(layout)
		}
	}
}

func (cl *ChartLines) SetHumanReadPrec(hrp int) {
	if cl.line != nil {
		for i := range cl.line.Line {
			cl.line.Line[i].SetHumanReadPrec(hrp)
		}
	}

	if cl.lines != nil {
		for i := range cl.lines.Lines {
			for j := range cl.lines.Lines[i] {
				cl.lines.Lines[i][j].SetHumanReadPrec(hrp)
			}
		}
	}

	if cl.oclh != nil {
		for i := range cl.oclh.OCLHs {
			cl.oclh.OCLHs[i][0].SetHumanReadPrec(hrp)
			cl.oclh.OCLHs[i][1].SetHumanReadPrec(hrp)
			cl.oclh.OCLHs[i][2].SetHumanReadPrec(hrp)
			cl.oclh.OCLHs[i][3].SetHumanReadPrec(hrp)
		}
		for i := range cl.oclh.Lines {
			for j := range cl.oclh.Lines[i] {
				cl.oclh.Lines[i][j].SetHumanReadPrec(hrp)
			}
		}
	}
}

func (cl *ChartLines) JSONAutoDetect() ([]byte, error) {
	layout := gtime.DetectBestLayout(cl.JSONTimes())
	fmt.Println(layout)
	cl.SetTimeLayout(layout)

	if cl.line != nil {
		return json.Marshal(cl.line)
	}
	if cl.lines != nil {
		return json.Marshal(cl.lines)
	}
	if cl.oclh != nil {
		return json.Marshal(cl.oclh)
	}
	return nil, gerrors.New("no valid member in ChartLines")
}

func (cl *ChartLines) Clone() *ChartLines {
	r := new(ChartLines)
	*r.line = *cl.line
	*r.lines = *cl.lines
	*r.oclh = *cl.oclh
	return r
}
