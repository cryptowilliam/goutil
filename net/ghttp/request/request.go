package request

// Detect as much information as possible from http client request

import (
	"github.com/cryptowilliam/goutil/net/ghttp"
	"github.com/cryptowilliam/goutil/net/gnet"
	"github.com/valyala/fasthttp"
	"strings"
)

type Kv struct {
	Key   string
	Value string
}

type RequestInfo struct {
	Host            string
	TlsUsed         bool
	Method          string
	Path            string // URI = Path + QueryArgs
	QueryArgs       string
	Headers         []Kv
	ClientIP        string
	ProxyDetected   bool
	ProxyIP         string
	UserAgentDetect *ghttp.UserAgnetInfo
}

/*
使用了非高匿代理的HTTP请求时，其Header中可能出现下述Key：
$ITEM
X-$ITEM
HTTP-X-$ITEM
上述所有可能中的'-'替换为'_'
*/
var proxyFlags = []string{
	"VIA",
	"FORWARDED",
	"FORWARDED-FOR",
	"FORWARDED-FOR-IP",
	"REAL-IP",
	"CLIENT-IP",
	"PROXY-ID",
	"PROXY-CONNECTION",
}

var realIpFlags = []string{
	"FORWARDED",
	"REAL-IP",
	"CLIENT-IP",
}

func getHeaders(r *fasthttp.RequestCtx) (result []Kv) {
	var item Kv

	r.Request.Header.VisitAll(func(k, v []byte) {
		item.Key = string(k)
		item.Value = string(v)
		result = append(result, item)

	})
	return result
}

func parseProxy(headers []Kv) (isViaProxy bool, clientRealIP string) {
	var k string

	// Check is via proxy
	isViaProxy = false
	for _, hd := range headers {
		for _, f := range proxyFlags {
			k = strings.ToLower(hd.Key)

			f = strings.ToLower(f)
			if k == f || k == "x-"+f || k == "http-x-"+f {
				isViaProxy = true
			}

			f = strings.Replace(f, "-", "_", -1)
			if k == f || k == "x_"+f || k == "http_x_"+f {
				isViaProxy = true
			}
		}
	}
	if !isViaProxy {
		return false, ""
	}

	// ParseIPString client real IP
	for _, hd := range headers {
		for _, f := range realIpFlags {
			k = strings.ToLower(hd.Key)

			f = strings.ToLower(f)
			if strings.Contains(k, f) {
				if gnet.IsIPString(hd.Value) {
					clientRealIP = hd.Value
					break
				}
			}

			f = strings.Replace(f, "-", "_", -1)
			if strings.Contains(k, f) {
				if gnet.IsIPString(hd.Value) {
					clientRealIP = hd.Value
					break
				}
			}
		}
		if len(clientRealIP) > 0 {
			break
		}
	}

	return isViaProxy, clientRealIP
}

func ParseRequest(ctx *fasthttp.RequestCtx) (*RequestInfo, error) {
	var ri RequestInfo
	var clientRealIp string

	ri.Host = string(ctx.Host())
	ri.TlsUsed = ctx.IsTLS()
	ri.Method = string(ctx.Method())
	ri.Headers = getHeaders(ctx)
	ri.Path = string(ctx.Path())
	ri.QueryArgs = ctx.QueryArgs().String()
	ri.ProxyDetected, clientRealIp = parseProxy(ri.Headers)
	if ri.ProxyDetected {
		if len(clientRealIp) > 0 {
			ri.ClientIP = clientRealIp
		} else {
			ri.ClientIP = "Undetected"
		}
	} else {
		ri.ClientIP = ctx.RemoteIP().String()
	}
	if ri.ProxyDetected {
		ri.ProxyIP = ctx.RemoteIP().String()
	} else {
		ri.ProxyIP = ""
	}

	ri.UserAgentDetect, _ = ghttp.ParseUserAgent(string(ctx.UserAgent()))

	return &ri, nil
}

/*
func handler(w http.ResponseWriter, r *http.Request) {
	var keys []string
	for k := range r.Header {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := strings.Join(r.Header[k], " ")
		fmt.Fprintf(w, "%s: %s\n", k, v)
	}
}
*/
