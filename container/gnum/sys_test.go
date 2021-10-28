package gnum

import (
	"github.com/cryptowilliam/goutil/basic/gtest"
	"github.com/cryptowilliam/goutil/container/grand"
	"math/big"
	"strconv"
	"testing"
)

func TestFormatAnySys(t *testing.T) {
	for system := 2; system <= 36; system++ {
		for i := 0; i < 100; i++ {

			zeroString := FormatAnySys(GetDefaultDigits()[:system], *big.NewInt(int64(0)))
			if zeroString != "0" {
				gtest.PrintlnExit(t, "under number system %d, number 0 expect 0 but got %s", system, zeroString)
			}

			pos1String := FormatAnySys(GetDefaultDigits()[:system], *big.NewInt(int64(1)))
			if pos1String != "1" {
				gtest.PrintlnExit(t, "under number system %d, number 1 expect 1 but got %s", system, pos1String)
			}

			neg1String := FormatAnySys(GetDefaultDigits()[:system], *big.NewInt(int64(-1)))
			if neg1String != "-1" {
				gtest.PrintlnExit(t, "under number system %d, number -1 expect -1 but got %s", system, neg1String)
			}

			randNum := grand.Int(-10000000, 10000000)
			got := FormatAnySys(GetDefaultDigits()[:system], *big.NewInt(int64(randNum)))
			expect := strconv.FormatInt(int64(randNum), system)
			if expect != got {
				gtest.PrintlnExit(t, "under number system %d, number %d expect %s but got %s", system, randNum, expect, got)
			}
		}
	}
}

func TestParseAnySys(t *testing.T) {
	for system := 2; system <= 36; system++ {
		for i := 0; i < 100; i++ {

			zero, err := ParseAnySys(GetDefaultDigits()[:system], "0")
			gtest.Assert(t, err)
			if zero.Int64() != 0 {
				gtest.PrintlnExit(t, "under number system %d, string '0' expect 0 but got %s", system, zero.String())
			}

			p1, err := ParseAnySys(GetDefaultDigits()[:system], "1")
			gtest.Assert(t, err)
			if p1.Int64() != 1 {
				gtest.PrintlnExit(t, "under number system %d, string '1' expect 1 but got %s", system, p1.String())
			}

			n1, err := ParseAnySys(GetDefaultDigits()[:system], "-1")
			gtest.Assert(t, err)
			if n1.Int64() != -1 {
				gtest.PrintlnExit(t, "under number system %d, string '-1' expect -1 but got %s", system, n1.String())
			}

			randNum := grand.Int(-10000000, 10000000)
			randStr := strconv.FormatInt(int64(randNum), system)
			got, err := ParseAnySys(GetDefaultDigits()[:system], randStr)
			gtest.Assert(t, err)
			expect := randNum
			if int64(expect) != got.Int64() {
				gtest.PrintlnExit(t, "under number system %d, string %s expect %d but got %s", system, randStr, expect, got.String())
			}
		}
	}
}