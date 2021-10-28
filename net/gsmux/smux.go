package gsmux

import (
	"github.com/cryptowilliam/goutil/net/gmux"
	"github.com/xtaci/smux"
	"net"
	"time"
)

func NewSmuxClient(conn net.Conn, muxBufferSize, muxkeepAliveSeconds int) (*smux.Session, error) {
	smuxConfig := smux.DefaultConfig()
	smuxConfig.MaxReceiveBuffer = muxBufferSize
	smuxConfig.KeepAliveInterval = time.Duration(muxkeepAliveSeconds) * time.Second

	return smux.Client(conn, smuxConfig)
}

func NewSmuxServer(conn net.Conn, muxBufferSize, muxkeepAliveSeconds int) (*smux.Session, error) {
	smuxConfig := smux.DefaultConfig()
	smuxConfig.MaxReceiveBuffer = muxBufferSize
	smuxConfig.KeepAliveInterval = time.Duration(muxkeepAliveSeconds) * time.Second

	return smux.Server(conn, smuxConfig)
}

func ConvertInternalError(err error) error {
	if err == smux.ErrInvalidProtocol {
		return gmux.ErrInvalidProtocol
	}
	return err
}
