package gnet

import (
	"fmt"
	"net"
)

type (
	Interface net.Interface
)

func WrapIfi(ifi net.Interface) Interface {
	return Interface(ifi)
}

func WrapIfiPtr(ifi *net.Interface) *Interface {
	return (*Interface)(ifi)
}

func Interfaces() ([]Interface, error) {
	ifis, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var result []Interface
	for _, v := range ifis {
		result = append(result, Interface(v))
	}
	return result, nil
}

func (ifi *Interface) Raw() *net.Interface {
	return (*net.Interface)(ifi)
}

// GetV4 returns v4 IP address and IP network of current network interface.
func (ifi *Interface) GetV4() (IP, *IPNet, error) {
	if addrList, err := ifi.Raw().Addrs(); err != nil {
		return nil, nil, err
	} else {
		for _, addr := range addrList {
			// Notice: sb use 'addr.(*net.IPNet)' convert addr to IPNet, this is DANGEROUS!
			// Yes 'addr' is a 'IPNet' structure inside golang source, but it is NOT
			// a valid IP network, its IP member is a specific IP address (like 192.168.9.123/24)
			// but not IP network required IP address (like 192.168.9.0/24). The correct IP address
			// in IP network is the head of IP ranges of the network, like 192.168.9.0/24.
			ip, ipNet, err := net.ParseCIDR(addr.String())
			if err != nil {
				return nil, nil, err
			}
			if WrapIP(ip).IsV4() {
				return IP(ip), WrapIPNetPtr(ipNet), nil
			}
		}
	}
	// Sanity-check that the interface has a good address.
	return nil, nil, fmt.Errorf("no IP4 network found")
}

func (ifi *Interface) IsUp() bool {
	return ifi.Raw().Flags&net.FlagUp > 0
}

func (ifi *Interface) IsLoopBack() bool {
	return ifi.Raw().Flags&net.FlagLoopback != 0
}

func (ifi *Interface) IsBroadcast() bool {
	return ifi.Raw().Flags&net.FlagBroadcast != 0
}

func (ifi *Interface) IsMulticast() bool {
	return ifi.Raw().Flags&net.FlagMulticast != 0
}

func (ifi *Interface) IsPointToPoint() bool {
	return ifi.Raw().Flags&net.FlagPointToPoint != 0
}
