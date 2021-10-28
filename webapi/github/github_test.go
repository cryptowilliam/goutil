package github

import "testing"

func TestClient_GetRepoFile(t *testing.T) {
	c, err := NewGithub("", "")
	if err != nil {
		t.Error(err)
		return
	}
	b, err := c.DownloadRepoFile("cryptowilliam/goutil/***/***/coin-tags.csv")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(b))
}
