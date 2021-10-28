package gfs

import (
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/google/uuid"
	"runtime"
)

// Generate a new temp filename for cache
func NewTempFilename() (string, error) {
	if runtime.GOOS == "windows" {
		return "", gerrors.New("Unsupport windows for now")
	} else {
		return "/tmp/" + uuid.New().String() + ".temp", nil
	}
}
