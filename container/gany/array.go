package gany

import (
	"math"
)

/*

https://github.com/golang/go/wiki/InterfaceSlice

错误示例
var dataSlice []int = foo()
var interfaceSlice []interface{} = dataSlice
This gets the error

cannot use dataSlice (type []int) as type []interface { } in assignment


正确示例
var dataSlice []int = foo()
var interfaceSlice []interface{} = make([]interface{}, len(dataSlice))
for i, d := range dataSlice {
	interfaceSlice[i] = d
}
*/

// split all interfaces to multiple []interface that count is blockCount
// sample:
// SplitByBlockCount([1, 2, 3, 4, 5], 2) = [ [1, 2, 3], [4, 5] ]
func SplitByBlockCount(all []interface{}, blockCount int) [][]interface{} {
	if len(all) == 0 || blockCount <= 0 {
		return [][]interface{}{}
	}
	blockSize := int(math.Ceil(float64(len(all)) / float64(blockCount)))
	return SplitByBlockSize(all, blockSize)
}

// split all interfaces to multiple []interface that each size <= blockSize
// sample:
// SplitByBlockSize([1, 2, 3, 4, 5], 2) = [ [1, 2], [3, 4], [5] ]
func SplitByBlockSize(all []interface{}, blockSize int) [][]interface{} {
	if len(all) == 0 {
		return [][]interface{}{}
	}
	if blockSize <= 0 || blockSize >= len(all) {
		return [][]interface{}{all}
	}

	var r [][]interface{}
	processedLen := 0
	for processedLen < len(all) {
		if processedLen+blockSize >= len(all) {
			blockSize = len(all) - processedLen
		}
		r = append(r, all[processedLen:processedLen+blockSize])
		processedLen += blockSize
	}
	return r
}
