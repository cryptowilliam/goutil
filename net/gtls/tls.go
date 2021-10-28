package gtls

// Wrap net.Conn and net.Listener with TLS.

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"net"
)

func NewTLSConn(conn net.Conn, certFile, keyFile string) (net.Conn, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}
	if len(cert.Certificate) != 2 {
		return nil, gerrors.New("client.crt should have 2 concatenated certificates: client + CA")
	}
	ca, err := x509.ParseCertificate(cert.Certificate[1])
	if err != nil {
		return nil, err
	}
	certPool := x509.NewCertPool()
	certPool.AddCert(ca)
	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      certPool,
	}

	return tls.Client(conn, config), nil
}

func NewTLSListener(lis net.Listener, certFile, keyFile string) (net.Listener, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}
	if len(cert.Certificate) != 2 {
		return nil, gerrors.New("server.crt should have 2 concatenated certificates: server + CA")
	}
	ca, err := x509.ParseCertificate(cert.Certificate[1])
	if err != nil {
		return nil, err
	}
	certPool := x509.NewCertPool()
	certPool.AddCert(ca)
	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	}
	config.Rand = rand.Reader
	return tls.NewListener(lis, config), nil
}
