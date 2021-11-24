package gmux

import (
	"io"
	"net"
)

// MuxConnIF is an interface for upper level multiplexing
// connection which based on underlying net.Conn.
type MuxConnIF interface {
	Open(streamName string) (io.ReadWriteCloser, error)
	Accept() (io.ReadWriteCloser, error)
	IsClosed() bool
	NumStreams() int
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	Close() error
}

// StreamIF is an interface for logical stream,
// it implements net.Conn.
type StreamIF interface {
	ID() uint32
	Name() string
	net.Conn
}
