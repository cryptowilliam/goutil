package github

import (
	"context"
	"fmt"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/container/gstring"
	"io/ioutil"
	"strings"
)

type (
	Client struct {
		spider *Spider
		api    *API
	}
)

func NewGithub(token, proxy string) (*Client, error) {
	api := new(API)
	spider := new(Spider)
	err := error(nil)

	if len(token) > 0 {
		api, err = newAPI(token, proxy)
	} else {
		spider, err = NewSpider(proxy)
	}

	if err != nil {
		return nil, err
	}
	return &Client{api: api, spider: spider}, nil
}

func (c *Client) GetAllReposInfo(user string) (*RepoInfoList, error) {
	if c.spider != nil {
		return c.spider.GetAllReposInfo(user)
	}
	if c.api != nil {
		return c.api.GetAllReposInfo(user)
	}
	return nil, gerrors.Errorf("Invalid github client")
}

func (c *Client) GetRateLimit() (RateLimit, error) {
	if c.spider != nil {
		return c.spider.GetRateLimit()
	}
	if c.api != nil {
		return c.api.GetRateLimit()
	}
	return RateLimit{}, gerrors.Errorf("Invalid github client")
}

// example
// filepath: cartman/cartman-repo/reports/coin-tags.csv
// user: cartman
// repo: cartman-repo
// filepath: reports/tags.csv
//
// filepath is different with public repo files' raw access url like:
// https://github.com/cartman/cartman-repo/raw/master/reports/tags.csv
func (c *Client) DownloadRepoFile(filepath string) ([]byte, error) {
	if c.api.cli == nil {
		return nil, gerrors.Errorf("nil api client, make sure you have a github token")
	}

	ss := strings.Split(filepath, "/")
	ss = gstring.RemoveByValue(ss, "")
	if len(ss) < 3 {
		return nil, gerrors.Errorf("invalid filepath %s", filepath)
	}
	user := ss[0]
	repo := ss[1]
	subPath := strings.Join(ss[2:], "/")
	fmt.Println(user, repo, subPath)
	r, err := c.api.cli.Repositories.DownloadContents(context.Background(), user, repo, subPath, nil)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return b, nil
}
