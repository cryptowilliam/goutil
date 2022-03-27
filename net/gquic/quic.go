package gquic

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"github.com/lucas-clemente/quic-go"
	"math/big"
	"net"
	"time"
)

type QuicConn struct {
	sess   quic.Session
	stream quic.Stream
}

type QuicListener struct {
	listener quic.Listener
	stream   quic.Stream
}

const (
	alpn = "idontknow"
)

func Dial(raddr string) (*QuicConn, error) {
	sess, err := quic.DialAddr(raddr, &tls.Config{InsecureSkipVerify: true, NextProtos: []string{alpn}}, nil)
	if err != nil {
		return nil, err
	}
	stream, err := sess.OpenStream()
	if err != nil {
		return nil, err
	}

	return &QuicConn{sess: sess, stream: stream}, nil
}

func Listen(laddr string) (net.Listener, error) {
	ln, err := quic.ListenAddr(laddr, generateTLSConfig(), nil)
	if err != nil {
		return nil, err
	}
	return &QuicListener{listener: ln}, nil
}

// Setup a bare-bones TLS config for the server
func generateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(any(err))
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(any(err))
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(any(err))
	}
	return &tls.Config{Certificates: []tls.Certificate{tlsCert}, NextProtos: []string{alpn}}
}

// Accept waits for and returns the next connection to the listener.
func (l *QuicListener) Accept() (net.Conn, error) {
	sess, err := l.listener.Accept(context.Background())
	if err != nil {
		return nil, err
	}
	stream, err := sess.AcceptStream(context.Background())
	if err != nil {
		return nil, err
	}

	return &QuicConn{sess: sess, stream: stream}, nil
}

// Closes the listener.
// Any blocked Accept operations will be unblocked and return gerrors.
func (l *QuicListener) Close() error {
	return l.listener.Close()
}

// Addr returns the listener's network address.
func (l *QuicListener) Addr() net.Addr {
	return l.Addr()
}

// Read reads data from the connection.
// Read can be made to time out and return an Error with Timeout() == true
// after a fixed time limit; see SetDeadline and SetReadDeadline.
func (c *QuicConn) Read(b []byte) (n int, err error) {
	return c.stream.Read(b)
}

// Write writes data to the connection.
// Write can be made to time out and return an Error with Timeout() == true
// after a fixed time limit; see SetDeadline and SetWriteDeadline.
func (c *QuicConn) Write(b []byte) (n int, err error) {
	return c.stream.Write(b)
}

// Closes the connection.
// Any blocked Read or Write operations will be unblocked and return gerrors.
func (c *QuicConn) Close() error {
	return c.sess.CloseWithError(0, "")
}

// LocalAddr returns the local network address.
func (c *QuicConn) LocalAddr() net.Addr {
	return c.sess.LocalAddr()
}

// RemoteAddr returns the remote network address.
func (c *QuicConn) RemoteAddr() net.Addr {
	return c.sess.RemoteAddr()
}

// SetDeadline sets the read and write deadlines associated
// with the connection. It is equivalent to calling both
// SetReadDeadline and SetWriteDeadline.
//
// A deadline is an absolute time after which I/O operations
// fail with a timeout (see type Error) instead of
// blocking. The deadline applies to all future and pending
// I/O, not just the immediately following call to Read or
// Write. After a deadline has been exceeded, the connection
// can be refreshed by setting a deadline in the future.
//
// An idle timeout can be implemented by repeatedly extending
// the deadline after successful Read or Write calls.
//
// A zero value for t means I/O operations will not time out.
func (c *QuicConn) SetDeadline(t time.Time) error {
	return c.stream.SetDeadline(t)
}

// SetReadDeadline sets the deadline for future Read calls
// and any currently-blocked Read call.
// A zero value for t means Read will not time out.
func (c *QuicConn) SetReadDeadline(t time.Time) error {
	return c.stream.SetReadDeadline(t)
}

// SetWriteDeadline sets the deadline for future Write calls
// and any currently-blocked Write call.
// Even if write times out, it may return n > 0, indicating that
// some of the data was successfully written.
// A zero value for t means Write will not time out.
func (c *QuicConn) SetWriteDeadline(t time.Time) error {
	return c.stream.SetWriteDeadline(t)
}
