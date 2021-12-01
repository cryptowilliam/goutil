// another implement github.com/liupeidong0620/gateway
// which could find gateway for each network interface card

package gnet

import (
	"github.com/jackpal/gateway"
	"net"
)

func DiscoverGateway() (net.IP, error) {
	return gateway.DiscoverGateway()
}
