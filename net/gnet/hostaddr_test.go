package gnet

import (
	"net"
	"testing"
)

func TestParseHostAddrOnline(t *testing.T) {
	var ip net.IP
	var port int
	var err error

	ip, port, err = ParseHostAddrOnline("localhost")
	if err != nil || ip.String() != "127.0.0.1" || port != -1 {
		t.Error("ParseAddrString failed") // 是不是LookupIP造成的？
	}

	ip, port, err = ParseHostAddrOnline("baidu.com")
	if err != nil || port != -1 {
		t.Error("ParseAddrString failed")
	}

	ip, port, err = ParseHostAddrOnline("127.0.0.1")
	if err != nil || ip.String() != "127.0.0.1" || port != -1 {
		t.Error("ParseAddrString failed")
	}

	ip, port, err = ParseHostAddrOnline(":8888")
	if err != nil || ip.String() != "0.0.0.0" || port != 8888 {
		t.Error("ParseAddrString failed")
	}

	ip, port, err = ParseHostAddrOnline("localhost:8888")
	if err != nil || ip.String() != "127.0.0.1" || port != 8888 {
		t.Error("ParseAddrString failed")
	}

	ip, port, err = ParseHostAddrOnline("localhost:88888")
	if err == nil {
		t.Error("ParseAddrString failed")
	}

	ip, port, err = ParseHostAddrOnline("127.0.0.1:8888")
	if err != nil || ip.String() != "127.0.0.1" || port != 8888 {
		t.Error("ParseAddrString failed")
	}

	ip, port, err = ParseHostAddrOnline("1127.0.0")
	if err == nil {
		t.Error("ParseAddrString failed")
	}

	ip, port, err = ParseHostAddrOnline("1127.0.0.1")
	if err == nil {
		t.Error("ParseAddrString failed")
	}

	ip, port, err = ParseHostAddrOnline("1127.0.0.1:8888")
	if err == nil {
		t.Error("ParseAddrString failed")
	}

	ip, port, err = ParseHostAddrOnline("8888")
	if err == nil {
		t.Error("ParseAddrString failed")
	}

	ip, port, err = ParseHostAddrOnline("abc")
	if err == nil {
		t.Error("ParseAddrString failed")
	}

	ip, port, err = ParseHostAddrOnline("google.hotel")
	if err == nil {
		t.Error("ParseAddrString failed")
	}
}
