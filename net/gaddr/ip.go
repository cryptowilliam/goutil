package gaddr

import (
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"net"
	"strings"
)

func ParseIP(s string) (net.IP, error) {
	ip := net.ParseIP(s)
	if ip == nil {
		return nil, gerrors.New("Invalid IP address string '" + s + "'")
	}
	return ip, nil
}

// get all my local IPs
func GetLanIps() ([]net.IP, error) {
	var result = make([]net.IP, 0)

	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}

		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}

		if strings.HasPrefix(iface.Name, "docker") || strings.HasPrefix(iface.Name, "w-") {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil || ip.IsLoopback() {
				continue
			}

			t := CheckIPString(ip.String())
			if t == IPv4_LAN || t == IPv6_LAN {
				result = append(result, ip)
			}
		}
	}

	return result, nil
}
