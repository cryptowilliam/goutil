package gaddr

// http://www.baidu.com/news
// socks5://username[:password]@host:port

import (
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/container/gstring"
	"github.com/cryptowilliam/goutil/encoding/gmultimedia"
	"github.com/goware/urlx"
	"net/url"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"v2ray.com/core/common/net"
)

// scheme:[//[user:password@]host[:port]][/]path[?query][#fragment]

const (
	/* public schemes */
	SchemeUnknown = iota
	SchemeHttp
	SchemeHttps
	SchemeFtp
	SchemeFtps
	SchemeMailto
	SchemeFile
	SchemeIdap
	SchemeNews
	SchemeGopher
	SchemeTelnet
	SchemeWais
	SchemeNntp
	SchemeData
	SchemeIrc
	SchemeIrcs
	SchemeWorldwind
	SchemeMms
	SchemeSocks4
	SchemeSocks4a
	SchemeSocks5
	SchemeSocks5s
	SchemeSocksHttp
	SchemeSocksHttps
	SchemeShadowsocks
	/* custom schemes */
	SchemeSvn
	SchemeHg
	SchemeGit
	SchemeThunder
	SchemeTencent
	SchemeEd2k
	SchemeMagnet
	SchemeTwitter
)

type Scheme int

func IsFilePath(s string) bool {
	if runtime.GOOS == "windows" {
		vn := filepath.VolumeName(s)
		return len(vn) > 0
	} else {
		return gstring.StartWith(s, "/")
	}
}

func IsUrl(s string) bool {
	_, err := url.Parse(s)
	return err == nil
}

// "http://bing.com/" is domain url, "http://bing.com/search" is not domain url
func IsAndOnlyDomain(str string) bool {
	u, err := url.Parse(str)
	if err != nil {
		return false
	}

	return (len(u.Path) == 0 || u.Path == "/") && len(u.RawQuery) == 0
}

// Combine absolute path and relative path to get a new absolute path
// If relUrl is absolute url, returns this relUrl
func Join(baseUrl string, relUrl string) (absUrl string, err error) {
	if len(baseUrl) == 0 || len(relUrl) == 0 {
		return "", gerrors.New("UrlJoin get invalid parameters")
	}
	base, err := url.Parse(baseUrl)
	if err != nil {
		return "", gerrors.New("baseUrl parse error: " + baseUrl)
	}
	if !base.IsAbs() {
		return "", gerrors.New("baseUrl is not absolute url: " + baseUrl)
	}
	rel, err := url.Parse(relUrl)
	if err != nil {
		return "", gerrors.New("relUrl parse error: " + relUrl)
	}
	return base.ResolveReference(rel).String(), nil
}

type UrlHost struct {
	Domain string // like "163.com"
	IP     string // like "8.8.8.8"
	Port   int // like "443"
}

type UrlAuth struct {
	User        string
	Password    string
	PasswordSet bool
}

type Path struct {
	Str    string
	Dirs   []string
	Params map[string][]string
}

func (ua *UrlAuth) String() string {
	res := ""
	if len(ua.User) > 0 {
		res += ua.User
	}
	if len(ua.Password) > 0 {
		res += ":" + ua.Password
	}
	return res
}

type AddrSlice struct {
	Scheme string // like "http", "ftp"
	Domain Domain // like "google.com"
	Auth   UrlAuth // like "usr:pwd"
	Host   UrlHost // like "google.com:443"
	Path   Path // like "?article=1260&lang=en#comment"
}

// Addr returns IP address or domain, without port.
func (uh *UrlHost) Addr() string {
	if uh.Domain != "" {
		return uh.Domain
	}
	return uh.IP
}

func (uh *UrlHost) String() string {
	var result string

	if len(uh.Domain) > 0 {
		result += uh.Domain
	} else {
		result += uh.IP
	}

	if IsValidPort(uh.Port) {
		result += ":" + strconv.FormatInt(int64(uh.Port), 10)
	}
	return result
}

func (us *AddrSlice) String() string {
	res := ""
	if len(us.Scheme) > 0 {
		res += us.Scheme + "://"
	}
	if len(us.Auth.String()) > 0 {
		res += us.Auth.String() + "@"
	}
	res += us.Host.String()
	if len(us.Path.Str) > 0 {
		res += "/" + us.Path.Str
	}
	return res
}

/*
type UrlChunk struct {
	Scheme   string
	User     string
	Password string
	Host     string
	Port     int64
	Path     string
	Param    string
}

// mongodb://user:password@127.0.0.1:27717
// ftp://user:password@127.0.0.1:1999/files/abc.mp4
func ParseUrlString(u string) (*UrlChunk, error) {
	r := UrlChunk{}
	//defErr := gerrors.Errorf("invalid url(%s)", u)

	// parse scheme
	uphppp := u // user password host port path param
	if ss := strings.Split(u, "://"); len(ss) == 2 {
		r.Scheme = ss[0]
		uphppp = ss[1]
	}

	// parse user password
	hppp := uphppp // host port path param
	if ss := strings.Split(uphppp, "@"); len(ss) == 2 {
		hppp = ss[1]
		usr_pwd := ss[0]
		ss := strings.Split(usr_pwd, ":")
		if len(ss) == 1 {
			if gstring.StartWith(usr_pwd, ":") {
				r.Password = usr_pwd
			}
			if gstring.EndWith(usr_pwd, ":") {
				r.User = usr_pwd
			}
		}
		if len(ss) == 2 {
			r.User = ss[0]
			r.Password = ss[1]
		}
	}

	// split host port & path param
	hp := hppp // host port
	pp := ""   // path param
	if ss := strings.Split(hppp, "/"); len(ss) == 2 {
		hp = ss[0]
		pp = ss[1]
	}

	// parse host, port
	r.Host = hp
	if ss := strings.Split(hp, ":"); len(ss) == 2 {
		r.Host = ss[0]
		r.Port, _ = strconv.ParseInt(ss[1], 10, 32)
	}

	// parse path, param
	if len(pp) > 0 {
		r.Path = pp
		if ss := strings.Split(pp, "?"); len(ss) == 2 {
			r.Path = ss[0]
			r.Param = ss[1]
		}
	}
	return &r, nil
}*/

// NOTICE
// url.Parse("192.168.1.1:80") reports error because RFC3986 says "192.168.1.1:80" is an invalid url, the correct way is "//192.168.1.1:80".
// In gaddr library and urlx library, "192.168.1.1:80" is a valid url because it is used a lot
// Reference: https://github.com/golang/go/issues/19297
func ParseUrl(urlStr string) (*AddrSlice, error) {
	u, err := urlx.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	var s AddrSlice

	// if there is no scheme in input url string, urlx.Parse will give default scheme "http://"
	// so I must check urlx.Parse return scheme
	if gstring.StartWith(strings.ToLower(urlStr), strings.ToLower(u.Scheme)+"://") {
		s.Scheme = u.Scheme
	}
	if u.User != nil {
		s.Auth.User = u.User.Username()
		s.Auth.Password, s.Auth.PasswordSet = u.User.Password()
	}
	if strings.Contains(u.Host, ":") {
		if host, portstr, err := net.SplitHostPort(u.Host); err != nil {
			return nil, err
		} else {
			if portstr != "" {
				if port, err := strconv.Atoi(portstr); err != nil {
					return nil, err
				} else {
					s.Host.Port = port
				}
			}
			if IsIPString(host) {
				s.Host.IP = host
			} else {
				s.Host.Domain = host
			}
		}
	} else {
		s.Host.Domain = u.Host
	}
	if s.Host.Domain != "" {
		domain, err := ParseDomain(s.Host.Domain)
		if err != nil {
			return nil, err
		}
		s.Domain = *domain
	}

	s.Path.Str = u.Path
	dirs := strings.Split(u.Path, "/")
	for _, v := range dirs {
		if v == "" {
			continue
		}
		s.Path.Dirs = append(s.Path.Dirs, v)
	}
	s.Path.Params = u.Query()

	return &s, nil
}

func IsImageUrl(url string) bool {
	_, err := ParseUrl(url)
	if err != nil {
		return false
	}
	url = strings.ToLower(url)
	for _, v := range gmultimedia.SuffixsOfImage {
		if gstring.EndWith(url, v) {
			return true
		}
	}
	return false
}

func IsVideoUrl(url string) bool {
	_, err := ParseUrl(url)
	if err != nil {
		return false
	}
	url = strings.ToLower(url)
	for _, v := range gmultimedia.SuffixsOfVideo {
		if gstring.EndWith(url, v) {
			return true
		}
	}
	return false
}

func IsAudioUrl(url string) bool {
	_, err := ParseUrl(url)
	if err != nil {
		return false
	}
	url = strings.ToLower(url)
	for _, v := range gmultimedia.SuffixsOfAudio {
		if gstring.EndWith(url, v) {
			return true
		}
	}
	return false
}

func LastPath(urlstr string) string {
	u, err := url.Parse(urlstr)
	if err != nil {
		return ""
	}

	idx := strings.LastIndex(u.Path, "/")
	if idx <= 0 || idx == (len(u.Path)-1) {
		return ""
	}
	return u.Path[idx+1:]
}

func RemoveDuplicateUrl(urls []string) []string {
	return gstring.RemoveDuplicate(urls)
}

// https://video-icn1-1.xx.fbcdn.net/v/t42.9040-2/58467180_2666273813399564_6679224605468524544_n.mp4?_nc_cat=100\u0026efg=eyJybHIiOjY5NCwicmxhIjo1MTIsInZlbmNvZGVfdGFnIjoic3ZlX3NkIn0=\u0026rl=694\u0026vabr=386\u0026_nc_ht=video-icn1-1.xx\u0026oh=881ead117c700970945a89716b3a0b54\u0026oe=5CB9BAA5
func DecodeURL(encodedUrl string) (string, error) {
	encodedUrl = strings.Replace(encodedUrl, "&amp;", "&", -1)

	return url.QueryUnescape(encodedUrl)
}

type AddrParser struct {
	addressString string
	as            *AddrSlice
	err           error
}

func NewParser(address_string string) *AddrParser {
	r := &AddrParser{}
	r.addressString = address_string
	r.err = nil
	r.as, r.err = ParseUrl(r.addressString)
	return r
}

func (ap *AddrParser) Verify() error {
	return ap.err
}

func (ap *AddrParser) Split() (*AddrSlice, error) {
	return ap.as, ap.err
}

func (ap *AddrParser) Scheme() string {
	if ap.err != nil {
		return ""
	}
	return ap.as.Scheme
}
