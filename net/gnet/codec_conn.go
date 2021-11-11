package gnet

import (
	"io"
	"net"
	"time"
)

type (
	// CodecConn is a generic stream-oriented network connection that can be compressed and decompressed.
	// It implements net.Conn interface. It is used to compress/decompress network connection stream,
	// for example, when transferring log files or JSON files over the Internet,
	// adding a compression algorithm like snappy can greatly improve the efficiency of data transfer.
	CodecConn struct {
		conn                 net.Conn           // original net.Conn used to access LocalAddr/RemoteAddr/SetDeadline/Close...
		codecReadWriteCloser io.ReadWriteCloser // implements data compress/decompress here
	}
)

// NewCodecConn create CodecConn with original connection and codec io.ReadWriteCloser.
// Note: you should implement data compress/decompress at `codecReadWriteCloser`
func NewCodecConn(conn net.Conn, codec io.ReadWriteCloser) *CodecConn {
	rst := new(CodecConn)
	rst.conn = conn
	rst.codecReadWriteCloser = codec
	return rst
}

// Read implements net.Conn.
func (c *CodecConn) Read(p []byte) (n int, err error) {
	return c.codecReadWriteCloser.Read(p)
}

// Write implements net.Conn.
func (c *CodecConn) Write(p []byte) (n int, err error) {
	return c.codecReadWriteCloser.Write(p)
}

// Close implements net.Conn.
func (c *CodecConn) Close() error {
	_ = c.codecReadWriteCloser.Close()
	return c.conn.Close()
}

// LocalAddr implements net.Conn.
func (c *CodecConn) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

// RemoteAddr implements net.Conn.
func (c *CodecConn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

// SetDeadline implements net.Conn.
func (c *CodecConn) SetDeadline(t time.Time) error {
	return c.conn.SetDeadline(t)
}

// SetReadDeadline implements net.Conn.
func (c *CodecConn) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

// SetWriteDeadline implements net.Conn.
func (c *CodecConn) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}
