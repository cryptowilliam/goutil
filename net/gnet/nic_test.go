package gnet

import (
	"fmt"
	"testing"
)

func TestNic(t *testing.T) {
	ns, _ := GetAllNicNames()
	for _, s := range ns {
		ni, _ := GetNicInfo(s)
		if !ni.IsPhysical {
			continue
		}
		fmt.Println("------------------")
		fmt.Println("name:", ni.Name)
		for _, v := range ni.Addrs {
			fmt.Println("CIDR:", v.CIDR, "ip:", v.IP, "mask", v.Netmask)
		}
		fmt.Println("desc:", ni.Desc)
		fmt.Println("physical:", ni.IsPhysical)
		fmt.Println("mac:", ni.MAC)
	}
}

func TestGetAllLANIPv4CIDRs(t *testing.T) {
	fmt.Println(GetAllLANIPv4CIDRs())
}

func TestGetOutboundIP(t *testing.T) {
	outbound, err := GetOutboundIP()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("outbound LAN addr", outbound)
}
