package gerrors

import (
	"github.com/cryptowilliam/goutil/basic/gtest"
	pkgErr "github.com/pkg/errors"
	"strings"
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

func TestJoin(t *testing.T) {
	errNil := error(nil)
	err1 := New("err1")
	err2 := New("err2")
	err3 := New("err3")

	cl := gtest.NewCaseList()

	cl.New().Input(errNil).Expect("")
	cl.New().Input(err1).Expect("err1")
	cl.New().Input(err2).Expect("err2")
	cl.New().Input(err3).Expect("err3")

	cl.New().Input(errNil).Input(errNil).Expect("")
	cl.New().Input(errNil).Input(err1).Expect("err1")
	cl.New().Input(errNil).Input(err2).Expect("err2")
	cl.New().Input(errNil).Input(err3).Expect("err3")
	cl.New().Input(err1).Input(errNil).Expect("err1")
	cl.New().Input(err1).Input(err1).Expect("error[1]:err1;error[2]:err1;")
	cl.New().Input(err1).Input(err2).Expect("error[1]:err1;error[2]:err2;")
	cl.New().Input(err1).Input(err3).Expect("error[1]:err1;error[2]:err3;")

	cl.New().Input(errNil).Input(err1).Input(err2).Expect("error[1]:err1;error[2]:err2;")
	cl.New().Input(errNil).Input(err2).Input(err3).Expect("error[1]:err2;error[2]:err3;")
	cl.New().Input(err1).Input(err2).Input(err3).Expect("error[1]:err1;error[2]:err2;error[3]:err3;")
	cl.New().Input(err3).Input(err2).Input(err1).Expect("error[1]:err3;error[2]:err2;error[3]:err1;")

	for _, v := range cl.Get() {
		if len(v.Inputs) == 0 {
			continue
		}
		var inputErrors []error
		var inputStrings []string
		for _, v := range v.Inputs {
			if v == nil {
				inputErrors = append(inputErrors, error(nil))
				inputStrings = append(inputStrings, "")
			} else {
				inputErrors = append(inputErrors, v.(error))
				inputStrings = append(inputStrings, v.(error).Error())
			}
		}

		errJoint := Join(inputErrors[0], inputErrors[1:]...)
		errJointString := ""
		if errJoint != nil {
			errJointString = errJoint.Error()
		}
		if errJointString != v.Expects[0].(string) {
			gtest.PrintlnExit(t, "inputs:[%s], expect:%s, but %s got", strings.Join(inputStrings, ","), v.Expects[0].(string), errJoint.Error())
		}
	}
}
