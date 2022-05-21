package interpreter

import (
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/container/gany"
)

type (
	// Callable represents a JavaScript function that can be called from Go.
	Callable func(args ...any) ([]any, error)

	Vm interface {
		RunScript(script string) (gany.Val, error)
		MapGoValueToScript(name string, value interface{}) error
		MapScriptFuncToGo(funcName string) (Callable, error)
	}
)

func NewVM(engine string) (Vm, error) {
	switch engine {
	case "goja":
		return newVMGoja()
	case "yaegi":
		return newVMYaegi()
	default:
		return nil, gerrors.New("unsupported interpreter engine %s", engine)
	}
}
