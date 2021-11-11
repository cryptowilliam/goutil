// IO enhancements to compression algorithms.

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
	// CompReadWriteCloser implements io.ReadWriteCloser
	CompReadWriteCloser struct {
		algo  CompAlgo
		param *CompParam
		rwc   io.ReadWriteCloser
		r     io.Reader
		w     gio.WriteFlusher
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
