package gaddr

import (
	"github.com/cryptowilliam/goutil/net/ghttp"
	"github.com/cryptowilliam/goutil/sys/gfs"
	"github.com/mohong122/ip2region/binding/golang/ip2region"
	"net"
	"os"
	"time"
)

// https://github.com/getlantern/geolookup 在线查询的，但是服务提供方需要翻墙

const (
	ip2RegionDbUrl = "https://github.com/lionsoul2014/ip2region/raw/master/data/ip2region.db"
)

type IpGeo struct {
	Country  string
	Region   string
	Province string
	City     string
	CityId   int64
	ISP      string
}

type GeoFinder struct {
	finder *ip2region.Ip2Region
}

func NewGeoFinderONLINE() (*GeoFinder, error) {
	var gf GeoFinder
	var err error

	// Download to local disk
	resp, err := ghttp.Get(ip2RegionDbUrl, "", time.Minute*2, true)
	if err != nil {
		return nil, err
	}
	buf, err := ghttp.ReadBodyBytes(resp)
	if err != nil {
		return nil, err
	}
	os.Remove("ip2region.db")
	gfs.BytesToFile(buf, "ip2region.db")

	gf.finder, err = ip2region.New("ip2region.db")
	if err != nil {
		return nil, err
	}
	return &gf, nil
}

func (gf *GeoFinder) GetByIP(ip net.IP) (*IpGeo, error) {
	var result IpGeo

	inf, err := gf.finder.MemorySearch(ip.String())
	if err != nil {
		return nil, err
	}
	result.Country = inf.Country
	result.Region = inf.Region
	result.Province = inf.Province
	result.City = inf.City
	result.CityId = inf.CityId
	result.ISP = inf.ISP

	return &result, nil
}

func (gf *GeoFinder) GetByIPString(s string) (*IpGeo, error) {
	ip, err := ParseIP(s)
	if err != nil {
		return nil, err
	}
	return gf.GetByIP(ip)
}

func (gf *GeoFinder) Close() {
	gf.finder.Close()
}
