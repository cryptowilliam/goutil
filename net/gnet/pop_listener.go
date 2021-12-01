package gnet

import (
	"container/list"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/safe/gwg"
	"github.com/cryptowilliam/goutil/sys/galloc"
	"net"
	"strings"
	"sync"
	"time"
)

type (
	// PopConn is packet-oriented protocols connection.
	// It implements net.Conn, makes PacketConn used like a net.Conn.
	PopConn struct {
		sync.RWMutex
		l          *PopListener
		localAddr  net.Addr
		remoteAddr net.Addr
		readBuf    *list.List
	}

	// PopListener is packet-oriented protocols listener.
	// It implements net.Listener, makes PacketConn used like a net.Listener.
	PopListener struct {
		network   string
		addr      string
		ls        net.PacketConn
		alloc     *galloc.Allocator
		connList  sync.Map
		chDie     chan struct{}
		chAccepts chan *PopConn
		chErr     chan error
		wg        *sync.WaitGroup
	}
)

func ListenPop(network, addr string) (*PopListener, error) {
	l := &PopListener{
		network:   network,
		addr:      addr,
		chDie:     make(chan struct{}),
		chAccepts: make(chan *PopConn, acceptBacklog),
		alloc:     galloc.NewAllocator(),
		wg:        gwg.New(),
	}
	err := error(nil)

	switch strings.ToLower(network) {
	case "udp", "udp4", "udp6", "unixgram", "ip:1", "ip:icmp":
		l.ls, err = net.ListenPacket(network, addr)
	default:
		err = gerrors.New("unsupported network %s", network)
	}
	if err != nil {
		return nil, err
	}

	l.wg.Add(1)
	go l.readRoutine()

	return l, nil
}

func (l *PopListener) readRoutine() {
	defer l.wg.Add(-1)

	for {
		select {
		case <-l.chDie:
			return
		default:
			buf := l.alloc.Get(2000)
			n, rmtAddr, err := l.ls.ReadFrom(buf)
			if err != nil {
				l.chErr <- err
				return
			}
			buf = buf[:n] // NOTICE: fix buf length to 'n'

			var conn *PopConn = nil
			if n > 0 {
				connIF, ok := l.connList.Load(rmtAddr.String())
				if ok {
					conn = connIF.(*PopConn)
				} else {
					conn = new(PopConn)
					conn.readBuf = list.New()
					conn.localAddr = l.Addr()
					conn.remoteAddr = rmtAddr
					conn.l = l
					l.chAccepts <- conn
				}

				conn.Lock()
				conn.readBuf.PushBack(buf)
				conn.Unlock()
			}

			if n > 0 {
				time.Sleep(time.Millisecond * 20)
			} else {
				time.Sleep(time.Millisecond * 100)
			}
		}
	}
}

func (l *PopListener) Accept() (net.Conn, error) {
	select {
	case <-l.chDie:
		return nil, nil
	case err := <- l.chErr:
		return nil, err
	case newConn := <-l.chAccepts:
		return newConn, nil
	}
}

func (l *PopListener) Close() error {
	close(l.chDie) // close read routine
	l.wg.Wait()    // wait read routine exit
	return l.ls.Close()
}

func (l *PopListener) Addr() net.Addr {
	return l.ls.LocalAddr()
}

func (c *PopConn) Read(b []byte) (int, error) {
	c.Lock()
	defer c.Unlock()

	if c.readBuf.Len() == 0 {
		return 0, nil
	}

	frontBuf := c.readBuf.Front().Value.([]byte)
	copyLen := copy(b, frontBuf)
	if copyLen >= len(frontBuf) {
		c.readBuf.Remove(c.readBuf.Front())
		return copyLen, nil
	} else {
		frontBuf = frontBuf[copyLen:]
		c.readBuf.Front().Value = frontBuf
		return copyLen, nil
	}
}

func (c *PopConn) Write(b []byte) (n int, err error) {
	return c.l.ls.WriteTo(b, c.remoteAddr)
}

func (c *PopConn) Close() error {
	c.Lock()
	defer c.Unlock()

	c.l.connList.Delete(c.remoteAddr)

	return nil
}

func (c *PopConn) LocalAddr() net.Addr {
	return c.localAddr
}

func (c *PopConn) RemoteAddr() net.Addr {
	return c.remoteAddr
}

// TODO
func (c *PopConn) SetDeadline(t time.Time) error {
	return nil
}

// TODO
func (c *PopConn) SetReadDeadline(t time.Time) error {
	return nil
}

// TODO
func (c *PopConn) SetWriteDeadline(t time.Time) error {
	return nil
}
