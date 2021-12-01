package gsocks5

import (
	"github.com/cryptowilliam/goutil/net/gnet"
	"github.com/cryptowilliam/goutil/net/gsocks5/socks5internal"
	"net"
)

type (
	Server struct {
		lis        net.Listener
		srv        *socks5internal.Server
		listenAddr string
		dialer     gnet.DialWithCtxFunc
		dnsResolver gnet.LookupIPWithCtxFunc
	}
)

// NewServer create new socks5 server.
// For now, non-tcp socks5 proxy server is not necessary, so there is no "network" param.
// listenAddr example: "127.0.0.1:8000"
func NewServer(listenAddr string) *Server {
	return &Server{listenAddr: listenAddr}
}

// SetCustomDialer sets custom dialer for requests.
// This operation is optional.
func (s *Server) SetCustomDialer(dialer gnet.DialWithCtxFunc) {
	s.dialer = dialer
}

// SetCustomDNSResolver sets custom DNS resolver for requests.
// This operation is optional.
func (s *Server) SetCustomDNSResolver(dnsResolver gnet.LookupIPWithCtxFunc) {
	s.dnsResolver = dnsResolver
}

// ListenAndServe start a tcp socks5 proxy server.
// For now, non-tcp socks5 proxy server is not necessary, so there is no "network" param.
// listenAddr example: "127.0.0.1:8000"
func (s *Server) ListenAndServe() error {
	if err := s.Listen(); err != nil {
		return err
	}
	return s.Serve()
}

// Listen start a tcp listener.
func (s *Server) Listen() error {
	// Create a SOCKS5 server
	conf := &socks5internal.Config{}
	if s.dialer != nil {
		conf.Dial = s.dialer
	}
	if s.dnsResolver != nil {
		conf.Resolver = s.dnsResolver
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
