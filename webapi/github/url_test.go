package github

import (
	"fmt"
	"testing"
)

func TestParseUrl(t *testing.T) {
	type testitem struct {
		url     string
		user    string
		repo    string
		subpath string
	}

	ts := []testitem{}
	ts = append(ts, testitem{url: "https://gist.github.com/jabbawookiees/2ffa21af85e637fbb1e3b661739cb96b", user: "jabbawookiees", repo: "", subpath: "2ffa21af85e637fbb1e3b661739cb96b"})
	ts = append(ts, testitem{url: "https://github.com/LinkEyeOfficial", user: "LinkEyeOfficial", repo: "", subpath: ""})
	ts = append(ts, testitem{url: "https://github.com/numuscrypto/numuscore", user: "numuscrypto", repo: "numuscore", subpath: ""})
	ts = append(ts, testitem{url: "https://github.com/MANOPlatform/manocoin/releases/tag/v1.0.0", user: "MANOPlatform", repo: "manocoin", subpath: "releases/tag/v1.0.0"})
	ts = append(ts, testitem{url: "https://github.com/cartman/cartman-repo/raw/master/reports/tags.csv", user: "cartman", repo: "cartman-repo", subpath: "raw/master/reports/tags.csv"})
	ts = append(ts, testitem{url: "https://github.com/cartman/cartman-repo/reports/tags.csv", user: "cartman", repo: "cartman-repo", subpath: "reports/tags.csv"})

	for _, v := range ts {
		user, repo, subpath, err := ParseUrl(v.url)
		if err != nil {
			t.Error(err)
			return
		}
		if user != v.user || repo != v.repo || subpath != v.subpath {
			t.Errorf("ParseUrl(%s) get (%s, %s, %s), but should get (%s, %s, %s)",
				v.url, user, repo, subpath, v.user, v.repo, v.subpath)
			return
		}
	}
}

func TestRepoInfoList_MainRepo(t *testing.T) {
	type testitem struct {
		user     string
		mainrepo string
	}

	ts := []testitem{}
	ts = append(ts, testitem{user: "LinkEyeOfficial", mainrepo: ""})
	ts = append(ts, testitem{user: "numuscrypto", mainrepo: "numuscrypto"})

	s, err := NewSpider("")
	if err != nil {
		t.Error(err)
		return
	}
	for _, v := range ts {
		repos, err := s.GetAllReposInfo(v.user)
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Println(repos)
		mainrepo, _ := repos.MainRepo()
		if mainrepo != nil && mainrepo.Name != v.mainrepo {
			t.Errorf("mainrepo of user (%s) get (%s), but should get (%s)", v.user, mainrepo.Name, v.mainrepo)
			return
		}
		if mainrepo == nil && v.mainrepo != "" {
			t.Errorf("mainrepo of user (%s) get (nil), but should get (%s)", v.user, v.mainrepo)
			return
		}
	}
}
