package facebook

import (
	"encoding/json"
	"fmt"
	"github.com/cryptowilliam/goutil/sys/gtime"
	"testing"
)

func TestGetIndex(t *testing.T) {
	dt, err := gtime.NewDate(2018, 10, 1)
	if err != nil {
		t.Error(err)
		return
	}
	tm := dt.ToTimeUTC()
	idx, err := GetIndex("https://www.facebook.com/justin.sun", "socks5://127.0.0.1:1086", &tm)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(idx)
}

func TestParseHtml(t *testing.T) {
	urls := []string{ //"https://www.facebook.com/CBSNews/posts/10156746780825950",
		"https://www.facebook.com/CBSNews/videos/2180728382257435",
	}
	for _, u := range urls {
		post, err := GetPost(u, "socks5://127.0.0.1:1086")
		if err != nil {
			t.Error(err)
			return
		}
		s, _ := json.Marshal(post)
		fmt.Println(string(s))
	}
}
