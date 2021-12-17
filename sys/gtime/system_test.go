package gtime

import (
	"fmt"
	"testing"
	"time"
)

func TestIsEpochBeginning(t *testing.T) {
	tm := EpochBeginTime
	if !IsEpochBeginning(tm) {
		t.Error("IsEpochBeginning test failed")
		return
	}
	tm = tm.In(time.UTC)
	if !IsEpochBeginning(tm) {
		t.Error("IsEpochBeginning test failed")
		return
	}
	tm = time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC)
	fmt.Println(tm.Second(), tm.Nanosecond())
	if IsEpochBeginning(tm) {
		t.Error("IsEpochBeginning test failed")
		return
	}
}

func TestUptime(t *testing.T) {
	fmt.Println(Uptime().String())
}
