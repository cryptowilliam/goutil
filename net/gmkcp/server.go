package gmkcp

// TODO
// Close mkcp's internal logging output.

import (
	"context"
	"github.com/cryptowilliam/goutil/net/gaddr"
	"net"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/transport/internet"
	mkcp "v2ray.com/core/transport/internet/kcp"
)

// mkcp server
type Server struct {
	ln      *mkcp.Listener
	accepts chan net.Conn
}

func Listen(listenAddr string) (*Server, error) {
	s := Server{}

	s.accepts = make(chan net.Conn, 1024)

	lnNetIP, lnPort, err := gaddr.ParseHostAddrOnline(listenAddr)
	if err != nil {
		return nil, err
	}
	lnIP := lnNetIP.String()

	config := mkcp.Config{
		Mtu:              &mkcp.MTU{Value: 1500},
		Tti:              &mkcp.TTI{Value: 10},
		Congestion:       true,
		UplinkCapacity:   &mkcp.UplinkCapacity{Value: 5},
		DownlinkCapacity: &mkcp.DownlinkCapacity{Value: 200},
		ReadBuffer:       &mkcp.ReadBuffer{Size: 10 * 1024 * 1024},
		WriteBuffer:      &mkcp.WriteBuffer{Size: 10 * 1024 * 1024},
	}
	s.ln, err = mkcp.NewListener(context.Background(), v2net.ParseAddress(lnIP), v2net.Port(lnPort),
		&internet.MemoryStreamConfig{
			ProtocolName:     "mkcp",
			ProtocolSettings: &config,
		}, func(conn internet.Connection) {
			s.accepts <- conn
		})
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (s *Server) Accept() (net.Conn, error) {
	return <-s.accepts, nil
}

func (s *Server) Close() error {
	return s.ln.Close()
}

func (s *Server) Addr() net.Addr {
	return s.ln.Addr()
}

func (s *Server) GetActiveConnCount() int {
	return s.ln.ActiveConnections()
}
