package gnum

import (
	"fmt"
	"github.com/cryptowilliam/goutil/encoding/gjson"
	"math"
	"testing"
)

func TestFloatAlmostEqual(t *testing.T) {
	type S struct {
		A float64
	}
	s := S{A: 2.345}
	fmt.Println(gjson.MarshalStringDefault(s, true))
	fmt.Println(fmt.Sprintf("%.5f", s.A))
}

func TestDetectMaxPrec(t *testing.T) {
	fmt.Println(math.IsNaN(PosInf))
	fmt.Println(fmt.Sprintf("%f", NegInf))
}
