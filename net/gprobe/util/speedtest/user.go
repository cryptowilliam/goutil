package speedtest

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"io/ioutil"
	"net/http"
)

// User information
type User struct {
	IP  string `xml:"ip,attr"`
	Lat string `xml:"lat,attr"`
	Lon string `xml:"lon,attr"`
	Isp string `xml:"isp,attr"`
}

// For decode xml
type Users struct {
	Users []User `xml:"client"`
}

func fetchUserInfo() (*User, error) {
	// Fetch xml user data
	resp, err := http.Get("http://speedtest.net/speedtest-config.php")
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Decode xml
	decoder := xml.NewDecoder(bytes.NewReader(body))
	users := Users{}
	for {
		t, _ := decoder.Token()
		if t == nil {
			break
		}
		switch se := t.(type) {
		case xml.StartElement:
			decoder.DecodeElement(&users, &se)
		}
	}
	if users.Users == nil {
		return nil, gerrors.New("Cannot fetch user information. http://www.speedtest.net/speedtest-config.php is temporarily unavailable.")
	}
	return &users.Users[0], nil
}

func (u User) String() string {
	if u.IP != "" {
		return fmt.Sprintf("IP: " + u.IP + " (" + u.Isp + ") [" + u.Lat + ", " + u.Lon + "]")
	}
	return ""
}
