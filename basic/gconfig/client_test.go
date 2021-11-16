package gconfig

import (
	"github.com/cryptowilliam/goutil/basic/gtest"
	"github.com/cryptowilliam/goutil/encoding/gjson"
	"github.com/google/uuid"
	"os"
	"testing"
)

func TestClient_Load(t *testing.T) {
	type Sample struct {
		AEncryptMe string
		B          string
	}

	s1 := Sample{AEncryptMe: "A", B: "B"}
	s2 := Sample{}
	randStr := uuid.New().String()+".json"

	cc, err := NewClient("")
	gtest.Assert(t, err)

	defer func() {
		_ = os.Remove(cc.getConfigFilePath(randStr))
	}()

	cc.SetPassword("pwd", "nonce")
	err = cc.Store(randStr, &s1)
	gtest.Assert(t, err)
	err = cc.Load(randStr, &s2, false)
	gtest.Assert(t, err)
	if gjson.MarshalStringDefault(s1, false) != gjson.MarshalStringDefault(s2, false) {
		gtest.PrintlnExit(t, "gconfig.Marshal != gconfig.Unmarshal")
		return
	}
}
