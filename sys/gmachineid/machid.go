package gmachineid

import (
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/container/gstring"
	"github.com/cryptowilliam/goutil/crypto/ghash"
	"github.com/cryptowilliam/goutil/net/gnet"
	"os/exec"
	"runtime"
	"strings"
)

// Get hardware UUID of MacOS
func MacosHardwareUUID() (string, error) {
	if runtime.GOOS != "darwin" {
		return "", gerrors.New("MacosHardwareUUID does not support " + runtime.GOOS)
	}
	output, err := exec.Command("system_profiler", "SPHardwareDataType").CombinedOutput()
	if err != nil {
		return "", err
	}
	uuid, err := gstring.SubstrBetween(string(output), "Hardware UUID:", "\n", true, true, false, false)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(uuid), nil
}

func nonMacosPhysicalMACs() (string, error) {
	ns, err := gnet.GetAllNicNames()
	if err != nil {
		return "", err
	}

	var macs string

	for _, s := range ns {
		ni, _ := gnet.GetNicInfo(s)
		if !ni.IsPhysical {
			continue
		}
		macs += ni.MAC
	}

	return ghash.Md5Str(macs)
}

func Get() (string, error) {
	var str string
	var err error

	if runtime.GOOS == "darwin" {
		str, err = MacosHardwareUUID()
	} else if runtime.GOOS == "linux" || runtime.GOOS == "windows" {
		str, err = nonMacosPhysicalMACs()
	} else {
		return "", gerrors.New("Unsupported OS " + runtime.GOOS)
	}

	if err != nil {
		return "", err
	}
	md5, err := ghash.Md5Str(str + "salt-duck-machid")
	if err != nil {
		return "", err
	}
	return md5, nil
}
