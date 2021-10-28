package gnum

/**
statistic functions
*/

import "math"

// Sum (和)
func Sum(values []float64) float64 {
	ret := 0.
	for _, v := range values {
		ret += v
	}
	return ret
}

// Mean (均值)
func Mean(values []float64) float64 {
	if len(values) == 0 {
		return math.NaN()
	}
	return Sum(values) / float64(len(values))
}

// Standard Deviation (均方差)
// ddof: Delta Degrees of Freedom. The divisor used in calculations is N - ddof, where N represents the number of elements. By default ddof is zero.
func Std(values []float64, ddof int) float64 {
	if len(values) == 0 {
		return math.NaN()
	}
	m := Mean(values)
	ss := 0.
	for _, v := range values {
		d := v - m
		ss += d * d
	}

	return math.Sqrt(ss / float64(len(values)-ddof))
}

// Variance (方差)
// func Var(sample []float64, wholePop bool) float64
