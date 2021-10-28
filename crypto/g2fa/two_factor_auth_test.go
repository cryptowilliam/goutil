package g2fa

import (
	"fmt"
	"github.com/cryptowilliam/goutil/basic/gtest"
	"testing"
)

func TestTwoFactorAuth(t *testing.T) {
	pwd, secondsRemaining, err := TwoFactorAuth("nzxxiidbebvwk6jb")
	gtest.Assert(t, err)
	if secondsRemaining < 0 || secondsRemaining > 30 {
		gtest.PrintlnExit(t, fmt.Sprintf("invalid secondsRemaining %d", secondsRemaining))
	}
	if len(pwd) != 6 {
		gtest.PrintlnExit(t, fmt.Sprintf("invalid 2FA password %s", pwd))
	}
}
