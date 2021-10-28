package gcron

import (
	"testing"
	"time"
)

func TestDurCron_CheckNowUnblock(t *testing.T) {
	dcc := NewDurCron(nil, false, time.Second)

	triggerCount := 0
	for i := 0; i < 26; i++ {
		time.Sleep(time.Second / 10)
		if _, triggered := dcc.CheckNowUnblock(); triggered {
			triggerCount++
		}
	}
	if triggerCount != 2 && triggerCount != 3 {
		t.Errorf("CheckNowUnblock error")
		return
	}
}
