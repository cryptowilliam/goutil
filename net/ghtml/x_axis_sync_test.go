package ghtml

import (
	"encoding/json"
	"fmt"
	"github.com/cryptowilliam/goutil/sys/gtime"
	"sort"
	"testing"
	"time"
)

type (
	_dot_ struct {
		Time  time.Time
		Value *float64
	}

	_dots4test_ struct {
		Items  []_dot_
		sorted bool
	}
)

func (ds *_dots4test_) Len() int {
	return len(ds.Items)
}

// Support sort
func (ds *_dots4test_) Swap(i, j int) {
	ds.Items[i], ds.Items[j] = ds.Items[j], ds.Items[i]
}

// Support sort
func (ds *_dots4test_) Less(i, j int) bool {
	return ds.Items[i].Time.Before(ds.Items[j].Time)
}

func (ds *_dots4test_) String() string {
	buf, err := json.Marshal(ds)
	if err != nil {
		return ""
	}
	return string(buf)
}

func (ds *_dots4test_) Add(d _dot_) {
	ds.Items = append(ds.Items, d)
}

func (ds *_dots4test_) Sort() {
	sort.Sort(ds)
	ds.sorted = true
}

func (ds *_dots4test_) Times() []time.Time {
	r := []time.Time{}
	for _, v := range ds.Items {
		r = append(r, v.Time)
	}
	return r
}

func (ds *_dots4test_) Values() []float64 {
	r := []float64{}
	for _, v := range ds.Items {
		r = append(r, *v.Value)
	}
	return r
}

type doubleDot struct {
	a *_dot_
	b *_dot_
}

type doubleDotList []doubleDot

func (dkl doubleDotList) Len() int {
	return len([]doubleDot(dkl))
}

func (dkl doubleDotList) Less(i, j int) bool {
	return []doubleDot(dkl)[i].a.Time.Before([]doubleDot(dkl)[j].a.Time)
}

func (dkl doubleDotList) Swap(i, j int) {
	[]doubleDot(dkl)[i], []doubleDot(dkl)[j] = []doubleDot(dkl)[j], []doubleDot(dkl)[i]
}

func TestLinesSync_Sync(t *testing.T) {
	ds1 := new(_dots4test_)
	ds2 := new(_dots4test_)

	newFloat := func(v float64) *float64 {
		return &v
	}

	date, _ := gtime.NewDate(2019, 1, 1)
	ds1.Items = append(ds1.Items, _dot_{Time: date.ToTime(0, 0, 0, 0, time.UTC), Value: newFloat(1)})
	date, _ = gtime.NewDate(2019, 1, 2)
	ds1.Items = append(ds1.Items, _dot_{Time: date.ToTime(0, 0, 0, 0, time.UTC), Value: newFloat(2)})
	date, _ = gtime.NewDate(2019, 1, 3)
	ds1.Items = append(ds1.Items, _dot_{Time: date.ToTime(0, 0, 0, 0, time.UTC), Value: newFloat(3)})

	date, _ = gtime.NewDate(2019, 1, 2)
	ds2.Items = append(ds2.Items, _dot_{Time: date.ToTime(0, 0, 0, 0, time.UTC), Value: newFloat(2)})
	date, _ = gtime.NewDate(2019, 1, 3)
	ds2.Items = append(ds2.Items, _dot_{Time: date.ToTime(0, 0, 0, 0, time.UTC), Value: newFloat(3)})
	date, _ = gtime.NewDate(2019, 1, 4)
	ds2.Items = append(ds2.Items, _dot_{Time: date.ToTime(0, 0, 0, 0, time.UTC), Value: newFloat(4)})

	ls := NewXAxisSync()
	if err := ls.AddFloats("dots1", ds1.Times(), ds1.Values()); err != nil {
		t.Log(err)
		return
	}
	if err := ls.AddFloats("dots2", ds2.Times(), ds2.Values()); err != nil {
		t.Log(err)
		return
	}
	ls.Sync()

	tms := ls.GetTimes()
	fmt.Println(tms)
	itfVals1 := ls.GetValues("dots1")
	itfVals2 := ls.GetValues("dots2")
	fmt.Println("\ndots1:")
	for i := range itfVals1 {
		if itfVals1[i] != nil {
			fmt.Println(tms[i], *itfVals1[i].(*float64))
		} else {
			fmt.Println(tms[i], "null")
		}
	}
	fmt.Println("\ndots2:")
	for i := range itfVals2 {
		if itfVals2[i] != nil {
			fmt.Println(tms[i], *itfVals2[i].(*float64))
		} else {
			fmt.Println(tms[i], "null")
		}
	}

	if itfVals1[3] != nil {
		t.Errorf("sync error")
	}
	if itfVals2[0] != nil {
		t.Errorf("sync error")
	}
}
