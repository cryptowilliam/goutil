package gsysinfo

// https://github.com/okelet/goutils/blob/4585eaae84d2dc05f02ea33fb585b78de299f041/proxy_utils.go
// https://stackoverflow.com/questions/36127254/change-proxy-on-mac-osx-programmatically
// https://developer.apple.com/documentation/networkextension/neproxysettings/1406766-proxyautoconfigurationjavascript
// https://github.com/universonic/shadowsocks-macos/blob/0d7d213b22864ca9712e23918341288efc072f22/Shadowsocks/ProxyConfHelper.m

/*
there are 3 method to configure system proxy
1: OS env
2: OS Pac proxy
3: OS global proxy

method 1:
On Ubuntu 16.04 LTS with Chrome v53 (64 bit), I had to set the http_proxy / HTTP_PROXY env variables to "http://proxyserver:port" for all users for Chrome to be able to communicate.
Modify /etc/profile
export {http,ftp,https,rsync}_proxy="socks5://proxyserver:port"
export {HTTP,FTP,HTTPS,RSYNC}_PROXY=$http_proxy
but under macOS, this seems works for curl/wget, doesn't works for browsers, some articles said "chrome --proxy-server=***" works
anyway, this is not a solid solution

method 2 and 3:
configure it in network configuration GUI window, or use networksetup command under macOS, use RegisterKey editor under windows
and, pac proxy operations using system api (but not command) works too, sample project https://github.com/getlantern/pac-cmd
so I think api to configure global proxy is possible too
*/
