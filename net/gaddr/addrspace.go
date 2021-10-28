package gaddr

import (
	"net"
	"strings"
)

// const for CheckIPString
const (
	NotIP         IpType = 0
	IPv4_LAN      IpType = 1
	IPv4_WAN      IpType = 2
	IPv6_LAN      IpType = 3
	IPv6_WAN      IpType = 4
	IPv4_LOOPBACK IpType = 5
	IPv6_LOOPBACK IpType = 6
	IPv4_ANY      IpType = 7
	IPv6_ANY      IpType = 8
)

type IpType int

var addrSpaceOfLoopbackIPv4 = [...]string{
	"127.0.0.0/8",
}

var addrSpaceOfLoopbackIPv6 = [...]string{
	"::1/128",
}

var addrSpaceOfAnyIPv6 = [...]string{
	"::/128",
}

var addrSpacesOfLANIPv4 = [...]string{
	"10.0.0.0/8",     // 10.*.*.*
	"172.16.0.0/12",  // 172.16.0.0 - 172.31.255.255
	"192.168.0.0/16", // 192.168.*.*
}

var addrSpacesOfLANIPv6 = [...]string{
	"fd00::/8",
}

func isLoopbackIPv4String(s string) bool {
	for _, it := range addrSpaceOfLoopbackIPv4 {
		_, cidrnet, err := net.ParseCIDR(it)
		if err != nil {
			panic(err) // assuming I did it right above
		}
		myaddr := net.ParseIP(strings.Split(s, "/")[0])

		if cidrnet.Contains(myaddr) {
			return true
		}
	}
	return false
}

func isLoopbackIPv6String(s string) bool {
	for _, it := range addrSpaceOfLoopbackIPv6 {
		_, cidrnet, err := net.ParseCIDR(it)
		if err != nil {
			panic(err) // assuming I did it right above
		}
		myaddr := net.ParseIP(strings.Split(s, "/")[0])

		if cidrnet.Contains(myaddr) {
			return true
		}
	}
	return false
}

func isAnyIPv4String(s string) bool {
	return s == "0.0.0.0"
}

func isAnyIPv6String(s string) bool {
	for _, it := range addrSpaceOfLoopbackIPv6 {
		_, cidrnet, err := net.ParseCIDR(it)
		if err != nil {
			panic(err) // assuming I did it right above
		}
		myaddr := net.ParseIP(strings.Split(s, "/")[0])

		if cidrnet.Contains(myaddr) {
			return true
		}
	}
	return false
}

func isLANIPv4String(s string) bool {
	for _, it := range addrSpacesOfLANIPv4 {
		_, cidrnet, err := net.ParseCIDR(it)
		if err != nil {
			panic(err) // assuming I did it right above
		}
		myaddr := net.ParseIP(strings.Split(s, "/")[0])

		if cidrnet.Contains(myaddr) {
			return true
		}
	}
	return false
}

func isLANIPv6String(s string) bool {
	for _, it := range addrSpacesOfLANIPv6 {
		_, cidrnet, err := net.ParseCIDR(it)
		if err != nil {
			panic(err) // assuming I did it right above
		}
		myaddr := net.ParseIP(strings.Split(s, "/")[0])

		if cidrnet.Contains(myaddr) {
			return true
		}
	}
	return false
}

// check is IPv4_LAN or IPv4_WAN or IPv6 or NotIP
func CheckIPString(s string) IpType {
	ip := net.ParseIP(s)
	return CheckIP(ip)
}

func CheckIP(ip net.IP) IpType {
	if ip == nil {
		return NotIP
	}
	s := ip.String()
	if ip.To4() != nil {
		if isLoopbackIPv4String(s) {
			return IPv4_LOOPBACK
		} else if isAnyIPv4String(s) {
			return IPv4_ANY
		} else if isLANIPv4String(s) {
			return IPv4_LAN
		} else {
			return IPv4_WAN
		}
	}
	if isLoopbackIPv6String(s) {
		return IPv6_LOOPBACK
	} else if isAnyIPv6String(s) {
		return IPv6_ANY
	} else if isLANIPv6String(s) {
		return IPv6_LAN
	} else {
		return IPv6_WAN
	}
}

func IsIPString(s string) bool {
	return CheckIPString(s) != NotIP
}
