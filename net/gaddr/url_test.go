package gaddr_test

import (
	"fmt"
	"github.com/cryptowilliam/goutil/net/gaddr"
	"testing"
)

// "mailto:thomas@gmail.com"
// "https://www.google.com"
// "ftp://user:pwd@12.103.38.24/files"
// "http://my-proxy.com:9090"
// "https://my-proxy.com:9090"
// "socks4://my-proxy.com:9090"
// "socks4a://my-proxy.com:9090"
// "socks5://my-proxy.com:9090"
// "ss://my-proxy.com:9090"
// "mongodb://user:pwd@192.168.3.12/mydb"
// "ed2k://..."
// "magnet://..."
// "ss://method:password@ip:port"
// "jet://method:password@ip:port"
func TestParseUrl(t *testing.T) {
	_, err := gaddr.ParseUrl("socks://127.0.0.1:1086")
	if err != nil {
		t.Error(err)
		return
	}
	_, err = gaddr.ParseUrl("ss://method:password@163.com:1633")
	if err != nil {
		t.Error(err)
		return
	}
	us, err := gaddr.ParseUrl("ss://admin@network:password@me@13.209.69.159:9292")
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(us.Auth.User, us.Auth.Password)
}
