// MultiListener is a Listener that listens on multiple addresses and ports,
// and it is used like one Listener.

package gnet

import (
	"fmt"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"net"
	"sync"
)

type (
	mlAccepted struct {
		network string
		addr    string
		conn    net.Conn
		err     error
	}

	Listener struct {
		network  string       // listening network
		addr     string       // listen address
		listener net.Listener // raw listener
		chDie    chan struct{}
	}

	MultiListener struct {
		lns       []Listener
		lnsMtx    sync.RWMutex
		chAccepts chan mlAccepted
		chDie     chan struct{}
		wgClose   *sync.WaitGroup
	}
)

const acceptBacklog = 4096

// NewMultiListener creates new MultiListener.
func NewMultiListener() *MultiListener {
	return &MultiListener{
		chAccepts: make(chan mlAccepted, acceptBacklog),
		chDie:     make(chan struct{}),
	}
}

func (ml *MultiListener) listenRoutine(ln Listener) {
	for {
		select {
		case <-ml.chDie:
			return
		case <-ln.chDie:
			return
		default:
			newConn, err := ln.listener.Accept()
			ml.chAccepts <- mlAccepted{
				network: ln.network,
				addr:    ln.addr,
				conn:    newConn,
				err:     err,
			}
			if err != nil {
				return
			}
		}
	}
}

// AddListen add new listen address.
func (ml *MultiListener) AddListen(network, addr string) error {
	lnRaw, err := Listen(network, addr)
	if err != nil {
		return err
	}

	ln := Listener{
		network:  network,
		addr:     addr,
		listener: lnRaw,
	}
	ml.lnsMtx.Lock()
	ml.lns = append(ml.lns, ln)
	ml.lnsMtx.Unlock()

	go ml.listenRoutine(ln)
	return nil
}

// Accept waits for and returns the next connection to the listener.
func (ml *MultiListener) Accept() (network string, addr string, conn net.Conn, err error) {
	select {
	case <-ml.chDie:
		return "", "", nil, nil
	case accepted := <-ml.chAccepts:
		return accepted.network, accepted.addr, accepted.conn, accepted.err
	}
}

// CloseOne closes one listener.
func (ml *MultiListener) CloseOne(network, addr string) error {
	err := error(nil)
	ml.lnsMtx.RLock()
	for _, v := range ml.lns {
		if v.network == network && v.addr == addr {
			close(v.chDie)
			err = v.listener.Close()
		}
	}
	ml.lnsMtx.RUnlock()
	return err
}

// Close closes all the listeners.
// Any blocked Accept operations will be unblocked and return errors.
func (ml *MultiListener) Close() error {
	close(ml.chDie)

	var errs []error
	ml.lnsMtx.RLock()
	for _, v := range ml.lns {
		close(v.chDie)
		err := v.listener.Close()
		if err != nil {
			errs = append(errs, err)
		}
	}
	ml.lnsMtx.RUnlock()

	return gerrors.JoinArray(errs)
}

// Addr returns all the listener's network address.
func (ml *MultiListener) Addr() []string {
	var res []string

	ml.lnsMtx.RLock()
	for _, v := range ml.lns {
		res = append(res, fmt.Sprintf("%s://%s", v.network, v.addr))
	}
	ml.lnsMtx.RUnlock()

	return res
}
