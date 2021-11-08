package gsecureconn

import (
	"net"
	"time"
)

/*
const (
	bufSize = 1024
)

var bpool sync.Pool

func init() {
	bpool.New = func() interface{} {
		return make([]byte, bufSize)
	}
}
func bufferPoolGet() []byte {
	return bpool.Get().([]byte)
}
func bufferPoolPut(b []byte) {
	bpool.Put(b)
}
*/
// 加密传输的 Socket
type SecureConn struct {
	conn   net.Conn
	cipher *Cipher
}

// see net.DialTCP
func WrapConn(rawConn net.Conn, cipher *Cipher) (net.Conn, error) {
	// Conn被关闭时直接清除所有数据 不管没有发送的数据
	//rawConn.SetLinger(0)

	return &SecureConn{
		conn:   rawConn,
		cipher: cipher,
	}, nil
}

// 从输入流里读取加密过的数据，解密后把原数据放到bs里
func (c *SecureConn) Read(b []byte) (n int, err error) {
	n, err = c.conn.Read(b)
	if err != nil {
		return
	}
	c.cipher.Decode(b[:n])
	return
}

// 把放在bs里的数据加密后立即全部写入输出流
func (c *SecureConn) Write(b []byte) (int, error) {
	c.cipher.Encode(b)
	return c.Write(b)
}

func (c *SecureConn) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *SecureConn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *SecureConn) SetDeadline(t time.Time) error {
	return c.conn.SetDeadline(t)
}

func (c *SecureConn) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

func (c *SecureConn) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}

func (c *SecureConn) Close() error {
	return c.conn.Close()
}

/*
// 从src中源源不断的读取原数据加密后写入到dst，直到src中没有数据可以再读取
func (c *SecureConn) EncodeCopy(dst net.Conn) error {
	buf := bufferPoolGet()
	defer bufferPoolPut(buf)
	for {
		readCount, errRead := c.conn.Read(buf)
		if errRead != nil {
			if errRead != io.EOF {
				return errRead
			} else {
				return nil
			}
		}
		if readCount > 0 {
			writeCount, errWrite := (&SecureConn{
				conn:   dst,
				cipher: c.cipher,
			}).Write(buf[0:readCount])
			if errWrite != nil {
				return errWrite
			}
			if readCount != writeCount {
				return io.ErrShortWrite
			}
		}
	}
}

// 从src中源源不断的读取加密后的数据解密后写入到dst，直到src中没有数据可以再读取
func (c *SecureConn) DecodeCopy(dst net.Conn) error {
	buf := bufferPoolGet()
	defer bufferPoolPut(buf)
	for {
		readCount, errRead := c.Read(buf)
		if errRead != nil {
			if errRead != io.EOF {
				return errRead
			} else {
				return nil
			}
		}
		if readCount > 0 {
			writeCount, errWrite := dst.Write(buf[0:readCount])
			if errWrite != nil {
				return errWrite
			}
			if readCount != writeCount {
				return io.ErrShortWrite
			}
		}
	}
}
*/
