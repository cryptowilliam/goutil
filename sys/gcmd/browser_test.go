package gcmd

import "testing"

func TestBrowserCMD_OpenURL(t *testing.T) {
	err := NewBrowser().OpenURL("http://127.0.0.1:8080")
	if err != nil {
		t.Error(err)
		return
	}
}
