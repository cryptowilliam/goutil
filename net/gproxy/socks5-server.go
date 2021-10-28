package gproxy

// https://github.com/Ichbinjoe/go-socks/
// https://github.com/kumakichi/go-socks5/ client.go shows how to send http request via socks5 protocol in TCP

type Socks5Server struct {
	addr                 string
	server               *socks5.Server
	network              string
	cascadingSocks5Proxy string
}

// Create a SOCKS5 server
func newServer(addr, username, password, cascadingSocks5Proxy string) (*Socks5Server, error) {
	s := Socks5Server{addr: addr}
	var err error

	conf := &socks5.Config{}
	s.server, err = socks5.New(conf)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (s *Socks5Server) ListenAndServe() {
	s.server.ListenAndServe("tcp", s.addr)
}
