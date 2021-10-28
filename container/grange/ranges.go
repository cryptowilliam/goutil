package grange

import (
	"fmt"
	"github.com/cryptowilliam/goutil/container/gnum"
	"sort"
	"strings"
)

// this is a data structure to describe a large amount
// and ALMOST consequent integer numbers, instead of huge integer array.
// replace huge memory cost with CPU.

type Range struct {
	Begin int64
	End   int64
}

func NewRange(a, b int64) Range {
	return Range{Begin: gnum.MinInt64(a, b), End: gnum.MaxInt64(a, b)}
}

// can number n can join current Range and become one Range but not 2 separate items
// ir [5,10]
// n = 3: false
// n = 4: true
// n = 5: true
// n = 10: true
// n = 11: true
// n = 12: false
func (ir *Range) IsConnected(cmp int64) bool {
	return cmp >= ir.Begin-1 && cmp <= ir.End+1
}

func (ir *Range) IsConnectedEx(cmp Range) bool {
	return ir.IsConnected(cmp.Begin) || ir.IsConnected(cmp.End)
}

// is number n in current range
// ir [5,10]
// n = 3: false
// n = 4: false
// n = 5: true
// n = 10: true
// n = 11: false
// n = 12: false
func (ir *Range) IsOverlap(cmp int64) bool {
	return cmp >= ir.Begin && cmp <= ir.End
}

func (ir *Range) IsOverlapEx(cmp Range) bool {
	return ir.IsOverlap(cmp.Begin) || ir.IsOverlap(cmp.End) || cmp.IsOverlap(ir.Begin) || cmp.IsOverlap(ir.End)
}

func (ir *Range) Equal(cmp Range) bool {
	return ir.Begin == cmp.Begin && ir.End == cmp.End
}

func (ir *Range) Len() int64 {
	return ir.End - ir.Begin + 1
}

func (ir *Range) String() string {
	return fmt.Sprintf("[%d,%d]", ir.Begin, ir.End)
}

type RangeFilter struct {
	Ranges []Range
}

func NewRangeFilter() *RangeFilter {
	return &RangeFilter{}
}

func NewRangeFilterEx(begin, end int64) *RangeFilter {
	r := &RangeFilter{}
	r.AddRange(NewRange(begin, end))
	return r
}

func (irf *RangeFilter) DataSize() int64 {
	sz := int64(0)
	for _, rg := range irf.Ranges {
		sz += rg.Len()
	}
	return sz
}

func (irf *RangeFilter) Len() int {
	return len(irf.Ranges)
}

func (irf *RangeFilter) Swap(i, j int) {
	irf.Ranges[i], irf.Ranges[j] = irf.Ranges[j], irf.Ranges[i]
}

func (irf *RangeFilter) Less(i, j int) bool {
	return irf.Ranges[i].Begin < irf.Ranges[j].Begin
}

// sort, try to join IntRanges to one if them can
func (irf *RangeFilter) rebuildRanges() {

	sort.Sort(irf)

	oldRanges := irf.Ranges
	var newRanges []Range

	for k, v := range oldRanges {
		if k == 0 {
			newRanges = append(newRanges, v)
		} else {
			// compare to last item
			if v.Begin <= newRanges[len(newRanges)-1].End+1 {
				newRanges[len(newRanges)-1].End = v.End
			} else {
				newRanges = append(newRanges, v)
			}
		}
	}

	irf.Ranges = newRanges
}

func (irf *RangeFilter) AddInt64(n int64) {
	// compare to exists items
	for k, v := range irf.Ranges {

		if n >= v.Begin && n <= v.End {
			return
		}

		if n == v.Begin-1 {
			irf.Ranges[k].Begin = n
			irf.rebuildRanges()
			return
		}

		if n == v.End+1 {
			irf.Ranges[k].End = n
			irf.rebuildRanges()
			return
		}
	}

	// no items or not belongs to any exists item
	irf.Ranges = append(irf.Ranges, Range{Begin: n, End: n})
	irf.rebuildRanges()
	return
}

func (irf *RangeFilter) AddRange(in Range) {
	irf.Ranges = append(irf.Ranges, in)
	irf.rebuildRanges()
}

func (irf *RangeFilter) SubRange(toSub Range) {

	oldRanges := irf.Ranges
	var newRanges []Range

	for _, v := range oldRanges {

		// do nothing and keep this
		if !v.IsOverlapEx(toSub) {
			newRanges = append(newRanges, v)
			continue
		}

		// ignore(remove) this one
		if toSub.Begin <= v.Begin && toSub.End >= v.End {
			continue
		}

		// sub and keep 2 sub results
		if toSub.Begin > v.Begin && toSub.End < v.End {
			newRanges = append(newRanges, NewRange(v.Begin, toSub.Begin-1))
			newRanges = append(newRanges, NewRange(toSub.End+1, v.End))
			continue
		}

		// sub and keep 1 sub result
		if toSub.Begin <= v.Begin {
			newRanges = append(newRanges, NewRange(toSub.End+1, v.End))
		} else {
			newRanges = append(newRanges, NewRange(v.Begin, toSub.Begin-1))
		}

	}

	irf.Ranges = newRanges
}

func (irf *RangeFilter) Add(in RangeFilter) {
	for _, v := range in.Ranges {
		irf.AddRange(v)
	}
}

func (irf *RangeFilter) Sub(in RangeFilter) {
	for _, v := range in.Ranges {
		irf.SubRange(v)
	}
}

func (irf *RangeFilter) MinMax() (int64, int64, bool) {
	if irf.Len() == 0 {
		return 0, 0, false
	}
	return irf.Ranges[0].Begin, irf.Ranges[irf.Len()-1].End, true
}

// get not linked/connected ranges of current RangeFilter
func (irf *RangeFilter) Lacks() *RangeFilter {
	minNum, maxNum, ok := irf.MinMax()
	if !ok {
		return NewRangeFilter()
	}

	full := NewRangeFilter()
	full.AddRange(NewRange(minNum, maxNum))
	full.Sub(*irf)
	return full
}

func (irf *RangeFilter) Has(n int64) bool {
	for _, v := range irf.Ranges {
		if n >= v.Begin && n <= v.End {
			return true
		}
	}
	return false
}

func (irf *RangeFilter) String() string {
	var ss []string
	for _, v := range irf.Ranges {
		ss = append(ss, fmt.Sprintf("[%d,%d]", v.Begin, v.End))
	}
	return strings.Join(ss, " ")
}
