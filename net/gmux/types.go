package gmux

import (
	"io"
	"net"
)

type Mux interface {
	Open(streamName string) (io.ReadWriteCloser, error)
	Accept() (io.ReadWriteCloser, error)
	IsClosed() bool
	NumStreams() int
	RemoteAddr() net.Addr
	Close() error
}

type StreamIF interface {
	io.ReadWriteCloser
	ID() uint32
	Name() string
	RemoteAddr() net.Addr
}
