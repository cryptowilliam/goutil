package github

import (
	"fmt"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/container/gstring"
	"net/url"
	"strings"
)

// support github.com/***, www.github.com/***, gist.github.com/***
func ParseUrl(githubUrl string) (user, repo, sub string, err error) {
	// If not begin with "http://" or "https://", add it.
	// url.Parse接口传入的url必须由Scheme开头，否则整个url都会被认为是Url中的Path
	s := strings.ToLower(githubUrl)
	if !gstring.StartWith(s, "http://") && !gstring.StartWith(s, "https://") {
		githubUrl = "https://" + githubUrl
	}

	// Parse
	URL, err := url.Parse(githubUrl)
	if err != nil {
		return "", "", "", err
	}

	// Check host.
	hoststr := strings.ToLower(URL.Host)
	if hoststr != "github.com" && !gstring.EndWith(hoststr, ".github.com") {
		return "", "", "", gerrors.Errorf("url(%s, host:%s) is not a github url", githubUrl, URL.Host)
	}

	// Check path.
	subPath := URL.Path
	if len(subPath) < 3 {
		return "", "", "", gerrors.Errorf("url(%s, subPath:%s) is not a github repository url", githubUrl, subPath)
	}
	// Remove "/" from head.
	if subPath[0:1] == "/" {
		subPath = subPath[1:]
	}
	items := strings.Split(subPath, "/")
	if len(items) == 0 {
		return "", "", "", gerrors.Errorf("%s is not a github repository url", githubUrl)
	}

	if hoststr == "www.github.com" || hoststr == "github.com" || hoststr == "gist.github.com" {
		user = items[0]
		items[0] = "" // clean user for make subpath
	}
	if (hoststr == "www.github.com" || hoststr == "github.com") && len(items) > 1 {
		repo = items[1]
		items[1] = "" // clean repo for make subpath
	}

	subPath = strings.Join(items, "/")
	for gstring.StartWith(subPath, "/") {
		subPath = subPath[1:]
	}

	return user, repo, subPath, nil
}

func BuildRepoUrl(user, repo string) string {
	return fmt.Sprintf("https://www.github.com/%s/%s", user, repo)
}
