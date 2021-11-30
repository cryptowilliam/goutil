package gnet

import (
	"fmt"
	"github.com/cryptowilliam/goutil/basic/gtest"
	"testing"
)

func TestDNSResolver(t *testing.T) {
	//stamp := "sdns://AgAAAAAAAAAACTIyMy41LjUuNSCoF6cUD2dwqtorNi96I2e3nkHPSJH1ka3xbdOglmOVkQ5kbnMuYWxpZG5zLmNvbQovZG5zLXF1ZXJ5"

	dc := NewDNSClient()
	err := dc.AddCustomDNSServer("127.0.0.1:8888")
	gtest.Assert(t, err)
	resp, err := dc.LookupIP("www.yahoo.com")
	gtest.Assert(t, err)
	fmt.Println(resp)

	addrList, err := dc.LookupIP("localhost")
	gtest.Assert(t, err)
	if !addrList[0].IsLoopback() {
		t.Fatalf("expected loopback")
	}
}
