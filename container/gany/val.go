package gany

import (
	"fmt"
	"reflect"
)

type (
	Val struct {any}

	TypeInfo struct {
		TypeName string
		IsNil bool
		IsPointer bool
		IsSlice bool
		IsMap bool
	}
)

var (
	ValNil = NewVal(nil)
)

func NewVal(v any) Val {
	// if type of `v` is reflect.Value, and cast it to any with `Val{v}`,
	// it's type and value info will be lost.
	if reflect.TypeOf(v) == reflect.TypeOf(reflect.Value{}) {
		return Val{v.(reflect.Value).Interface()}
	} else {
		return Val{v}
	}
}

func (v Val) Value() reflect.Value {
	return reflect.ValueOf(v.Any())
}

func (v Val) Any() any {
	return v.any
}

func (v Val) IsNil() bool {
	return v.Any() == nil
}

func (v Val) IsPointer() bool {
	t := reflect.TypeOf(v.Any())
	return t.Kind() == reflect.Ptr
}

func (v Val) IsSlice() bool {
	return IsSlice(v.Any())
}

func (v Val) IsMap() bool {
	return IsMap(v.Any())
}

func (v Val) TypeName() string {
	if v.IsNil() {
		return "nil" // 'nil' means uninitialized interface{}.
	}

	return reflect.TypeOf(v.Any()).String()
}

func (v Val) TypeReflect() reflect.Type {
	return reflect.TypeOf(v.Any())
}

func (v Val) TypeInfo() TypeInfo {
	res := TypeInfo{}

	res.IsNil = v.Any() == nil
	if t := reflect.TypeOf(v.Any()); t.Kind() == reflect.Ptr {
		res.TypeName = t.Elem().Name()
		res.IsPointer = true
	} else {
		res.TypeName = t.Name()
		res.IsPointer = false
	}
	res.IsMap = IsMap(v.Any())
	res.IsSlice = IsSlice(v.Any())

	return res
}
/*
func (v Val) CastToByte() byte {
	return v.any.(byte)
}

func (v Val) CastToRune() rune {
	return v.any.(rune)
}

func (v Val) CastToUint() uint {
	return v.any.(uint)
}

func (v Val) CastToUint8() uint8 {
	return v.any.(uint8)
}

func (v Val) CastToUint16() uint16 {
	return v.any.(uint16)
}

func (v Val) CastToUint32() uint32 {
	return v.any.(uint32)
}

func (v Val) CastToUint64() uint64 {
	return v.any.(uint64)
}

func (v Val) CastToInt() int {
	return v.any.(int)
}

func (v Val) CastToInt8() int8 {
	return v.any.(int8)
}

func (v Val) CastToInt16() int16 {
	return v.any.(int16)
}

func (v Val) CastToInt32() int32 {
}

func (v Val) CastToInt64() int64 {

}

func (v Val) CastToFloat32() float32 {

}

func (v Val) CastToFloat64() float64 {

}

func (v Val) CastToBigInt() big.Int {

}

func (v Val) CastToBigFloat() big.Float {

}

func (v Val) CastToBool() bool {

}

func (v Val) CastToBytes() []byte {

}

func (v Val) CastToString() string {

}

func (v Val) Equals(Val) bool {

}

func (v Val) StrictEquals(v2 Val) bool {
	return v.TypeName() == v2.TypeName() &&
		v.
}*/


func (v Val) String() string {
	return fmt.Sprintf("%v", v.Any())
}


