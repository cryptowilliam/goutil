package gaddr

//github.com/joeguo/tldextract
//github.com/weppos/publicsuffix-go
//判断是不是公共域名后缀
//判断两个域名是不是同一个所有人，比如news.baidu和www.baidu.com就是同一个所有者

import (
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/container/gstring"
	"github.com/cryptowilliam/goutil/sys/gfs"
	"github.com/cryptowilliam/goutil/sys/gsysinfo"
	"github.com/domainr/whois"
	"github.com/globalsign/publicsuffix"
	"github.com/joeguo/tldextract"
	"github.com/liamcurry/domains"
	gowhois "github.com/likexian/whois"
	"os"
	"strings"
	"time"
)

type Domain struct {
	TLD        string // "com" | "com.cn"
	SLD_ROOT   string // "baidu"
	TRD_SUB    string // "www"
	SiteDomain string // "baidu.com"
}

var extractor *tldextract.TLDExtract = nil

// 当缓存文件不存在，或者修改时间在24小时以前，则重新下载缓存文件
func updateLtdListAndExtractor() error {
	var err error
	home, err := gsysinfo.GetHomeDir()
	if err != nil {
		return err
	}
	cacheFilename := home + "/.ltd.suffix.cache.dat"
	pi, err := gfs.GetPathInfo(cacheFilename)
	if err != nil {
		extractor = nil
		return err
	}
	if !pi.Exist {
		extractor = nil
	}
	if pi.Exist && time.Now().Sub(pi.ModifiedTime).Hours() > 24 {
		os.Remove(cacheFilename)
		extractor = nil
	}
	if extractor == nil {
		extractor, err = tldextract.New(cacheFilename, false)
		if err != nil {
			return err
		}
	}
	return nil
}

// NOTICE
// 优点: 从权威网站下载TLD列表，判断结果准确
// 缺点: 初始化或者更新时必须在线工作，下载期间接口响应慢
func ParseONLINE(domain string) (*Domain, error) {
	var result Domain

	if err := updateLtdListAndExtractor(); err != nil {
		return nil, err
	}

	if extractor != nil {
		data := extractor.Extract(domain)
		if data != nil && len(data.Tld) > 0 && len(data.Root) > 0 {
			result.TLD = data.Tld
			result.SLD_ROOT = data.Root
			result.TRD_SUB = data.Sub
			result.SiteDomain = result.SLD_ROOT + "." + result.TLD
			return &result, nil
		} else {
			return nil, gerrors.New("Illegal domain")
		}
	}
	return nil, gerrors.New("Nil extractor")
}

// NOTICE
// This an offline domain parse function, please update source repo often
func ParseDomain(domain string) (*Domain, error) {
	result := Domain{}

	ret := false
	result.TLD, ret = publicsuffix.PublicSuffix(domain)
	if !ret || result.TLD == "" {
		return nil, gerrors.Errorf("%s is not a valid domain", domain)
	}
	s := gstring.RemoveTail(domain, len(result.TLD))
	if len(s) > 0 && s[len(s)-1] == '.' {
		s = gstring.RemoveTail(s, 1)
	}
	if len(s) == 0 {
		return &result, nil
	}
	ss := strings.Split(s, ".")
	if len(ss) > 0 {
		result.SLD_ROOT = ss[len(ss)-1]
		result.SiteDomain = result.SLD_ROOT + "." + result.TLD
		ss = ss[:len(ss)-1]
		if len(ss) > 0 {
			result.TRD_SUB = strings.Join(ss, ".")
		}
	}
	return &result, nil
}

func IsDomainONLINE(domain string) bool {
	_, err := ParseONLINE(domain)
	if err == nil {
		return true
	} else {
		return false
	}
}

func IsDomain(domain string) bool {
	_, err := ParseDomain(domain)
	if err == nil {
		return true
	} else {
		return false
	}
}

// Cloned from github.com/domainr/whois
// Whois response represents a whois response from a server.
type Whois struct {
	// Query and Host are copied from the Request.
	// Query string
	Host string

	// FetchedAt is the date and time the response was fetched from the server.
	FetchedAt time.Time

	// MediaType and Charset hold the MIME-type and character set of the response body.
	//MediaType string
	//Charset   string

	// Body contains the raw bytes of the network response (minus HTTP headers).
	//Body []byte
	Body string
}

// Principle: WHOIS information of domains which are not taken include "No match".
func IsRegistrable(domain string) bool {
	c := domains.NewChecker()
	return !c.IsTaken(domain)
}

func GetWhoisWithDomain(domain string) (*Whois, error) {
	request, err := whois.NewRequest(domain)
	if err != nil {
		return nil, err
	}
	response, err := whois.DefaultClient.Fetch(request)
	if err != nil {
		return nil, err
	}
	w := Whois{
		//Query:     response.Query,
		Host:      response.Host,
		FetchedAt: response.FetchedAt,
		//MediaType: response.MediaType,
		//Charset:   response.Charset,
		Body: string(response.Body),
	}
	return &w, nil
}

func GetWhoisWithIP(ip string) (*Whois, error) {
	result, err := gowhois.Whois(ip)
	if err != nil {
		return nil, err
	}
	w := Whois{
		Host:      ip,
		FetchedAt: time.Now(),
		Body:      result,
	}
	return &w, nil
}
