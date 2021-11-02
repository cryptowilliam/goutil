package github

import (
	"context"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/net/ghttp"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"io/ioutil"
	"time"
)

type (
	GitDrive struct {
		gc *github.Client
	}

	GitPartition struct {
		drv       *GitDrive
		partition string
	}
)

func NewGitDrive(token, proxy string) (*GitDrive, error) {
	r := &GitDrive{}

	if token == "" {
		return nil, gerrors.Errorf("token required")
	}
	tc := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}))
	r.gc = github.NewClient(tc)

	if len(proxy) > 0 {
		if err := ghttp.SetProxy(tc, proxy); err != nil {
			return nil, err
		}
	}
	return r, nil
}

func (gd *GitDrive) MyName() (string, error) {
	return "", nil
}

func (gd *GitDrive) Partitions() ([]string, error) {
	repos, _, err := gd.gc.Repositories.List(context.Background(), "", nil)
	if err != nil {
		return nil, err
	}
	var r []string
	for _, v := range repos {
		if *v.Private == true {
			r = append(r, *v.Name)
		}
	}
	return r, nil
}

func (gd *GitDrive) CreatePartition(name string) error {
	private := true
	gitignore := "go"
	_, _, err := gd.gc.Repositories.Create(context.Background(), "", &github.Repository{Name: &name, Private: &private, GitignoreTemplate: &gitignore})
	return err
}

func (gd *GitDrive) Partition(name string) *GitPartition {
	r := &GitPartition{
		drv:       gd,
		partition: name,
	}
	return r
}

func (gp *GitPartition) LastUpdateTime() (time.Time, error) {
	statuses, _, err := gp.drv.gc.Repositories.ListStatuses(context.Background(), "nifflerfox", gp.partition, "master", &github.ListOptions{Page: 1, PerPage: 1})
	if err != nil {
		return time.Time{}, err
	}
	if len(statuses) != 1 {
		return time.Time{}, gerrors.Errorf("statuses length is %d", len(statuses))
	}
	return *statuses[0].UpdatedAt, nil
}

// Download repository as zip file, and unzip it to specify save directory
func (gp *GitPartition) ReadAsZip() ([]byte, error) {
	r, err := gp.drv.gc.Repositories.DownloadContents(context.Background(), "", gp.partition, "", nil)
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

func (gp *GitPartition) ReadFile(file string) ([]byte, error) {
	r, err := gp.drv.gc.Repositories.DownloadContents(context.Background(), "", gp.partition, file, nil)
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

func (gp *GitPartition) WriteFile(file string, data []byte) error {
	_, _, err := gp.drv.gc.Repositories.CreateFile(context.Background(), "", gp.partition, file, &github.RepositoryContentFileOptions{Content: data})
	if err == nil {
		return nil
	}
	_, _, err = gp.drv.gc.Repositories.UpdateFile(context.Background(), "", gp.partition, file, &github.RepositoryContentFileOptions{Content: data})
	return err
}
