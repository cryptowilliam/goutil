package gprobe

import (
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/container/gnum"
	"github.com/cryptowilliam/goutil/net/gaddr"
	"net"
	"time"
)

func TCPing(host string, port int, timeout time.Duration) (opened bool, duration time.Duration, err error) {
	if !gaddr.IsValidPort(port) {
		return false, 0, gerrors.Errorf("Invalid port " + gnum.ToString(port))
	}

	ip := ""
	if gaddr.IsIPString(host) {
		ip = host
	} else {
		ipArr, err := gaddr.LookupIP(host)
		if err != nil {
			return false, 0, err
		}
		if len(ipArr) == 0 {
			return false, 0, gerrors.New("ip addresses length is zero")
		}
		ip = ipArr[0].String()
	}

	startTime := time.Now()
	conn, err := net.DialTimeout("tcp", ip+":"+gnum.ToString(port), timeout)
	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()
	if err != nil {
		return false, time.Now().Sub(startTime), nil // Maybe Closed
	}
	return true, time.Now().Sub(startTime), nil // Opened
}
