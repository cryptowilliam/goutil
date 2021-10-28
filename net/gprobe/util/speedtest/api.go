package speedtest

import (
	"github.com/cryptowilliam/goutil/container/gspeed"
)

func TestDownloadSpeed() (gspeed.Speed, error) {
	// Get and select test server(s)
	u, err := fetchUserInfo()
	if err != nil {
		return gspeed.Speed(0), err
	}
	allServers, err := getAllTestServers(u) // return list has already been sorted by distance
	if err != nil {
		return gspeed.Speed(0), err
	}
	selectServers, err := allServers.selectNearestTestServers(1) // return nearest test server
	if err != nil {
		return gspeed.Speed(0), err
	}

	// Test all selected servers and calculate average download speed
	avg := 0.0
	for _, s := range selectServers.Servers {
		dlSpeed, err := s.testDownload()
		if err != nil {
			return gspeed.Speed(0), err
		}
		avg += dlSpeed
	}
	avg = avg / float64(len(selectServers.Servers))
	return gspeed.FromBits(avg * float64(gspeed.Mb))
}

func TestUploadSpeed() (gspeed.Speed, error) {
	// Get and select test server(s)
	u, err := fetchUserInfo()
	if err != nil {
		return gspeed.Speed(0), err
	}
	allServers, err := getAllTestServers(u) // return list has already been sorted by distance
	if err != nil {
		return gspeed.Speed(0), err
	}
	selectServers, err := allServers.selectNearestTestServers(1) // return most nearest test server
	if err != nil {
		return gspeed.Speed(0), err
	}

	// Test all selected servers and calculate average upload speed
	avg := 0.0
	for _, s := range selectServers.Servers {
		ulSpeed, err := s.testUpload()
		if err != nil {
			return gspeed.Speed(0), err
		}
		avg += ulSpeed
	}
	avg = avg / float64(len(selectServers.Servers))
	return gspeed.FromBits(avg * float64(gspeed.Mb))
}
