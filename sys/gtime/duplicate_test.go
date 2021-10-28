package gtime

import (
	"testing"
	"time"
)

func TestHasDuplicated(t *testing.T) {
	now := time.Now()
	var times []time.Time
	for i := 0; i < 1000; i++ {
		now = now.Add(time.Nanosecond)
		times = append(times, now)
	}

	has, err := HasDuplicated(times, time.Nanosecond)
	if err != nil {
		t.Error(err)
		return
	}
	if has != false {
		t.Errorf("HasDuplicated Nanosecond should be false")
		return
	}

	has, err = HasDuplicated(times, time.Millisecond)
	if err != nil {
		t.Error(err)
		return
	}
	if has != true {
		t.Errorf("HasDuplicated Millisecond should be true")
		return
	}
}
