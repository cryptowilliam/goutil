// +build !windows,!linux

package gmtu

import "github.com/cryptowilliam/goutil/basic/gerrors"

func check(addr string, size int) (bool, int, error) {
	return false, 0, gerrors.ErrNotImplemented
}
