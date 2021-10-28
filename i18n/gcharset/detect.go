package gcharset

// github.com/chrisport/go-lang-detector

import (
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/johngb/langreg" // ISO-639-1 Code
	"github.com/taruti/langdetect"
)

func DetectNatrualLanguage(utf8Slice []byte) (lang langdetect.Language, err error) {
	if len(utf8Slice) == 0 {
		return langdetect.Language{}, gerrors.Errorf("Empty input")
	}
	l := langdetect.DetectLanguage(utf8Slice, "")
	return l, nil
}

func IsChinese(s string) (bool, error) {
	lang, err := DetectNatrualLanguage([]byte(s))
	if err != nil {
		return false, err
	}
	return lang == langdetect.Zh, nil
}

func IsValidISO_693_1_Code(code string) bool {
	return langreg.IsValidLanguageCode(code)
}
