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
)

type (
	WriterFlusher interface {
		io.Writer
		gio.Flusher
	}

	CompStream struct {
		algo  CompAlgo
		param *CompParam
		conn  io.ReadWriteCloser
		w     WriterFlusher
		r     io.Reader
	}
)

func (c *CompStream) Read(p []byte) (n int, err error) {
	return c.r.Read(p)
}

func (c *CompStream) Write(p []byte) (n int, err error) {
	n, err = c.w.Write(p)
	err = c.w.Flush()
	return n, err
}

func (c *CompStream) Close() error {
	return c.conn.Close()
}

func NewStream(compAlgo CompAlgo, param *CompParam, conn io.ReadWriteCloser) (io.ReadWriteCloser, error) {
	rst := new(CompStream)
	rst.algo = compAlgo
	if param != nil {
		*rst.param = *param
		if err := param.Verify(compAlgo); err != nil {
			return nil, err
		}
	}
	rst.conn = conn
	switch compAlgo {
	case CompAlgoSnappy:
		rst.r = snappy.NewReader(conn)
		rst.w = snappy.NewBufferedWriter(conn)
	case CompAlgoS2:
		rst.r = s2.NewReader(conn)
		rst.w = s2.NewWriter(conn)
	case CompAlgoGzip:
		err := error(nil)
		rst.r, err = gzip.NewReader(conn)
		if err != nil {
			return nil, err
		}
		rst.w = gzip.NewWriter(conn)
	case CompAlgoPgZip:
		err := error(nil)
		rst.r, err = pgzip.NewReader(conn)
		if err != nil {
			return nil, err
		}
		rst.w = pgzip.NewWriter(conn)
	case CompAlgoZStd:
		err := error(nil)
		rst.r, err = zstd.NewReader(conn)
		if err != nil {
			return nil, err
		}
		rst.w, err = zstd.NewWriter(conn)
		if err != nil {
			return nil, err
		}
	case CompAlgoZLib:
		err := error(nil)
		rst.r, err = zlib.NewReader(conn)
		if err != nil {
			return nil, err
		}
		rst.w = zlib.NewWriter(conn)
	case CompAlgoFlate:
		rst.r = flate.NewReader(conn)
		level := -1 // default level: -1
		if param != nil {
			level = param.Level
		}
		err := error(nil)
		rst.w, err = flate.NewWriter(conn, level)
		if err != nil {
			return nil, err
		}
	default:
		return nil, gerrors.New("unsupported compress algorithm %s", compAlgo)
	}
	return rst, nil
}
