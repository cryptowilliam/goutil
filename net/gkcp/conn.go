package gkcp

import (
	"github.com/xtaci/kcp-go"
	"io"
	"net"
	"time"
)

type KcpConn struct {
	server *Server
	rw     io.ReadWriteCloser
	sess   *kcp.UDPSession
	addr   string
}

func (c *KcpConn) Read(buf []byte) (n int, err error) {
	return c.rw.Read(buf)
}

func (c *KcpConn) Write(data []byte) (n int, err error) {
	return c.rw.Write(data)
}

func (c *KcpConn) SetDeadline(t time.Time) error {
	err := c.SetReadDeadline(t)
	if err != nil {
		return err
	}
	return c.SetWriteDeadline(t)
}

func (c *KcpConn) SetReadDeadline(t time.Time) error {
	return c.sess.SetReadDeadline(t)
}

func (c *KcpConn) SetWriteDeadline(t time.Time) error {
	return c.sess.SetWriteDeadline(t)
}

func (c *KcpConn) LocalAddr() net.Addr {
	return c.sess.LocalAddr()
}

func (c *KcpConn) RemoteAddr() net.Addr {
	return c.sess.RemoteAddr()
}

func (c *KcpConn) Close() error {
	_ = c.sess.Close()
	return c.rw.Close()
}
