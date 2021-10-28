package gsort

import (
	"fmt"
	"github.com/cryptowilliam/goutil/basic/gtest"
	"testing"
)

type (
	testDataArray struct {
		data [29]rune
	}
)

var (
	testArray = testDataArray{data: [29]rune{'b', 'c', 'a', 'v', 'j', 'w', 'p', 'o', 'g', 'h', 's', 'e', 'r', 'x', 'q', 'z', 'f', 'k', 'l', 'd', 'm', 't', 'n', 'i', 'u', 'y', 'U', 'S', 'A'}}
)

func (td *testDataArray) Len() int {
	return len(td.data)
}

func (td *testDataArray) Less(i, j int) bool {
	return td.data[i] < td.data[j]
}

func (td *testDataArray) Swap(i, j int) {
	td.data[i], td.data[j] = td.data[j], td.data[i]
}

func TestSort(t *testing.T) {
	err := Sort(QuickSort, &testArray)
	gtest.Assert(t, err)
	fmt.Println(testArray)
}
