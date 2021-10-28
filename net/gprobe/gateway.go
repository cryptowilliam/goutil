package gprobe

import (
	"github.com/jackpal/gateway"
	"net"
)

func DiscoverGateway() (net.IP, error) {
	return gateway.DiscoverGateway()
}
