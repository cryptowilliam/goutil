package interpreter

import (
	"github.com/cryptowilliam/goutil/container/gany"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"reflect"
)

type (
	VmYaegj struct {
		vmYaegi *interp.Interpreter
		customLib map[string]map[string]reflect.Value
	}
)

func valYaegi2Comm(yaegiVal reflect.Value) gany.Val {
	return gany.NewVal(yaegiVal)
}

func newVMYaegi() (*VmYaegj, error) {
	res := &VmYaegj{vmYaegi: interp.New(interp.Options{})}
	res.customLib = make(map[string]map[string]reflect.Value)
	res.customLib["custom/custom"] = make(map[string]reflect.Value)

	if err := res.vmYaegi.Use(stdlib.Symbols); err != nil {
		return nil, err
	}

	return res, nil
}

func (vm *VmYaegj) RunScript(script string) (gany.Val, error) {
	val, err := vm.vmYaegi.Eval(script)
	if err != nil {
		return gany.ValNil, err
	}
	return valYaegi2Comm(val), nil
}

// TODO 确定重复调用这个接口，比如映射两个不同的name、value对， 是可以正常执行的
// 引入的指针和相关的自定义类型（非基础类型）都可以存取，但是默认不可以声明
// 如果需要在脚本中主动声明Go Runtime中的自定义类型，目前知道的方法是在脚本中import该类型所在的包，
// MapGoValueToScript allows script to access go runtime `value` with `name`.
// Usually `value` is a pointer in the go runtime.
func (vm *VmYaegj) MapGoValueToScript(name string, value interface{}) error {
	vm.customLib["custom/custom"][name] = reflect.ValueOf(value)
	if err := vm.vmYaegi.Use(vm.customLib); err != nil {
		return err
	}

	vm.customLib["custom/custom"]["ctx"] = reflect.ValueOf(value)
	if err := vm.vmYaegi.Use(vm.customLib); err != nil {
		return err
	}

	if _, err := vm.vmYaegi.Eval(`import . "custom"`); err != nil {
		return err
	}
	return nil
}

// MapScriptFuncToGo allows go runtime to access script function.
func (vm *VmYaegj) MapScriptFuncToGo(funcName string) (Callable, error) {
	fn, err := vm.vmYaegi.Eval(funcName)
	if err != nil {
		return nil, err
	}

	return func(args ...any) ([]gany.Val, error) {
		var argsSlice []reflect.Value
		for _, item := range args {
			argsSlice = append(argsSlice, reflect.ValueOf(item))
		}

		resVals := fn.Call(argsSlice)
		var res []gany.Val
		for _, item := range resVals {
			res = append(res, gany.NewVal(item))
		}
		return res, nil
	}, nil
}