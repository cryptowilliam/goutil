package twitter

// This is twitter spider API.
// If you need history data without limit of twitter API limit,
// please check https://github.com/zxl777/NewsCrawler.

// TODO: 确定GIF链接、视频链接和音频链接的获取
// TODO: 确定官方接口在返回用户历史数据时的最大限制是200还是3200 -> https://dev.twitter.com/rest/reference/get/statuses/user_timeline
// TODO: 如果网络有异常，NewApi和SubKeyWords都不返回错误，这个bug要修正

/*
"vggavuMh904DH1LHgZ6Amg0QQ",
"ZKl6EfLKrIUv3UPUj1ETWxcu0I2N1TI3qcsfVUPR9BVLmyW2AI",
"2946383622-ZkS2t3B0ZZnngEUAGSaIWLDTyk1HrAo02gUb7h7",
"S68WqCuICbHLaIrf1poXHr84WG9uwH3J5XnKPtVzl5DSV"
*/

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/cryptowilliam/goutil/container/gnum"
	"github.com/cryptowilliam/goutil/container/gstring"
	"github.com/taruti/langdetect"
	"regexp"
	"strings"
	"time"

	"github.com/cryptowilliam/goutil/net/ghttp"
)

type Api struct {
	inApi *anaconda.TwitterApi
}

const (
	// https://developer.twitter.com/en/docs/tweets/timelines/api-reference/get-statuses-user_timeline.html
	MaxRetrieveTweets = 200
)

type SimpleUser struct {
	Id             int64  // Global unique int64 id given by twitter, user can't modify it
	UserName       string // Global unique string user name in whole twitter site, users can modify it
	DisplayName    string // Display name whatever user likes
	FollowersCount int
	FollowingCount int
	TweetsCount    int
	JoinedTime     time.Time
	Birthday       time.Time
	Location       string
	Lang           langdetect.Language
	Timezone       time.Location
}

type SimpleTweet struct {
	Id              int64  // Tweet Id
	UserName        string // For generate tweet url
	FullText        string
	UserId          int64
	Time            time.Time
	PhotoUrls       []string
	AnimatedGifUrls []string
	VideoUrls       []string
	ReTweetCount    int
	FavoriteCount   int
	CommentCount    int
}

func (t *SimpleTweet) GetUrl() string {
	return fmt.Sprintf("https://twitter.com/%s/status/%s", t.UserName, gnum.FormatInt64(t.Id))
}

// Removes hashtags and @s from tweet, preserves content after # character.
// Reference: https://github.com/Conorbro/twitter-sentiment-analysis/blob/master/utils.go
func (t *SimpleTweet) GetTextWithoutAtHashtag() string {
	clean := t.GetCleanText()

	// Remove @...
	atUserRegex := regexp.MustCompile("@[A-Za-z]*")
	clean = atUserRegex.ReplaceAllString(clean, "")

	// Remove #...
	atUserRegex = regexp.MustCompile("#[A-Za-z]*")
	clean = atUserRegex.ReplaceAllString(clean, "")
	return clean
}

// 去除末尾大量的不构成语句组成部分的 @... #...
// example: https://twitter.com/H_O_L_O_/status/1046193782033731584/photo/1
// https://twitter.com/H_O_L_O_/status/1045495983126327297
func (t *SimpleTweet) GetCleanText() string {
	words := strings.Split(t.FullText, " ")
	if len(words) == 0 {
		return strings.Join(words, " ")
	}

	// Remove @... #... if their count more than 3,
	// because if the count is less than 3,
	// they may be part of the text and don't need to remove.
	cutPos := len(words) // 需要保留的单词的位置
	for i := len(words) - 1; i >= 0; i-- {
		if !(gstring.StartWith(words[i], "@") ||
			gstring.StartWith(words[i], "#")) {
			cutPos = i
			break
		}
	}
	if len(words)-cutPos >= 3 {
		words = words[0 : cutPos+1]
	}
	return strings.Join(words, " ")
}

func New(consumerKey, consumerSecret, accessToken, accessTokenSecret string, proxyUrl string) (*Api, error) {
	a := Api{}
	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	a.inApi = anaconda.NewTwitterApi(accessToken, accessTokenSecret)
	if len(proxyUrl) > 0 {
		if err := ghttp.SetProxy(a.inApi.HttpClient, proxyUrl); err != nil {
			return nil, err
		}
	}
	return &a, nil
}

func (a *Api) RateLimit() time.Duration {
	return time.Second * 60
}
