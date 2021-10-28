package gnum

import "math/bits"

const _MaxUint_ = ^uint(0)
const _MinUint_ = 0
const _MaxInt_ = int(_MaxUint_ >> 1)
const _MinInt_ = -_MaxInt_ - 1

const (
	_MaxUint_2_ uint = (1 << bits.UintSize) - 1
	_MaxInt_2_  int  = (1<<bits.UintSize)/2 - 1
	_MinInt_2_  int  = (1 << bits.UintSize) / -2
)
