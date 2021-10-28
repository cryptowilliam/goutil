package ghttp

import (
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	resp, err := Get("http://whatismyip.akamai.com", "socks5://127.0.0.1:1086", time.Second*10, true)
	if err != nil {
		t.Error(err)
		return
	}
	ipstr, err := ReadBodyString(resp)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(ipstr)
}

func TestSetProxy(t *testing.T) {
	client := http.DefaultClient

	if SetProxy(client, "socks://127.0.0.1:1086") == nil {
		t.Errorf("socks://*** is a invalid HTTP proxy scheme")
		return
	}

	resp, err := client.Get("http://whatismyip.akamai.com")
	if err != nil {
		t.Error(err)
		return
	}
	original_ipstr, _ := ReadBodyString(resp)
	resp.Body.Close()
	original_ipstr = strings.Trim(original_ipstr, "\r") // icanhazip.com 的返回结果会带换行符
	original_ipstr = strings.Trim(original_ipstr, "\n")
	t.Log(original_ipstr)

	if err := SetProxy(client, "socks5://127.0.0.1:1086"); err != nil {
		t.Error(err)
		return
	}
	resp, err = client.Get("http://whatismyip.akamai.com")
	if err != nil {
		t.Error(err)
		return
	}
	proxyed_ipstr, _ := ReadBodyString(resp)
	resp.Body.Close()
	proxyed_ipstr = strings.Trim(proxyed_ipstr, "\r") // icanhazip.com 的返回结果会带换行符
	proxyed_ipstr = strings.Trim(proxyed_ipstr, "\n")
	t.Log(proxyed_ipstr)

	if original_ipstr == proxyed_ipstr {
		t.Error("SetProxy doesn't work")
	}
}

func TestGetBigFile(t *testing.T) {
	_, err := GetBigFile("http://smartmesh.io/SmartMeshWhitePaperEN.pdf", "test.pdf")
	if err != nil {
		t.Error(err)
		return
	}
}
