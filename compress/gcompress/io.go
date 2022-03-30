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
		algo  Comp
		param *CompParam
		rwc   io.ReadWriteCloser
		r     io.Reader
		w     gio.WriteFlusher
	}
)

// Read implements io.ReadWriteCloser.
func (c *CompReadWriteCloser) Read(p []byte) (n int, err error) {
	return c.r.Read(p)
}

// Write implements io.ReadWriteCloser.
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

// Close implements io.ReadWriteCloser.
func (c *CompReadWriteCloser) Close() error {
	return c.rwc.Close()
}

// NewCompReadWriteCloser create compress-orient ReadWriteCloser.
func NewCompReadWriteCloser(compAlgo Comp, param *CompParam, rwc io.ReadWriteCloser) (*CompReadWriteCloser, error) {
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
	case CompNone:
		return nil, gerrors.New("can't create CompReadWriteCloser for 'none' algo")
	case CompSnappy:
		rst.r = snappy.NewReader(rwc)
		rst.w = snappy.NewBufferedWriter(rwc)
	case CompS2:
		rst.r = s2.NewReader(rwc)
		rst.w = s2.NewWriter(rwc)
	case CompGzip:
		err := error(nil)
		rst.r, err = gzip.NewReader(rwc)
		if err != nil {
			return nil, err
		}
		rst.w = gzip.NewWriter(rwc)
	case CompPgZip:
		err := error(nil)
		rst.r, err = pgzip.NewReader(rwc)
		if err != nil {
			return nil, err
		}
		rst.w = pgzip.NewWriter(rwc)
	case CompZStd:
		err := error(nil)
		rst.r, err = zstd.NewReader(rwc)
		if err != nil {
			return nil, err
		}
		rst.w, err = zstd.NewWriter(rwc)
		if err != nil {
			return nil, err
		}
	case CompZLib:
		err := error(nil)
		rst.r, err = zlib.NewReader(rwc)
		if err != nil {
			return nil, err
		}
		rst.w = zlib.NewWriter(rwc)
	case CompFlate:
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
