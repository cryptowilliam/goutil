package gio

import (
	"fmt"
	"github.com/cryptowilliam/goutil/basic/gtest"
	"net"
	"testing"
	"time"
)

func TestReadFull(t *testing.T) {
	lis, err := net.Listen("tcp", "127.0.0.1:12345")
	gtest.Assert(t, err)
	go func() {
		for {
			lis.Accept()
		}
	}()

	conn, err := net.Dial("tcp", "127.0.0.1:12345")
	gtest.Assert(t, err)

	buf := make([]byte, 10)
	timeout := 10 * time.Second
	n, err := ReadFull(conn, buf, &timeout)
	gtest.Assert(t, err)
	fmt.Println(n)
}
