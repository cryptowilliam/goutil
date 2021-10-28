package ginterface

import (
	"encoding/json"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/encoding/gjson"
	"log"
	"math/big"
	"reflect"
	"time"
)

var (
	// Any initialized interface's type is no longer interface, but the type of init value,
	// var itfc interface{}
	// Now type of 'itfc' is 'nil'
	// var itfc interface{} = interface{}("str")
	// Now type of 'itfc' becomes string, not type of uninitialized interface anymore.
	UnInitIF interface{} = nil
)

func Parse(x interface{}) (typeName string, isPointer bool) {
	if t := reflect.TypeOf(x); t.Kind() == reflect.Ptr {
		return t.Elem().Name(), true
	} else {
		return t.Name(), false
	}
}

func Type(x interface{}) string {
	if x == nil || reflect.TypeOf(x) == nil {
		return "nil" // 'nil' means uninitialized interface{}.
	}

	return reflect.TypeOf(x).String()
}

// Try to detect slice member type if x is a slice.
func TypeEx(x interface{}) (string, error) {
	// uninitialized interface{}
	if x == nil || reflect.TypeOf(x) == nil {
		return "nil", nil // 'nil' means uninitialized interface{}.
	}

	s := reflect.TypeOf(x).String()

	// not slice
	if reflect.TypeOf(x).Kind() != reflect.Slice {
		return s, nil
	}

	// Slice member type already detected.
	if s /*like '[]string'*/ != reflect.TypeOf([]interface{}{}).String() /*'[]interface {}'*/ {
		return s, nil
	}

	// Detect slice member type with json.Unmarshal.
	buf, err := gjson.MarshalBuffer(x, false)
	if err != nil {
		return "", err
	}
	var dst interface{}
	err = json.Unmarshal(buf, &dst)
	if err != nil {
		return "", err
	}

	// if x is uninitialized interface slice, dst will be nil
	if dst == nil {
		return s, nil
	}

	var sliceMemType *string = nil
	for _, v := range dst.([]interface{}) {
		if sliceMemType == nil {
			smt := reflect.TypeOf(v).String()
			sliceMemType = &smt
		} else {
			if reflect.TypeOf(v).String() != *sliceMemType {
				return s /*[]interface {}*/, nil
			}
		}
	}

	if sliceMemType != nil {
		return "[]" + *sliceMemType, nil
	}else{
		return s, nil
	}
}

func IsSlice(x interface{}) bool {
	return reflect.TypeOf(x) != nil && reflect.TypeOf(x).Kind() == reflect.Slice
}

func IsMap(x interface{}) bool {
	return reflect.TypeOf(x) != nil && reflect.TypeOf(x).Kind() == reflect.Map
}

func IsIntNum(x interface{}) bool {
	switch Type(x) {
	case "bool", "rune", "uint", "uint8", "uint16", "uint32", "uint64", "int", "int8", "int16", "int32", "int64", "big.Int":
		return true
	default:
		return false
	}
}

func ConvertToIntNum(x interface{}) (*big.Int, error) {
	res := &big.Int{}
	switch Type(x) {
	case "bool":
		if x.(bool) {
			res.SetInt64(1)
		} else {
			res.SetInt64(0)
		}
	case "rune":
		res.SetUint64(uint64(x.(rune)))
	case "uint8":
		res.SetUint64(uint64(x.(uint8)))
	case "uint16":
		res.SetUint64(uint64(x.(uint16)))
	case "uint32":
		res.SetUint64(uint64(x.(uint32)))
	case "uint64":
		res.SetUint64(x.(uint64))
	case "int8":
		res.SetInt64(int64(x.(uint8)))
	case "int16":
		res.SetInt64(int64(x.(uint16)))
	case "int32":
		res.SetInt64(int64(x.(uint32)))
	case "int64":
		res.SetInt64(int64(x.(uint64)))
	case "int":
		res.SetInt64(int64(x.(int)))
	default:
		res = nil
	}
	if res == nil {
		return nil, gerrors.New("x type %s is not a integer number", Type(x))
	}
	return res, nil
}

func IsFloatNum(x interface{}) bool {
	switch Type(x) {
	case "float32", "float64", "big.Float":
		return true
	default:
		return false
	}
}

func ConvertToFloatNum(x interface{}) (*big.Float, error) {
	res := &big.Float{}
	switch Type(x) {
	case "float32":
		res.SetInt64(int64(x.(uint32)))
	case "float64":
		res.SetInt64(int64(x.(uint64)))
	default:
		res = nil
	}
	if res == nil {
		return nil, gerrors.New("x type %s is not a float number", Type(x))
	}
	return res, nil
}

//   -1 if x <  y
//    0 if x == y
//   +1 if x >  y
func Cmp(x, y interface{}) (int, error) {
	if Type(x) != Type(y) {
		return 0, gerrors.New("x type %s != y type %s", Type(x), Type(y))
	}

	if IsIntNum(x) {
		xNum, err := ConvertToIntNum(x)
		if err != nil {
			return 0, err
		}
		yNum, err := ConvertToIntNum(y)
		if err != nil {
			return 0, err
		}
		return xNum.Cmp(yNum), nil
	}

	if IsFloatNum(x) {
		xNum, err := ConvertToFloatNum(x)
		if err != nil {
			return 0, err
		}
		yNum, err := ConvertToFloatNum(y)
		if err != nil {
			return 0, err
		}
		return xNum.Cmp(yNum), nil
	}

	if Type(x) == "time.Time" {
		xTime := x.(time.Time)
		yTime := y.(time.Time)
		if xTime.Before(yTime) {
			return -1, nil
		}
		if xTime.Equal(yTime) {
			return 0, nil
		}
		if xTime.After(yTime) {
			return 1, nil
		}
	}

	if Type(x) == "string" {
		xString := x.(string)
		yString := y.(string)
		if xString < yString {
			return -1, nil
		}
		if xString == yString {
			return 0, nil
		}
		if xString > yString {
			return 1, nil
		}
	}

	return 0, gerrors.New("ginterface.Cmp doesn't support %s", Type(x))
}

func TypeEqualValueEqual(x, y interface{}) bool {
	res, err := Cmp(x, y)
	if err != nil {
		return false
	}
	return res == 0
}

func TypeEqualValueLess(x, y interface{}) bool {
	res, err := Cmp(x, y)
	if err != nil {
		return false
	}
	return res < 0
}

func TypeEqualValueLessOrEqual(x, y interface{}) bool {
	res, err := Cmp(x, y)
	if err != nil {
		return false
	}
	return res <= 0
}

func TypeEqualValueGreater(x, y interface{}) bool {
	res, err := Cmp(x, y)
	if err != nil {
		return false
	}
	return res > 0
}

func TypeEqualValueGreaterOrEqual(x, y interface{}) bool {
	res, err := Cmp(x, y)
	if err != nil {
		return false
	}
	return res >= 0
}

type CommonFunc struct{}

var cf CommonFunc

func (c *CommonFunc) Merge2(s ...[]interface{}) (result []interface{}) {
	switch len(s) {
	case 0:
		break
	case 1:
		result = s[0]
		break
	default:
		s1 := s[0]
		s2 := cf.Merge2(s[1:]...) //...将数组元素打散
		result = make([]interface{}, len(s1)+len(s2))
		copy(result, s1)
		copy(result[len(s1):], s2)
		break
	}

	return
}

func Merge(s ...[]interface{}) (result []interface{}) {
	switch len(s) {
	case 0:
		break
	case 1:
		result = s[0]
		break
	default:
		s1 := s[0]
		s2 := Merge(s[1:]...) // ...可以将数组元素打散
		result = make([]interface{}, len(s1)+len(s2))
		copy(result, s1)
		copy(result[len(s1):], s2)
		break
	}

	return result
}

/**
  @retry  重试次数
  @method 调用的函数，比如: api.GetTicker ,注意：不是api.GetTicker(...)
  @params 参数,顺序一定要按照实际调用函数入参顺序一样
  @return 返回
*/
func ReCallItfc(retry int, method interface{}, params ...interface{}) interface{} {

	invokeM := reflect.ValueOf(method)
	if invokeM.Kind() != reflect.Func {
		panic("method not a function")
		return nil
	}

	var value = make([]reflect.Value, len(params))
	var i = 0
	for ; i < len(params); i++ {
		value[i] = reflect.ValueOf(params[i])
	}

	var retV interface{}
	var retryC = 0
_CALL:
	if retryC > 0 {
		log.Println("sleep....", time.Duration(retryC*200*int(time.Millisecond)))
		time.Sleep(time.Duration(retryC * 200 * int(time.Millisecond)))
	}

	retValues := invokeM.Call(value)

	for _, vl := range retValues {
		if vl.Type().String() == "error" {
			if !vl.IsNil() {
				log.Println(vl)
				retryC++
				if retryC <= retry {
					log.Printf("Invoke Method(%s) Error , Begin Retry Call [%d] ...", invokeM.String(), retryC)
					goto _CALL
				} else {
					panic("Invoke Method Fail ???" + invokeM.String())
				}
			}
		} else {
			retV = vl.Interface()
		}
	}

	return retV
}
