package gnet

import (
	"fmt"
	"testing"
)

func TestParseDomain(t *testing.T) {
	_, err := ParseDomain("")
	if err == nil {
		t.Error("ParseDomain nil string test failed")
	}

	_, err = ParseDomain("google.hotels2")
	if err == nil {
		t.Error("ParseDomain google.hotels2 test failed")
		return
	}

	dm, err := ParseDomain(".co")
	if err != nil {
		t.Error("ParseDomain test failed, .co")
		return
	}
	fmt.Println(dm)

	dm, err = ParseDomain("google.co")
	if err != nil {
		t.Error("ParseDomain google.co test failed")
		return
	}
	fmt.Println(dm)

	dm, err = ParseDomain("cn.groups.google.co")
	if err != nil {
		t.Error("ParseDomain cn.groups.google.co test failed")
		return
	}
	fmt.Println(dm)
}

func TestGetWhoisWithIP(t *testing.T) {
	fmt.Println(GetWhoisWithIP("1.1.1.1"))
}
