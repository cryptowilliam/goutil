package gbarcode

import (
	"github.com/cryptowilliam/goutil/sys/gfs"
	"os"
	"testing"
)

func TestDecodeQrCode(t *testing.T) {
	defContent := "0x1234567890abcdef"
	defFilename := "qrcode.gif"

	w, err := os.Create(defFilename)
	if err != nil {
		t.Error(err)
		return
	}
	defer os.Remove(defFilename)

	err = EncodeQrCode(defContent, 200, "gif", w)
	if err != nil {
		t.Error(err)
		return
	}

	r, f, err := gfs.FilenameToReader(defFilename)
	if err != nil {
		t.Error(err)
		return
	}
	result, err := DecodeQrCode(r)
	f.Close()
	if result != defContent {
		t.Error("EncodeQrCode/DecodeQrCode failed")
	}
}
