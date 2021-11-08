package netdrive

import (
	"github.com/cryptowilliam/goutil/basic/gtest"
	"testing"
)

func TestGoogleDrive_DownloadFile(t *testing.T) {
	gd, err := newGoogleDrive("")
	gtest.Assert(t, err)

	b, err := gd.DownloadFile("/abc.txt")
	gtest.Assert(t, err)
	gtest.PrintlnExit(t, string(b))
}
