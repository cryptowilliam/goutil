package youtube

import (
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"strings"
)

/*
ChannelId / Username / ShowName(Title)
同一个账户可能有多个命名形式，其中，第一种方式中是username，第二种方式的最后一段就是channelId，第三种就是账号Title(或者叫show name)，显示在网页账号Logo旁的
Title: Mel-P
Username: https://www.youtube.com/user/mttp88/videos
ChannelId: https://www.youtube.com/channel/UCAXi6QwUzcX9Kp_-MCZ_xFw
*/

var (
	UrlTypeUnknown UrlType = ""
	UrlTypeHome    UrlType = "home"
	UrlTypeUser    UrlType = "user"
	UrlTypeChannel UrlType = "channel"
	UrlTypeVideo   UrlType = "video"
)

type (
	Url struct {
		Type      UrlType
		ChannelId string
		Username  string
		VideoId   string
	}

	UrlType string
)

// user,channel,video
func ParseUrl(url string) (*Url, error) {
	userDelimiter := "/user/"
	channelDelimiter := "/channel/"
	videoDelimiter := "/watch?v="
	r := &Url{}
	defErr := gerrors.Errorf("invalid youtube url(%s)", url)

	if strings.Contains(url, userDelimiter) {
		ss := strings.Split(url, userDelimiter)
		if len(ss) != 2 {
			return nil, defErr
		}
		sss := strings.Split(ss[1], "/")
		r.Type = UrlTypeUser
		r.Username = sss[0]
		return r, nil
	}
	if strings.Contains(url, channelDelimiter) {
		ss := strings.Split(url, channelDelimiter)
		if len(ss) != 2 {
			return nil, defErr
		}
		sss := strings.Split(ss[1], "/")
		r.Type = UrlTypeChannel
		r.ChannelId = sss[0]
		return r, nil
	}
	if strings.Contains(url, videoDelimiter) {
		ss := strings.Split(url, videoDelimiter)
		if len(ss) != 2 {
			return nil, defErr
		}
		r.Type = UrlTypeVideo
		r.VideoId = ss[1]
		return r, nil
	}
	return nil, defErr
}
