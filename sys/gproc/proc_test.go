package gproc

import (
	"fmt"
	"testing"
)

func TestGetProcInfo(t *testing.T) {
	pid := GetPidOfMyself()
	pi, err := GetProcInfo(pid)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(pi)
}

func TestGetPidByProcName(t *testing.T) {
	ids, err := GetPidByProcName("fontd")
	if err != nil {
		t.Error(err)
		return
	}
	for k, v := range ids {
		fmt.Println(k, v)
	}
}

func TestGetPidCreateTime(t *testing.T) {
	idtms, err := GetProcCreateTime("fontd")
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(idtms)
}
