package translate

import (
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/zgs225/youdao"
	"regexp"
	"strings"
	"time"
)

const (
	DefaultYoudaoAppId     = "2f871f8481e49b4c"
	DefaultYoudaoAppSecret = "CQFItxl9hPXuQuVcQa5F2iPmZSbN0hYS"
)

func GoogleTranslateByWeb(SourceLang, TargetLang, Text, Proxy string, timeout *time.Duration) (string, error) {

	resp, err := NewClient(SourceLang, TargetLang).Translate(Text, Proxy, timeout).Get().Do()
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", gerrors.Errorf("Failed, HTTP status %d", resp.StatusCode)
	}
	respHtml := string(resp.ResponseBody)

	re := regexp.MustCompile(`class="t0">(.*?)<`)
	match := re.FindStringSubmatch(respHtml)
	if len(match) != 2 {
		return "", gerrors.New("Failed to translate")
	}

	translated := strings.Replace(match[1], "&quot;", "", -1)

	// Google会把http://中的引号翻译成中文全角引号，下面把它恢复
	translated = strings.Replace(translated, "http：//", "http://", -1)
	translated = strings.Replace(translated, "https：//", "https://", -1)

	// 人工修正不恰当的翻译
	/*if gstring.StartWith(translated, "RT ") {
		translated = "转推 " + translated[3:]
	}
	translated = strings.Replace(translated, "分散应用程序", "分布式应用程序", -1)
	translated = strings.Replace(translated, "街区之外", "区块之外", -1)*/
	return translated, nil
}

func YoudaoTranslateEnCnAuto(appId, appSecret, s string) (string, error) {
	c := &youdao.Client{
		AppID:     appId,
		AppSecret: appSecret,
	}
	r, err := c.Query(s)
	if err != nil {
		return "", err
	}
	return (*r.Translation)[0], nil
}
