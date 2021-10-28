package facebook

import (
	"fmt"
	"github.com/cryptowilliam/goutil/container/gstring"
	"github.com/cryptowilliam/goutil/net/gaddr"
	"strings"
)

type UrlType string

const (
	UrlTypeNotFacebook UrlType = "notFacebook"
	UrlTypeUnknown     UrlType = "unknown"
	UrlTypeHome        UrlType = "home"
	UrlTypeAccount     UrlType = "account"
	UrlTypePost        UrlType = "post"
)

func ParseAccountInUrl(urlstr string) string {
	us, err := gaddr.ParseUrl(urlstr)
	if err != nil || strings.ToLower(us.Domain.SiteDomain) != "facebook.com" {
		fmt.Println(err)
		return ""
	}
	if len(us.Path.Dirs) == 0 {
		return ""
	}
	return us.Path.Dirs[0]
}

// https://www.facebook.com/pages/Gaza-City/107601032603195?__xts__%5B0%5D=68.ARBlLNlHuQoI3d3kal378TdFDPzLWTXtmCzZ8fvLJDhuYRhPueFv41ZbDwODUB9UjknyvB0Tx1touVXErBmiWkILMeYGjSqsFWU7hLHWFG0e21XYPKkAB7b2oXfV_TNEO3AXzNynaxdfA-BQkLq5t4tyYWHRLAGlpNQuhe6cDnwgADlsOuGgqvUcRgny850vffNzHpMnJuNJ2lATjuaNL4EiUig9u4_On5j37X1YLSOM01cnntuNydB_9anjR3q2psrCb_2OZHZ-iRZwnK6p1Pkl5N58j9OVAEmMFM2txVJ8pX2WMUfnEzJk1Ch0ptCGL7RChdwUFZXq3p6kp87ehIw3HX8xH7Cn5GI&__tn__=-R
func ParseUrl(urlstr string) (ut UrlType, accountName string) {
	if gstring.StartWith(urlstr, "https://www.facebook.com/pages/") {
		return UrlTypeUnknown, ""
	}
	if gstring.StartWith(urlstr, "https://www.facebook.com/photo.php") {
		return UrlTypePost, ""
	}

	us, err := gaddr.ParseUrl(urlstr)
	if err != nil || strings.ToLower(us.Domain.SiteDomain) != "facebook.com" {
		fmt.Println(err)
		return UrlTypeUnknown, ""
	}
	if len(us.Path.Dirs) == 0 {
		return UrlTypeHome, ""
	}
	if len(us.Path.Dirs) == 1 {
		return UrlTypeAccount, us.Path.Dirs[0]
	}
	return UrlTypePost, ""
}

// change https://www.facebook.com/BBCChinese into https://www.facebook.com/pg/BBCChinese/posts
func GetCrawlerFriendlyAccountAliasURL(rawUrl string) (newUrl string) {
	t, accountName := ParseUrl(rawUrl)

	if t == UrlTypeAccount {
		return fmt.Sprintf("https://www.facebook.com/pg/%s/posts", accountName)
	}

	return rawUrl
}

// change https://www.facebook.com/washingtonpost/posts/10158279295492293?__xts__%B0%D=68.ARD3s7zDKJU1jcn0-IG5HFLCPQ3VQy6uRylPCHHYWUsyXCJy3nYdv1Z_r4xqeh8h--Cn6oycJf5AK8_Fl3UAw2xODedWwsjJOOcPY_3rCI1zCWGlf-lz3a6KrJOxrZGiktib-Oi1VqDgYEZk54DfYzH6VN8X8D84y34IGlVmObwN12aevk--kIzKTZ3JRo8paQSWC92BQOnGsASTT7FD292P_K7dTUcNq3y1VtFI-xOCHpeK6-2S4bISnJs3H-u-F3STKxFcd3HxY5E-nsK7_8A9en4dxwu_-rQkpmohtV1XI8V5z1pHvI_Z8R5D-yfKyRuadKT5jdRDHr4&__tn__=-R
// into
// https://www.facebook.com/washingtonpost/posts/10158279295492293
func FixFacebookPostURL(rawUrl string) (newUrl string) {
	t, _ := ParseUrl(rawUrl)

	if t == UrlTypePost {
		if strings.Contains(rawUrl, "facebook.com/photo.php?fbid=") {
			ss := strings.Split(rawUrl, "&")
			newUrl = ss[0]
			for gstring.EndWith(newUrl, "/") {
				newUrl = gstring.RemoveTail(newUrl, 1)
			}
			return newUrl
		}
		if strings.Contains(rawUrl, "?") {
			ss := strings.Split(rawUrl, "?")
			newUrl = ss[0]
			for gstring.EndWith(newUrl, "/") {
				newUrl = gstring.RemoveTail(newUrl, 1)
			}
			return newUrl
		}
	}

	return rawUrl
}
