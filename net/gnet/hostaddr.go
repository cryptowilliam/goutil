package gnet

// Email such like "HOSTNAME" / "DOMAIN" / "IP" / ":PORT" / "HOSTNAME:PORT" / "DOMAIN:PORT" / "IP:PORT"

import (
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"net"
	"strconv"
	"strings"
)

// "HOSTNAME" / "DOMAIN" / "IP" / ":PORT" / "HOSTNAME:PORT" / "DOMAIN:PORT" / "IP:PORT" -> net.IP, port
// return port maybe -1, this is NOT an error
// NOTICE:
// ResolveIPAddr() & LookupIP() API can't recognize "1127.0.0.1" or "abc127.0.0.1" style illegal IP string,
// they still returns a IP address and nil error
func ParseHostAddrOnline(addr string) (IP net.IP, port int, err error) {
	if len(addr) == 0 {
		return nil, -1, gerrors.New("Empty address string")
	}
	addr = strings.ToLower(addr)

	var addrWithoutPort = addr
	var resultPort = -1

	// Try parse PORT
	// Infact, ResolveTCPAddr() will get a correct port from any string looks like "ANYTHING:PORT-NUMBER"
	a, e := net.ResolveTCPAddr("tcp", addr) // Must be "DOMAIN:PORT" / "IP:PORT" style (Must have PORT)
	if e == nil && a.Port > 0 {
		resultPort = a.Port

		// Remove port from address
		end := strings.LastIndex(addr, ":")
		if end < 0 {
			return nil, -1, gerrors.New(addr + " is not a legal host address")
		}
		addrWithoutPort = addr[0:end]
	}

	// Try parse address which removed port
	if len(addrWithoutPort) == 0 {
		return net.ParseIP("0.0.0.0"), resultPort, nil
	} else if addrWithoutPort == "localhost" {
		return net.ParseIP("127.0.0.1"), resultPort, nil
	} else if IsIPString(addrWithoutPort) { // "IP" style
		return net.ParseIP(addrWithoutPort), resultPort, nil
	} else if IsDomainONLINE(addrWithoutPort) == true { // "DOMAIN" style
		ips, e := net.LookupIP(addrWithoutPort)
		if e != nil {
			return nil, -1, gerrors.New("Can't lookup IP from address string " + addr)
		}
		return ips[0], resultPort, nil
	} else {
		return nil, -1, gerrors.New(addr + " is not a legal socket address")
	}
}

func ParseAddr(addr net.Addr) (IP net.IP, port int, err error) {
	// parse IP
	switch raw := addr.(type) {
	case *net.UDPAddr:
		IP = raw.IP
	case *net.TCPAddr:
		IP = raw.IP
	default:
		host, _, err := net.SplitHostPort(addr.String())
		if err != nil {
			return nil, 0, err
		}
		IP = net.ParseIP(host)
	}

	// parse port
	switch raw := addr.(type) {
	case *net.UDPAddr:
		port = raw.Port
	default:
		_, portStr, err := net.SplitHostPort(addr.String())
		if err != nil {
			return nil, 0, err
		}
		i64, err := strconv.ParseInt(portStr, 0, 0)
		if err != nil {
			return nil, 0, err
		}
		port = int(i64)
	}

	return IP, port, nil
}
