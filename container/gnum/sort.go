package gnum

import (
	"math/big"
	"sort"
)

type uintSlice []uint

func (us uintSlice) Len() int           { return len(([]uint)(us)) }
func (us uintSlice) Less(i, j int) bool { return ([]uint)(us)[i] < ([]uint)(us)[j] }
func (us uintSlice) Swap(i, j int) {
	([]uint)(us)[i], ([]uint)(us)[j] = ([]uint)(us)[j], ([]uint)(us)[i]
}

type int64Slice []int64

func (i64s int64Slice) Len() int           { return len(([]int64)(i64s)) }
func (i64s int64Slice) Less(i, j int) bool { return ([]int64)(i64s)[i] < ([]int64)(i64s)[j] }
func (i64s int64Slice) Swap(i, j int) {
	([]int64)(i64s)[i], ([]int64)(i64s)[j] = ([]int64)(i64s)[j], ([]int64)(i64s)[i]
}

type bigIntSlice []big.Int

func (bis bigIntSlice) Len() int { return len(([]big.Int)(bis)) }
func (bis bigIntSlice) Less(i, j int) bool {
	return ([]big.Int)(bis)[i].Cmp(&([]big.Int)(bis)[j]) == -1
}
func (bis bigIntSlice) Swap(i, j int) {
	([]big.Int)(bis)[i], ([]big.Int)(bis)[j] = ([]big.Int)(bis)[j], ([]big.Int)(bis)[i]
}

type bigFloatSlice []big.Float

func (bfs bigFloatSlice) Len() int { return len(([]big.Float)(bfs)) }
func (bfs bigFloatSlice) Less(i, j int) bool {
	return ([]big.Float)(bfs)[i].Cmp(&([]big.Float)(bfs)[j]) == -1 // works if number is +Inf -Inf
}
func (bfs bigFloatSlice) Swap(i, j int) {
	([]big.Float)(bfs)[i], ([]big.Float)(bfs)[j] = ([]big.Float)(bfs)[j], ([]big.Float)(bfs)[i]
}

func SortUints(in []uint)             { sort.Sort((uintSlice)(in)) }
func SortInts(in []int)               { sort.Ints(in) }
func SortInt64s(in []int64)           { sort.Sort((int64Slice)(in)) }
func SortFloats(in []float64)         { sort.Float64s(in) }
func SortBigInts(in []big.Int)        { sort.Sort((bigIntSlice)(in)) }
func SortBigFloats(in []big.Float)    { sort.Sort((bigFloatSlice)(in)) }
func ReverseUints(in []uint)          { sort.Sort(sort.Reverse(uintSlice(in))) }
func ReverseInts(in []int)            { sort.Sort(sort.Reverse(sort.IntSlice(in))) }
func ReverseInt64s(in []int64)        { sort.Sort(sort.Reverse(int64Slice(in))) }
func ReverseFloats(in []float64)      { sort.Sort(sort.Reverse(sort.Float64Slice(in))) }
func ReverseBigInts(in []big.Int)     { sort.Sort(sort.Reverse(bigIntSlice(in))) }
func ReverseBigFloats(in []big.Float) { sort.Sort(sort.Reverse(bigFloatSlice(in))) }
