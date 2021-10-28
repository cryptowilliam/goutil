package gtime

import (
	"sort"
	"time"
)

type timeSet struct {
	values []time.Time
}

func (ta *timeSet) Len() int {
	return len(ta.values)
}

// Support sort
func (ta *timeSet) Swap(i, j int) {
	ta.values[i], ta.values[j] = ta.values[j], ta.values[i]

}

// Support sort
func (ta *timeSet) Less(i, j int) bool {
	return ta.values[i].Before(ta.values[j])
}

func SortTimes(in []time.Time) []time.Time {
	tmp := &timeSet{}
	tmp.values = in
	sort.Sort(tmp)
	return tmp.values
}
