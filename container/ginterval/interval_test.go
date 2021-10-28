package ginterval

import (
	"github.com/cryptowilliam/goutil/basic/gtest"
	"github.com/cryptowilliam/goutil/container/gdecimal"
	"testing"
)

func TestParse(t *testing.T) {
	cl := gtest.NewCaseList()

	cl.New().Input("[-2.5,-1]").Expect("[-2.5,-1]").Expect(true)
	cl.New().Input("[-1,  100.5]").Expect("[-1,100.5]").Expect(true)
	cl.New().Input("[2, 1000]").Expect("[2,1000]").Expect(true)
	cl.New().Input("[-2, -1)").Expect("[-2,-1)").Expect(true)
	cl.New().Input("[-1.5, 100)").Expect("[-1.5,100)").Expect(true)
	cl.New().Input("[2,  1000)").Expect("[2,1000)").Expect(true)
	cl.New().Input("(-2, -1)").Expect("(-2,-1)").Expect(true)
	cl.New().Input("(-1, 100.24)").Expect("(-1,100.24)").Expect(true)
	cl.New().Input("(2.14,1000)").Expect("(2.14,1000)").Expect(true)
	cl.New().Input("(-2, -1]").Expect("(-2,-1]").Expect(true)
	cl.New().Input("(-1,100]").Expect("(-1,100]").Expect(true)
	cl.New().Input("(2,    1000]").Expect("(2,1000]").Expect(true)
	cl.New().Input("(-∞, -100)").Expect("(-∞,-100)").Expect(true)
	cl.New().Input("(10, +∞)").Expect("(10,+∞)").Expect(true)
	cl.New().Input("(-∞, +∞)").Expect("(-∞,+∞)").Expect(true)

	cl.New().Input("[-100, -200]").Expect("").Expect(false)
	cl.New().Input("[1000, 100]").Expect("").Expect(false)
	cl.New().Input("[100,-200]").Expect("").Expect(false)
	cl.New().Input("[-∞, +∞]").Expect("").Expect(false)
	cl.New().Input("[-∞, -100]").Expect("").Expect(false)
	cl.New().Input("[-100, -∞]").Expect("").Expect(false)
	cl.New().Input("[10, +∞]").Expect("").Expect(false)
	cl.New().Input("[-100, -200)").Expect("").Expect(false)
	cl.New().Input("[1000, 100)").Expect("").Expect(false)
	cl.New().Input("[100,-200)").Expect("").Expect(false)
	cl.New().Input("[-∞, +∞)").Expect("").Expect(false)
	cl.New().Input("[-∞,  -100)").Expect("").Expect(false)
	cl.New().Input("[-100, -∞)").Expect("").Expect(false)
	cl.New().Input("[-100, -200)").Expect("").Expect(false)
	cl.New().Input("[1000, 100)").Expect("").Expect(false)
	cl.New().Input("[100, -200)").Expect("").Expect(false)
	cl.New().Input("[-∞,+∞)").Expect("").Expect(false)
	cl.New().Input("[-∞, -100)").Expect("").Expect(false)
	cl.New().Input("[-100, -∞)").Expect("").Expect(false)
	cl.New().Input("(-100, -200)").Expect("").Expect(false)
	cl.New().Input("(1000, 100)").Expect("").Expect(false)
	cl.New().Input("(100,-200)").Expect("").Expect(false)
	cl.New().Input("(-100, -∞)").Expect("").Expect(false)
	cl.New().Input("[3.14, 0.68]").Expect("").Expect(false)
	cl.New().Input("[+-100, 200]").Expect("").Expect(false)
	cl.New().Input("[yes,OK]").Expect("").Expect(false)

	for _, v := range cl.Get() {
		intvl, err := Parse(v.Inputs[0].(string))
		gotIntvl := ""
		if intvl != nil {
			gotIntvl = intvl.String()
		}
		gotOK := err == nil

		expectIntvl, expectOK := v.Expects[0].(string), v.Expects[1].(bool)

		if gotIntvl != expectIntvl || gotOK != expectOK {
			gtest.PrintlnExit(t, "string '%s', expect '%s' '%v', got '%s' '%v'", v.Inputs[0].(string), expectIntvl, expectOK, gotIntvl, gotOK)
		}
	}
}


func TestInterval_IsOverlap(t *testing.T) {
	cl := gtest.NewCaseList()

	// A
	cl.New().Input("[-2,2]").Input("[-2, 2]").Expect(true)
	cl.New().Input("[-2,2]").Input("[-4,-2]").Expect(true)
	cl.New().Input("[-2,2]").Input("[ 2, 4]").Expect(true)
	cl.New().Input("[-2,2]").Input("[-4,-1]").Expect(true)
	cl.New().Input("[-2,2]").Input("[ 1, 4]").Expect(true)
	cl.New().Input("[-2,2]").Input("[-4, 4]").Expect(true)
	cl.New().Input("[-2,2]").Input("[-1, 1]").Expect(true)
	cl.New().Input("[-2,2]").Input("[-6,-4]").Expect(false)
	cl.New().Input("[-2,2]").Input("[ 4, 6]").Expect(false)

	cl.New().Input("[-2,2]").Input("(-2, 2]").Expect(true)
	cl.New().Input("[-2,2]").Input("(-4,-2]").Expect(true)
	cl.New().Input("[-2,2]").Input("( 2, 4]").Expect(false)
	cl.New().Input("[-2,2]").Input("(-4,-1]").Expect(true)
	cl.New().Input("[-2,2]").Input("( 1, 4]").Expect(true)
	cl.New().Input("[-2,2]").Input("(-4, 4]").Expect(true)
	cl.New().Input("[-2,2]").Input("(-1, 1]").Expect(true)
	cl.New().Input("[-2,2]").Input("(-6,-4]").Expect(false)
	cl.New().Input("[-2,2]").Input("( 4, 6]").Expect(false)

	cl.New().Input("[-2,2]").Input("[-2, 2)").Expect(true)
	cl.New().Input("[-2,2]").Input("[-4,-2)").Expect(false)
	cl.New().Input("[-2,2]").Input("[ 2, 4)").Expect(true)
	cl.New().Input("[-2,2]").Input("[-4,-1)").Expect(true)
	cl.New().Input("[-2,2]").Input("[ 1, 4)").Expect(true)
	cl.New().Input("[-2,2]").Input("[-4, 4)").Expect(true)
	cl.New().Input("[-2,2]").Input("[-1, 1)").Expect(true)
	cl.New().Input("[-2,2]").Input("[-6,-4)").Expect(false)
	cl.New().Input("[-2,2]").Input("[ 4, 6)").Expect(false)

	cl.New().Input("[-2,2]").Input("(-2, 2)").Expect(true)
	cl.New().Input("[-2,2]").Input("(-4,-2)").Expect(false)
	cl.New().Input("[-2,2]").Input("( 2, 4)").Expect(false)
	cl.New().Input("[-2,2]").Input("(-4,-1)").Expect(true)
	cl.New().Input("[-2,2]").Input("( 1, 4)").Expect(true)
	cl.New().Input("[-2,2]").Input("(-4, 4)").Expect(true)
	cl.New().Input("[-2,2]").Input("(-1, 1)").Expect(true)
	cl.New().Input("[-2,2]").Input("(-6,-4)").Expect(false)
	cl.New().Input("[-2,2]").Input("( 4, 6)").Expect(false)


	// B
	cl.New().Input("(-2,2]").Input("[-2, 2]").Expect(true)
	cl.New().Input("(-2,2]").Input("[-4,-2]").Expect(false)
	cl.New().Input("(-2,2]").Input("[ 2, 4]").Expect(true)
	cl.New().Input("(-2,2]").Input("[-4,-1]").Expect(true)
	cl.New().Input("(-2,2]").Input("[ 1, 4]").Expect(true)
	cl.New().Input("(-2,2]").Input("[-4, 4]").Expect(true)
	cl.New().Input("(-2,2]").Input("[-1, 1]").Expect(true)
	cl.New().Input("(-2,2]").Input("[-6,-4]").Expect(false)
	cl.New().Input("(-2,2]").Input("[ 4, 6]").Expect(false)

	cl.New().Input("(-2,2]").Input("(-2, 2]").Expect(true)
	cl.New().Input("(-2,2]").Input("(-4,-2]").Expect(false)
	cl.New().Input("(-2,2]").Input("( 2, 4]").Expect(false)
	cl.New().Input("(-2,2]").Input("(-4,-1]").Expect(true)
	cl.New().Input("(-2,2]").Input("( 1, 4]").Expect(true)
	cl.New().Input("(-2,2]").Input("(-4, 4]").Expect(true)
	cl.New().Input("(-2,2]").Input("(-1, 1]").Expect(true)
	cl.New().Input("(-2,2]").Input("(-6,-4]").Expect(false)
	cl.New().Input("(-2,2]").Input("( 4, 6]").Expect(false)

	cl.New().Input("(-2,2]").Input("[-2, 2)").Expect(true)
	cl.New().Input("(-2,2]").Input("[-4,-2)").Expect(false)
	cl.New().Input("(-2,2]").Input("[ 2, 4)").Expect(true)
	cl.New().Input("(-2,2]").Input("[-4,-1)").Expect(true)
	cl.New().Input("(-2,2]").Input("[ 1, 4)").Expect(true)
	cl.New().Input("(-2,2]").Input("[-4, 4)").Expect(true)
	cl.New().Input("(-2,2]").Input("[-1, 1)").Expect(true)
	cl.New().Input("(-2,2]").Input("[-6,-4)").Expect(false)
	cl.New().Input("(-2,2]").Input("[ 4, 6)").Expect(false)

	cl.New().Input("(-2,2]").Input("(-2, 2)").Expect(true)
	cl.New().Input("(-2,2]").Input("(-4,-2)").Expect(false)
	cl.New().Input("(-2,2]").Input("( 2, 4)").Expect(false)
	cl.New().Input("(-2,2]").Input("(-4,-1)").Expect(true)
	cl.New().Input("(-2,2]").Input("( 1, 4)").Expect(true)
	cl.New().Input("(-2,2]").Input("(-4, 4)").Expect(true)
	cl.New().Input("(-2,2]").Input("(-1, 1)").Expect(true)
	cl.New().Input("(-2,2]").Input("(-6,-4)").Expect(false)
	cl.New().Input("(-2,2]").Input("( 4, 6)").Expect(false)

	// C
	cl.New().Input("[-2,2)").Input("[-2, 2]").Expect(true)
	cl.New().Input("[-2,2)").Input("[-4,-2]").Expect(true)
	cl.New().Input("[-2,2)").Input("[ 2, 4]").Expect(false)
	cl.New().Input("[-2,2)").Input("[-4,-1]").Expect(true)
	cl.New().Input("[-2,2)").Input("[ 1, 4]").Expect(true)
	cl.New().Input("[-2,2)").Input("[-4, 4]").Expect(true)
	cl.New().Input("[-2,2)").Input("[-1, 1]").Expect(true)
	cl.New().Input("[-2,2)").Input("[-6,-4]").Expect(false)
	cl.New().Input("[-2,2)").Input("[ 4, 6]").Expect(false)

	cl.New().Input("[-2,2)").Input("(-2, 2]").Expect(true)
	cl.New().Input("[-2,2)").Input("(-4,-2]").Expect(true)
	cl.New().Input("[-2,2)").Input("( 2, 4]").Expect(false)
	cl.New().Input("[-2,2)").Input("(-4,-1]").Expect(true)
	cl.New().Input("[-2,2)").Input("( 1, 4]").Expect(true)
	cl.New().Input("[-2,2)").Input("(-4, 4]").Expect(true)
	cl.New().Input("[-2,2)").Input("(-1, 1]").Expect(true)
	cl.New().Input("[-2,2)").Input("(-6,-4]").Expect(false)
	cl.New().Input("[-2,2)").Input("( 4, 6]").Expect(false)

	cl.New().Input("[-2,2)").Input("[-2, 2)").Expect(true)
	cl.New().Input("[-2,2)").Input("[-4,-2)").Expect(false)
	cl.New().Input("[-2,2)").Input("[ 2, 4)").Expect(false)
	cl.New().Input("[-2,2)").Input("[-4,-1)").Expect(true)
	cl.New().Input("[-2,2)").Input("[ 1, 4)").Expect(true)
	cl.New().Input("[-2,2)").Input("[-4, 4)").Expect(true)
	cl.New().Input("[-2,2)").Input("[-1, 1)").Expect(true)
	cl.New().Input("[-2,2)").Input("[-6,-4)").Expect(false)
	cl.New().Input("[-2,2)").Input("[ 4, 6)").Expect(false)

	cl.New().Input("[-2,2)").Input("(-2, 2)").Expect(true)
	cl.New().Input("[-2,2)").Input("(-4,-2)").Expect(false)
	cl.New().Input("[-2,2)").Input("( 2, 4)").Expect(false)
	cl.New().Input("[-2,2)").Input("(-4,-1)").Expect(true)
	cl.New().Input("[-2,2)").Input("( 1, 4)").Expect(true)
	cl.New().Input("[-2,2)").Input("(-4, 4)").Expect(true)
	cl.New().Input("[-2,2)").Input("(-1, 1)").Expect(true)
	cl.New().Input("[-2,2)").Input("(-6,-4)").Expect(false)
	cl.New().Input("[-2,2)").Input("( 4, 6)").Expect(false)

	// D
	cl.New().Input("(-2,2)").Input("[-2, 2]").Expect(true)
	cl.New().Input("(-2,2)").Input("[-4,-2]").Expect(false)
	cl.New().Input("(-2,2)").Input("[ 2, 4]").Expect(false)
	cl.New().Input("(-2,2)").Input("[-4,-1]").Expect(true)
	cl.New().Input("(-2,2)").Input("[ 1, 4]").Expect(true)
	cl.New().Input("(-2,2)").Input("[-4, 4]").Expect(true)
	cl.New().Input("(-2,2)").Input("[-1, 1]").Expect(true)
	cl.New().Input("(-2,2)").Input("[-6,-4]").Expect(false)
	cl.New().Input("(-2,2)").Input("[ 4, 6]").Expect(false)

	cl.New().Input("(-2,2)").Input("(-2, 2]").Expect(true)
	cl.New().Input("(-2,2)").Input("(-4,-2]").Expect(false)
	cl.New().Input("(-2,2)").Input("( 2, 4]").Expect(false)
	cl.New().Input("(-2,2)").Input("(-4,-1]").Expect(true)
	cl.New().Input("(-2,2)").Input("( 1, 4]").Expect(true)
	cl.New().Input("(-2,2)").Input("(-4, 4]").Expect(true)
	cl.New().Input("(-2,2)").Input("(-1, 1]").Expect(true)
	cl.New().Input("(-2,2)").Input("(-6,-4]").Expect(false)
	cl.New().Input("(-2,2)").Input("( 4, 6]").Expect(false)

	cl.New().Input("(-2,2)").Input("[-2, 2)").Expect(true)
	cl.New().Input("(-2,2)").Input("[-4,-2)").Expect(false)
	cl.New().Input("(-2,2)").Input("[ 2, 4)").Expect(false)
	cl.New().Input("(-2,2)").Input("[-4,-1)").Expect(true)
	cl.New().Input("(-2,2)").Input("[ 1, 4)").Expect(true)
	cl.New().Input("(-2,2)").Input("[-4, 4)").Expect(true)
	cl.New().Input("(-2,2)").Input("[-1, 1)").Expect(true)
	cl.New().Input("(-2,2)").Input("[-6,-4)").Expect(false)
	cl.New().Input("(-2,2)").Input("[ 4, 6)").Expect(false)

	cl.New().Input("(-2,2)").Input("(-2, 2)").Expect(true)
	cl.New().Input("(-2,2)").Input("(-4,-2)").Expect(false)
	cl.New().Input("(-2,2)").Input("( 2, 4)").Expect(false)
	cl.New().Input("(-2,2)").Input("(-4,-1)").Expect(true)
	cl.New().Input("(-2,2)").Input("( 1, 4)").Expect(true)
	cl.New().Input("(-2,2)").Input("(-4, 4)").Expect(true)
	cl.New().Input("(-2,2)").Input("(-1, 1)").Expect(true)
	cl.New().Input("(-2,2)").Input("(-6,-4)").Expect(false)
	cl.New().Input("(-2,2)").Input("( 4, 6)").Expect(false)

	// 1
	cl.New().Input("[-2,2]").Input("[-4, +∞)").Expect(true)
	cl.New().Input("[-2,2]").Input("[-2, +∞)").Expect(true)
	cl.New().Input("[-2,2]").Input("[ 0, +∞)").Expect(true)
	cl.New().Input("[-2,2]").Input("[ 2, +∞)").Expect(true)
	cl.New().Input("[-2,2]").Input("[ 4, +∞)").Expect(false)

	cl.New().Input("[-2,2]").Input("(-4, +∞)").Expect(true)
	cl.New().Input("[-2,2]").Input("(-2, +∞)").Expect(true)
	cl.New().Input("[-2,2]").Input("( 0, +∞)").Expect(true)
	cl.New().Input("[-2,2]").Input("( 2, +∞)").Expect(false)
	cl.New().Input("[-2,2]").Input("( 4, +∞)").Expect(false)

	// 2
	cl.New().Input("(-2,2]").Input("[-4, +∞)").Expect(true)
	cl.New().Input("(-2,2]").Input("[-2, +∞)").Expect(true)
	cl.New().Input("(-2,2]").Input("[ 0, +∞)").Expect(true)
	cl.New().Input("(-2,2]").Input("[ 2, +∞)").Expect(true)
	cl.New().Input("(-2,2]").Input("[ 4, +∞)").Expect(false)

	cl.New().Input("(-2,2]").Input("(-4, +∞)").Expect(true)
	cl.New().Input("(-2,2]").Input("(-2, +∞)").Expect(true)
	cl.New().Input("(-2,2]").Input("( 0, +∞)").Expect(true)
	cl.New().Input("(-2,2]").Input("( 2, +∞)").Expect(false)
	cl.New().Input("(-2,2]").Input("( 4, +∞)").Expect(false)

	// 3
	cl.New().Input("[-2,2)").Input("[-4, +∞)").Expect(true)
	cl.New().Input("[-2,2)").Input("[-2, +∞)").Expect(true)
	cl.New().Input("[-2,2)").Input("[ 0, +∞)").Expect(true)
	cl.New().Input("[-2,2)").Input("[ 2, +∞)").Expect(false)
	cl.New().Input("[-2,2)").Input("[ 4, +∞)").Expect(false)

	cl.New().Input("[-2,2)").Input("(-4, +∞)").Expect(true)
	cl.New().Input("[-2,2)").Input("(-2, +∞)").Expect(true)
	cl.New().Input("[-2,2)").Input("( 0, +∞)").Expect(true)
	cl.New().Input("[-2,2)").Input("( 2, +∞)").Expect(false)
	cl.New().Input("[-2,2)").Input("( 4, +∞)").Expect(false)

	// 4
	cl.New().Input("(-2,2)").Input("[-4, +∞)").Expect(true)
	cl.New().Input("(-2,2)").Input("[-2, +∞)").Expect(true)
	cl.New().Input("(-2,2)").Input("[ 0, +∞)").Expect(true)
	cl.New().Input("(-2,2)").Input("[ 2, +∞)").Expect(false)
	cl.New().Input("(-2,2)").Input("[ 4, +∞)").Expect(false)

	cl.New().Input("(-2,2)").Input("(-4, +∞)").Expect(true)
	cl.New().Input("(-2,2)").Input("(-2, +∞)").Expect(true)
	cl.New().Input("(-2,2)").Input("( 0, +∞)").Expect(true)
	cl.New().Input("(-2,2)").Input("( 2, +∞)").Expect(false)
	cl.New().Input("(-2,2)").Input("( 4, +∞)").Expect(false)

	// a
	cl.New().Input("[-2,2]").Input("(-∞, -4]").Expect(false)
	cl.New().Input("[-2,2]").Input("(-∞, -2]").Expect(true)
	cl.New().Input("[-2,2]").Input("(-∞,  0]").Expect(true)
	cl.New().Input("[-2,2]").Input("(-∞,  2]").Expect(true)
	cl.New().Input("[-2,2]").Input("(-∞,  4]").Expect(true)

	cl.New().Input("[-2,2]").Input("(-∞, -4)").Expect(false)
	cl.New().Input("[-2,2]").Input("(-∞, -2)").Expect(false)
	cl.New().Input("[-2,2]").Input("(-∞,  0)").Expect(true)
	cl.New().Input("[-2,2]").Input("(-∞,  2)").Expect(true)
	cl.New().Input("[-2,2]").Input("(-∞,  4)").Expect(true)

	// b
	cl.New().Input("(-2,2]").Input("(-∞, -4]").Expect(false)
	cl.New().Input("(-2,2]").Input("(-∞, -2]").Expect(false)
	cl.New().Input("(-2,2]").Input("(-∞,  0]").Expect(true)
	cl.New().Input("(-2,2]").Input("(-∞,  2]").Expect(true)
	cl.New().Input("(-2,2]").Input("(-∞,  4]").Expect(true)

	cl.New().Input("(-2,2]").Input("(-∞, -4)").Expect(false)
	cl.New().Input("(-2,2]").Input("(-∞, -2)").Expect(false)
	cl.New().Input("(-2,2]").Input("(-∞,  0)").Expect(true)
	cl.New().Input("(-2,2]").Input("(-∞,  2)").Expect(true)
	cl.New().Input("(-2,2]").Input("(-∞,  4)").Expect(true)

	// c
	cl.New().Input("[-2,2)").Input("(-∞, -4]").Expect(false)
	cl.New().Input("[-2,2)").Input("(-∞, -2]").Expect(true)
	cl.New().Input("[-2,2)").Input("(-∞,  0]").Expect(true)
	cl.New().Input("[-2,2)").Input("(-∞,  2]").Expect(true)
	cl.New().Input("[-2,2)").Input("(-∞,  4]").Expect(true)

	cl.New().Input("[-2,2)").Input("(-∞, -4)").Expect(false)
	cl.New().Input("[-2,2)").Input("(-∞, -2)").Expect(false)
	cl.New().Input("[-2,2)").Input("(-∞,  0)").Expect(true)
	cl.New().Input("[-2,2)").Input("(-∞,  2)").Expect(true)
	cl.New().Input("[-2,2)").Input("(-∞,  4)").Expect(true)

	// d
	cl.New().Input("(-2,2)").Input("(-∞, -4]").Expect(false)
	cl.New().Input("(-2,2)").Input("(-∞, -2]").Expect(false)
	cl.New().Input("(-2,2)").Input("(-∞,  0]").Expect(true)
	cl.New().Input("(-2,2)").Input("(-∞,  2]").Expect(true)
	cl.New().Input("(-2,2)").Input("(-∞,  4]").Expect(true)

	cl.New().Input("(-2,2)").Input("(-∞, -4)").Expect(false)
	cl.New().Input("(-2,2)").Input("(-∞, -2)").Expect(false)
	cl.New().Input("(-2,2)").Input("(-∞,  0)").Expect(true)
	cl.New().Input("(-2,2)").Input("(-∞,  2)").Expect(true)
	cl.New().Input("(-2,2)").Input("(-∞,  4)").Expect(true)

	// I
	cl.New().Input("(-∞,+∞)").Input("(-∞, -2)").Expect(true)
	cl.New().Input("(-∞,+∞)").Input("(-∞,  0)").Expect(true)
	cl.New().Input("(-∞,+∞)").Input("(-∞,  2)").Expect(true)
	cl.New().Input("(-∞,+∞)").Input("(-∞, -2]").Expect(true)
	cl.New().Input("(-∞,+∞)").Input("(-∞,  0]").Expect(true)
	cl.New().Input("(-∞,+∞)").Input("(-∞,  2]").Expect(true)

	cl.New().Input("(-∞,+∞)").Input("(-2, +∞)").Expect(true)
	cl.New().Input("(-∞,+∞)").Input("( 0, +∞)").Expect(true)
	cl.New().Input("(-∞,+∞)").Input("( 2, +∞)").Expect(true)
	cl.New().Input("(-∞,+∞)").Input("[-2, +∞)").Expect(true)
	cl.New().Input("(-∞,+∞)").Input("[ 0, +∞)").Expect(true)
	cl.New().Input("(-∞,+∞)").Input("[ 2, +∞)").Expect(true)

	cl.New().Input("(-∞,+∞)").Input("(-∞, +∞)").Expect(true)

	for _, v := range cl.Get() {
		s1 := v.Inputs[0].(string)
		int1, err := Parse(s1)
		gtest.Assert(t, err)
		s2 := v.Inputs[1].(string)
		int2, err := Parse(s2)
		gtest.Assert(t, err)
		expect := v.Expects[0].(bool)

		got := int1.IsOverlap(*int2)

		if got != expect {
			gtest.PrintlnExit(t, "interval %s vs %s, overlap check expect %v but got %v", s1, s2, expect, got)
		}
	}
}


func TestInterval_Contains(t *testing.T) {
	cl := gtest.NewCaseList()

	cl.New().Input("[-2,2]").Input(-3).Expect(false)
	cl.New().Input("[-2,2]").Input(-2).Expect(true)
	cl.New().Input("[-2,2]").Input(0).Expect(true)
	cl.New().Input("[-2,2]").Input(2).Expect(true)
	cl.New().Input("[-2,2]").Input(3).Expect(false)

	cl.New().Input("(-2,2]").Input(-3).Expect(false)
	cl.New().Input("(-2,2]").Input(-2).Expect(false)
	cl.New().Input("(-2,2]").Input(0).Expect(true)
	cl.New().Input("(-2,2]").Input(2).Expect(true)
	cl.New().Input("(-2,2]").Input(3).Expect(false)

	cl.New().Input("[-2,2)").Input(-3).Expect(false)
	cl.New().Input("[-2,2)").Input(-2).Expect(true)
	cl.New().Input("[-2,2)").Input(0).Expect(true)
	cl.New().Input("[-2,2)").Input(2).Expect(false)
	cl.New().Input("[-2,2)").Input(3).Expect(false)

	cl.New().Input("(-2,2)").Input(-3).Expect(false)
	cl.New().Input("(-2,2)").Input(-2).Expect(false)
	cl.New().Input("(-2,2)").Input(0).Expect(true)
	cl.New().Input("(-2,2)").Input(2).Expect(false)
	cl.New().Input("(-2,2)").Input(3).Expect(false)

	cl.New().Input("(-∞,0]").Input(-1).Expect(true)
	cl.New().Input("(-∞,0]").Input(0).Expect(true)
	cl.New().Input("(-∞,0]").Input(1).Expect(false)

	cl.New().Input("(-∞,0)").Input(-1).Expect(true)
	cl.New().Input("(-∞,0)").Input(0).Expect(false)
	cl.New().Input("(-∞,0)").Input(1).Expect(false)

	cl.New().Input("[0,+∞)").Input(-1).Expect(false)
	cl.New().Input("[0,+∞)").Input(0).Expect(true)
	cl.New().Input("[0,+∞)").Input(1).Expect(true)

	cl.New().Input("(0,+∞)").Input(-1).Expect(false)
	cl.New().Input("(0,+∞)").Input(0).Expect(false)
	cl.New().Input("(0,+∞)").Input(1).Expect(true)

	cl.New().Input("(-∞,+∞)").Input(0).Expect(true)

	for _, v := range cl.Get() {
		s1 := v.Inputs[0].(string)
		int1, err := Parse(s1)
		gtest.Assert(t, err)
		n := v.Inputs[1].(int)
		nDec := gdecimal.NewFromInt(n)
		expect := v.Expects[0].(bool)
		got := int1.Contains(nDec)
		if got != expect {
			gtest.PrintlnExit(t, "interval %s, contains check %d expect %v but got %v", s1, n, expect, got)
		}
	}
}