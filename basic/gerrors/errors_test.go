package gerrors

import (
	pkgErr "github.com/pkg/errors"
	"testing"
)

func TestErrDetail_Details(t *testing.T) {
	gerr := New("test error %s", "abc")
	if gerr.Details().Stack == "" {
		t.Failed()
		return
	}
}

func TestErrDetail_Wrap(t *testing.T) {
	pkgErr := pkgErr.Errorf("test error %s", "abc")
	gerr := Wrap(pkgErr, "")
	if !IsGerror(gerr) {
		t.Failed()
		return
	}
}
