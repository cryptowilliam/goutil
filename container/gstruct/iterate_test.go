package gstruct

import (
	"fmt"
	"github.com/cryptowilliam/goutil/basic/gtest"
	"reflect"
	"testing"
)

func TestIterate(t *testing.T) {
	type (
		In struct {
			S1 string `crypt:"true"`
			S2 string `crypt:"false"`
			I1 int    `crypt:"true"`
			I2 int    `crypt:"false"`
		}

		Out struct {
			In In
		}
	)

	s := Out{
		In{
			S1: "s1",
			S2: "s2",
			I1: 1,
			I2: 2,
		},
	}

	modifyFn := func(v reflect.Value) (newVal reflect.Value, modified bool, err error) {
		if v.Kind() == reflect.String {
			return reflect.ValueOf("NEWS1"), true, nil
		}
		fmt.Println(v.Kind(), v)
		return reflect.ValueOf(0), false, nil
	}
	err := Iterate(&s, "crypt", "true", modifyFn)
	gtest.Assert(t, err)
	fmt.Println(s)
}
