package github

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/cryptowilliam/goutil/net/ghtml"
	"github.com/cryptowilliam/goutil/net/ghttp"
	"github.com/cryptowilliam/goutil/sys/gtime"
	"strconv"
	"strings"
	"time"
)

type (
	Spider struct {
		proxy string
	}
)

func NewSpider(proxy string) (*Spider, error) {
	return &Spider{}, nil
}

func (s *Spider) GetAllReposInfo(user string) (*RepoInfoList, error) {
	res := RepoInfoList{}

	// begin: loop for all pages
	for page := 0; ; page++ {

		// build request path
		path := ""
		if page == 0 {
			path = fmt.Sprintf("https://github.com/%s?tab=repositories", user)
		} else {
			path = fmt.Sprintf("https://github.com/%s?tab=repositories&page=%d", user, page)
		}

		// get page html
		html, err := ghttp.GetString(path, s.proxy, time.Second*10)
		if err != nil {
			return nil, err
		}
		doc, err := ghtml.NewDocFromHtmlSrc(&html)
		if err != nil {
			return nil, err
		}

		// begin: loop for all repos of current page
		sel := doc.Find("h3")
		for _, v := range sel.Nodes {

			// get name
			link, ok := goquery.NewDocumentFromNode(v).Find("a").Attr("href")
			if !ok {
				continue
			}
			items := strings.Split(link, "/")
			if len(items) == 0 {
				continue
			}
			name := items[len(items)-1]

			// get repo "<li" node
			if v.Parent == nil || v.Parent.Parent == nil || v.Parent.Parent.Data != "li" {
				continue
			}
			li := v.Parent.Parent
			lidoc := goquery.NewDocumentFromNode(li)

			// get last updated datetime
			lastUpdatedStr, ok := lidoc.Find("relative-time").Attr("datetime")
			if !ok {
				continue
			}
			lastUpdated, err := gtime.ParseDatetimeStringFuzz(lastUpdatedStr)
			if err != nil {
				continue
			}

			// get stars and forks
			stars := int64(0)
			forks := int64(0)
			lidoc.Find(".muted-link").Each(func(i int, s *goquery.Selection) {
				html, err := s.Html()
				if err != nil {
					return
				}
				if strings.Contains(html, `aria-label="star"`) {
					s := strings.Replace(s.Text(), "\r", "", -1)
					s = strings.Replace(s, "\n", "", -1)
					s = strings.Replace(s, " ", "", -1)
					stars, err = strconv.ParseInt(s, 10, 32)
					if err != nil {
						return
					}
				}
				if strings.Contains(html, `aria-label="fork"`) {
					s := strings.Replace(s.Text(), "\r", "", -1)
					s = strings.Replace(s, "\n", "", -1)
					s = strings.Replace(s, " ", "", -1)
					forks, err = strconv.ParseInt(s, 10, 32)
					if err != nil {
						return
					}
				}
			})

			// get programing language
			s := lidoc.Find(`[itemprop='programmingLanguage']`).Text()
			s = strings.Replace(s, "\r", "", -1)
			s = strings.Replace(s, "\n", "", -1)
			s = strings.TrimLeft(s, " ")
			s = strings.TrimRight(s, " ")

			// append
			item := RepoInfo{Name: name, StarsCount: int(stars), ForkCount: int(forks), RepoLastUpdateTime: lastUpdated, Language: s}
			res.Items = append(res.Items, item)
		} // end: loop for all repos of current page

		// check next page
		if len(doc.Find(".next_page").Nodes) == 0 ||
			len(doc.Find(".next_page.disabled").Nodes) > 0 { // last page
			return &res, nil
		}
		// has next page

	} // end: loop for all pages

	return &res, nil
}

func (s *Spider) GetRateLimit() (RateLimit, error) {
	res := RateLimit{CoreReqPerHour: 5000,
		SearchReqPerHour: 100}
	res.CoreInterval = time.Hour / time.Duration(res.CoreReqPerHour)
	res.SearchInterval = time.Hour / time.Duration(res.SearchReqPerHour)
	return res, nil
}
