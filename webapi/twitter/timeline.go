package twitter

import (
	"github.com/ChimeraCoder/anaconda"
	"github.com/cryptowilliam/goutil/container/gnum"
	"github.com/cryptowilliam/goutil/net/ghttp"
	"net/url"
	"strings"
	"time"
)

func (a *Api) GetMentions(tweetId int64, timeout *time.Duration) (*SimpleTweet, error) {
	return nil, nil
}

func (a *Api) GetTweet(tweetId int64, timeout *time.Duration) (*SimpleTweet, error) {
	tt, err := a.inApi.GetTweet(tweetId, nil)
	if err != nil {
		return nil, err
	}
	return readSimpleTweet(&tt), nil
}

func (a *Api) GetTweets(userId int64, maxCount int64, sinceId *int64, timeout *time.Duration) ([]SimpleTweet, error) {
	if maxCount <= 0 || maxCount > MaxRetrieveTweets {
		maxCount = MaxRetrieveTweets
	}

	v := url.Values{}
	// if you don't set since_id, will get maxCount latest tweets
	if sinceId != nil {
		v.Set("since_id", gnum.FormatInt64(*sinceId))
	}
	v.Set("user_id", gnum.FormatInt64(userId))
	v.Set("count", gnum.FormatInt64(maxCount))

	if timeout != nil {
		ghttp.SetTimeout(a.inApi.HttpClient, timeout, nil, nil, nil, nil)
	}
	tweets, err := a.inApi.GetUserTimeline(v)
	if err != nil {
		return nil, err
	}

	return readSimpleTweetArray(tweets), nil
}

func (a *Api) GetLastTweet(uid int64) (*SimpleTweet, error) {
	timeout := time.Minute
	tws, err := a.GetTweets(uid, 1, nil, &timeout)
	if err != nil {
		return nil, err
	}
	if len(tws) == 0 {
		return nil, nil
	}
	return &tws[0], nil
}

func (a *Api) GetTweetsSinceTime(uid int64, sinceTime time.Time, timeout *time.Duration) ([]SimpleTweet, error) {
	var sinceId *int64 = nil
	var r []SimpleTweet
	for {
		tws, err := a.GetTweets(uid, MaxRetrieveTweets, sinceId, timeout)
		if err != nil {
			return nil, err
		}
		if len(tws) == 0 {
			break
		}
		sinceId = &tws[0].Id // 第一个Id就是最大最新的
		reachSinceTime := false
		for _, tw := range tws {
			if tw.Time.Before(sinceTime) {
				reachSinceTime = true
				break
			}
			r = append(r, tw)
		}
		if reachSinceTime {
			break
		}
		if len(tws) < MaxRetrieveTweets {
			break
		}

	}
	return r, nil
}

/*
func (a *Api) GetTweetsSinceId(userId int64, sinceId *int64, timeout *time.Duration) ([]SimpleTweet, error) {
	v := url.Values{}
	v.Set("user_id", numeric.FormatInt64(userId))
	if sinceId != nil {
		v.Set("since_id", numeric.FormatInt64(*sinceId))
	}
	if timeout != nil {
		httpz.SetTimeout(a.inApi.HttpClient, *timeout)
	}
	tweets, err := a.inApi.GetUserTimeline(v)
	if err != nil {
		return nil, err
	}

	return readSimpleTweetArray(tweets), nil
}*/

// FIXME： 哪个是评论量
func readSimpleTweet(tw *anaconda.Tweet) *SimpleTweet {
	if tw == nil {
		return nil
	}
	result := SimpleTweet{}
	result.UserId = tw.User.Id
	result.Time, _ = tw.CreatedAtTime()
	// FIXME
	// 部分加长tweet无法获取全文，比如 https://twitter.com/holochain/status/1045290743869386754
	result.FullText = tw.FullText
	// 修正"&"被转换成"&amp;"的情况，
	// 这类转换通常没必要，即使是为了方便在Html中显示，大可在放进Html之前再转换不迟
	// example: https://twitter.com/marycamacho/status/1045333515536150529
	result.FullText = strings.Replace(result.FullText, "&amp;", "&", -1)
	result.Id = tw.Id
	result.UserName = tw.User.ScreenName
	result.ReTweetCount = tw.RetweetCount
	result.FavoriteCount = tw.FavoriteCount

	for _, v := range tw.ExtendedEntities.Media {
		switch v.Type {
		case "photo":
			result.PhotoUrls = append(result.PhotoUrls, v.Media_url)
		case "multi_photos":
			result.PhotoUrls = append(result.PhotoUrls, v.Media_url)
		case "animated_gif":
			result.AnimatedGifUrls = append(result.AnimatedGifUrls, v.Media_url)
		case "video":
			result.VideoUrls = append(result.VideoUrls, v.Media_url)
		}
	}
	return &result
}

func readSimpleTweetArray(tws []anaconda.Tweet) []SimpleTweet {
	if len(tws) == 0 {
		return nil
	}
	sts := []SimpleTweet{}
	for _, v := range tws {
		sts = append(sts, *readSimpleTweet(&v))
	}
	return sts
}
