package interpreter

import (
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/container/gany"
)

type (
	// Callable represents a JavaScript function that can be called from Go.
	Callable func(args ...any) ([]gany.Val, error)

	// Vm is script interpreter.
	Vm interface {
		// RunScript runs script.
		RunScript(script string) (gany.Val, error)

		// MapGoValueToScript allows script to access go runtime `value` with `name`.
		// Usually `value` is a pointer in the go runtime.
		MapGoValueToScript(name string, value interface{}) error

		// MapScriptFuncToGo allows go runtime to access script function.
		MapScriptFuncToGo(funcName string) (Callable, error)
	}
)

// NewVM creates interpreter.
// It supports Golang script and ECMAScript languages like Javascript, TypeScripts.
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
