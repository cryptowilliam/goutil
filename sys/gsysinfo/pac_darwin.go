package gsysinfo

import (
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/container/gstring"
	"os/exec"
	"strings"
)

func GetPacProxy() (pacUrl string, enabled bool, err error) {
	itfc, err := GetCurrentNetworkInterface()
	if err != nil {
		return "", false, err
	}

	b, err := exec.Command("networksetup", "-getautoproxyurl", itfc).CombinedOutput()
	if err != nil {
		return "", false, gerrors.New(string(b) + err.Error())
	}
	ss := strings.Split(string(b), "\n")
	if len(ss) < 2 {
		return "", false, gerrors.Errorf("invalid return(%s)", string(b))
	}

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
		if gstring.StartWith(v, "URL: ") {
			pacUrl = strings.Replace(v, "URL: ", "", 1)
		}
	}

	return pacUrl, enabled, nil
}

// pac更新后怎样刷新： https://www.zhihu.com/question/19947389
func SetPacProxyOn(pacUrl string) error {
	itfc, err := GetCurrentNetworkInterface()
	if err != nil {
		return err
	}

	b, err := exec.Command("networksetup", "-setautoproxyurl", itfc, pacUrl).CombinedOutput()
	if err != nil {
		return gerrors.New(string(b) + err.Error())
	}

	b, err = exec.Command("networksetup", "-setautoproxystate", itfc, "on").CombinedOutput()
	if err != nil {
		return gerrors.New(string(b) + err.Error())
	}

	return nil
}

func SetPacProxyOff() error {
	itfc, err := GetCurrentNetworkInterface()
	if err != nil {
		return err
	}

	b, err := exec.Command("networksetup", "-setautoproxystate", itfc, "off").CombinedOutput()
	if err != nil {
		return gerrors.New(string(b) + err.Error())
	}

	return nil
}
