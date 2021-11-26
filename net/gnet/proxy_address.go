package gnet

import (
	"strings"
)

/*
 Proxy address example
 "http://my-proxy.com:9090"
 "https://my-proxy.com:9090"
 "socks4://my-proxy.com:9090"
 "socks4a://my-proxy.com:9090"
 "socks5://my-proxy.com:9090"
 "ss://my-proxy.com:9090"
*/

func ParseProxyAddr(address string) (proxyType, host string, err error) {
	us, err := ParseUrl(address)
	us.Scheme = strings.ToLower(us.Scheme)

	return us.Scheme, us.Host.String(), nil
}
