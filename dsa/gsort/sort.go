package gsort

import (
	"fmt"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"sort"
)

/**
sort —— 排序算法
Go内置的sort包实现了四种基本排序算法：插入排序、归并排序、堆排序和快速排序。
但是这四种排序方法是不公开的，它们只被用于sort包内部使用。所以在对数据集合排序时不必考虑应当选择哪一种排序方法，
只要实现了sort.Interface定义的三个方法：获取数据集合长度的Len()方法、比较两个元素大小的Less()方法和交换两个元素位置的Swap()方法，就可以顺利对数据集合进行排序。
sort包会根据实际数据自动选择高效的排序算法。 除此之外，为了方便对常用数据类型的操作，sort包提供了对[]int切片、[]float64切片和[]string切片完整支持
*/

type (
	SortAlgo string
)

const (
	BubbleSort    = SortAlgo("bubbleSort")
	SelectionSort = SortAlgo("selectionSort")
	InsertionSort = SortAlgo("insertionSort")
	ShellSort     = SortAlgo("shellSort")
	QuickSort     = SortAlgo("quickSort")
	HeapSort      = SortAlgo("heapSort")
	MergeSort     = SortAlgo("mergeSort")
	CountSort     = SortAlgo("countSort")
	BinSort       = SortAlgo("binSort")
	RadixSort     = SortAlgo("radixSort")
)

func Sort(algo SortAlgo, data sort.Interface) error {
	switch algo {
	case BubbleSort:
		bubbleSort(data)
		return nil
	case SelectionSort:
		selectionSort(data)
		return nil
	case InsertionSort:
		insertionSort(data)
		return nil
	case ShellSort:
		shellSort(data)
		return nil
	case QuickSort:
		quickSort(data)
		return nil
	default:
		return gerrors.New("unsupported sort %s", algo)
	}
}

// The bubble sort repeatedly walks through the column of elements to be sorted, comparing two adjacent elements in turn.
// If the order is wrong, swap these two elements.
// The work of visiting elements is repeated until there are no neighboring elements to swap, i.e. the element column is already sorted.
// The name of the algorithm comes from the fact that the smaller elements are swapped and slowly "float" to the top of the sequence.
// Just as the bubbles of carbon dioxide in carbonated beverages eventually rise to the top, hence the name "bubble sort".
// The complexity of bubble sort is O(n) ~ O(n^2).
func bubbleSort(data sort.Interface) {
	for {
		swapCount := 0
		for i := 0; i < data.Len()-1; i++ {
			if data.Less(i+1, i) {
				data.Swap(i, i+1)
				swapCount++
			}
		}
		if swapCount == 0 {
			return
		}
	}
}

// Selection sort works by selecting the smallest element from the data elements to be sorted, store it at the beginning of the sequence.
// Then find the smallest element from the remaining unsorted elements and place it at the end of the sorted sequence.
// And so on, until count of data elements to be sorted are zero.
// The complexity of selection sort is O(n) ~ O(n^2).
func selectionSort(data sort.Interface) {
	for i := 0; i < data.Len(); i++ {
		minEleIdx := i
		for j := i; j < data.Len(); j++ {
			if data.Less(j, minEleIdx) {
				minEleIdx = j
			}
		}
		data.Swap(minEleIdx, i)
	}
}

// Insert sort, traversing backwards from the second element, for each element E, the following operation is performed:
// For all elements in front of element E, from back to front, compare the value of E and each element I in the list,
// if E is smaller than I, swap E and I, otherwise end the operation.
func insertionSort(data sort.Interface) {
	for i := 1; i < data.Len(); i++ {
		for j := i; j >= 1; j-- {
			if data.Less(j, j-1) {
				data.Swap(j, j-1)
			} else {
				break
			}
		}
	}
}

// Shell sort
// The length of the data to be sorted is L, then there are the following gaps L/2, L/4, L/8 .... 1,
// perform the following operations on the above gaps in order:
// Branches the data in gap, similar to a matrix, if the length is odd then the last element is in a separate row,
// treat each column of data as a separate piece of data to be sorted by insertion sort algorithm.
func shellSort(data sort.Interface) {
	insertionSortWithGap := func(data sort.Interface, startIdx int, gap int) {
		for i := gap; i < data.Len(); i += gap {
			for j := i; j >= gap; j -= gap {
				if data.Less(j, j-gap) {
					data.Swap(j, j-gap)
				} else {
					break
				}
			}
		}
	}

	for gap := data.Len() / 2; gap > 0; gap /= 2 {
		for i := 0; i < gap; i++ {
			insertionSortWithGap(data, i, gap)
		}
	}
}

// Quick sort
// 将两个指针i，j分别指向区间的起始和最后位置，取区间的首元素作为基准值T，反复操作以下步骤S:
// (1) j逐渐减小，并逐次比较j指向的元素和T的大小，若p(j)<T则交换位置，也就是放到T的右边
// (2) i逐渐增大，并逐次比较i指向的元素和T的大小，若p(i)>T则交换位置，也就是放到T的左边
// (3) 直到i，j指向同一个值，循环结束
// 若待排序数据的长度为L，那么按如下区间[L,L/2,L/4...1)，对区间内数据进行排序，若L为奇数，则最后一个元素拼接放入最后一个区间
func quickSort(data sort.Interface) {
	quickSortWithRegion := func(data sort.Interface, regionLeft, regionRight int) {
		fmt.Println(regionLeft, regionRight)
		i := regionLeft
		j := regionRight
		baselineIndex := regionLeft
		rightTurn := true
		for i < j {
			if rightTurn {
				if data.Less(j, baselineIndex) {
					data.Swap(j, baselineIndex)
					baselineIndex = j
					rightTurn = !rightTurn
				}
				j--
			} else {
				if data.Less(baselineIndex, i) {
					data.Swap(baselineIndex, i)
					baselineIndex = i
					rightTurn = !rightTurn
				}
				i++
			}
		}
	}

	for regionLen := data.Len(); regionLen > 1; regionLen /= 2 {
		for regionSN := 0; regionSN < data.Len()/regionLen; regionSN++ {
			regionLeft := regionSN * regionLen
			regionRight := regionLeft + regionLen - 1
			if regionLen < data.Len() && regionSN == data.Len()/regionLen-1 && data.Len()%2 > 0 { // 如果是最后一个区间，而且待排序数据长度为奇数
				regionRight++
			}
			quickSortWithRegion(data, regionLeft, regionRight)
		}
	}
}
