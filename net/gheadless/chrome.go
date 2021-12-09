package gheadless

// This module has been tested successfully.

import (
	"context"
	"github.com/chromedp/chromedp"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/basic/glog"
	"io/ioutil"
	"strings"
	"time"
)

type (
	ChromeHeadless struct {
	}
)

func NewChrome() *ChromeHeadless {
	return &ChromeHeadless{}
}

func (ch *ChromeHeadless) Screenshot(urlStr, pathToSavePNG, proxy string, log glog.Interface, timeout time.Duration) error {
	// fix urlStr, chromedp doesn't accept URL without "http://" or "https://"
	if !strings.HasPrefix(strings.ToLower(urlStr), "http://") && !strings.HasPrefix(strings.ToLower(urlStr), "https://") {
		urlStr = "http://"+urlStr
	}

	// create context
	opts := chromedp.DefaultExecAllocatorOptions[:]
	if proxy != "" { // add proxy
		opts = append(opts,
			// 1) specify the proxy server.
			// Note that the username/password is not provided here.
			// Check the link below for the description of the proxy settings:
			// https://www.chromium.org/developers/design-documents/network-settings
			chromedp.ProxyServer(proxy),
			// By default, Chrome will bypass localhost.
			// The test server is bound to localhost, so we should add the
			// following flag to use the proxy for localhost URLs.
			chromedp.Flag("proxy-bypass-list", "<-loopback>"),
		)
	}
	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// log the protocol messages to understand how it works.
	debugFunc := func(string, ...interface{}) {}
	if log != nil {
		debugFunc = log.Debgf
	}
	// Note:
	// If remove two lines code below, Screenshot() will report "invalid context",
	// so if input "log" is nil, we give it a fake debug callback,
	// I don't know why.
	ctx, cancel = chromedp.NewContext(ctx, chromedp.WithDebugf(debugFunc))
	defer cancel()

	// capture entire browser viewport, returning png with quality=90
	chDone := make(chan error, 1)
	go func() {
		var buf []byte
		if err := chromedp.Run(ctx, fullScreenshot(urlStr, 90, &buf)); err != nil {
			chDone <- err
			return
		}
		if err := ioutil.WriteFile(pathToSavePNG, buf, 0o644); err != nil {
			chDone <- err
			return
		}
		close(chDone) // "chan error" returns nil after close() action
	}()

	// wait screenshot result
	if timeout > 0 {
		ticker := time.NewTicker(timeout)
		select {
		case <-ticker.C:
			return gerrors.ErrTimeout
		case err := <-chDone:
			return err
		}
	} else {
		return <-chDone
	}
}

// fullScreenshot takes a screenshot of the entire browser viewport.
//
// Note: chromedp.FullScreenshot overrides the device's emulation settings. Use
// device.Reset to reset the emulation and viewport settings.
func fullScreenshot(urlStr string, quality int, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlStr),
		chromedp.FullScreenshot(res, quality),
	}
}