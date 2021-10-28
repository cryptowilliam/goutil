package gfs

import (
	"fmt"
	"os"
	"testing"
)

func TestMakeDir(t *testing.T) {
	if err := MakeDir("abc"); err != nil {
		t.Error(err)
		return
	}
	if err := RemoveDir("abc"); err != nil {
		t.Error(err)
		return
	}
}

func TestDirSize(t *testing.T) {
	s, err := os.Stat("")
	fmt.Println(err)
	fmt.Println(s.Size())
}
