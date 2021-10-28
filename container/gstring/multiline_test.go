package gstring

import "testing"

func TestGetLastLines(t *testing.T) {
	src := "abc\ndef\nghi\njkl"
	result := GetLastLines(src, 2)
	if result != "ghi\njkl" {
		t.Error("GetLastLines failed")
	}
}
