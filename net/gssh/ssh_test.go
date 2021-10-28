package gssh

import (
	"fmt"
	"github.com/cryptowilliam/goutil/basic/gtest"
	"testing"
)

func TestClient_RunCommand(t *testing.T) {
	ssh, err := Dial("xx.xx.xx.xx", "root", "", "*.pem", "")
	gtest.Assert(t, err)
	out, err := ssh.RunCommand("ls")
	gtest.Assert(t, err)
	fmt.Print(out)
}
