package gtime

import (
	"github.com/cryptowilliam/goutil/basic/gtest"
	"testing"
)

func TestParseYearMonthInt(t *testing.T) {
	cl := gtest.NewCaseList()
	cl.New().Input(200007).Expect("200007")
	cl.New().Input(-200007).Expect("-200007")

	for _, v := range cl.Get() {
		output, err := ParseYearMonthInt(v.Inputs[0].(int))
		if err != nil {
			t.Error(err)
			return
		}
		if output.StringYYYYMM() != v.Expects[0].(string) {
			t.Errorf("output %s, expected output %s", output, v.Expects[0].(string))
			return
		}
	}
}
