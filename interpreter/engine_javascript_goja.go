package interpreter

import (
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/container/gany"
	"github.com/dop251/goja"
)

type (
	VmGoja struct {
		vmGoja *goja.Runtime
	}
)

func valGoja2Comm(gojaVal goja.Value) any {
	return gojaVal.Export()
}

func newVMGoja() (*VmGoja, error) {
	res := &VmGoja{vmGoja: goja.New()}
	return res, nil
}

func (vm *VmGoja) toGojaValue(anyValue any) goja.Value {
	return vm.vmGoja.ToValue(anyValue)
}

func (vm *VmGoja) RunScript(script string) (gany.Val, error) {
	val, err := vm.vmGoja.RunString(script)
	if err != nil {
		return gany.ValNil, err
	}
	return gany.NewVal(val.Export()), nil
}

func (vm *VmGoja) MapGoValueToScript(name string, value interface{}) error {
	return vm.vmGoja.Set(name, value)
}

func (vm *VmGoja) MapScriptFuncToGo(funcName string) (Callable, error) {
	callable, ok := goja.AssertFunction(vm.vmGoja.Get(funcName))
	if !ok {
		return nil, gerrors.New("%s is not a valid function", funcName)
	}

	return func(args ...any) ([]gany.Val, error) {
		var items []goja.Value
		for _, item := range args {
			items = append(items, vm.toGojaValue(item))
		}
		retGoja, err := callable(goja.Undefined(), items...)
		if err != nil {
			return nil, err
		}
		return []gany.Val{gany.NewVal(retGoja.Export())}, nil
	}, nil
}