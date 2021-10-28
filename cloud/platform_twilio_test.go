package cloud

import (
	"github.com/cryptowilliam/goutil/basic/gtest"
	"testing"
	"time"
)

func TestTwilio_SmsSendMsg(t *testing.T) {
	tw, err := newTwilio("", "")
	gtest.Assert(t, err)
	err = tw.SmsSendMsg("", "", "test msg from twilio api at "+time.Now().String())
	gtest.Assert(t, err)
}
