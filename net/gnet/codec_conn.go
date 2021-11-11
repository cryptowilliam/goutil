package gnet

import (
	"io"
	"net"
	"time"
)

type (
	// CodecConn implements net.Conn, used to compress/decompress network connection.
	// For example, transferring log files or JSON files over the Internet,
	// adding a compression algorithm like snappy can greatly improve the efficiency of data transfer.
	CodecConn struct {
		conn                 net.Conn           // original net.Conn used to access LocalAddr/RemoteAddr/SetDeadline/Close...
		codecReadWriteCloser io.ReadWriteCloser // implements data compress/decompress here
	}
)

// NewCodecConn create CodecConn with original connection and codec io.ReadWriteCloser.
// Note: you should implement data compress/decompress at `codecReadWriteCloser`
func NewCodecConn(conn net.Conn, codecReadWriteCloser io.ReadWriteCloser) *CodecConn {
	rst := new(CodecConn)
	rst.conn = conn
	rst.codecReadWriteCloser = codecReadWriteCloser
	return rst
}

func (c *CodecConn) Read(p []byte) (n int, err error) {
	return c.codecReadWriteCloser.Read(p)
}

func (c *CodecConn) Write(p []byte) (n int, err error) {
	return c.codecReadWriteCloser.Write(p)
}

func (c *CodecConn) Close() error {
	_ = c.codecReadWriteCloser.Close()
	return c.conn.Close()
}

func (c *CodecConn) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *CodecConn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *CodecConn) SetDeadline(t time.Time) error {
	return c.conn.SetDeadline(t)
}

func (c *CodecConn) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

func (c *CodecConn) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}
