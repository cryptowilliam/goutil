package gexcel

import (
	"github.com/cryptowilliam/goutil/encoding/gjson"
	"github.com/cryptowilliam/goutil/sys/gfs"
	"testing"
)

func TestMemDoc_ToXlsx(t *testing.T) {
	s := "test.xlsx"
	xd, err := OpenPath(s)
	if err != nil {
		t.Error(err)
	}

	b, err := xd.ToMemDoc(true).ToXlsx()
	if err != nil {
		t.Error(err)
	}
	gfs.BytesToFile(b, s+"2.xlsx")
}

func TestXlsDoc_ToMemDoc(t *testing.T) {
	s := "test.xlsx"
	xd, err := OpenPath(s)
	if err != nil {
		t.Error(err)
	}

	b := xd.ToMemDoc(true)
	if err != nil {
		t.Error(err)
	}
	gfs.StringToFile(gjson.MarshalStringDefault(b, true), "test_in_memory.json")
}
