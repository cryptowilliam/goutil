package gchart

import (
	"bytes"
	"github.com/wcharczuk/go-chart" //exposes "chart"
	"github.com/wcharczuk/go-chart/drawing"
	"time"
)

type TimeSeriesDot struct {
	Time   time.Time
	YValue float64
}

type TimeSeriesLine struct {
	Dots      []TimeSeriesDot
	LineColor drawing.Color
}

func DrawTimeSeriesLines(xAxisLabel, yAxisLabel string, width, height int, lines []TimeSeriesLine) (PNG []byte, err error) {
	graph := chart.Chart{
		Canvas: chart.Style{
			FillColor: drawing.ColorWhite,
		},
		Width:  width,
		Height: height,
		XAxis: chart.XAxis{
			Style: chart.Style{
				Hidden: false,
			},
		},
		YAxis: chart.YAxis{
			Style: chart.Style{
				Hidden: false,
			},
		},
	}
	if len(xAxisLabel) > 0 {
		graph.XAxis.Name = xAxisLabel
	}
	if len(yAxisLabel) > 0 {
		graph.YAxis.Name = yAxisLabel
	}

	// Draw multiple lines.
	for _, line := range lines {
		tslist := []time.Time{}
		ylist := []float64{}
		for _, dot := range line.Dots {
			tslist = append(tslist, dot.Time)
			ylist = append(ylist, dot.YValue)
		}
		graph.Series = []chart.Series{
			chart.TimeSeries{
				XValues: tslist,
				YValues: ylist,
				Style: chart.Style{
					Hidden:      false,                        // note: if we set ANY other properties, we must set this to true.
					StrokeColor: line.LineColor,               // will supercede defaults
					FillColor:   line.LineColor.WithAlpha(80), // will supercede defaults
				},
			}}
	}

	// Render.
	buffer := bytes.NewBuffer([]byte{})
	if err := graph.Render(chart.PNG, buffer); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
