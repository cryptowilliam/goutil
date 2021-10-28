package gaddr

import (
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"net"
)

// Extracts IP mask from CIDR address.
func CIDRToMask(cidr string) (net.IPMask, error) {
	_, ip, err := net.ParseCIDR(cidr)
	return ip.Mask, err
}

func CidrContainsIP(cidr, ip string) (bool, error) {
	if !IsIPString(ip) {
		return false, gerrors.Errorf("%s is not ip address", ip)
	}

	if _, in, err := net.ParseCIDR(cidr); err != nil {
		return false, err
	} else {
		ipa := net.ParseIP(ip)
		return in.Contains(ipa), nil
	}
}

//  http://play.golang.org/p/m8TNTtygK0
func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

// List all ip address contained in this cidr
func CIDRListAll(cidr string) ([]string, error) {

	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}

	// remove network address and broadcast address
	return ips[1 : len(ips)-1], nil
}
