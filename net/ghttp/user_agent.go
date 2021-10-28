package ghttp

import (
	"github.com/avct/uasurfer" // 准确性很高
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/container/gnum"
	"github.com/mssola/user_agent" // 实测准确性较差
	"strings"
)

type UserAgnetInfo struct {
	Platform       string
	OsName         string
	OsVersion      string
	DeviceType     string
	IsMobile       bool
	IsBot          bool
	EngineName     string
	EngineVersion  string
	BrowserName    string
	BrowserVersion string
}

func verToString(ver uasurfer.Version) (result string) {
	if ver.Major >= 0 {
		result += gnum.ToString(ver.Major)
		if ver.Minor >= 0 {
			result += "." + gnum.ToString(ver.Minor)
			if ver.Patch >= 0 {
				result += "." + gnum.ToString(ver.Patch)
			}
		}
	}
	return result
}

func ParseUserAgent(uaString string) (*UserAgnetInfo, error) {
	if len(uaString) == 0 {
		return nil, gerrors.New("Nil uaString")
	}
	var uai UserAgnetInfo
	ua := user_agent.New(uaString)
	ua2 := uasurfer.Parse(uaString)

	uai.Platform = ua2.OS.Platform.String()
	uai.Platform = strings.Replace(uai.Platform, "Platform", "", 1) // "PaltformMac" -> "Mac"
	uai.OsName = ua2.OS.Name.String()
	uai.OsName = strings.Replace(uai.OsName, "OS", "", 1) // "OSMacOSX" -> "MacOSX"
	uai.OsVersion = verToString(ua2.OS.Version)
	uai.DeviceType = ua2.DeviceType.String()
	uai.DeviceType = strings.Replace(uai.DeviceType, "Device", "", 1) // "DeviceComputer" -> "Computer"
	uai.IsMobile = uai.DeviceType == "Phone" || uai.DeviceType == "Tablet" || uai.DeviceType == "Wearable"
	uai.IsBot = ua.Bot() // uasurfer库de代码中出现了Bot，但README中没有提及，所以还是用的ua.Bot()接口
	uai.EngineName, uai.EngineVersion = ua.Engine()
	uai.BrowserName = ua2.Browser.Name.String()
	uai.BrowserName = strings.Replace(uai.BrowserName, "Browser", "", 1) // "BrowserChrome" -> "Chrome"
	uai.BrowserVersion = verToString(ua2.Browser.Version)

	return &uai, nil
}
