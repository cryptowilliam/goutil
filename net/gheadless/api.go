package gheadless

import (
	"github.com/goware/urlx"
)

func ScreenShot(urlStr, path string) error {
	url, err := urlx.Parse(urlStr)
	if err != nil {
		return err
	}

	cli := New()
	cli.Setup()
	if err := cli.ScreenshotURL(url, path); err != nil {
		return err
	}
	return nil
}
