package gnum

import (
	"github.com/markcheno/go-talib"
)

// Avg - Vector average
func Avg(inReal []float64, timePeriod int) []float64 {
	r := talib.Sum(inReal, timePeriod)
	for i := range r {
		r[i] = r[i] / float64(timePeriod)
	}
	return r
}

func MinSS(inReal0, inReal1 []float64) []float64 {
	if len(inReal0) == 0 || len(inReal1) == 0 {
		return nil
	}
	if len(inReal0) != len(inReal1) {
		return nil
	}
	var min []float64
	for i := 0; i < len(inReal0); i++ {
		min = append(min, MinFloat(inReal0[i], inReal1[i]))
	}
	return min
}

func MinSV(inReal []float64, inTimePeriod int) []float64 {
	return talib.Min(inReal, inTimePeriod)
}

func MaxSS(inReal0, inReal1 []float64) []float64 {
	if len(inReal0) == 0 || len(inReal1) == 0 {
		return nil
	}
	if len(inReal0) != len(inReal1) {
		return nil
	}
	var max []float64
	for i := 0; i < len(inReal0); i++ {
		max = append(max, MaxFloat(inReal0[i], inReal1[i]))
	}
	return max
}

func MaxSV(inReal []float64, inTimePeriod int) []float64 {
	return talib.Max(inReal, inTimePeriod)
}

func AddSS(inReal0, inReal1 []float64) []float64 {
	return talib.Add(inReal0, inReal1)
}

func AddSV(inReal []float64, v float64) []float64 {
	var r []float64
	for i := 0; i < len(inReal); i++ {
		r = append(r, inReal[i]+v)
	}
	return r
}

func SubSS(inReal0, inReal1 []float64) []float64 {
	return talib.Sub(inReal0, inReal1)
}

func SubSV(inReal []float64, v float64) []float64 {
	var r []float64
	for i := 0; i < len(inReal); i++ {
		r = append(r, inReal[i]-v)
	}
	return r
}

func MulSS(inReal0, inReal1 []float64) []float64 {
	return talib.Mult(inReal0, inReal1)
}

func MulSV(inReal []float64, v float64) []float64 {
	var r []float64
	for i := 0; i < len(inReal); i++ {
		r = append(r, inReal[i]*v)
	}
	return r
}

func SetHead0(inReal []float64, count int) {
	if count > len(inReal) {
		count = len(inReal)
	}
	for i := 0; i < count; i++ {
		inReal[i] = 0
	}
}
