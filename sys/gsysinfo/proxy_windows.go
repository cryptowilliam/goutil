package gsysinfo

import (
	"golang.org/x/sys/windows/registry"
)

const (
	keyProxyEnable = "ProxyEnable"
	keyProxyServer = "ProxyServer"
)

// reference: https://github.com/andreyvit/systemproxy/blob/master/sysproxy_windows.go

func GetGlobalSocks5Proxy() (string, bool, error) {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Internet Settings`, registry.QUERY_VALUE)
	if err != nil {
		return "", false, err
	}
	defer k.Close()

	en, _, err := k.GetIntegerValue(keyProxyEnable)
	if err != nil && err != registry.ErrNotExist {
		return "", false, err
	}
	enabled := (en != 0)

	defaultServer, _, err := k.GetStringValue(keyProxyServer)
	if err != nil && err != registry.ErrNotExist {
		return "", false, err
	}

	return defaultServer, enabled, nil
}

func SetGlobalSocks5Proxy(defaultServer string, enabled bool) error {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Internet Settings`, registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer k.Close()

	var en uint32
	if enabled {
		en = 1
	}
	err = k.SetDWordValue(keyProxyEnable, en)
	if err != nil {
		return err
	}

	err = k.SetStringValue(keyProxyServer, defaultServer)
	if err != nil {
		return err
	}

	return nil
}
