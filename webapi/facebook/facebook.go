package facebook

import (
	"fmt"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/container/gnum"
	"github.com/cryptowilliam/goutil/container/gstring"
	"github.com/cryptowilliam/goutil/net/ghtml"
	"github.com/cryptowilliam/goutil/net/ghttp"
	"github.com/cryptowilliam/goutil/net/gnet"
	"github.com/cryptowilliam/goutil/sys/gtime"
	"strconv"
	"time"
)

// github.com/henrykhadass/gowebhook
// github.com/huandu/facebook

type Index struct {
}

type Post struct {
	Url          string
	Account      string
	CrawlTime    time.Time
	PublishTime  time.Time
	Content      string
	Forward      bool
	LikeCount    int
	CommentCount int
	ShareCount   int
	ImageUrls    []string
	VideoUrls    []string
}

func GetIndex(accountUrl, proxy string, since *time.Time) ([]Index, error) {

	return nil, nil
}

// remove all html tags to pure visible string in browser
func html2Text(html *string) string {
	doc, err := ghtml.NewDocFromHtmlSrc(html)
	if err != nil {
		return ""
	}
	return doc.Text()
}

func fbPostGetContent(html *string) (string, error) {
	// https://www.facebook.com/CBSNews/posts/10156746780825950
	// https://www.facebook.com/CBSNews/videos/2180728382257435
	htmlSlice, err := gstring.SubstrBetweenUTF8(*html, " userContent ", "</html>", true, false, false, true)
	if err == nil {
		content, err := gstring.SubstrBetweenUTF8(htmlSlice, "<p>", "</p>", true, false, false, false)
		if err == nil {
			return html2Text(&content), nil
		}
	}

	return "", gerrors.Errorf("can't find content")
}

func fbPostGetTime(html *string) (time.Time, error) {
	// https://www.facebook.com/CBSNews/posts/10156746780825950
	// https://www.facebook.com/CBSNews/videos/2180728382257435

	timeString, err := gstring.SubstrBetweenUTF8(*html, `data-utime="`, `"`, true, true, false, false)
	if err == nil {
		if gnum.IsDigit(timeString) {
			if ts, err := strconv.ParseInt(timeString, 10, 64); err == nil {
				return gtime.EpochSecToTime(ts), nil
			}
		}
	}

	return gtime.ZeroTime, gerrors.Errorf("can't parse time")
}

func fbPostGetLikeCount(html *string) (int, error) {
	// https://www.facebook.com/CBSNews/posts/10156746780825950
	// https://www.facebook.com/CBSNews/videos/2180728382257435
	if likeCountString, err := gstring.SubstrBetweenUTF8(*html, `reaction_count:{count:`, `}`, false, true, false, false); err == nil {
		if likeCount, err := strconv.ParseInt(likeCountString, 10, 64); err == nil {
			return int(likeCount), nil
		}
	}
	return 0, gerrors.Errorf("can't parse like count")
}

func fbPostGetCommentCount(html *string) (int, error) {
	// https://www.facebook.com/CBSNews/posts/10156746780825950
	// https://www.facebook.com/CBSNews/videos/2180728382257435
	if commentCountString, err := gstring.SubstrBetweenUTF8(*html, `comment_count:{total_count:`, `}`, false, true, false, false); err == nil {
		if commentCount, err := strconv.ParseInt(commentCountString, 10, 64); err == nil {
			return int(commentCount), nil
		}
	}
	return 0, gerrors.Errorf("can't parse comment count")
}

func fbPostGetShareCount(html *string) (int, error) {
	// https://www.facebook.com/CBSNews/posts/10156746780825950
	// https://www.facebook.com/CBSNews/videos/2180728382257435
	if shareCountString, err := gstring.SubstrBetweenUTF8(*html, `share_count:{count:`, `}`, false, true, false, false); err == nil {
		if shareCount, err := strconv.ParseInt(shareCountString, 10, 64); err == nil {
			return int(shareCount), nil
		}
	}
	return 0, gerrors.Errorf("can't parse share count")
}

// TODO 清除特殊字符、Unicode转换成本来的字符
func fbPostGetImageUrls(html *string) ([]string, error) {
	r := []string{}
	// https://www.facebook.com/CBSNews/posts/10156746780825950
	if imageUrl, err := gstring.SubstrBetweenUTF8(*html, `<img class="scaledImageFitWidth img" src="`, `""`, false, true, false, false); err == nil {
		if len(imageUrl) > 0 {
			if imageUrl, err = gnet.DecodeURL(imageUrl); err == nil {
				r = append(r, imageUrl)
			}
		}
	}
	return r, gerrors.Errorf("can't parse image urls")
}

// TODO 清除特殊字符、Unicode转换成本来的字符
func fbPostGetVideo(html *string) ([]string, error) {
	r := []string{}
	// https://www.facebook.com/CBSNews/videos/2180728382257435
	if videoUrl, err := gstring.SubstrBetweenUTF8(*html, `<meta property="og:video" content="`, `" />`, false, true, false, false); err == nil {
		if len(videoUrl) > 0 {
			if videoUrl, err = gnet.DecodeURL(videoUrl); err == nil {
				r = append(r, videoUrl)
			}
		}
	}
	return r, gerrors.Errorf("can't parse video urls")
}

func ParseHtml(html *string) (*Post, error) {
	_, err := ghtml.NewDocFromHtmlSrc(html)
	if err != nil {
		return nil, err
	}
	r := new(Post)
	if r.Content, err = fbPostGetContent(html); err != nil {
		return nil, err
	}
	if r.PublishTime, err = fbPostGetTime(html); err != nil {
		return nil, err
	}
	r.CrawlTime = time.Now().In(time.UTC)
	if r.LikeCount, err = fbPostGetLikeCount(html); err != nil {
		return nil, err
	}
	if r.CommentCount, err = fbPostGetCommentCount(html); err != nil {
		return nil, err
	}
	if r.ShareCount, err = fbPostGetShareCount(html); err != nil {
		return nil, err
	}
	r.ImageUrls, _ = fbPostGetImageUrls(html)
	r.VideoUrls, _ = fbPostGetVideo(html)
	return r, nil
}

func GetPost(postUrl, proxy string) (*Post, error) {
	html, err := ghttp.GetString(postUrl, proxy, time.Minute)
	if err != nil {
		return nil, err
	}
	fmt.Println(html)
	return ParseHtml(&html)
}
