package gany

import (
	"fmt"
	"reflect"
	"testing"
)

func TestVal_TypeName(t *testing.T) {
	fmt.Println(NewVal(1).TypeName())
	fmt.Println(NewVal(nil).TypeName())
	fmt.Println(NewVal(nil).String())
}

func TestNewVal(t *testing.T) {
	n := 1
	refVal := reflect.ValueOf(n)
	fmt.Println(refVal.Kind())

	val := NewVal(refVal)
	fmt.Println(val.Value().Interface())
	fmt.Println("String: " + val.String())
}
