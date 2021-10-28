package gsysinfo

import "github.com/getlantern/pac"

func SetPacProxyOn(pacUrl string) error {
	return pac.On(pacUrl)
}

func SetPacProxyOff() error {
	return pac.Off("old-pac-url-prefix")
}
