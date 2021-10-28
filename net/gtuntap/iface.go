package gtuntap

import (
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"strconv"

	"github.com/songgao/water"
)

type Iface struct {
	name string
	ip   string
	mtu  int
	ifce *water.Interface
}

func NewIface(name, ip string, mtu int) *Iface {
	return &Iface{
		name: name,
		ip:   ip,
		mtu:  mtu,
	}
}

func (i *Iface) Start() error {
	ip, netIP, err := net.ParseCIDR(i.ip)
	if err != nil {
		return err
	}

	if runtime.GOOS == "darwin" {
		config := water.Config{
			DeviceType: water.TUN,
		}
		//config.Name = i.name
		i.ifce, err = water.New(config)
		if err != nil {
			return err
		}
		mask := netIP.Mask
		netmask := fmt.Sprintf("%d.%d.%d.%d", mask[0], mask[1], mask[2], mask[3])
		cmd := exec.Command("ifconfig", i.Name(),
			ip.String(), "netmask", netmask,
			"mtu", strconv.Itoa(i.mtu), "up")
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("err: %s %s", err, string(output))
		}
	}

	if runtime.GOOS == "linux" || runtime.GOOS == "windows" {
		config := water.Config{
			DeviceType: water.TAP,
		}

		i.ifce, err = water.New(config)
		if err != nil {
			return err
		}
	}

	return nil
}

func (i *Iface) Name() string {
	return i.ifce.Name()
}

func (i *Iface) Read(pkt PacketIP) (int, error) {
	return i.ifce.Read(pkt)
}

func (i *Iface) Write(pkt PacketIP) (int, error) {
	return i.ifce.Write(pkt)
}

type PacketIP []byte

func NewPacketIP(size int) PacketIP {
	return PacketIP(make([]byte, size))
}

func (p PacketIP) GetSourceIP() net.IP {
	return net.IP(p[12:16])
}

func (p PacketIP) GetDestinationIP() net.IP {
	return net.IP(p[16:20])
}
