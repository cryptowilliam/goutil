package gprobe

import (
	"fmt"
	"github.com/Cubox-/libping"
	"net"
	"time"
)

// https://github.com/ipsecdiagtool/ipsecdiagtool ?

func DiscoverMtu() (int, error) {
	chann := make(chan libping.Response, 100)
	go libping.Pinguntil("192.168.100.100", 10, chann, time.Second)
	for i := range chann {
		if ne, ok := i.Error.(net.Error); ok && ne.Timeout() {
			fmt.Printf("Request timeout for icmp_seq %d\n", i.Seq)
			continue
		} else if i.Error != nil {
			fmt.Println(i.Error)
		} else {
			fmt.Printf("%d bytes from %s: icmp_seq=%d time=%s\n", i.Readsize, i.Destination, i.Seq, i.Delay)
		}
	}
	return 0, nil
}
