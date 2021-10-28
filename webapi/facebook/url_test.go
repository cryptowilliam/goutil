package facebook

import (
	"fmt"
	"testing"
)

func TestFixFacebookPostURL(t *testing.T) {
	s := "https://www.facebook.com/zaobaosg/posts/2325781157493039?__xts__%!B(MISSING)0%!D(MISSING)=68.ARAES5NElT7Xm-z7tdoSU5Mj12lFO3e1uZII6DClAHPbHUv6f7qowNncQIBbFMCAnt4dwISEdquSVdOC3Aik1gL0cWVaXejC2_tRySnDIu7zDMrORuhnNzehbyDOCE5mrX5pfSVVLZAPGGkL-bbSq-orZM409s1SMADK3jdWn8s2QM-znxsbT1C8XQaTeGMW458-rl-o25G1BbqHFPX8IxontiKzL9-mHugFFplPmaqlVTvmvYyP9Vj_ox9tCNYjNwfhnc1oZmouMhg9EO75Srin-Giw9sLkimGdYMATrKusJ0-hiJyyc4PYentoBl4wpo1pW5pN7fyfmwU2bYkB88zgpA&__tn__=-R"
	d := "https://www.facebook.com/zaobaosg/posts/2325781157493039"
	if FixFacebookPostURL(s) != d {
		t.Error("TestFixFacebookPostURL error")
	}

	s = "https://www.facebook.com/photo.php?fbid=2302280916525473&set=a.350429498377301&type=3&__xts__%5B0%5D=68.ARD_A37Id7QWDYWwANJcYkMY2HKo4J_VhBBirPlkNkzYDdkCwjDyPAflu_QSbij2JDigpdsrlw3C1JkS8cHLEU9kLSVvEf2oRCESZKr0IN_6rnXvoCGF0iFQqJy1Ifq0BQFPxjhOk6SiDx4LLHhzanUcRdZnX2ujHgGR9hxOtuKluzl9MGkjmw0wSgzNR2H24lPk1MjxSRVMtUEaXSNnR5xAhAwldFHM4Bb9h6YdciItndm_dYB1CjE5iu3fB3gyC-yK11ydf8kt1hCFFwwKm-6Cz89DcJW7e9YVmVttKpJSa3LQEgWYoEI8k6-5q35iIyyjdJImGGL8eSmaRzNteJAVuA&__tn__=H-R"
	//d = "https://www.facebook.com/zaobaosg/posts/2325781157493039"
	fmt.Println(FixFacebookPostURL(s))
}

func TestParseUrl(t *testing.T) {
	url := "https://www.facebook.com/photo.php?fbid=10114721051227214&set=a.745277455924&type=3&__xts__%!B(MISSING)0%!D(MISSING)=68.ARCU0Qln27IQtoFXxlmSLaCw4lnTCSHyOUO-3km4PinXbxSgJ4QSA-kglZew3Wi_dOpMc53Z59NsIiPSM5knf8myHOc5qBNjy0rGA0xFLujdzBJwqQgRb_Ifl9FlIMRmbSKUzzFGcHUNBldlIhNEjzffp_rWZenYej3WgEBMP132XLLaabSTV7FRyuikd5sKWkMSKf6h1DCPYHZFFk_jIypDbVU2xibMuG6UBDt3JGZnomp6V8brjfNcBET3Euo6aKzbsB4V7WGaFo_swqKpYH3kSrtZGGgpzcRiOYWGoa85dbp_u4F9euq75T8v_V9E3SIZ9s4wZJYLt_SL9Q&__tn__=H-R"
	ut, _ := ParseUrl(url)
	if ut != UrlTypePost {
		t.Errorf("ParseUrl error")
		return
	}
}
