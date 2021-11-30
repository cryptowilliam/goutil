package gnet

import (
	"fmt"
	"github.com/cryptowilliam/goutil/basic/gtest"
	"testing"

	"golang.org/x/net/context"
)

func TestDNSResolver(t *testing.T) {
	stampStr := "sdns://AgAAAAAAAAAACTIyMy41LjUuNSCoF6cUD2dwqtorNi96I2e3nkHPSJH1ka3xbdOglmOVkQ5kbnMuYWxpZG5zLmNvbQovZG5zLXF1ZXJ5"

	dc := NewDNSClient()
	resp, err := dc.LookupIP("www.yahoo.com", stampStr, "")
	gtest.Assert(t, err)
	fmt.Println(resp)
	return

	d := NewDNSClient().SetDefaultServers()
	ctx := context.Background()

	addr, err := d.ResolveIPAddrWithCtx(ctx, "tcp", "localhost")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if !addr.IP.IsLoopback() {
		t.Fatalf("expected loopback")
	}
}
