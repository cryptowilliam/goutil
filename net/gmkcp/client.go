package gmkcp

// TODO:
// Cant's detect connection established or closed.

import (
	"context"
	"github.com/cryptowilliam/goutil/net/gnet"
	"net"
	"time"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/transport/internet"
	. "v2ray.com/core/transport/internet/kcp"
)

// mkcp client
type Client struct {
	serverAddr string
	mkcpConn   internet.Connection
}

// Create new client connection and dial
func Dial(raddr string) (*Client, error) {
	var mc Client

	serverNetIP, serverPort, err := gnet.ParseHostAddrOnline(raddr)
	if err != nil {
		return nil, err
	}
	serverIP := serverNetIP.String()

	config := Config{
		Mtu:              &MTU{Value: 1500},
		Tti:              &TTI{Value: 10},
		Congestion:       true,
		UplinkCapacity:   &UplinkCapacity{Value: 10},
		DownlinkCapacity: &DownlinkCapacity{Value: 10},
		ReadBuffer:       &ReadBuffer{Size: 10 * 1024 * 1024},
		WriteBuffer:      &WriteBuffer{Size: 10 * 1024 * 1024},
	}
	clientConn, err := DialKCP(context.Background(), v2net.UDPDestination(v2net.ParseAddress(serverIP), v2net.Port(serverPort)),
		&internet.MemoryStreamConfig{
			ProtocolName:     "mkcp",
			ProtocolSettings: &config,
		})

	if err != nil {
		return nil, err
	}
	mc.mkcpConn = clientConn
	return &mc, nil
}

func (c *Client) State() int {
	return int(c.mkcpConn.(*Connection).State())
}

func (c *Client) Read(b []byte) (n int, err error) {
	return c.mkcpConn.Read(b)
}

func (c *Client) Write(b []byte) (n int, err error) {
	return c.mkcpConn.Write(b)
}

func (c *Client) Close() error {
	return c.mkcpConn.Close()
}

func (c *Client) LocalAddr() net.Addr {
	return c.mkcpConn.LocalAddr()
}

func (c *Client) RemoteAddr() net.Addr {
	return c.mkcpConn.RemoteAddr()
}

func (c *Client) SetDeadline(t time.Time) error {
	return c.mkcpConn.SetDeadline(t)
}

func (c *Client) SetReadDeadline(t time.Time) error {
	return c.SetReadDeadline(t)
}

func (c *Client) SetWriteDeadline(t time.Time) error {
	return c.SetWriteDeadline(t)
}
