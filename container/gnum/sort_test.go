package gnum

import (
	"fmt"
	"math"
	"math/big"
	"strings"
	"testing"
)

func TestSortUints(t *testing.T) {
	in := []uint{1, 3, 5, 7, 9, 8, 0, 6, 4, 2} //{1.1, 3.3, 5.5, 7.7, 9.9, 8.8, 0.0, 6.6, 4.4, 2.2}
	SortUints(in)
	if fmt.Sprintf("%v", in) != "[0 1 2 3 4 5 6 7 8 9]" {
		t.Errorf("SortUints failed")
		return
	}
}

func TestSortBigInts(t *testing.T) {
	in := []big.Int{*big.NewInt(1), *big.NewInt(3), *big.NewInt(5), *big.NewInt(7), *big.NewInt(9), *big.NewInt(8), *big.NewInt(0), *big.NewInt(6), *big.NewInt(4), *big.NewInt(2)}
	SortBigInts(in)
	var ss []string
	for _, v := range in {
		ss = append(ss, v.String())
	}
	if strings.Join(ss, ", ") != "0, 1, 2, 3, 4, 5, 6, 7, 8, 9" {
		t.Errorf("SortBigInts failed")
		return
	}
}

func TestSortBigFloats(t *testing.T) {
	in := []big.Float{*big.NewFloat(1.1), *big.NewFloat(3.3), *big.NewFloat(5.5), *big.NewFloat(7.7), *big.NewFloat(9.9), *big.NewFloat(8.8), *big.NewFloat(0.0), *big.NewFloat(6.6), *big.NewFloat(4.4), *big.NewFloat(2.2)}
	SortBigFloats(in)
	var ss []string
	for _, v := range in {
		ss = append(ss, v.String())
	}
	if strings.Join(ss, ", ") != "0, 1.1, 2.2, 3.3, 4.4, 5.5, 6.6, 7.7, 8.8, 9.9" {
		t.Errorf("SortBigFloats failed")
		return
	}

	in = []big.Float{*big.NewFloat(math.Inf(1)), *big.NewFloat(2), *big.NewFloat(math.Inf(-1))}
	SortBigFloats(in)
	ss = nil
	for _, v := range in {
		ss = append(ss, v.String())
	}
	if strings.Join(ss, ", ") != "-Inf, 2, +Inf" {
		t.Errorf("SortBigFloats failed2")
		return
	}
}
