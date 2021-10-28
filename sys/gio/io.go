package gio

import (
	"bytes"
	"io"
	"io/ioutil"
	"strings"
	"sync"
	"time"
)

type SetDeadlineCallback func(t time.Time) error
type CopiedSizeCallback func(size int64)

// Forked from standard library io.Copy
func CopyTimeout(dst io.Writer, dstWriteCb SetDeadlineCallback, src io.Reader, srcReadCb SetDeadlineCallback, timeout time.Duration) (written int64, err error) {
	buf := make([]byte, 32*1024)
	var nr, nw int
	var er, ew error

	/*
		// If the reader has a WriteTo method, use it to do the copy.
		// Avoids an allocation and a copy.
		if wt, ok := src.(WriterTo); ok {
			return wt.WriteTo(dst)
		}
		// Similarly, if the writer has a ReadFrom method, use it to do the copy.
		if rt, ok := dst.(ReaderFrom); ok {
			return rt.ReadFrom(src)
		}
		if buf == nil {
			buf = make([]byte, 32*1024)
		}
	*/

	for {
		if timeout > 0 {
			srcReadCb(time.Now().Add(timeout))
		}
		nr, er = src.Read(buf)
		if nr > 0 {
			if timeout > 0 {
				dstWriteCb(time.Now().Add(timeout))
			}
			nw, ew = dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er == io.EOF {
			break
		}
		if er != nil {
			err = er
			break
		}
	}
	return written, err
}

// A pool for stream copying
var xmitBuf sync.Pool

func init() {
	xmitBuf.New = func() interface{} {
		return make([]byte, 32768)
	}
}

// https://github.com/xtaci/kcptun/blob/master/server/main.go
func CopyStream(dst io.Writer, src io.ReadCloser) chan struct{} {
	die := make(chan struct{})
	go func() {
		buf := xmitBuf.Get().([]byte)
		genericCopyBuffer(dst, src, buf)
		xmitBuf.Put(buf)
		close(die)
	}()
	return die
}

// https://github.com/xtaci/kcptun/blob/86cc46f437592e88b2504c79ef1ecfea37bb3cbb/generic/copy.go
// io.CopyBuffer has extra tests for interface like io.ReaderFrom and io.WriterTo
// which is not efficient in memory management from tests
func genericCopyBuffer(dst io.Writer, src io.Reader, buf []byte) (written int64, err error) {
	if buf != nil && len(buf) == 0 {
		panic("empty buffer in copyBuffer")
	}

	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}
	return written, err
}

// https://github.com/xtaci/kcptun/blob/master/server/main.go
// io.CopyBuffer has extra tests for interface like io.ReaderFrom and io.WriterTo
// which is not efficient in memory management from tests
func CopyBuffer(dst io.Writer, src io.Reader, buf []byte) (written int64, err error) {
	if buf != nil && len(buf) == 0 {
		panic("empty buffer in copyBuffer")
	}

	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}
	return written, err
}

// https://github.com/efarrer/iothrottler/
// https://github.com/jwkohnen/bwio/
// limit copy speed / duration, read copied size
func CopyEx(dst io.Writer, dstWriteCb SetDeadlineCallback, src io.Reader, srcReadCb SetDeadlineCallback, timeout time.Duration, sizeCallback CopiedSizeCallback) (written int64, err error) {
	buf := make([]byte, 32*1024)
	var nr, nw int
	var er, ew error
	lastNotify := time.Time{}

	/*
		// If the reader has a WriteTo method, use it to do the copy.
		// Avoids an allocation and a copy.
		if wt, ok := src.(WriterTo); ok {
			return wt.WriteTo(dst)
		}
		// Similarly, if the writer has a ReadFrom method, use it to do the copy.
		if rt, ok := dst.(ReaderFrom); ok {
			return rt.ReadFrom(src)
		}
		if buf == nil {
			buf = make([]byte, 32*1024)
		}
	*/

	for {
		if timeout > 0 {
			srcReadCb(time.Now().Add(timeout))
		}
		nr, er = src.Read(buf)
		if nr > 0 {
			if timeout > 0 {
				dstWriteCb(time.Now().Add(timeout))
			}
			nw, ew = dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)

				// notify callback
				if sizeCallback != nil {
					if time.Now().Sub(lastNotify) > time.Second {
						sizeCallback(written)
					}
				}
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er == io.EOF {
			break
		}
		if er != nil {
			err = er
			break
		}
	}
	return written, err
}

func ReadCloserToBytes(rd io.ReadCloser) ([]byte, error) {
	b, err := ioutil.ReadAll(rd)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func ReadCloserToInterface(rd io.ReadCloser) (interface{}, error) {
	b, err := ReadCloserToBytes(rd)
	if err != nil {
		return nil, err
	}
	return interface{}(b), nil
}

func ReadCloserToString(rd io.ReadCloser) (string, error) {
	b, err := ReadCloserToBytes(rd)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func BytesToReadCloser(b []byte) io.ReadCloser {
	return ioutil.NopCloser(bytes.NewReader(b))
}

func StringToReadCloser(s string) io.ReadCloser {
	return ioutil.NopCloser(strings.NewReader(s))
}
