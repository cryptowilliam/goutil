package gnet

import (
	"context"
	"net"
)

type (
	Dialer interface {
		Dial(network, remoteAddr string) (net.Conn, error)
	}

	DialerWithCtx interface {
		DialWithCtx(ctx context.Context, network, remoteAddr string) (net.Conn, error)
	}

	DialFunc = func(network, remoteAddr string) (net.Conn, error)

	DialWithCtxFunc = func(ctx context.Context, network, remoteAddr string) (net.Conn, error)
)
