package youtube

import (
	"google.golang.org/api/youtube/v3"
)

type (
	CommentThread youtube.CommentThread
)

// warning: 只是一部分评论
func (yt Youtube) GetComments(videoID string) ([]*CommentThread, error) {
	response, err := yt.service.CommentThreads.List([]string{"snippet"}).VideoId(videoID).TextFormat("plainText").Order("relevance").MaxResults(200).Do()
	if err != nil {
		return nil, err
	}
	var r []*CommentThread
	for _, v := range response.Items {
		r = append(r, (*CommentThread)(v))
	}
	return r, nil
}
