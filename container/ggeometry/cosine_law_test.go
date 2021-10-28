package ggeometry

import (
	"fmt"
	"math"
	"testing"
)

func TestAngleAB(t *testing.T) {
	fmt.Println(math.Cos(AngleToRadian(45)))
	fmt.Println(AngleAB(4.242640687119285, 3, 3))
}
