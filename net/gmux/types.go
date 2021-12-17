package gmux

import (
	"io"
	"net"
)

type (
	CloseNotifier func(stream StreamIF, ctx interface{})

	// MuxConnIF is an interface for upper level multiplexing
	// connection which based on underlying net.Conn.
	MuxConnIF interface {
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
	StreamIF interface {
		ID() uint32
		Name() string
		SetCloseNotifier(notifier CloseNotifier, ctx interface{})
		net.Conn
	}
)
