package gnum

import "math"

const (
	GoldenRatio = 0.618
)

func Fibonacci(n int) int {
	if n <= 0 {
		return 0
	}
	if n == 1 || n == 2 {
		return 1
	}

	a := 1
	b := 1
	num := 0
	for i := 3; i <= n; i++ {
		num = a + b
		a = b
		b = num
	}
	return b
}

func RealGoldenRatio() float64 {
	return (math.Sqrt(5) - 1) / 2
}
