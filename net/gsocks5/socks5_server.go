package gsocks5

import (
	"github.com/cryptowilliam/goutil/net/gnet"
	socks5internal "github.com/cryptowilliam/goutil/net/gsocks5/socks5internal"
	"net"
)

type (
	Server struct {
		lis net.Listener
		srv *socks5internal.Server
		listenAddr string
		dialer     gnet.DialWithCtxFunc
	}
)

// NewServer create new socks5 server.
// For now, non-tcp socks5 proxy server is not necessary, so there is no "network" param.
// listenAddr example: "127.0.0.1:8000"
func NewServer(listenAddr string) *Server {
	return &Server{listenAddr: listenAddr}
}

// SetCustomDialer set custom dialer for requests.
// This operation is optional.
func (s *Server) SetCustomDialer(dialer gnet.DialWithCtxFunc) {
	s.dialer = dialer
}

// ListenAndServe start a tcp socks5 proxy server.
// For now, non-tcp socks5 proxy server is not necessary, so there is no "network" param.
// listenAddr example: "127.0.0.1:8000"
func (s *Server) ListenAndServe() error {
	// Create a SOCKS5 server
	conf := &socks5internal.Config{}
	if s.dialer != nil {
		conf.Dial = s.dialer
	}
	server, err := socks5internal.New(conf)
	if err != nil {
		return err
	}

	// Create SOCKS5 proxy on localhost port 8000
	if err := server.ListenAndServe("tcp", s.listenAddr); err != nil {
		return err
	}
	return nil
}

// Listen start a tcp listener.
func (s *Server) Listen() error {
	// Create a SOCKS5 server
	conf := &socks5internal.Config{}
	if s.dialer != nil {
		conf.Dial = s.dialer
	}
	err := error(nil)
	s.srv, err = socks5internal.New(conf)
	if err != nil {
		return err
	}
	s.lis, err = net.Listen("tcp", s.listenAddr)
	return err
}

// Serve socks5 proxy server and wait until it returns error.
func (s *Server) Serve() error {
	return s.srv.Serve(s.lis)
}