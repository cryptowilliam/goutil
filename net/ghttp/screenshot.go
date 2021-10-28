package ghttp

// Command screenshot is a chromedp example demonstrating how to take a
// screenshot of a specific element and of the entire browser viewport.

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/dom"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"io/ioutil"
	"math"
	"time"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

const (
	CloudflareSel = "div.cf-browser-verification"
)

// elementScreenshot takes a screenshot of a specific element.
func ScreenshotElement(urlStr, sel string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlStr),
		chromedp.WaitVisible(sel, chromedp.ByID),
		chromedp.Screenshot(sel, res, chromedp.NodeVisible, chromedp.ByID),
	}
}

// fullScreenshot takes a screenshot of the entire browser viewport.
// Note: this will override the viewport emulation settings.
func Screenshot(urlStr, proxyServer, waitVisibleByQuery string, timeout time.Duration) (snapshotImage []byte, fullHtml string, err error) {
	var buf []byte
	html := ""
	// quality compression quality from range [0..100] (jpeg only).
	quality := int64(100)

	var waitVisibleAction chromedp.Action = nil
	if waitVisibleByQuery != "" {
		waitVisibleAction = chromedp.WaitReady(waitVisibleByQuery, chromedp.ByQuery)
	}

	t := chromedp.Tasks{
		chromedp.Navigate(urlStr),
		chromedp.WaitNotPresent(CloudflareSel, chromedp.ByQuery), // cross cloudflare DDos protecting page
		waitVisibleAction,
		chromedp.Sleep(1 * time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {

			// get full html
			node, err := dom.GetDocument().Do(ctx)
			if err != nil {
				return err
			}
			html, err = dom.GetOuterHTML().WithNodeID(node.NodeID).Do(ctx)
			if err != nil {
				return err
			}

			// get layout metrics
			_, _, contentSize, _, _, _, err := page.GetLayoutMetrics().Do(ctx)
			if err != nil {
				return err
			}

			width, height := int64(math.Ceil(contentSize.Width)), int64(math.Ceil(contentSize.Height))

			// force viewport emulation
			err = emulation.SetDeviceMetricsOverride(width, height, 1, false).
				WithScreenOrientation(&emulation.ScreenOrientation{
					Type:  emulation.OrientationTypePortraitPrimary,
					Angle: 0,
				}).
				Do(ctx)
			if err != nil {
				return err
			}

			// capture screenshot
			buf, err = page.CaptureScreenshot().
				WithQuality(quality).
				WithClip(&page.Viewport{
					X:      contentSize.X,
					Y:      contentSize.Y,
					Width:  contentSize.Width,
					Height: contentSize.Height,
					Scale:  1,
				}).Do(ctx)
			if err != nil {
				return err
			}
			return nil
		}),
	}

	// create context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	if ctx == nil {
		return nil, "", gerrors.New("nil chromedp ctx")
	}

	// force max timeout for retrieving and processing the data
	if timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, timeout)
	}

	// capture entire browser viewport, returning png with quality=90
	if proxyServer != "" {
		chromedp.ProxyServer(proxyServer)
	}

	// run
	if err := chromedp.Run(ctx, t); err != nil {
		return nil, "", err
	}
	return buf, html, nil
}

func ScreenshotEx(urlstr, proxyServer, waitVisibleByQuery string, timeout time.Duration, savePath string) error {
	data, html, err := Screenshot(urlstr, proxyServer, waitVisibleByQuery, timeout)
	if err != nil {
		return err
	}
	fmt.Println(html)

	if err := ioutil.WriteFile(savePath, data, 0644); err != nil {
		return err
	}
	return nil
}
