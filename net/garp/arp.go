package garp

import (
	"fmt"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/basic/glog"
	"github.com/cryptowilliam/goutil/net/gnet"
	"github.com/cryptowilliam/goutil/safe/gwg"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"net"
	"sync"
	"time"
)

type (
	Device struct {
		Ifi  string
		MAC  string
		IP   net.IP
		Name string
	}

	ScanErr struct {
		Ifi string
		Err error
	}

	Scanner struct {
		chErr    chan ScanErr
		chResult chan Device
		results  sync.Map
	}
)

// NewScanner create arp scanner.
func NewScanner() *Scanner {
	return &Scanner{
		chErr:    make(chan ScanErr, 4096),
		chResult: make(chan Device, 4096),
	}
}

// Scan requires ROOT privilege.
func (s *Scanner) Scan() ([]Device, error) {
	// clean up result cache
	s.results = sync.Map{}

	// Get a list of all interfaces.
	ifs, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	scanWg := gwg.New()
	for _, ifc := range ifs {
		scanWg.Add(1)
		// Start up a scan on each interface.
		go func(iface net.Interface) {
			defer scanWg.Add(-1)
			if err := s.scanIfi(iface); err != nil {
				s.chErr <- ScanErr{
					Ifi: iface.Name,
					Err: err,
				}
				glog.Erro(err)
			}

		}(ifc)
	}
	// Wait for all interfaces' scans to complete.
	scanWg.Wait()

	var result []Device
	s.results.Range(func(key, value interface{}) bool {
		result = append(result, value.(Device))
		return true
	})
	return result, nil
}

// scanIfi scans an individual interface's local network for machines using ARP requests/replies.
//
// scan loops forever, sending packets out regularly.  It returns an error if
// it's ever unable to write a packet.
func (s *Scanner) scanIfi(ifi net.Interface) error {
	// We just look for IPv4 addresses, so try to find if the interface has one.
	ipAddr, ipNet, err := gnet.WrapIfiPtr(&ifi).GetV4()
	if err != nil {
		return err
	}

	if gnet.WrapIfiPtr(&ifi).IsLoopBack() {
		return nil
	}

	// Open up a pcap handle for packet reads/writes.
	handle, err := pcap.OpenLive(ifi.Name, 65536, true, pcap.BlockForever)
	if err != nil {
		return err
	}
	defer handle.Close()

	// Start up a goroutine to read in packet data.
	chExit := make(chan struct{})
	go s.receiveARP(handle, ifi, chExit)
	defer close(chExit)

	// Looping multiple times because of the responses & devices being unreliable.
	for count := 1; count <= 6; count++ {
		// Write our scan packets out to the handle.
		if err := s.sendARP(handle, ifi, ipAddr.Raw(), *ipNet.Raw()); err != nil {
			return gerrors.Wrap(err, fmt.Sprintf("error writing packets on %s", ifi.Name))
		}

		// Sleep here to wait for arp responses - 3 secs seems to be the most reliable
		time.Sleep(3 * time.Second)
	}

	return nil
}

// receiveARP watches a handle for incoming ARP responses we might care about, and prints them.
// receiveARP loops until 'chExit' is closed.
func (s *Scanner) receiveARP(handle *pcap.Handle, ifi net.Interface, chExit chan struct{}) {
	lookupAddrFirst := func(ip string) string {
		addr, err := net.LookupAddr(ip)
		if err != nil {
			glog.Erro(err)
			return ""
		}
		if len(addr) > 0 {
			return addr[0]
		}
		return ""
	}

	src := gopacket.NewPacketSource(handle, layers.LayerTypeEthernet)
	in := src.Packets()
	for {
		var packet gopacket.Packet
		select {
		case <-chExit:
			return
		case packet = <-in:
			if packet == nil {
				continue
			}
			arpLayer := packet.Layer(layers.LayerTypeARP)
			if arpLayer == nil {
				continue
			}
			arp := arpLayer.(*layers.ARP)
			if arp.Operation != layers.ARPReply {
				// This is a packet I sent.
				continue
			}
			newItem := Device{
				Ifi:  ifi.Name,
				Name: lookupAddrFirst(net.IP(arp.SourceProtAddress).String()),
				IP:   arp.SourceProtAddress,
				MAC:  net.HardwareAddr(arp.SourceHwAddress).String(),
			}
			s.chResult <- newItem
			s.results.Store(newItem.MAC, newItem)
		}
	}
}

// sendARP writes an ARP request for each address on our local network to the pcap handle.
func (s *Scanner) sendARP(handle *pcap.Handle, ifi net.Interface, ipAddr net.IP, ipNet net.IPNet) error {
	// Set up all the layers' fields we can.
	eth := layers.Ethernet{
		SrcMAC:       ifi.HardwareAddr,
		DstMAC:       net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		EthernetType: layers.EthernetTypeARP,
	}
	arp := layers.ARP{
		AddrType:        layers.LinkTypeEthernet,
		Protocol:        layers.EthernetTypeIPv4,
		HwAddressSize:   6,
		ProtAddressSize: 4,
		Operation:       layers.ARPRequest,
		SourceHwAddress: []byte(ifi.HardwareAddr),
		// For example:
		// Sender's IP address is "192.168.3.123/24", its IP
		// network is "192.168.3.0/24",
		// Notice: use sending device's IP address like 192.168.3.123 here,
		// don't use IP network's IP head like 192.168.3.0, otherwise sender
		// won't receive ARP response.
		// Why invoke ".To4()"?
		// IPv4 []byte length could be 16 too, not only IPv6 has 16 bytes.
		SourceProtAddress: []byte(ipAddr.To4()),
		DstHwAddress:      []byte{0, 0, 0, 0, 0, 0},
	}
	// Set up buffer and options for serialization.
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}
	// Send one packet for every address.
	for _, ip := range gnet.WrapIPNetPtr(&ipNet).ListAll() {
		arp.DstProtAddress = ip
		if err := gopacket.SerializeLayers(buf, opts, &eth, &arp); err != nil {
			return err
		}
		if err := handle.WritePacketData(buf.Bytes()); err != nil {
			return err
		}
	}
	return nil
}
