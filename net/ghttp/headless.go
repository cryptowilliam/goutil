package ghttp

import (
	"context"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
	phantomgo "github.com/cryptowilliam/goutil/net/ghttp/phantom"
	"gopkg.in/headzoo/surf.v1"
	"io/ioutil"
)

func ChromeDriverGet(urlStr, proxy string) (string, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var res string
	chromedp.ProxyServer(proxy)
	err := chromedp.Run(ctx,
		chromedp.Navigate(urlStr),
		chromedp.ActionFunc(func(ctx context.Context) error {
			node, err := dom.GetDocument().Do(ctx)
			if err != nil {
				return err
			}
			res, err = dom.GetOuterHTML().WithNodeID(node.NodeID).Do(ctx)
			return err
		}),
	)
	if err != nil {
		return "", err
	}
	return res, nil
}

/*
func HttpHeadlessGet(url, proxy, sel string, timeout time.Duration) (string, error) {
	var err error

	// create context
	ctxt, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create chrome instance
	c, err := chromedp.New(ctxt, chromedp.WithLog(log.Printf))
	if err != nil {
		return "", err
	}

	// run task list
	var buf []byte
	err = c.Run(ctxt, headlessGet(url, sel, &buf))
	if err != nil {
		return "", err
	}

	// shutdown chrome
	err = c.Shutdown(ctxt)
	if err != nil {
		return "", err
	}

	// wait for chrome to finish
	err = c.Wait()
	if err != nil {
		return "", err
	}

	return str, nil

}

func headlessGet(urlstr, sel string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.Sleep(2 * time.Second),
		chromedp.WaitVisible(sel, chromedp.ByID),
	}
}*/

func HeadlessGetPhantom(url, proxy string) (string, error) {
	p := &phantomgo.Param{
		Method:       "GET", //POST or GET ..
		Url:          url,
		UsePhantomJS: true,
	}
	brower := phantomgo.NewPhantom()
	if proxy != "" {
		brower.SetProxy(proxy)
	}
	resp, err := brower.Download(p)
	if err != nil {
		return "", err
	}
	if body, err := ioutil.ReadAll(resp.Body); err != nil {
		return "", err
	} else {
		return string(body), nil
	}
}

/*
func HeadlessGetCdp(url string) (string, error) {
	var err error

	// create context
	ctxt, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create chrome instance
	c, err := chromedp.New(ctxt, chromedp.WithLog(log.Printf))
	if err != nil {
		return "", err
	}

	// run task list
	var res []string
	err = c.Run(ctxt, chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.WaitReady(`html`, chromedp.ByID),
		chromedp.Evaluate(`document.documentElement.innerHTML`, &res),
	})



	if err != nil {
		return "", err
	}

	// shutdown chrome
	err = c.Shutdown(ctxt)
	if err != nil {
		return "", err
	}

	// wait for chrome to finish
	err = c.Wait()
	if err != nil {
		return "", err
	}

	fmt.Println(res)

	return res[0], nil
}*/

// bow is a headless browser with some browser behave like web browser, includes: cookie, submit forms,
// but JavaScript NOT supported
func HeadlessGetBow(url string) (string, error) {
	bow := surf.NewBrowser()
	err := bow.Open("http://www.yahoo.com")
	if err != nil {
		return "", err
	}
	return bow.Dom().Html()
}
