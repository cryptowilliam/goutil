package gkcp

import (
	"github.com/xtaci/kcp-go"
	"net"
)

// kcp server type
type Server struct {
	listener      *kcp.Listener
	udpListenAddr string
	conn          kcp.UDPSession
	opt           Option // 实际上应该命名为kcpOption
}

type Listener struct{}

func (l Listener) Listen(listenAddr string) (net.Listener, error) {
	return Listen(listenAddr)
}

func Listen(listenAddr string) (*Server, error) {
	return ListenWithOptions(listenAddr, DefaultOption(false))
}

func ListenWithOptions(listenAddr string, opt Option) (*Server, error) {
	var s Server

	s.opt = opt

	block, _ := s.opt.getCryptBlock()
	ln, err := kcp.ListenWithOptions(listenAddr, block, s.opt.DataShard, s.opt.ParityShard)
	if err != nil {
		return nil, err
	}
	if err := ln.SetDSCP(s.opt.DSCP); err != nil {
		return nil, err
	}
	if err := ln.SetReadBuffer(s.opt.SockBuf); err != nil {
		return nil, err
	}
	if err := ln.SetWriteBuffer(s.opt.SockBuf); err != nil {
		return nil, err
	}
	s.listener = ln

	return &s, nil
}

func (s *Server) Accept() (net.Conn, error) {
	c, err := s.listener.AcceptKCP() // AcceptKCP
	if err != nil {
		return nil, err
	}

	c.SetStreamMode(true)
	c.SetWriteDelay(false)
	c.SetNoDelay(s.opt.NoDelay, s.opt.Interval, s.opt.Resend, s.opt.NoCongestion)
	c.SetMtu(s.opt.MTU)
	c.SetWindowSize(s.opt.SndWnd, s.opt.RcvWnd)
	c.SetACKNoDelay(s.opt.AckNodelay)

	var uc KcpConn
	uc.server = s
	uc.addr = c.RemoteAddr().String()
	uc.sess = c
	uc.rw = uc.sess
	return &uc, nil
}

func (s *Server) Addr() net.Addr {
	return s.listener.Addr()
}

func (s *Server) Close() error {
	return s.conn.Close()
}
