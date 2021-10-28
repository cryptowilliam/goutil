package gnum

import (
	"math/big"
)

func BigIntToFloat64(it *big.Int) float64 {
	num, _ := new(big.Float).SetInt(it).Float64()
	return num
}

func BiggerThan(left, right *big.Int) bool {
	return left.Cmp(right) > 0
}

func BiggerEqualThan(left, right *big.Int) bool {
	return left.Cmp(right) >= 0
}

func SmallerThan(left, right *big.Int) bool {
	return left.Cmp(right) < 0
}

func SmallerEqualThan(left, right *big.Int) bool {
	return left.Cmp(right) <= 0
}

func IsZero(val *big.Int) bool {
	return val.Cmp(big.NewInt(0)) == 0
}
