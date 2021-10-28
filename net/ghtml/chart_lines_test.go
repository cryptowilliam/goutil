package ghtml

import (
	"testing"
	"time"
)

func TestChartLines_JSONAutoDetect(t *testing.T) {
	tmsRain := []time.Time{
		time.Date(2017, 1, 1, 0, 0, 1, 0, time.UTC),
		time.Date(2019, 1, 2, 3, 0, 0, 0, time.UTC),
		time.Date(2019, 1, 3, 3, 0, 0, 0, time.Local),
	}
	valRain := []float64{1.1, 2.2, 3.3}

	tmsWind := []time.Time{
		time.Date(2017, 2, 1, 0, 0, 1, 0, time.UTC),
		time.Date(2019, 2, 2, 3, 0, 0, 0, time.UTC),
		time.Date(2019, 2, 3, 3, 0, 0, 0, time.Local),
	}
	valWind := []float64{2.2, 3.3, 4.4}

	tmsTemp := []time.Time{
		time.Date(2017, 3, 1, 0, 0, 1, 0, time.UTC),
		time.Date(2019, 3, 2, 3, 0, 0, 0, time.UTC),
		time.Date(2019, 3, 3, 3, 0, 0, 0, time.Local),
	}
	valsTemp := []float64{28.5, 27, 29.1}

	cl := NewChartLines()
	cl.LineMode().Set("rain", tmsRain, valRain)
	b, err := cl.JSONAutoDetect()
	if err != err {
		t.Error(err)
		return
	}
	t.Log(string(b))

	cl.LinesMode().Append("rain", tmsRain, valRain)
	cl.LinesMode().Append("wind", tmsWind, valWind)
	cl.LinesMode().Append("temperature", tmsTemp, valsTemp)
	b, err = cl.JSONAutoDetect()
	if err != err {
		t.Error(err)
		return
	}
	t.Log(string(b))
	if string(b) != `{"Times":["2017-01-01 00:00:01","2017-02-01 00:00:01","2017-03-01 00:00:01","2019-01-02 03:00:00","2019-01-02 19:00:00","2019-02-02 03:00:00","2019-02-02 19:00:00","2019-03-02 03:00:00","2019-03-02 19:00:00"],"LineNames":["rain","wind","temperature"],"Lines":[[1.1,null,null,2.2,3.3,null,null,null,null],[null,2.2,null,null,null,3.3,4.4,null,null],[null,null,28.5,null,null,null,null,27,29.1]]}` {
		t.Errorf("chart lines sync error")
		return
	}
}
