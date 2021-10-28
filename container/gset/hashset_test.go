package gset

import (
	"testing"
)

func TestNewHashSet4String(t *testing.T) {
	s := NewHashSet()
	s.Add("abc")
	if !s.Contains("abc") {
		t.Error("Contains failed")
	}
	if s.Contains("abcd") {
		t.Error("Contains failed")
	}
}
