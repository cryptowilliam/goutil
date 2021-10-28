package gnum

import (
	"math"
	"strconv"
)

/*
个位(数) = units (digit)
十位(数) = tens (digit)
百位(数) = hundreds (digit)
千位(数) = thousands (digit)
万位(数) = ten thousands (digit)
十万位(数) = hundred thousands (digit)
*/

// 199012 -> 2
func GetLast1(v int) int {
	return v % 10
}

// 199012 -> 12
func GetLast2(v int) int {
	return v % 100
}

// 199012 -> 19901
func RemoveLast1(v int) int {
	return v / 10
}

// 199012 -> 1990
func RemoveLast2(v int) int {
	return v / 100
}

// SliceBetweenEqual(19901201, 4, 2) -> 12
// SliceBetweenEqual(-19901201, 4, 2) -> 12
func Slice(v int64, begin, len int) (int64, error) {
	if v < 0 {
		v = -v
	}
	s := strconv.FormatInt(v, 10)
	s = s[begin:len]
	return strconv.ParseInt(s, 10, 64)
}

func ContainsInt(src []int, toFind int) bool {
	for _, v := range src {
		if v == toFind {
			return true
		}
	}
	return false
}

func IsOddInt64(n int64) bool {
	absNum := int(math.Abs(float64(n)))
	return absNum%2 != 0
}

func IsEvenInt64(n int64) bool {
	return !IsOddInt64(n)
}

func IsOddUint64(n uint64) bool {
	absNum := int(math.Abs(float64(n)))
	return absNum%2 != 0
}

func IsEvenUint64(n uint64) bool {
	return !IsOddUint64(n)
}
