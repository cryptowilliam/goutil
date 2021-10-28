package gaddr

import (
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/container/gstring"
	"github.com/cryptowilliam/goutil/net/ghtml"
	"github.com/cryptowilliam/goutil/net/ghttp"
	"github.com/cryptowilliam/goutil/net/gprobe/xonline"
	"net"
	"strings"
	"time"
)

// get my wan IPs by 3rd party service
func GetWanIpOL(proxy string) (net.IP, error) {
	var ipstr string
	var exist bool
	var firstCheckWanOnline = true

	// best choices: plain text of ip string
	endpoints := []string{
		"http://whatismyip.akamai.com",
		"http://ident.me",
		"http://myip.dnsomatic.com",
		"http://icanhazip.com",
		"http://ifconfig.co/ip"}
	eps := gstring.Shuffle(endpoints)
	for _, url := range eps {
		resp, err := ghttp.Get(url, proxy, time.Second*3, true)
		if err != nil {
			if firstCheckWanOnline {
				if !xonline.IsWanOnline(proxy) {
					return nil, gerrors.New("Can't get WAN ip because of internet offline ")
				}
				firstCheckWanOnline = false
			}
			continue
		}
		ipstr, _ = ghttp.ReadBodyString(resp)
		resp.Body.Close()
		ipstr = strings.Trim(ipstr, "\r") // icanhazip.com 的返回结果会带换行符
		ipstr = strings.Trim(ipstr, "\n")
		t := CheckIPString(ipstr)
		if t == IPv4_WAN {
			return ParseIP(ipstr)
		}
	}

	// backup choices
	htmlString, err := ghttp.GetString("http://bot.whatismyipaddress.com", proxy, time.Second*5)
	if err != nil {
		return nil, err
	}
	doc, err := ghtml.NewDocFromHtmlSrc(&htmlString)
	if err == nil {
		ipstr = doc.Text()
		t := CheckIPString(ipstr)
		if t == IPv4_WAN {
			return ParseIP(ipstr)
		}
	}
	htmlString, err = ghttp.GetString("http://network-tools.com", proxy, time.Second*5)
	if err != nil {
		return nil, err
	}
	doc, err = ghtml.NewDocFromHtmlSrc(&htmlString)
	if err == nil {
		ipstr, exist = doc.Find("#field").First().Attr("value")
		if exist {
			t := CheckIPString(ipstr)
			if t == IPv4_WAN {
				return ParseIP(ipstr)
			}
		}
	}

	if !xonline.IsWanOnline(proxy) {
		return nil, gerrors.New("Can't get WAN ip because of internet offline ")
	} else {
		return nil, gerrors.New("Can't get WAN ip, unknown error")
	}
}
