package gnet

import (
	"container/list"
	"fmt"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/safe/gwg"
	"github.com/cryptowilliam/goutil/sys/galloc"
	"github.com/cryptowilliam/goutil/sys/gtime"
	"github.com/pkg/errors"
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
		l                *PopListener
		localAddr        net.Addr
		remoteAddr       net.Addr
		readBuf          *list.List
		chReadEvent      chan struct{}
		readDeadline     time.Time
		// Note:
		// When no data timeout error is triggered in the Read or Write function and been returned,
		// io.Copy will also end because of Read/Write/WriteTo returns error, thus freeing up more resources.
		noDataTimeout time.Duration // long time no data read/write timeout duration, it is different from read write deadline
		noDataTicker  *time.Ticker // long time no data read/write timeout ticker
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
		chReadErr chan error
		wg        *sync.WaitGroup
	}
)

var (
	monoClock = gtime.NewMonoClock()
)

func ListenPop(network, addr string) (*PopListener, error) {
	l := &PopListener{
		network:   network,
		addr:      addr,
		chDie:     make(chan struct{}),
		chAccepts: make(chan *PopConn, acceptBacklog),
		chReadErr: make(chan error, 1),
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

func (l *PopListener) notifyReadEvent(pop *PopConn) {
	select {
	case pop.chReadEvent <- struct{}{}:
	default:
	}
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
				l.chReadErr <- err
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
					conn.chReadEvent = make(chan struct{}, 1)
					conn.l = l
					conn.noDataTimeout = 5 * time.Minute
					conn.noDataTicker = time.NewTicker(conn.noDataTimeout)
					l.chAccepts <- conn
				}

				conn.Lock()
				conn.readBuf.PushBack(buf)
				conn.Unlock()
				l.notifyReadEvent(conn)
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
	case err := <-l.chReadErr:
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

	for {
		// deadline for current reading operation
		var timeout *time.Timer
		var chDeadlineTimeout <-chan time.Time
		if !c.readDeadline.IsZero() {
			if time.Now().After(c.readDeadline) {
				return 0, errors.WithStack(fmt.Errorf("read timeout"))
			}
			delay := time.Until(c.readDeadline)
			timeout = time.NewTimer(delay)
			chDeadlineTimeout = timeout.C
		}

		select {
		case err := <-c.l.chReadErr:
			return 0, err
		case <-c.chReadEvent:
			if c.readBuf.Len() == 0 {
				return 0, nil
			}
			frontBuf := c.readBuf.Front().Value.([]byte)
			copyLen := copy(b, frontBuf)
			if copyLen > 0 {
				c.noDataTicker.Reset(c.noDataTimeout)
			}
			if copyLen >= len(frontBuf) {
				c.readBuf.Remove(c.readBuf.Front())
				return copyLen, nil
			} else {
				frontBuf = frontBuf[copyLen:]
				c.readBuf.Front().Value = frontBuf
				return copyLen, nil
			}
		case <-chDeadlineTimeout:
			return 0, errors.WithStack(fmt.Errorf("deadline timeout"))
		case <-c.noDataTicker.C: // it is different from read write deadline timeout
			return 0, errors.WithStack(fmt.Errorf("no data timeout"))
		}
	}
}

func (c *PopConn) Write(b []byte) (n int, err error) {
	select {
	case <-c.noDataTicker.C:
		return 0, errors.WithStack(fmt.Errorf("no data timeout"))
	default:
		n, err = c.l.ls.WriteTo(b, c.remoteAddr)
		if n > 0 && err == nil {
			c.noDataTicker.Reset(c.noDataTimeout)
		}
		return n, err
	}
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

func (c *PopConn) SetReadDeadline(t time.Time) error {
	c.Lock()
	defer c.Unlock()
	c.readDeadline = t
	return nil
}

// TODO
func (c *PopConn) SetWriteDeadline(t time.Time) error {
	return nil
}
