package nlp

// client library: https://cloud.google.com/natural-language/docs/reference/libraries
// offline library: https://github.com/cdipaolo/sentiment

import (
	"bytes"
	"encoding/json"
	"fmt"
	gs "github.com/bitly/go-simplejson"
	"github.com/cdipaolo/sentiment"
	"github.com/cryptowilliam/goutil/net/ghttp"
	"github.com/taruti/langdetect"
	"io/ioutil"
	"net/http"
	"time"
)

var url = "https://language.googleapis.com/v1beta2/documents:analyzeSentiment"

type document struct {
	Type     string `json:"type"`
	Language string `json:"language"`
	Content  string `json:"content"`
}

type requestBody struct {
	Document     document `json:"document"`
	EncodingType string   `json:"encodingType"`
}

func newRequestBody(text, lang string) *requestBody {
	return &requestBody{
		Document: document{
			Type:     "PLAIN_TEXT",
			Language: lang,
			Content:  text,
		},
		EncodingType: "UTF8",
	}
}

func AnalyzeSentimentONLINE(text string, language langdetect.Language, proxyUrl *string) (score float64, err error) {
	key := "AIzaSyArmOhGM6MOGJAyRrrCvDmhHD2dPQ6oPsE"

	// Ready for request
	client := &http.Client{}
	if proxyUrl != nil && len(*proxyUrl) > 0 {
		err = ghttp.SetProxy(client, *proxyUrl)
		if err != nil {
			return 0, err
		}
	}
	to := time.Second * 5
	if err := ghttp.SetTimeout(client, &to, nil, nil, nil, nil); err != nil {
		return 0, err
	}
	body, err := json.Marshal(newRequestBody(text, language.String()))
	if err != nil {
		return 0, err
	}

	// SendEth request google cloud service
	req, err := http.NewRequest("POST", fmt.Sprintf("%s?key=%s", url, key), bytes.NewBuffer(body))
	if err != nil {
		return 0, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	// Parse response
	j, err := gs.NewJson(result)
	if err != nil {
		return 0, err
	}
	res := j.Get("documentSentiment").Get("score")
	return res.Float64()
}

func AnalyzeSentiment(text string, language langdetect.Language) (score float64, err error) {
	model, err := sentiment.Restore()
	if err != nil {
		return 0, err
	}
	analysis := model.SentimentAnalysis(text, sentiment.English)
	total := float64(0)
	for _, v := range analysis.Words {
		total += float64(v.Score)
	}
	avg := total / float64(len(analysis.Words))
	return avg, nil
}
