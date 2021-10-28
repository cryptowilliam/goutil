package ggeometry

import (
	"math"
)

// get angle degree of angle built by a & b in triangle (a, b, c)
func AngleAB(a, b, c float64) float64 {
	return RadianToAngle(math.Acos((math.Pow(a, 2) + math.Pow(b, 2) - math.Pow(c, 2)) / (2 * a * b)))
}
