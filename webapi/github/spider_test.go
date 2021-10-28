package github

import (
	"fmt"
	"testing"
)

func TestSpider_GetAllReposInfo(t *testing.T) {
	c, err := NewSpider("")
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println("------------")
	repos, err := c.GetAllReposInfo("PascalLite")
	if err != nil {
		t.Error(err)
		return
	}
	for _, v := range repos.Items {
		fmt.Println(v)
	}
	if mainrepo, err := repos.MainRepo(); err == nil {
		fmt.Println("main repo:", *mainrepo)
	}

	fmt.Println("------------")
	repos, err = c.GetAllReposInfo("holochain")
	if err != nil {
		t.Error(err)
		return
	}
	for _, v := range repos.Items {
		fmt.Println(v)
	}
	if mainrepo, err := repos.MainRepo(); err == nil {
		fmt.Println("main repo:", *mainrepo)
	}
}
