package grpcs

// A rpc proxy used to cache(save the network traffic) and filter(only allow whitelist functions).

import (
	"fmt"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/basic/glog"
	"sort"
	"strings"
	"sync"
)

type (
	Proxy struct {
		c            *Client
		s            *Server
		allowFns     map[string]struct{}
		cacheFns     map[string]bool
		cacheReplies map[string]*Reply
		mu           sync.RWMutex
	}
)

// build new proxy
func NewProxy(clientRpcType RpcType, clientNetwork, clientAddress string,
	serverRpcType RpcType, serverNetwork, serverAddress string) (*Proxy, error) {
	c, err := Dial(clientRpcType, clientNetwork, clientAddress, *NewParamChecker())
	if err != nil {
		return nil, err
	}
	s, err := Listen(serverRpcType, serverNetwork, serverAddress)
	if err != nil {
		return nil, err
	}

	return &Proxy{
		c:            c,
		s:            s,
		cacheReplies: map[string]*Reply{},
	}, nil
}

func (p *Proxy) AddAllowFunc(fn string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.allowFns[fn] = struct{}{}
}

// Add function names which need to cache.
func (p *Proxy) AddCacheFunc(fn string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.cacheReplies[fn] = nil
}

// Format request into string.
func (p *Proxy) format(req Request) string {
	var ss []string
	for k, v := range req.Params {
		ss = append(ss, k+"="+fmt.Sprintf("%v", v))
	}
	// "Request.Params" is a map, and it is disorderly.
	// Sort it and make that same "Request" has same format() output.
	sort.Strings(ss)
	key := req.Func + "(" + strings.Join(ss, ",") + ")"
	return key
}

func (p *Proxy) onReq(in Request, out *Reply) error {
	p.mu.RLock()
	_, allow := p.allowFns[in.Func]
	_, needCache := p.cacheFns[p.format(in)]
	nowCache, hasCache := p.cacheReplies[p.format(in)]
	p.mu.RUnlock()

	if !allow {
		return gerrors.New("Func(%s) now allowed", in.Func)
	}

	if needCache {
		glog.Debgf("%s need cache", p.format(in))
		if hasCache {
			glog.Debgf("%s has cache", p.format(in))
			*out = *nowCache
			return nil
		} else {
			glog.Debgf("%s doesn't has cache, send request and cache it if request success", p.format(in))
			if err := p.c.Call(in.Func, in, out); err != nil {
				return err
			}
			p.mu.Lock()
			p.cacheReplies[p.format(in)] = out
			p.mu.Unlock()
			return nil
		}
	} else {
		glog.Debgf("func(%s) doesn't need cache", in.Func)
		return p.c.Call(in.Func, in, out)
	}
}

func (p *Proxy) Run() error {
	pc := NewParamChecker()
	return p.s.Run(*pc, p.onReq)
}

func (p *Proxy) Close() error {
	if err := p.c.Close(); err != nil {
		return err
	}
	if err := p.s.Close(); err != nil {
		return err
	}
	return nil
}
