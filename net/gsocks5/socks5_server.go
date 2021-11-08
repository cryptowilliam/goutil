package gsocks5

import (
	socks5internal "github.com/cryptowilliam/goutil/net/gsocks5/socks5internal"
)

// ListenAndServe start a tcp socks5 proxy server.
// For now, non-tcp socks5 proxy server is not necessary, so there is no "network" param.
// listenAddr example: "127.0.0.1:8000"
func ListenAndServe(listenAddr string) error {
	// Create a SOCKS5 server
	conf := &socks5internal.Config{}
	server, err := socks5internal.New(conf)
	if err != nil {
		return err
	}

	// Create SOCKS5 proxy on localhost port 8000
	if err := server.ListenAndServe("tcp", listenAddr); err != nil {
		return err
	}
	return nil
}
