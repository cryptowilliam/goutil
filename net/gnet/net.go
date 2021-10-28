package gnet

import (
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/net/gkcp"
	"github.com/cryptowilliam/goutil/net/gquic"
	"net"
	"strings"
)

// Dial multiple protocols.
func Dial(network, address string) (net.Conn, error) {
	switch strings.ToLower(network) {
	case "tcp":
		return net.Dial("tcp", address)
	case "kcp":
		return gkcp.Dial(address)
	case "quic":
		return gquic.Dial(address)
	default:
		return nil, gerrors.New("unsupported network %s", network)
	}
}

// Listen multiple protocols.
func Listen(network, address string) (net.Listener, error) {
	switch strings.ToLower(network) {
	case "tcp":
		return net.Listen("tcp", address)
	case "kcp":
		return gkcp.Listen(address)
	case "quic":
		return gquic.Listen(address)
	default:
		return nil, gerrors.New("unsupported network %s", network)
	}
}
