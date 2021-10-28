package gsysinfo

/*
  Get system version information
*/

import (
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/container/gstring"
	"github.com/cryptowilliam/goutil/sys/gsysinfo/util"
	"github.com/getlantern/osversion"
	"runtime"
	"strconv"
	"strings"
)

type SysVer struct {
	PlatformName    string // windows / linux / darwin
	PlatformVer     string // Under linux this is kernel version
	Arch            string // i386 / AMD64 / ARM / MIPS
	LinuxDistroName string // "CentOS" / "Debian"
	LinuxDistroVer  string
}

func init() {
	ForbiddenCentOs5x()
}

func (sv SysVer) String() string {
	return "Platform:" + sv.PlatformName + "\n" + "PlatformVer:" + sv.PlatformVer + "\n" + "Arch:" + sv.Arch + "\n" + "LinuxDistroName:" + sv.LinuxDistroName + "\n" + "LinuxDistroVer:" + sv.LinuxDistroVer + "\n"
}

func Get() (*SysVer, error) {
	var sv SysVer
	var err error

	sv.PlatformName = runtime.GOOS
	sv.PlatformVer, err = osversion.GetHumanReadable()
	if err != nil {
		return nil, err
	}
	if runtime.GOOS == "linux" {
		cutPos := strings.Index(sv.PlatformVer, "kernel: ")
		if cutPos >= 0 {
			// "kernel: 2.6.32-431.el6.x86_64" -> "2.6.32-431.el6.x86_64"
			// "Ubuntu 16.04.1 LTS kernel: 4.4.0-57-generic" -> "4.4.0-57-generic"
			tmp, err := gstring.SubstrAscii(sv.PlatformVer, cutPos+8, len(sv.PlatformVer)-1)
			if err == nil {
				sv.PlatformVer = tmp
			}
		}
		// sv.PlatformVer = strings.Trim(sv.PlatformVer, "kernel: ")
	} else if runtime.GOOS == "darwin" {
		sv.PlatformVer = strings.Trim(sv.PlatformVer, "OS X ") // "OS X 10.11.0 El Capitan" -> "10.11.0 El Capitan"
	} else if runtime.GOOS == "windows" {
		sv.PlatformVer = strings.Trim(sv.PlatformVer, "Windows ")
	} else {
		return nil, gerrors.New(runtime.GOOS + " unsupported")
	}
	sv.Arch = runtime.GOARCH

	if runtime.GOOS == "linux" {

		do := util.LinuxDistribution(nil)
		if do == nil {
			return nil, gerrors.New("LinuxDistribution error")
		}
		sv.LinuxDistroName = do.Name(true)
		sv.LinuxDistroName = strings.TrimSpace(sv.LinuxDistroName)
		sv.LinuxDistroName = strings.ToLower(sv.LinuxDistroName)

		// "CentOS 6.5" -> "CentOS"
		spacePos := strings.Index(sv.LinuxDistroName, " ")
		if spacePos > 0 {
			tmp, err := gstring.SubstrAscii(sv.LinuxDistroName, 0, spacePos)
			if err == nil {
				sv.LinuxDistroName = tmp
			}
		}

		// Config pretty:
		// Example in ubuntu, pretty = true, return "16.04", pretty = false, return "16.04 (Xenial Xerus)
		sv.LinuxDistroVer = do.Version(true, true)
		sv.LinuxDistroVer = strings.TrimSpace(sv.LinuxDistroVer)
		sv.LinuxDistroVer = strings.ToLower(sv.LinuxDistroVer)
	}

	return &sv, nil
}

// Golang does not support CentOS 5.x
func ForbiddenCentOs5x() error {
	if runtime.GOOS != "linux" {
		return gerrors.Errorf("linux supported only")
	}
	sv, err := Get()
	if err != nil {
		return err
	}
	if sv.LinuxDistroName == "centos" && len(sv.LinuxDistroVer) > 0 {
		majorVerString, err := gstring.SubstrAscii(sv.LinuxDistroVer, 0, 1)
		if err != nil {
			return err
		}
		majorVerNum, err := strconv.ParseInt(majorVerString, 10, 32)
		if err != nil {
			return err
		}
		if majorVerNum > 1 && majorVerNum <= 5 {
			return gerrors.New(sv.LinuxDistroName + " " + sv.LinuxDistroVer + " unsupported")
		}
	}
	return nil
}
