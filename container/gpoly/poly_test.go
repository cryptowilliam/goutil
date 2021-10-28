package gpoly

import (
	"fmt"
	"testing"
)

func TestPolyAngle(t *testing.T) {
	ys := []float64{2, 3, 4, 5, 6}
	fmt.Println(PolyAngle(ys))

	ys = []float64{4, 3, 2, 1, 0}
	fmt.Println(PolyAngle(ys))
}
