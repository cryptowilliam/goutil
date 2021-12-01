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
	case "udp":
		return net.Dial("udp", address)
	case "kcp":
		return gkcp.Dial(address)
	case "quic":
		return gquic.Dial(address)
	default:
		return nil, gerrors.New("unsupported network %s", network)
	}
}

// ListenCop listens multiple connection-oriented protocols.
func ListenCop(network, address string) (net.Listener, error) {
	switch strings.ToLower(network) {
	case "tcp", "tcp4", "tcp6":
		return net.Listen(network, address)
	case "kcp":
		return gkcp.Listen(address)
	case "quic":
		return gquic.Listen(address)
	default:
		return nil, gerrors.New("unsupported network %s", network)
	}
}

// ListenAny listens any supported protocols.
func ListenAny(network, address string) (net.Listener, error) {
	switch strings.ToLower(network) {
	case "udp", "udp4", "udp6", "unixgram", "ip:1", "ip:icmp":
		return ListenPop(network, address)
	default:
		return ListenCop(network, address)
	}
}
