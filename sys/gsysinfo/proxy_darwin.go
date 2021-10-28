package gsysinfo

import (
	"fmt"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/container/gstring"
	"os/exec"
	"strings"
)

// from: github.com/txthinking/brook/sysproxy
// GetNetworkInterface returns default interface name, not dev name.
func GetCurrentNetworkInterface() (string, error) {
	c := exec.Command("sh", "-c", "networksetup -listnetworkserviceorder | grep $(route -n get default | grep interface | awk '{print $2}') | awk 'BEGIN {FS=\",\"}; {print $1}' | awk 'BEGIN {FS=\": \"}; {print $2}'")
	out, err := c.CombinedOutput()
	if err != nil {
		return "", gerrors.New(string(out) + err.Error())
	}
	return strings.TrimSpace(string(out)), nil
}

func GetAllNetworkInterfaces() ([]string, error) {
	b, err := exec.Command("networksetup", "-listallnetworkservices").CombinedOutput()
	if err != nil {
		return nil, err
	}
	ss := strings.Split(string(b), "\n")
	if len(ss) > 0 && strings.Contains(ss[0], "*") {
		ss = ss[1:]
	}
	return ss, nil
}

func GetGlobalSocks5Proxy() (server string, enabled bool, err error) {
	itfc, err := GetCurrentNetworkInterface()
	if err != nil {
		return "", false, err
	}
	b, err := exec.Command("networksetup", "-getsocksfirewallproxy", itfc).CombinedOutput()
	if err != nil {
		return "", false, err
	}
	ss := strings.Split(string(b), "\n")
	if len(ss) < 4 {
		return "", false, gerrors.Errorf("invalid return(%s)", string(b))
	}

	svr := ""
	port := ""
	for _, v := range ss {
		if gstring.StartWith(v, "Enabled: ") {
			v = strings.Replace(v, "Enabled: ", "", 1)
			v = strings.ToLower(v)
			if v == "yes" {
				enabled = true
			} else if v == "no" {
				enabled = false
			} else {
				return "", false, gerrors.Errorf("invalid enabled flag(%s)", v)
			}
		}
		if gstring.StartWith(v, "Server: ") {
			svr = strings.Replace(v, "Server: ", "", 1)
		}
		if gstring.StartWith(v, "Port: ") {
			port = strings.Replace(v, "Port: ", "", 1)
		}
	}
	if enabled {
		server = fmt.Sprintf("%s:%s", svr, port)
	}
	return server, enabled, nil
}

func SetGlobalSocks5ProxyOn(socks5Proxy string) error {
	itfc, err := GetCurrentNetworkInterface()
	if err != nil {
		return err
	}

	ss := strings.Split(socks5Proxy, ":")
	if len(ss) != 2 {
		return gerrors.Errorf("invalid socks5 proxy(%s)", socks5Proxy)
	}

	// networksetup -setsocksfirewallproxy Wi-Fi 127.0.0.1 1088
	b, err := exec.Command("networksetup", "-setsocksfirewallproxy", itfc, ss[0], ss[1]).CombinedOutput()
	if err != nil {
		return gerrors.New(string(b) + err.Error())
	}

	b, err = exec.Command("networksetup", "-setsocksfirewallproxystate", itfc, "on").CombinedOutput()
	if err != nil {
		return gerrors.New(string(b) + err.Error())
	}

	return nil
}

func SetGlobalSocks5ProxyOff() error {
	itfc, err := GetCurrentNetworkInterface()
	if err != nil {
		return err
	}

	b, err := exec.Command("networksetup", "-setsocksfirewallproxystate", itfc, "off").CombinedOutput()
	if err != nil {
		return gerrors.New(string(b) + err.Error())
	}

	return nil
}
