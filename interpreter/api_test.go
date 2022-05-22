package interpreter

import (
	"github.com/cryptowilliam/goutil/basic/gtest"
	"testing"
)

type (
	ctx struct {}

	bar struct {Field int}
)

func (f *ctx) Bar(field int) *bar {
	return &bar{Field: field}
}

func (b *bar) Double() int {
	return b.Field * 2
}

func TestVm_GetScriptFunc(t *testing.T) {
	testData := map[string]string{}

	testData["goja"] = `
	function onReply(a, b) {
		return ctx.Bar(a + b).Double();
	}`

	testData["yaegi"] = `
	func onReply(a int, b int) int {
		return ctx.Bar(a + b).Double()
	}`

	for engine, script := range testData {
		vm, err := NewVM(engine)
		if err != nil {
			gtest.PrintlnExit(t, err.Error())
		}

		// set context
		err = vm.MapGoValueToScript("ctx", &ctx{})
		if err != nil {
			gtest.PrintlnExit(t, err.Error())
		}

		// load javascript
		_, err = vm.RunScript(script)
		if err != nil {
			gtest.PrintlnExit(t, "RunScript for engine %s error: %s", engine, err.Error())
		}

		// register javascript function
		onReply, err := vm.MapScriptFuncToGo("onReply")
		if err != nil {
			gtest.PrintlnExit(t, err.Error())
		}

		// call javascript function
		res, err := onReply(40, 2)
		if err != nil {
			gtest.PrintlnExit(t, err.Error())
		}

		numStr := res[0].String()
		if numStr != "84" {
			gtest.PrintlnExit(t, "engine %s result should be 84 but not %s", engine, numStr)
		}
	}
}
