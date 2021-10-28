package github

import (
	"context"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/net/ghttp"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"net/http"
	"runtime"
	"strings"
	"time"
)

type (
	ReleaseAsset struct {
		Name        string
		ByteSize    int
		DownloadUrl string
	}

	RateLimit struct {
		CoreReqPerHour, SearchReqPerHour int
		CoreInterval, SearchInterval     time.Duration
	}

	Release struct {
		Version       string
		DescTitle     string
		DescContent   string
		PublishedTime time.Time
		Assets        []ReleaseAsset
	}
)

// 根据文件名尝试分析当前平台对应的条目。
// 主要就是分析文件名中的"windows"/"darwin"/"linux","i386"/"arm64"/"amd64"等字样
func (r *Release) ParseCurrentPlatform() (*ReleaseAsset, error) {
	os := strings.ToLower(runtime.GOOS)
	arch := strings.ToLower(runtime.GOARCH)

	for _, v := range r.Assets {
		name := strings.ToLower(v.Name)
		if strings.Contains(name, os) && strings.Contains(name, arch) {
			return &v, nil
		}
	}

	return nil, gerrors.Errorf("Can't find release for %s-%s", os, arch)
}

func GetLatestRelease(user, repo string) (*Release, error) {
	res := Release{}

	c := github.NewClient(http.DefaultClient)
	release, _, err := c.Repositories.GetLatestRelease(context.Background(), user, repo)
	if err != nil {
		return nil, err
	}
	res.Version = *release.TagName
	res.DescTitle = *release.Name
	res.DescContent = *release.Body
	res.PublishedTime = release.PublishedAt.Time

	for _, v := range release.Assets {
		asset := ReleaseAsset{}
		asset.Name = *v.Name
		asset.ByteSize = *v.Size
		asset.DownloadUrl = *v.BrowserDownloadURL
		res.Assets = append(res.Assets, asset)
	}

	return &res, nil
}

type API struct {
	cli *github.Client
}

func newAPI(token, proxy string) (*API, error) {
	res := API{}
	tc := new(http.Client)

	if token == "" {
		res.cli = github.NewClient(nil)
	} else {
		tc = oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		))
		res.cli = github.NewClient(tc)
	}

	if len(proxy) > 0 {
		if err := ghttp.SetProxy(tc, proxy); err != nil {
			return nil, err
		}
	}
	return &res, nil
}

func (c *API) GetRateLimit() (RateLimit, error) {
	limits, _, err := c.cli.RateLimits(context.Background())
	if err != nil {
		return RateLimit{}, err
	}
	res := RateLimit{CoreReqPerHour: limits.GetCore().Limit,
		SearchReqPerHour: limits.GetSearch().Limit}
	res.CoreInterval = time.Hour / time.Duration(res.CoreReqPerHour)
	res.SearchInterval = time.Hour / time.Duration(res.SearchReqPerHour)
	return res, nil
}

func (c *API) GetRepoInfo(user, repo string) (*RepoInfo, error) {
	res := RepoInfo{}

	r, _, err := c.cli.Repositories.Get(context.Background(), user, repo)
	if err != nil {
		return nil, err
	}
	res.Name = repo
	res.CreateTime = r.GetCreatedAt().Time
	res.RepoLastUpdateTime = r.GetUpdatedAt().Time
	res.CodeLastPushedTime = r.GetPushedAt().Time
	res.StarsCount = r.GetStargazersCount()
	res.ForkCount = r.GetForksCount()
	res.Language = r.GetLanguage()
	return &res, nil
}

func (c *API) GetRepoBasicInfoWithUrl(repoUrl string) (*RepoInfo, error) {
	user, repo, _, err := ParseUrl(repoUrl)
	if err != nil {
		return nil, err
	}
	return c.GetRepoInfo(user, repo)
}

func (c *API) GetAllReposInfo(user string) (*RepoInfoList, error) {
	repos, _, err := c.cli.Repositories.List(context.Background(), user, nil)
	if err != nil {
		return nil, err
	}

	info := RepoInfo{}
	infos := RepoInfoList{}
	for _, v := range repos {
		info.CreateTime = v.GetCreatedAt().Time
		info.CodeLastPushedTime = v.GetPushedAt().Time
		info.RepoLastUpdateTime = v.GetUpdatedAt().Time
		info.StarsCount = v.GetStargazersCount()
		info.ForkCount = v.GetForksCount()
		info.Language = v.GetLanguage()
		infos.Items = append(infos.Items, info)
	}
	return &infos, nil
}

/*
func (c *API) GetMostStarsRepo(user string) (repo string, err error) {
	repos, err := c.ListAllRepos(user)
	if err != nil {
		return "", err
	}
	if len(repos) == 0 {
		return "", gerrors.Errorf("No repository")
	}

	maxStarsCount := 0
	for _, v := range repos {
		if v.StarsCount > maxStarsCount {
			maxStarsCount = v.StarsCount
		}
	}
	if maxStarsCount == 0 {
		return "", gerrors.Errorf("No repository has any star")
	}
	for _, v := range repos {
		if v.StarsCount == maxStarsCount {
			return v.Name, nil
		}
	}
	return "", gerrors.Errorf("Unknown error")
}*/
