package gcharset

import (
	"fmt"
	"testing"
)

func TestDetectNatrualLanguage(t *testing.T) {
	fmt.Println(DetectNatrualLanguage([]byte("中国")))
}
