package gaddr

import (
	"net"
	"strings"
)

type NicAddr struct {
	IP      net.IP
	Netmask net.IPMask
	CIDR    string
}

type NicInfo struct {
	Name       string
	Desc       string
	MAC        string
	Addrs      []NicAddr
	IsPhysical bool // Is physical network card or virtual network card
	Flags      string
	MTU        int
}

func GetNicInfo(name string) (NicInfo, error) {
	var ni NicInfo
	var na NicAddr
	inf, err := net.InterfaceByName(name)
	if err != nil {
		return ni, err
	}

	// Get MAC address / mtu / Flags
	ni.Name = name
	ni.MAC = strings.TrimSuffix(inf.HardwareAddr.String(), ":00:00")
	ni.MTU = inf.MTU
	ni.Flags = inf.Flags.String()

	// Get IPs
	addrs, err := inf.Addrs()
	if err != nil {
		return ni, err
	}
	for _, addr := range addrs {
		switch v := addr.(type) {
		case *net.IPNet:
			na.IP = v.IP
			na.Netmask = v.Mask
		}
		na.CIDR = addr.String()
		ni.Addrs = append(ni.Addrs, na)
	}

	ni.IsPhysical = isPhysical(ni)
	return ni, nil
}

func GetAllNicNames() ([]string, error) {
	var names []string
	infs, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, inf := range infs {
		names = append(names, inf.Name)
	}
	return names, nil
}

func GetAllLANIPv4CIDRs() ([]string, error) {
	nics, err := GetAllNicNames()
	if err != nil {
		return nil, err
	}

	cidrs := []string{}
	for _, v := range nics {
		ni, err := GetNicInfo(v)
		if err != nil {
			return nil, err
		}
		for _, v := range ni.Addrs {
			if isLANIPv4String(v.IP.String()) {
				cidrs = append(cidrs, v.CIDR)
			}
		}
	}
	return cidrs, nil
}

func isPhysical(ni NicInfo) bool {
	if len(ni.MAC) != 17 {
		return false
	}
	if len(ni.Name) == 0 {
		return false
	}

	nameLower := strings.ToLower(ni.Name)
	descLower := strings.ToLower(ni.Desc)

	if nameLower == "lo0" || len(ni.MAC) != 17 {
		return false
	}

	var virtualFeatures = []string{
		"virtual",
		"vmware",
		"vmnet",
		"oraybox",
		"pseudo",
		"bridge",
		"loopback",
		"vpn",
		"p2p",
		"{",
		"."}

	for _, feature := range virtualFeatures {
		if strings.Contains(nameLower, feature) || strings.Contains(descLower, feature) || strings.Contains(ni.Flags, feature) {
			return false
		}
	}

	return true
}

// get preferred outbound ip of this machine
// it will fail if device is not connected to LAN router
func GetOutboundIP() (net.IP, error) {
	// 1.1.1.1:1 is a fake target address, you can use anyone instead
	conn, err := net.Dial("udp", "1.1.1.1:1")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP, nil
}
