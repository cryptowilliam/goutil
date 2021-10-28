package github

import (
	"fmt"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/encoding/gzip"
	"github.com/cryptowilliam/goutil/net/ghttp"
	"github.com/cryptowilliam/goutil/sys/gtime"
	"os"
	"path/filepath"
	"time"
)

type (
	RepoInfo struct {
		Name               string
		CreateTime         time.Time
		RepoLastUpdateTime time.Time // Repository last updated time, include issue/wiki....
		CodeLastPushedTime time.Time // Source code last pushed time.
		StarsCount         int
		ForkCount          int
		Language           string
	}

	RepoInfoList struct {
		Items []RepoInfo
	}
)

func (rs *RepoInfoList) MainRepo() (*RepoInfo, error) {
	mostStars := 0
	mostStarRepo := ""
	for _, v := range rs.Items {
		if v.StarsCount > mostStars {
			mostStars = v.StarsCount
			mostStarRepo = v.Name
		}
	}
	if mostStars > 0 {
		for _, v := range rs.Items {
			if v.Name == mostStarRepo {
				return &v, nil
			}
		}
	}

	mostFork := 0
	mostForkRepo := ""
	for _, v := range rs.Items {
		if v.ForkCount > mostFork {
			mostFork = v.ForkCount
			mostForkRepo = v.Name
		}
	}
	if mostFork > 0 {
		for _, v := range rs.Items {
			if v.Name == mostForkRepo {
				return &v, nil
			}
		}
	}

	latestUpdate := gtime.EpochBeginTime
	latestUpdateRepo := ""
	for _, v := range rs.Items {
		if v.RepoLastUpdateTime.After(latestUpdate) {
			latestUpdate = v.RepoLastUpdateTime
			latestUpdateRepo = v.Name
		}
	}
	if latestUpdate != gtime.EpochBeginTime {
		for _, v := range rs.Items {
			if v.Name == latestUpdateRepo {
				return &v, nil
			}
		}
	}

	if len(rs.Items) > 0 {
		return &rs.Items[0], nil
	}

	return nil, gerrors.Errorf("No main repo")
}

// Download repository as zip file, and unzip it to specify save directory
func DownloadRepo(user, repo, saveDir string, proxy string, timeout time.Duration) error {
	// make HTTP request path
	path := fmt.Sprintf("https://github.com/%s/%s/archive/master.zip", user, repo)

	// Get zip file and unzip to temp dir
	resp, err := ghttp.Get(path, proxy, timeout, true)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Clean unzip dir before UnZip
	os.RemoveAll(filepath.Join(saveDir, repo+"-master"))
	if err := gzip.UnZip(resp.Body, saveDir); err != nil {
		return err
	}

	return nil
}
