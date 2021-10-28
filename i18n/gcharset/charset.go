package gcharset

import (
	"bytes"
	"github.com/richardlehane/characterize"
	"github.com/saintfish/chardet"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
)

func DetectCharset(str string, isHtml bool) (string, error) {
	var d *chardet.Detector
	if isHtml {
		d = chardet.NewHtmlDetector()
	} else {
		d = chardet.NewTextDetector()
	}
	r, e := d.DetectBest([]byte(str))
	if e != nil {
		return "", e
	} else {
		return r.Charset, nil
	}
}

// characterize seems great
func DetectCharsetEx(in []byte) (string, error) {
	ct := characterize.Detect(in)
	return ct.String(), nil
}

/*
// has cross compile problem
func ConvCharset(src string, fromCharset string, toCharset string) (string, error) {
	fromCharset = strings.ToLower(fromCharset)
	toCharset = strings.ToLower(toCharset)
	if fromCharset == "gb-18030" {
		fromCharset = "gb18030"
	}
	if toCharset == "gb-18030" {
		toCharset = "gb18030"
	}
	if fromCharset == toCharset {
		return src, nil
	}
	rst, e2 := iconv.ConvertString(src, fromCharset, toCharset)
	return rst, e2
}

func ConvCharsetAuto(src string, isHtml bool, toCharset string) (string, error) {
	fromCharset, e := DetectCharset(src, isHtml)
	if e != nil {
		return "", e
	}
	return ConvCharset(src, fromCharset, toCharset)
}*/

func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func Utf8ToGbk(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}
