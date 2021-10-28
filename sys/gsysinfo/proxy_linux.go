package gsysinfo

import (
	"github.com/cryptowilliam/goutil/basic/gerrors"
)

func GetSystemProxy() (string, bool, error) {
	// TODO
	return "", false, gerrors.Errorf("unsupported for now")
}

func SetSystemProxy(defaultServer string, enabled bool) error {
	// TODO
	return gerrors.Errorf("unsupported for now")
}
