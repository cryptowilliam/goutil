package gnum

import (
	"fmt"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"strings"
)

/*
func Min(firstNum interface{}, args... interface{}) (interface{}, error) {

	switch fv := firstNum.(type) {
	case int:
		for _ , v := range args{
			if int(fv) > int(v) {
				fv = int(v)
			}
		}
		return fv, nil
	case int8:
		for _ , v := range args{
			if int8(fv) > int8(v) {
				fv = int8(v)
			}
		}
		return fv, nil
	case int16:
		for _ , v := range args{
			if int16(fv) > int16(v) {
				fv = int16(v)
			}
		}
		return fv, nil
	case int32:
		for _ , v := range args{
			if int32(fv) > int32(v) {
				fv = int32(v)
			}
		}
		return fv, nil
	case int64:
		for _ , v := range args{
			if int64(fv) > int64(v) {
				fv = int64(v)
			}
		}
		return fv, nil
	case uint:
		for _ , v := range args{
			if uint(fv) > uint(v) {
				fv = uint(v)
			}
		}
		return fv, nil
	case uint8:
		for _ , v := range args{
			if uint8(fv) > uint8(v) {
				fv = uint8(v)
			}
		}
		return fv, nil
	case uint16:
		for _ , v := range args{
			if uint16(fv) > uint16(v) {
				fv = uint16(v)
			}
		}
		return fv, nil
	case uint32:
		for _ , v := range args{
			if uint32(fv) > uint32(v) {
				fv = uint32(v)
			}
		}
		return fv, nil
	case uint64:
		for _ , v := range args{
			if uint64(fv) > uint64(v) {
				fv = uint64(v)
			}
		}
		return fv, nil
	case big.Int:
		for _ , v := range args{
			if big.Int(fv).Cmp(&big.Int(v)) == 1 {
				fv = big.Int(v)
			}
		}
		return fv, nil
	default:
		return "", gerrors.New("Unsupported type")
	}
	return nil, gerrors.New("Unknown error")
}*/

func FloatsEqual(a, b []float64) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func MinInt(first int, args ...int) int {
	for _, v := range args {
		if first > v {
			first = v
		}
	}
	return first
}

func MinInt64(first int64, args ...int64) int64 {
	for _, v := range args {
		if first > v {
			first = v
		}
	}
	return first
}

func MinUint32(first uint32, args ...uint32) uint32 {
	for _, v := range args {
		if first > v {
			first = v
		}
	}
	return first
}

func MinUintArray(args []uint) uint {
	if len(args) == 0 {
		panic(gerrors.Errorf("MinUintArray nil input args"))
		return 0
	}

	min := args[0]
	for _, v := range args {
		if v < min {
			min = v
		}
	}
	return min
}

func MinUint64(first uint64, args ...uint64) uint64 {
	for _, v := range args {
		if first > v {
			first = v
		}
	}
	return first
}

func MinFloat64(first float64, args ...float64) float64 {
	for _, v := range args {
		if first > v {
			first = v
		}
	}
	return first
}

func MinFloat(args ...float64) float64 {
	if len(args) == 0 {
		panic("nil args input")
	}
	first := args[0]
	for _, v := range args {
		if first > v {
			first = v
		}
	}
	return first
}

func MaxInt(first int, args ...int) int {
	for _, v := range args {
		if first < v {
			first = v
		}
	}
	return first
}

func MaxInt64(first int64, args ...int64) int64 {
	for _, v := range args {
		if first < v {
			first = v
		}
	}
	return first
}

func MaxUint32(first uint32, args ...uint32) uint32 {
	for _, v := range args {
		if first < v {
			first = v
		}
	}
	return first
}

func MaxUint64(first uint64, args ...uint64) uint64 {
	for _, v := range args {
		if first < v {
			first = v
		}
	}
	return first
}

func MaxFloat64(first float64, args ...float64) float64 {
	for _, v := range args {
		if first < v {
			first = v
		}
	}
	return first
}

func MaxFloat(args ...float64) float64 {
	if len(args) == 0 {
		panic("nil args input")
	}
	first := args[0]
	for _, v := range args {
		if first < v {
			first = v
		}
	}
	return first
}

func SumFloat(args ...float64) float64 {
	if len(args) == 0 {
		panic("nil args input")
	}
	sum := float64(0)
	for _, v := range args {
		sum += v
	}
	return sum
}

// 给toBound划界，不超过[min, max]的范围
func BoundUint32(min, toBound, max uint32) uint32 {
	return MinUint32(MaxUint32(min, toBound), max)
}

func RemoveDuplicate(elements []int) []int {
	// Use map to record duplicates as we find them.
	encountered := map[int]bool{}
	result := []int{}

	for v := range elements {
		if encountered[elements[v]] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v]] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}

func JoinFormatInt(sep string, num ...int) string {
	tmp := []string{}
	ns := append([]int{}, num...)
	for _, v := range ns {
		tmp = append(tmp, fmt.Sprintf("%d", v))
	}
	return strings.Join(tmp, sep)
}

func FloatsToFloatsPtr(in []float64) []*float64 {
	r := []*float64{}
	for _, v := range in {
		p := new(float64)
		*p = v
		r = append(r, p)
	}
	return r
}

func FloatsPtrToFloats(in []*float64) []float64 {
	r := []float64{}
	for _, v := range in {
		r = append(r, *v)
	}
	return r
}

func CmpInt64Array(a, b []int64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func RemoveUint(src []uint, toRemove uint, removeCount int) []uint {
	removed := 0
	var r []uint
	for _, v := range src {
		if removed < removeCount && v == toRemove {
			removed++
		} else {
			r = append(r, v)
		}
	}
	return r
}

/*
type (
	MultiFloats struct {
		ns []big.Float
	}
)

func NewMultiFloats(first big.Float, others ...big.Float) *MultiFloats {
	mf := new(MultiFloats)
	mf.ns = append(mf.ns, first)
	for _, v := range others {
		mf.ns = append(mf.ns, v)
	}
	return mf
}

func (mf *MultiFloats) Min() big.Float {
	min := mf.ns[0]

	for _, v := range mf.ns {

	}
}

func (mf *MultiFloats) Max() big.Float {

}

func (mf *MultiFloats) Avg() big.Float {

}*/
