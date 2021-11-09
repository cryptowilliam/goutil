package gcompress

import (
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/sys/gio"
	"github.com/golang/snappy"
	"github.com/klauspost/compress/s2"
	"github.com/klauspost/compress/zstd"
	"github.com/klauspost/pgzip"
	"io"
	"net"
	"time"
)

type (
	WriterFlusher interface {
		io.Writer
		gio.Flusher
	}

	CompIOStream struct {
		algo  CompAlgo
		param *CompParam
		rwc   io.ReadWriteCloser
		r     io.Reader
		w     WriterFlusher
	}

	CompNetStream struct {
		conn     net.Conn
		ioStream *CompIOStream
	}
)

func (c *CompIOStream) Read(p []byte) (n int, err error) {
	return c.r.Read(p)
}

func (c *CompIOStream) Write(p []byte) (n int, err error) {
	n, err = c.w.Write(p)
	if err != nil {
		return 0, err
	}
	if err := c.w.Flush(); err != nil {
		return 0, err
	}

	return n, nil
}

func (c *CompIOStream) Close() error {
	return c.rwc.Close()
}

func NewIOStream(compAlgo CompAlgo, param *CompParam, rwc io.ReadWriteCloser) (*CompIOStream, error) {
	rst := new(CompIOStream)
	rst.algo = compAlgo
	if param != nil {
		*rst.param = *param
		if err := param.Verify(compAlgo); err != nil {
			return nil, err
		}
	}
	rst.rwc = rwc
	switch compAlgo {
	case CompAlgoSnappy:
		rst.r = snappy.NewReader(rwc)
		rst.w = snappy.NewBufferedWriter(rwc)
	case CompAlgoS2:
		rst.r = s2.NewReader(rwc)
		rst.w = s2.NewWriter(rwc)
	case CompAlgoGzip:
		err := error(nil)
		rst.r, err = gzip.NewReader(rwc)
		if err != nil {
			return nil, err
		}
		rst.w = gzip.NewWriter(rwc)
	case CompAlgoPgZip:
		err := error(nil)
		rst.r, err = pgzip.NewReader(rwc)
		if err != nil {
			return nil, err
		}
		rst.w = pgzip.NewWriter(rwc)
	case CompAlgoZStd:
		err := error(nil)
		rst.r, err = zstd.NewReader(rwc)
		if err != nil {
			return nil, err
		}
		rst.w, err = zstd.NewWriter(rwc)
		if err != nil {
			return nil, err
		}
	case CompAlgoZLib:
		err := error(nil)
		rst.r, err = zlib.NewReader(rwc)
		if err != nil {
			return nil, err
		}
		rst.w = zlib.NewWriter(rwc)
	case CompAlgoFlate:
		rst.r = flate.NewReader(rwc)
		level := -1 // default level: -1
		if param != nil {
			level = param.Level
		}
		err := error(nil)
		rst.w, err = flate.NewWriter(rwc, level)
		if err != nil {
			return nil, err
		}
	default:
		return nil, gerrors.New("unsupported compress algorithm %s", compAlgo)
	}
	return rst, nil
}

func (c *CompNetStream) Read(p []byte) (n int, err error) {
	return c.ioStream.Read(p)
}

func (c *CompNetStream) Write(p []byte) (n int, err error) {
	return c.ioStream.Write(p)
}

func (c *CompNetStream) Close() error {
	return c.ioStream.Close()
}

func (c *CompNetStream) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *CompNetStream) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *CompNetStream) SetDeadline(t time.Time) error {
	return c.conn.SetDeadline(t)
}

func (c *CompNetStream) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

func (c *CompNetStream) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}

func NewNetStream(compAlgo CompAlgo, param *CompParam, conn net.Conn) (*CompNetStream, error) {
	rst := new(CompNetStream)
	ioStream, err := NewIOStream(compAlgo, param, conn)
	if err != nil {
		return nil, err
	}
	rst.ioStream = ioStream
	rst.conn = conn
	return rst, nil
}
