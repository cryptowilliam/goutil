package grand

import (
	"fmt"
	"github.com/cryptowilliam/goutil/basic/gtest"
	"testing"
)

func TestRanges_Generate(t *testing.T) {
	rr := NewRanges().Allow('a', 'z').Allow('A', 'Z').Allow('0', '9')
	for i := 0; i < 62; i++ {
		fmt.Println(rr.Generate(32))
	}
}

func TestInt(t *testing.T) {
	for i := 0; i < 100; i++ {
		fmt.Println(Int(0, 2))
	}
}

func TestRandomString(t *testing.T) {
	s := RandomString(9)
	if len(s) != 9 {
		t.Errorf("RandomString error")
		return
	}
}

func TestRandomBuffer(t *testing.T) {
	buf := make([]byte, 10)
	_, err := RandomBuffer(buf)
	gtest.Assert(t, err)
}