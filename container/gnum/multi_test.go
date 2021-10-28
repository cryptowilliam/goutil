package gnum

import "testing"

func TestMin(t *testing.T) {
	var a, b, c int
	a = 1
	b = 2
	c = 3
	if MinInt(a, b, c) != 1 {
		t.Error("MinInt test failed")
	}
}
