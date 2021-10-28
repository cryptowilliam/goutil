package github

import (
	"github.com/cryptowilliam/goutil/sys/gfs"
	"github.com/cryptowilliam/goutil/sys/gsysinfo"
	"testing"
	"time"
)

func TestParseRepoUrl(t *testing.T) {
	errUrls := []string{
		"xgithub.com/myname/myrepo",
		"https://xgithub.com/myname/myrepo",
	}

	for _, s := range errUrls {
		_, _, _, err := ParseUrl(s)
		if err == nil {
			t.Errorf("invalid github url %s parse error", s)
		}
	}

	okUrls := []string{
		"github.com/myname/myrepo",
		"www.github.com/myname/myrepo",
		"www.github.com/myname/myrepo/",
		"https://github.com/myname/myrepo",
		"https://www.github.com/myname/myrepo",
		"https://www.github.com/myname/myrepo/releases",
	}

	for _, s := range okUrls {
		user, repo, _, err := ParseUrl(s)
		if err != nil {
			t.Error("correct github url " + s + " parse error, error msg is " + err.Error())
		}
		if user != "myname" || repo != "myrepo" {
			t.Error("correct github url " + s + " parse error")
		}
	}
}

func TestDownloadRepo(t *testing.T) {
	savedir, err := gsysinfo.GetSharedUserDir()
	if err != nil {
		t.Error(err)
		return
	}
	savedir += gfs.DirSlash() + "tokenbase"
	t.Log("Save dir", savedir)
	err = DownloadRepo("forkdelta", "tokenbase", savedir, "", time.Minute)
	if err != nil {
		t.Error(err)
		return
	}
	if !gfs.DirExits(savedir) {
		t.Errorf("Dir %s is empty", savedir)
		return
	}
}
