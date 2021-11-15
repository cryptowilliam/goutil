package gerrors

// Reference
// https://github.com/aletheia7/gerrors/blob/master/e.go

import (
	stdErr "errors"
	"fmt"
	extErr "github.com/go-errors/errors"
	pkgErr "github.com/pkg/errors"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

type (
	// Implementation structure of gerror.
	GErr struct {
		IsFatal bool
		Num     string
		Msg     string
		Stack   string
	}

	// gerror interface
	gerror interface {
		error
		Details() GErr
		SetErrNum(errNum error)
		SetFatal()
	}
)

// Error number.
var (
	ErrNil            = error(nil)
	ErrNotExist       = stdErr.New("not exist")
	ErrAlreadyExist   = stdErr.New("already exist")
	ErrNotFound       = stdErr.New("not found") // This is not a really run error, it means Database/Collection not exist in mongodb.
	ErrNotSupport     = stdErr.New("not support")
	ErrNotImplemented = stdErr.New("not implemented")
)

// Implements error interface, stack will not display in output.
func (ge *GErr) Error() string {
	errNumEnding := ""
	if ge.Num != "" {
		errNumEnding = "\n"
	}
	/*errMsgEnding := ""
	if ge.Num != "" {
		errMsgEnding = "\n"
	}*/
	return ge.Num + errNumEnding + ge.Msg // + errMsgEnding + ge.Stack
}

// Implements gerror interface.
func (ge *GErr) Details() GErr {
	return *ge
}

// Implements gerror interface.
func (ge *GErr) SetErrNum(errNum error) {
	ge.Num = errNum.Error()
}

// Implements gerror interface.
func (ge *GErr) SetFatal() {
	ge.IsFatal = true
}

func IsGerror(err error) bool {
	return strings.Contains(itfType(err), "GErr")
}

// New error.
func New(format string, args ...interface{}) gerror {
	res := &GErr{
		IsFatal: false,
		Num:     "",
		Msg:     fmt.Sprintf(format, args...),
		Stack:   extErr.Errorf(format, args...).ErrorStack(),
	}
	return res
}

// New error with details.
func NewExt(isFatal bool, errNum string, format string, args ...interface{}) gerror {
	res := &GErr{
		IsFatal: isFatal,
		Num:     errNum,
		Msg:     fmt.Sprintf(format, args...),
		Stack:   extErr.Errorf(format, args...).ErrorStack(),
	}
	return res
}

func Errorf(format string, args ...interface{}) gerror {
	return New(format, args...)
}

// Wrap error to GErr.
func Wrap(err error, message string) gerror {
	if IsGerror(err) {
		return err.(gerror)
	} else {
		return &GErr{
			IsFatal: false,
			Num:     "",
			Msg:     pkgErr.Wrap(err, message).Error(),
			Stack:   GetStack(err),
		}
	}
}

// Join combine multiple errors to one error.
func Join(err error, errs ...error) error {
	errCount := 0
	errJoin := ""
	errLast := error(nil)
	if err != nil {
		errCount++
		errJoin += New("error[%d]:%s;", errCount, err.Error()).Error()
		errLast = err
	}
	for _, v := range errs {
		if v != nil {
			errCount++
			errJoin += New("error[%d]:%s;", errCount, v.Error()).Error()
			errLast = v
		}
	}

	if errCount == 0 {
		return nil
	}
	if errCount == 1 {
		return errLast
	}
	return New(errJoin)
}

func JoinArray(errs []error) error {
	var errsNotNil []error
	for _, v := range errs {
		if v == nil {
			continue
		}
		errsNotNil = append(errsNotNil, v)
	}

	if len(errsNotNil) == 0 {
		return nil
	}
	return Join(errsNotNil[0], errsNotNil[1:]...)
}

func removeFirstLines(src string, count int) string {
	if count <= 0 {
		return src
	}
	sa := strings.Split(src, "\n")
	if len(sa) <= count {
		return ""
	}
	sa = sa[count:]
	return strings.Join(sa, "\n")
}

func endWith(s, toFind string) bool {
	if len(s) == 0 || len(toFind) == 0 {
		return false
	}

	pos := strings.LastIndex(s, toFind)

	// pos < 0: can't find toFind
	// if don't add pos >= 0, there will be a bug if EndWith("astring", "*astring")
	return pos >= 0 && pos == len(s)-len(toFind)
}

func itfType(x interface{}) string {
	if x == nil {
		return "nil"
	}
	return reflect.TypeOf(x).String()
}

func interfaceType(x interface{}) string {
	if x == nil {
		return "nil"
	}
	return reflect.TypeOf(x).String()
}

// Get stack for 3 types of error.
func GetStack(err error) string {
	if err == nil {
		return ""
	}
	stack := ""

	// GErr
	if strings.Contains(interfaceType(err), "GErr") {
		// Stack has been generated.
		return err.(*GErr).Stack
	}

	// pkg/errors
	// Support "github/pkg/errors".
	// Not support standard "github.com/cryptowilliam/goutil/basic/gerrors".
	// Sometimes even pkg/errors used, can't get stack for error created by gerrors.New(). But errors.New works always.
	stack = removeFirstLines(fmt.Sprintf("%+v\n", err), 1)

	// standard error
	if len(stack) == 0 || len(strings.Replace(stack, " ", "", -1)) == 0 {
		// Get call stack.
		for i := 2; i < 10; i++ {
			pc, file, line, _ := runtime.Caller(i) // Caller filename and line number.
			if len(file) == 0 {
				break
			}
			if endWith(file, "runtime.main") {
				break
			}
			f := runtime.FuncForPC(pc) // Caller package name and function name.
			stack = stack + "    -> " + file + ":" + strconv.FormatInt(int64(line), 10) + " " + f.Name()
		}
	}

	return stack
}
