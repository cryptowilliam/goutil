package twitter

import (
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"strconv"
	"strings"
)

var (
	UrlTypeUnknown UrlType = ""
	UrlTypeHome    UrlType = "home"
	UrlTypeUser    UrlType = "user"
	UrlTypeTweet   UrlType = "tweet"
)

type (
	Url struct {
		Type     UrlType
		Username string
		TweetId  int64
	}

	UrlType string
)

// user,tweet
func ParseUrl(url string) (*Url, error) {
	userDelimiter := "twitter.com/"
	tweetDelimiter := "/status/"
	//videoDelimiter := "/watch?v="
	r := &Url{}
	defErr := gerror.Errorf("invalid youtube url(%s)", url)

	if strings.Contains(url, tweetDelimiter) {
		ss := strings.Split(url, tweetDelimiter)
		if len(ss) != 2 {
			return nil, defErr
		}
		r.Type = UrlTypeTweet
		err := error(nil)
		r.TweetId, err = strconv.ParseInt(ss[1], 10, 64)
		if err != nil {
			return nil, err
		}
		return r, nil
	}

	if strings.Contains(url, userDelimiter) {
		ss := strings.Split(url, userDelimiter)
		if len(ss) != 2 {
			return nil, defErr
		}
		if len(strings.Split(ss[1], "/")) != 1 {
			return nil, defErr
		}
		r.Type = UrlTypeUser
		r.Username = ss[1]
		return r, nil
	}

	return nil, defErr
}
