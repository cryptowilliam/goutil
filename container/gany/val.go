package gany

import (
	"fmt"
	"reflect"
)

type (
	// TODO change to reflect.Value
	Val reflect.Value

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
	return Val(reflect.ValueOf(v))
}

func (v Val) Value() reflect.Value {
	return (reflect.Value)(v)
}

func (v Val) Any() any {
	return (reflect.Value)(v).Interface()
}
/*
func (v Val) ToByte() byte {

}

func (v Val) ToRune() rune {

}

func (v Val) ToUint() uint {

}

func (v Val) ToUint8() uint8 {

}

func (v Val) ToUint16() uint16 {

}

func (v Val) ToUint32() uint32 {

}

func (v Val) ToUint64() uint64 {

}

func (v Val) ToInt() int {

}

func (v Val) ToInt8() int8 {

}

func (v Val) ToInt16() int16 {

}

func (v Val) ToInt32() int32 {
}

func (v Val) ToInt64() int64 {

}

func (v Val) ToFloat32() float32 {

}

func (v Val) ToFloat64() float64 {

}

func (v Val) ToBigInt() big.Int {

}

func (v Val) ToBigFloat() big.Float {

}

func (v Val) ToBool() bool {

}

func (v Val) ToBytes() []byte {

}

func (v Val) ToString() string {

}

func (v Val) Equals(Val) bool {

}

func (v Val) StrictEquals(Val) bool {

}


*/

// FIXME
func (v Val) String() string {
	switch v.TypeName() {
	case "int":
		return fmt.Sprintf("%d", v.Any().(int))
	case "int64":
		return fmt.Sprintf("%d", v.Any().(int64))
	default:
		return fmt.Sprintf("unsupported type %s", v.TypeName())
	}
}

func (v Val) TypeReflect() reflect.Type {
	return reflect.TypeOf(v.Any())
}

func (v Val) TypeName() string {
	return Type(v.Any())
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

