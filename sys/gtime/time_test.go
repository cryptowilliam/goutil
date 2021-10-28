package gtime

import (
	"testing"
	"time"
)

func TestJSONTime_DetectBestLayout(t *testing.T) {
	tm := time.Date(2017, 1, 1, 0, 0, 1, 1, time.UTC)
	jtm := NewElegantTime(tm, "")
	t.Log(jtm.DetectBestLayout())
	b, _ := jtm.JSONAutoDetect()
	t.Log(string(b))

	tm = time.Date(2019, 1, 1, 3, 0, 0, 0, time.UTC)
	jtm = NewElegantTime(tm, "")
	t.Log(jtm.DetectBestLayout())
	b, _ = jtm.JSONAutoDetect()
	t.Log(string(b))

	tm = time.Date(2019, 1, 1, 3, 0, 0, 0, time.Local)
	jtm = NewElegantTime(tm, "")
	t.Log(jtm.DetectBestLayout())
	b, _ = jtm.JSONAutoDetect()
	t.Log(string(b))
}

func TestDetectBestLayout(t *testing.T) {
	tm1 := time.Date(2017, 1, 1, 0, 0, 1, 0, time.UTC)
	tm2 := time.Date(2019, 1, 1, 3, 0, 0, 0, time.UTC)
	tm3 := time.Date(2019, 1, 1, 3, 0, 0, 0, time.Local)
	jtm1 := NewElegantTime(tm1, "")
	jtm2 := NewElegantTime(tm2, "")
	jtm3 := NewElegantTime(tm3, "")

	t.Log(DetectBestLayout([]ElegantTime{jtm1, jtm2, jtm3}))
}
