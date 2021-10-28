package ghtml

import (
	"github.com/PuerkitoBio/goquery"
	"strings"
)

func NewDocFromHtmlSrc(htmlSrc *string) (*goquery.Document, error) {
	sr := strings.NewReader(*htmlSrc)
	doc, err := goquery.NewDocumentFromReader(sr)
	if err != nil {
		return nil, err
	}
	return doc, nil
}
