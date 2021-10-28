package gtime

import (
	"github.com/beevik/ntp"
	"github.com/cryptowilliam/goutil/container/gstring"
	"time"
)

// Get network time in UTC
func GetNetTimeInUTCONLINE() (time.Time, error) {
	ntpServers := []string{
		"ntp1.aliyun.com",
		"ntp2.aliyun.com",
		"ntp3.aliyun.com",
		"ntp4.aliyun.com",
		"ntp5.aliyun.com",
		"ntp6.aliyun.com",
		"ntp7.aliyun.com",
		"time1.cloud.tencent.com",
		"time2.cloud.tencent.com",
		"time3.cloud.tencent.com",
		"time4.cloud.tencent.com",
		"time5.cloud.tencent.com",
		"time1.apple.com",
		"time2.apple.com",
		"time3.apple.com",
		"time4.apple.com",
		"time5.apple.com",
		"time6.apple.com",
		"time7.apple.com",
		"time.apple.com",
		"3.asia.pool.ntp.org",
		"0.hk.pool.ntp.org",
		"0.jp.pool.ntp.org",
		"1.jp.pool.ntp.org",
		"2.jp.pool.ntp.org",
		"3.jp.pool.ntp.org"}
	var tm time.Time
	var err error
	ntpServers = gstring.Shuffle(ntpServers)
	for _, server := range ntpServers {
		tm, err = ntp.Time(server)
		if err == nil {
			return tm, nil
		}
	}
	return tm, err
}

// Get network time in local machine timezone
func GetNetTimeInLocalONLINE() (time.Time, error) {
	tm, err := GetNetTimeInUTCONLINE()
	if err == nil {
		return tm.In(time.Local), nil
	} else {
		return tm, err
	}
}

// Get network time and update time for local machine
// This API must run as root/admin， 但不知道Windows下是不是也是如此
// 有时警告有时不警告，内容为，sudo: timestamp too far in the future
func SyncNetTimeONLINEROOT() error {
	var tm time.Time
	var err error
	if tm, err = GetNetTimeInLocalONLINE(); err != nil {
		return err
	}

	return SetSystemTimeROOT(tm)
}
