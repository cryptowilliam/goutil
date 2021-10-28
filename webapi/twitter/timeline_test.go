package twitter

import (
	"fmt"
	"github.com/cryptowilliam/goutil/encoding/gjson"
	"github.com/cryptowilliam/goutil/sys/gtime"
	"testing"
	"time"
)

func TestApi_GetTweetsEx(t *testing.T) {
	//urls := `https://twitter.com/asahi_kokusai`

	tw, err := New("",
		"",
		"",
		"",
		"socks5://127.0.0.1:1086")
	if err != nil {
		t.Error(err)
		return
	}

	// https://twitter.com/asahi_kokusai/status/1224163734752378881
	// https://twitter.com/inoko1102/status/1223455869633028096
	fmt.Println(tw.GetTweet(1211546208516182024, nil))
	return

	user, err := tw.GetUserByUsername("asahi_kokusai")
	if err != nil {
		t.Error(err)
		return
	}

	sinceTime := gtime.Sub(time.Now(), gtime.Day*10)
	timeout := time.Minute * 2
	tws, err := tw.GetTweetsSinceTime(user.Id, sinceTime, &timeout)
	if err != nil {
		t.Error(err)
		return
	}
	for _, v := range tws {
		fmt.Println(gjson.MarshalStringDefault(v, false))
	}
}
