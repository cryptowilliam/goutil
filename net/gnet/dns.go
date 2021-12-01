package gnet

import (
	"context"
	"github.com/AdguardTeam/dnsproxy/upstream"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/container/ginterface"
	"github.com/cryptowilliam/goutil/container/gstring"
	"github.com/miekg/dns"
	"net"
	"strings"
	"sync"
	"time"
)

type (
	// DNSResolver is used to implement custom name resolution
	DNSResolver interface {
		LookupIP(host string) ([]net.IP, error)
		LookupAddr(addr string) (names []string, err error)
	}

	dnsSrv struct {
		host string
		u    upstream.Upstream
	}

	DNSClient struct {
		customDNSServers    map[string]dnsSrv // specified DNS servers
		useSysDNSIfNoCustom bool
		sync.RWMutex
	}

	LookupIPWithCtxFunc = func(ctx context.Context, host string) ([]net.IP, error)
)

var (
	// SysDNSResolver uses the system DNS to resolve host names
	SysDNSResolver = NewDNSClient().CleanCustomServers().UseSysDNSIfNoCustom(true)
)

func NewDNSClient() *DNSClient {
	return &DNSClient{
		customDNSServers: map[string]dnsSrv{},
	}
}

// AddCustomDNSServer adds new custom DNS server.
// host samples:
// Plain: 8.8.8.8:53
// Plain: 8.8.4.4:53
// Plain: 1.0.0.1:53
// Plain: 1.1.1.1:53
// DNS-over-TLS: tls://dns.adguard.com
// DNS-over-HTTPS: https://8.8.8.8/dns-query
// DNS-over-HTTPS: https://dns.adguard.com/dns-query
// DNSCrypt-stamp: sdns://AQIAAAAAAAAAFDE3Ni4xMDMuMTMwLjEzMDo1NDQzINErR_JS3PLCu_iZEIbq95zkSV2LFsigxDIuUso_OQhzIjIuZG5zY3J5cHQuZGVmYXVsdC5uczEuYWRndWFyZC5jb20
// DNS-over-QUIC: quic://dns.adguard.com
//
// More public DNS servers:
// https://github.com/DNSCrypt/dnscrypt-resolvers/blob/master/v3/public-resolvers.md
func (dl *DNSClient) AddCustomDNSServer(host string) error {
	opts := &upstream.Options{
		Timeout:            time.Duration(10) * time.Second,
		InsecureSkipVerify: false, // don't set it true
	}

	if IsIPString(host) {
		netIP := net.ParseIP(host)
		if netIP == nil {
			return gerrors.New("invalid IP %s", host)
		}
		opts.ServerIPAddrs = []net.IP{netIP}
	}

	u, err := upstream.AddressToUpstream(host, opts)
	if err != nil {
		return gerrors.New("Cannot create an upstream: %s", err.Error())
	}

	dl.customDNSServers[host] = dnsSrv{
		host: host,
		u:    u,
	}

	return nil
}

func (dl *DNSClient) RemoveCustomDNSServer(host, ip string) {
	delete(dl.customDNSServers, host)
	delete(dl.customDNSServers, ip)
}

// CleanCustomServers uses the system DNS to resolve host names.
func (dl *DNSClient) CleanCustomServers() *DNSClient {
	dl.Lock()
	defer dl.Unlock()
	dl.customDNSServers = make(map[string]dnsSrv)
	return dl
}

func (dl *DNSClient) UseSysDNSIfNoCustom(use bool) *DNSClient {
	dl.useSysDNSIfNoCustom = use
	return dl
}

func (dl *DNSClient) LookupIP(host string) ([]net.IP, error) {
	dl.RLock()
	defer dl.RUnlock()

	if IsIPString(host) {
		ip := net.ParseIP(host)
		if ip != nil {
			return []net.IP{ip}, nil
		}
	}

	if strings.ToLower(host) == "localhost" {
		return []net.IP{net.ParseIP("127.0.0.1")}, nil
	}

	if strings.ToLower(host) == "::1" {
		return []net.IP{net.ParseIP("::1")}, nil
	}

	if len(dl.customDNSServers) == 0 {
		if dl.useSysDNSIfNoCustom {
			return net.LookupIP(host)
		} else {
			return nil, gerrors.New("no DNS servers")
		}
	}

	// '.' means root in domain standard.
	if !gstring.EndWith(host, ".") {
		host += "."
	}
	req := dns.Msg{}
	req.Id = dns.Id()
	req.RecursionDesired = true
	req.Question = []dns.Question{
		{Name: host, Qtype: dns.TypeA, Qclass: dns.ClassINET},
	}

	var u upstream.Upstream = nil
	for _, srv := range dl.customDNSServers {
		u = srv.u
	}
	reply, err := u.Exchange(&req)
	if err != nil {
		return nil, gerrors.New("Cannot make the DNS request: %s", err.Error())
	}
	var ipMap = make(map[string]net.IP)
	for _, v := range reply.Answer {
		if ginterface.Type(v) != ginterface.Type(&dns.A{}) {
			continue
		}
		answer := v.(*dns.A)
		if answer.A != nil {
			ipMap[answer.A.String()] = answer.A
		}
	}

	var result []net.IP
	for _, v := range ipMap {
		result = append(result, v)
	}
	return result, nil
}

// LookupAddr looks up host names or domains by ip address from DNS server,
// if no DNS server is set up locally, then a query request will be sent to
// the default gateway.
func (dl *DNSClient) LookupAddr(addr string) (names []string, err error) {
	dl.RLock()
	defer dl.RUnlock()

	if len(dl.customDNSServers) == 0 {
		return net.LookupAddr(addr)
	}

	// TODO: to finish it
	m := new(dns.Msg)
	m.SetQuestion(addr, dns.TypeA)
	return nil, nil
}
