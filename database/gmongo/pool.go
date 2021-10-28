package gmongo

import (
	"context"
	"sync"
	"time"
)

// fixed size mongodb connection pool

type ConnFromPool struct {
	pool      *Pool
	rawConn   *Conn
	pos       int
	available bool
}

func (ce *ConnFromPool) RawConn() *Conn {
	return ce.rawConn
}

type Pool struct {
	checkAlive *Conn
	list       map[int]*ConnFromPool
	stateMu    sync.Mutex // mutex for available state
	size       int
}

func DialPool(mongodbUrl string, size int) (*Pool, error) {
	kc, err := Dial(mongodbUrl)
	if err != nil {
		return nil, err
	}

	r := &Pool{}
	r.list = make(map[int]*ConnFromPool)
	r.size = 0
	r.checkAlive = kc

	for i := 0; i < size; i++ {
		rc, err := Dial(mongodbUrl)
		if err != nil {
			r.Close()
			return nil, err
		}
		cex := ConnFromPool{rawConn: rc, pos: i, available: true, pool: r}
		r.list[i] = &cex
		r.size++
	}

	return r, nil
}

func (p *Pool) GetAvailable() (*ConnFromPool, bool) {
	p.stateMu.Lock()
	defer p.stateMu.Unlock()

	for k, v := range p.list {
		if v.available {
			p.list[k].available = false
			cp := *p.list[k] // returns copy, then I can destroy the pointer when PutBack it
			return &cp, true
		}
	}
	return nil, false
}

func (p *Pool) WaitAvailable() *ConnFromPool {
	for {
		p.stateMu.Lock()
		for k, v := range p.list {
			if v.available {
				p.list[k].available = false
				cp := *p.list[k] // returns copy, then I can destroy the pointer when PutBack it
				p.stateMu.Unlock()
				return &cp
			}
		}
		p.stateMu.Unlock()
		time.Sleep(time.Second)
	}
}

func (p *Pool) PutBack(cex **ConnFromPool) {
	p.stateMu.Lock()
	defer p.stateMu.Unlock()

	p.list[(*cex).pos].available = false

	// destroy copy
	*cex = nil
}

func (p *Pool) Ping() error {
	return p.checkAlive.inCli.Ping(context.Background(), nil)
}

func (p *Pool) Close() {
	if p.checkAlive != nil {
		_ = p.checkAlive.Close()
	}

	for i := 0; i < p.size; i++ {
		_ = p.list[i].rawConn.Close()
	}
}
