package gcmd

import "github.com/pkg/browser"
import "github.com/StevenZack/openurl"

type BrowserCMD struct {
}

func NewBrowser() *BrowserCMD {
	return &BrowserCMD{}
}

func (b *BrowserCMD) OpenURL(url string) error {
	return browser.OpenURL(url)
}

func (b *BrowserCMD) OpenFile(path string) error {
	return browser.OpenFile(path)
}

// open in Chrome APP mode, no toolbar/buttons, just web page
func (b *BrowserCMD) OpenURLAppMode(url string) error {
	return openurl.OpenApp(url)
}
