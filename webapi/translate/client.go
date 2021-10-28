package translate

import (
	"fmt"
	"github.com/cryptowilliam/goutil/net/ghttp"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// 从 github.com/aerokite/go-google-translate/pkg fork来的，做了结构简化，增加了proxy支持

const (
	googleTranslate = `https://translate.google.com/m?client=m&oe=UTF-8&ie=UTF-8&text=%s&sl=%s&tl=%s`
)

type Client struct {
	SourceLang string
	TargetLang string
}

type translator struct {
	client *http.Client
	path   string
	err    error
	req    *http.Request
}

type Response struct {
	Err          error
	Status       string
	StatusCode   int
	ResponseBody []byte
}

func NewClient(SourceLang, TargetLang string) *Client {
	c := &Client{
		SourceLang: SourceLang,
		TargetLang: TargetLang,
	}
	return c
}

func (c *Client) Translate(text, proxy string, timeout *time.Duration) *translator {
	/*client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}*/
	client := http.DefaultClient
	if len(proxy) > 0 {
		ghttp.SetProxy(client, proxy)
	}
	ghttp.SetInsecureSkipVerify(client, true)
	if timeout != nil {
		if err := ghttp.SetTimeout(client, timeout, nil, nil, nil, nil); err != nil {
			return nil
		}
	}

	escapesText := url.QueryEscape(text)
	path := fmt.Sprintf(googleTranslate, escapesText, c.SourceLang, c.TargetLang)
	return &translator{
		path:   path,
		client: client,
	}
}

func (t *translator) Get() *translator {
	t.req, t.err = http.NewRequest("GET", t.path, nil)
	return t
}

func (t *translator) Do() (*Response, error) {
	t.req.Header.Set("Accept", "application/json")
	resp, err := t.client.Do(t.req)
	if err != nil {
		return nil, err
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return &Response{
		StatusCode:   resp.StatusCode,
		Status:       resp.Status,
		ResponseBody: responseBody,
	}, nil
}
