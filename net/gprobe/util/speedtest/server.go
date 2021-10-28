package speedtest

import (
	"encoding/xml"
	"fmt"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/container/ggeo"
	"github.com/cryptowilliam/goutil/container/gnum"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
)

// TestServer information
type TestServer struct {
	URL      string `xml:"url,attr"`
	Lat      string `xml:"lat,attr"`
	Lon      string `xml:"lon,attr"`
	Name     string `xml:"name,attr"`
	Country  string `xml:"country,attr"`
	Sponsor  string `xml:"sponsor,attr"`
	ID       string `xml:"id,attr"`
	URL2     string `xml:"url2,attr"`
	Host     string `xml:"host,attr"`
	Distance float64
}

// TestServerList : List of TestServer
type TestServerList struct {
	Servers []TestServer `xml:"servers>server"`
}

// ByDistance : For sorting Servers.
type ByDistance struct {
	TestServerList
}

func getAllTestServers(user *User) (*TestServerList, error) {
	// Fetch xml server data
	resp, err := http.Get("http://www.speedtest.net/speedtest-servers-static.php")
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if len(body) == 0 {
		resp, err = http.Get("http://c.speedtest.net/speedtest-servers-static.php")
		if err != nil {
			return nil, err
		}
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
	}

	if resp.StatusCode != 200 {
		return nil, gerrors.New("Fail to get speedtest server list, status code " + strconv.FormatInt(int64(resp.StatusCode), 10))
	}
	if len(body) == 0 {
		return nil, gerrors.New("Fail to get speedtest server list, empty response!")
	}

	// 为什么这里读不到任何东西
	//fmt.Println(ioutil.ReadAll(resp.Body))

	// Decode xml
	var l TestServerList
	xml.Unmarshal(body, &l)
	if len(l.Servers) == 0 {
		return nil, gerrors.New("No speedtest server got")
	}

	// Calculate distance
	// Notice:
	// If you want write value to list item struct, don't use "for _, v := range mylist", but do it like below
	// Because in "for _, v := range mylist", v is a clone of real item struct, not itself
	for i := range l.Servers {
		server := &l.Servers[i]
		sLat, _ := strconv.ParseFloat(server.Lat, 64)
		sLon, _ := strconv.ParseFloat(server.Lon, 64)
		uLat, _ := strconv.ParseFloat(user.Lat, 64)
		uLon, _ := strconv.ParseFloat(user.Lon, 64)
		server.Distance = ggeo.GeoDistance(sLat, sLon, uLat, uLon)
	}

	return &l, nil
}

// Get server information
func (s TestServer) GetInfo() string {
	return fmt.Sprintf("[%4s] %8.2fkm ", s.ID, s.Distance) + fmt.Sprintf(s.Name+" ("+s.Country+") by "+s.Sponsor)
}

func (s TestServerList) Len() int {
	return len(s.Servers)
}

// Perp : swap i-th and j-th. For sorting Servers.
func (s TestServerList) Swap(i, j int) {
	s.Servers[i], s.Servers[j] = s.Servers[j], s.Servers[i]
}

// Less : compare the distance. For sorting Servers.
func (b ByDistance) Less(i, j int) bool {
	return b.Servers[i].Distance < b.Servers[j].Distance
}

// selectNearestTestServers : find server by serverID
func (l *TestServerList) selectNearestTestServers(selectCount int) (*TestServerList, error) {
	if selectCount <= 0 {
		return nil, gerrors.New("Can't select " + strconv.FormatInt(int64(selectCount), 10) + " servers")
	}
	if len(l.Servers) == 0 {
		return nil, gerrors.New("Can't search from empty server list")
	}

	// Sort by distance
	sort.Sort(ByDistance{*l})

	var result TestServerList
	selectCount = gnum.MinInt(len(l.Servers), selectCount)
	for i := 0; i < selectCount; i++ {
		result.Servers = append(result.Servers, l.Servers[i])
		fmt.Println("Test server: " + result.Servers[i].GetInfo())
	}

	return &result, nil
}
