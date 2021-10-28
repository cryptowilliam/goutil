package youtube

import (
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"google.golang.org/api/youtube/v3"
	"time"
)

type (
	Channel youtube.Channel
)

func (ch *Channel) ChannelId() string {
	return ch.Id // this is channelId
}

// ch.Snippet.Title is account title or showName, not username(username is like https://www.youtube.com/user/mttp88/videos)
func (ch *Channel) Title() string {
	return ch.Snippet.Title
}

// 根据用户名查询Channel信息
// example: "TIME"
// The getChannelInfo uses forUsername
// to get info (id, tittle, totalViews and description)
func (yt Youtube) GetChannelByUsername(forUsername string) (*Channel, error) {
	if forUsername == "" {
		return nil, gerrors.Errorf("empty forUsername")
	}
	response, err := yt.service.Channels.List(snippetContentDetailsStatistics).ForUsername(forUsername).Do()
	if err != nil {
		return nil, err
	}
	if len(response.Items) == 0 {
		return nil, gerrors.Errorf("nil response length in GetChannelByUsername")
	}
	return (*Channel)(response.Items[0]), nil
}

func (yt Youtube) GetChannelById(channelId string) (*Channel, error) {
	if channelId == "" {
		return nil, gerrors.Errorf("empty channelId")
	}
	response, err := yt.service.Channels.List(snippetContentDetailsStatistics).Id(channelId).Do()
	if err != nil {
		return nil, err
	}
	if len(response.Items) == 0 {
		return nil, gerrors.Errorf("nil response length in GetChannelById")
	}
	return (*Channel)(response.Items[0]), nil
}

// 不可以检索完整的视频信息，要详细的视频信息还是要调用GetVideoInfo接口
func (yt Youtube) GetVideosByChannelId(channelId string, publishAfter *time.Time) ([]string, error) {
	pageToken := ""
	var result []string

	for {
		call := yt.service.Search.List([]string{"snippet", "id"}).Type("video").ChannelId(channelId).MaxResults(50).Order("date")
		if pageToken != "" {
			call = call.PageToken(pageToken)
		}
		if publishAfter != nil {
			call = call.PublishedAfter(publishAfter.Format(time.RFC3339))
		}
		response, err := call.Do()
		if err != nil {
			return nil, err
		}

		// Iterate through each item and add it to the correct list.
		for _, item := range response.Items {
			switch item.Id.Kind {
			case "youtube#video":
				// snippet.channelTitle就是账号的username
				result = append(result, item.Id.VideoId)
			}
		}

		if response.NextPageToken != "" {
			pageToken = response.NextPageToken
		} else {
			break
		}
	}

	return result, nil
}
