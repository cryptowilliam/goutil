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
	WriteFlusher interface {
		io.Writer
		gio.Flusher
	}

	CompReadWriteCloser struct {
		algo  CompAlgo
		param *CompParam
		rwc   io.ReadWriteCloser
		r     io.Reader
		w     WriteFlusher
	}

	CompStream struct {
		conn            net.Conn
		readWriteCloser *CompReadWriteCloser
	}
)

func (c *CompReadWriteCloser) Read(p []byte) (n int, err error) {
	return c.r.Read(p)
}

func (c *CompReadWriteCloser) Write(p []byte) (n int, err error) {
	n, err = c.w.Write(p)
	if err != nil {
		return 0, err
	}
	if err := c.w.Flush(); err != nil {
		return 0, err
	}

	return n, nil
}

func (c *CompReadWriteCloser) Close() error {
	return c.rwc.Close()
}

func NewCompReadWriteCloser(compAlgo CompAlgo, param *CompParam, rwc io.ReadWriteCloser) (*CompReadWriteCloser, error) {
	rst := new(CompReadWriteCloser)
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

func (c *CompStream) Read(p []byte) (n int, err error) {
	return c.readWriteCloser.Read(p)
}

func (c *CompStream) Write(p []byte) (n int, err error) {
	return c.readWriteCloser.Write(p)
}

func (c *CompStream) Close() error {
	return c.readWriteCloser.Close()
}

func (c *CompStream) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *CompStream) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *CompStream) SetDeadline(t time.Time) error {
	return c.conn.SetDeadline(t)
}

func (c *CompStream) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

func (c *CompStream) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}

func NewCompStream(compAlgo CompAlgo, param *CompParam, conn net.Conn) (*CompStream, error) {
	rst := new(CompStream)
	readWriteCloser, err := NewCompReadWriteCloser(compAlgo, param, conn)
	if err != nil {
		return nil, err
	}
	rst.readWriteCloser = readWriteCloser
	rst.conn = conn
	return rst, nil
}
