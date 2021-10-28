package youtube

import (
	"fmt"
	"github.com/cryptowilliam/goutil/encoding/gjson"
	"testing"
)

func TestNew(t *testing.T) {
	videoID := "VX0_LbHaRqc"
	yt, err := NewWithKey("", "socks5://127.0.0.1:1086")
	if err != nil {
		t.Error(err)
		return
	}

	ci, err := yt.GetChannelById("UCAXi6QwUzcX9Kp_-MCZ_xFw")
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(gjson.MarshalStringDefault(ci, true))
	return

	fmt.Println(yt.GetVideosByChannelId("UCAXi6QwUzcX9Kp_-MCZ_xFw", nil))
	return

	ci, err = yt.GetChannelByUsername("mttp88")
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(gjson.MarshalStringDefault(ci, true))
	return

	vi, err := yt.GetVideo(videoID)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(gjson.MarshalStringDefault(vi, true))
	return

	comments, err := yt.GetComments(videoID)
	if err != nil {
		t.Error(err)
		return
	}
	for _, v := range comments {
		fmt.Println(gjson.MarshalStringDefault(v, false))
	}

}
