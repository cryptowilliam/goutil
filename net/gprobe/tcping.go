package gprobe

import (
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/container/gnum"
	"github.com/cryptowilliam/goutil/net/gaddr"
	"net"
	"time"
)

func TcpingOnline(host string, port int, timeout time.Duration) (opened bool, err error) {
	if !gaddr.IsValidPort(port) {
		return false, gerrors.Errorf("Invalid port " + gnum.ToString(port))
	}

	ip, _, err := gaddr.ParseHostAddrOnline(host)
	if err != nil {
		return false, err
	}

	conn, err := net.DialTimeout("tcp", ip.String()+":"+gnum.ToString(port), timeout)
	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()
	if err != nil {
		return false, nil // Maybe Closed
	}
	return true, nil // Opened
}
