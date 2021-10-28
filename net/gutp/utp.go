package gutp

import (
	"context"
	"github.com/anacrolix/utp"
	"net"
	"time"
)

func Dial(addr string) (net.Conn, error) {
	return utp.Dial(addr)
}

func DialTimeout(addr string, timeout time.Duration) (net.Conn, error) {
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	return utp.DialContext(ctx, addr)
}

func Listen(laddr string) (net.Listener, error) {
	return utp.Listen(laddr)
}
