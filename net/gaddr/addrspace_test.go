package gaddr

import (
	"fmt"
	"testing"
)

func TestCheckIPString(t *testing.T) {
	ip4lanlist := []string{"10.78.1.2", "192.168.1.12"}
	for _, v := range ip4lanlist {
		if CheckIPString(v) != IPv4_LAN {
			t.Error(fmt.Sprintf("%s test failed", v))
		}
	}

	ip4loopbacklist := []string{"127.0.0.1"}
	for _, v := range ip4loopbacklist {
		if CheckIPString(v) != IPv4_LOOPBACK {
			t.Error(fmt.Sprintf("%s test failed", v))
		}
	}

	ip4anylist := []string{"0.0.0.0"}
	for _, v := range ip4anylist {
		if CheckIPString(v) != IPv4_ANY {
			t.Error(fmt.Sprintf("%s test failed", v))
		}
	}
}
