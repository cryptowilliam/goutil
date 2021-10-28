package gnum

import "testing"

func TestParseFloat64(t *testing.T) {

	f, err := ParseFloat64("2.07539829195e-05")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(f)
}

func TestIsDigit(t *testing.T) {
	ss := []string{"123", "00123", "+123", "-123", "123.456", "2.07539829195e-05"}
	for _, s := range ss {
		if !IsDigit(s) {
			t.Errorf("%s is a number", s)
			return
		}
	}

	sserr := []string{"a123", "00-123", "10+123", "2.07.539829195e-05"}
	for _, s := range sserr {
		if IsDigit(s) {
			t.Errorf("%s is NOT a number", s)
			return
		}
	}
}
