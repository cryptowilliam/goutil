package ghdd

import (
	"encoding/json"
	"github.com/cryptowilliam/goutil/sys/gproc"
	"log"
	"testing"
)

func TestGetVolumeInfo(t *testing.T) {
	mydir, err := gproc.SelfPath()
	if err != nil {
		t.Error(err)
		return
	}

	vi, err := GetVolumeInfo(mydir)
	if err != nil {
		t.Error(err)
		return
	}
	buf, err := json.Marshal(vi)
	if err != nil {
		t.Error(err)
		return
	}
	log.Print(string(buf))
}
