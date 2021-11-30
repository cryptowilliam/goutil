package gnet

import (
	"encoding/binary"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"net"
	"strings"
)

type (
	// IP address.
	IP net.IP

	// IPNet defines IP network, or IP range.
	// Notice:
	// Valid IPNet samples: 192.168.7.0/24
	// Invalid IPNet samples: 192.168.7.123/24
	IPNet net.IPNet
)

var (
	addrSpaceOfLoopBackIPv4 = [...]string{
		"127.0.0.0/8",
	}

	addrSpaceOfLoopBackIPv6 = [...]string{
		"::1/128",
	}

	addrSpaceOfAnyIPv6 = [...]string{
		"::/128",
	}

	addrSpacesOfLANIPv4 = [...]string{
		"10.0.0.0/8",     // 10.*.*.*
		"172.16.0.0/12",  // 172.16.0.0 - 172.31.255.255
		"192.168.0.0/16", // 192.168.*.*
	}

	addrSpacesOfLANIPv6 = [...]string{
		"fd00::/8",
	}
)

func WrapIP(ip net.IP) IP {
	return IP(ip)
}

func ParseIP(s string) (IP, error) {
	ip := net.ParseIP(s)
	if ip == nil {
		return nil, gerrors.New("Invalid IP address string '" + s + "'")
	}
	return IP(ip), nil
}

func WrapIPNet(ipNet net.IPNet) IPNet {
	return IPNet(ipNet)
}

func WrapIPNetPtr(ipNet *net.IPNet) *IPNet {
	return (*IPNet)(ipNet)
}

func (ip IP) Raw() net.IP {
	return net.IP(ip)
}

func (ip IP) String() string {
	return net.IP(ip).String()
}

func (ip IP) IsLoopBack() bool {
	for _, it := range addrSpaceOfLoopBackIPv4 {
		_, cidrNet, err := net.ParseCIDR(it)
		if err != nil {
			panic(err) // assuming I did it right above
		}
		myAddr := net.ParseIP(strings.Split(ip.String(), "/")[0])

		if cidrNet.Contains(myAddr) {
			return true
		}
	}

	for _, it := range addrSpaceOfLoopBackIPv6 {
		_, cidrNet, err := net.ParseCIDR(it)
		if err != nil {
			panic(err) // assuming I did it right above
		}
		myAddr := net.ParseIP(strings.Split(ip.String(), "/")[0])

		if cidrNet.Contains(myAddr) {
			return true
		}
	}

	return false
}

func (ip IP) IsV4() bool {
	return ip.Raw().To4() != nil
}

func (ip IP) IsV6() bool {
	return !ip.IsV4()
}

func (ip IP) IsPublic() bool {
	return !ip.IsPrivate()
}

func (ip IP) IsPrivate() bool {
	return ip.Raw().IsPrivate()
	/*for _, it := range addrSpacesOfLANIPv4 {
		_, cidrNet, err := net.ParseCIDR(it)
		if err != nil {
			panic(err) // assuming I did it right above
		}
		myAddr := net.ParseIP(strings.Split(ip.String(), "/")[0])

		if cidrNet.Contains(myAddr) {
			return true
		}
	}

	for _, it := range addrSpacesOfLANIPv6 {
		_, cidrNet, err := net.ParseCIDR(it)
		if err != nil {
			panic(err) // assuming I did it right above
		}
		myAddr := net.ParseIP(strings.Split(ip.String(), "/")[0])

		if cidrNet.Contains(myAddr) {
			return true
		}
	}

	return false*/
}

func (ip IP) IsAny() bool {
	if ip.String() == "0.0.0.0" {
		return true
	}

	for _, it := range addrSpaceOfLoopBackIPv6 {
		_, cidrNet, err := net.ParseCIDR(it)
		if err != nil {
			panic(err) // assuming I did it right above
		}
		myAddr := net.ParseIP(strings.Split(ip.String(), "/")[0])

		if cidrNet.Contains(myAddr) {
			return true
		}
	}

	return false
}

func IsIPString(s string) bool {
	_, err := ParseIP(s)
	return err == nil
}

func LookupIP(host string) ([]net.IP, error) {
	return net.LookupIP(host)
}

// get all my local IPs
func GetPrivateIPs() ([]net.IP, error) {
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

			if WrapIP(ip).IsPrivate() {
				result = append(result, ip)
			}
		}
	}

	return result, nil
}

func (in *IPNet) Raw() *net.IPNet {
	return (*net.IPNet)(in)
}

func (in *IPNet) String() string {
	return in.Raw().String()
}

// Verify checks if IPNet is a valid IP network or not.
// Valid IPNet samples: 192.168.7.0/24
// Invalid IPNet samples: 192.168.7.12/24
func (in *IPNet) Verify() error {
	_, parsedIpNet, err := net.ParseCIDR(in.String())
	if err != nil {
		return err
	}
	if parsedIpNet.String() != in.String() {
		return gerrors.New("IPNet(%s) is not a valid IP network", in.String())
	}
	return nil
}

// ListAll returns all IPs of current IP network.
// FIXME: IPv6 not supported for now.
func (in *IPNet) ListAll() []net.IP {
	var result []net.IP
	num := binary.BigEndian.Uint32(in.IP)
	mask := binary.BigEndian.Uint32(in.Mask)
	num &= mask
	for mask < 0xffffffff {
		var buf [4]byte
		binary.BigEndian.PutUint32(buf[:], num)
		result = append(result, buf[:])
		mask += 1
		num += 1
	}
	return result
}
