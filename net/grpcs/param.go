package grpcs

import (
	"encoding/gob"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/container/gany"
	"github.com/cryptowilliam/goutil/container/gstring"
)

func init() {
	gob.Register(NewRequest())
	gob.Register(NewReply())
}

var (
	In  = true
	Out = false

	TypeBool        = true
	TypeString      = ""
	TypeStringSlice = []string{""}
	TypeByteSlice   = []byte("")
	TypeIF          = gany.UnInitIF
	TypeIFSlice     []interface{}
	StringSlice     []string
)

type (
	Request struct {
		Func   string
		Params map[string]interface{}
	}

	Reply map[string]interface{}

	ParamChecker struct {
		InRequirements  map[string]map[string]interface{}
		OutRequirements map[string]map[string]interface{}
	}
)

func NewRequest() Request {
	return Request{
		Func:   "",
		Params: map[string]interface{}{},
	}
}

func (p *Request) Get(key string) interface{} {
	return p.Params[key]
}

func (p *Request) Set(key string, val interface{}) {
	p.Params[key] = val
}

func NewReply() Reply {
	return map[string]interface{}{}
}

func (p *Reply) Get(key string) interface{} {
	return (*p)[key]
}

func (p *Reply) Set(key string, val interface{}) {
	(*p)[key] = val
}

func NewParamChecker() *ParamChecker {
	return &ParamChecker{
		// map[functionName]map[paramName]paramTypeSample,
		// if paramTypeSample is nil, means interface{} param, doesn't need to check type.
		InRequirements:  map[string]map[string]interface{}{},
		OutRequirements: map[string]map[string]interface{}{},
	}
}

// Require param with specified type.
// Type of initialized interface{} like interface{}('') is string, not interface{}.
// Type of uninitialized interface{} is 'nil'.
// If requires uninitialized interface{} type to allow multiple param types, set typeSample as 'nil',
func (pc *ParamChecker) Require(function string, isIn bool, name string, typeSample interface{}) {
	if isIn {
		existedMap, ok := pc.InRequirements[function]
		if !ok {
			existedMap = map[string]interface{}{}
		}
		existedMap[name] = typeSample
		pc.InRequirements[function] = existedMap
	} else {
		existedMap, ok := pc.OutRequirements[function]
		if !ok {
			existedMap = map[string]interface{}{}
		}
		existedMap[name] = typeSample
		pc.OutRequirements[function] = existedMap
	}
}

// Require interface{} param.
// 'nil' means uninitialized interface{} param, interface{} type param doesn't need to check type.
func (pc *ParamChecker) RequireIF(function string, isIn bool, name string) {
	pc.Require(function, isIn, name, nil)
}

func (pc *ParamChecker) VerifyIn(function string, in Request) error {
	for requireKey, requireVal := range pc.InRequirements[function] {
		paramVal, ok := in.Params[requireKey]
		if !ok {
			return gerrors.New("function %s requires input param %s", function, requireKey)
		}
		// 'nil' means interface{} param, interface{} type param doesn't need to check type.
		if gany.Type(requireVal) != gany.Type(nil) && gany.Type(requireVal) != gany.Type(paramVal) {
			return gerrors.New("function %s input param %s require type %s, but %s got", function, requireKey, gany.Type(requireVal), gany.Type(paramVal))
		}
	}
	return nil
}

// If set 'fuzzyNumSlice' true, consider []uint8 / []uint16 ... []int64 as []float64
// because after json.Unmarshal([]byte, &interface{}), output interface is set as map[string]interface{},
// and any type number map value is float64 type.
func (pc *ParamChecker) VerifyOut(function string, out *Reply, fuzzyNumSlice bool) error {
	for requireKey, requireVal := range pc.OutRequirements[function] {
		paramVal, ok := (*out)[requireKey]
		if !ok {
			return gerrors.New("function %s requires output param %s", function, requireKey)
		}
		// 'nil' means interface{} param, interface{} type param doesn't need to check type.
		if gany.Type(requireVal) == gany.Type(nil) {
			return nil
		}

		requireType, err := gany.TypeEx(requireVal)
		if err != nil {
			return err
		}
		paramType, err := gany.TypeEx(paramVal)
		if err != nil {
			return err
		}

		// In jsonrpc, after json.Unmarshal to map[string]interface{}, all numbers in map's value will becomes float64 type.
		if fuzzyNumSlice {
			numTypes := []string{
				"[]uint",
				"[]uint8",
				"[]uint16",
				"[]uint32",
				"[]uint64",
				"[]int",
				"[]int8",
				"[]int16",
				"[]int32",
				"[]int64",
				"[]float32",
			}
			for _, v := range numTypes {
				if requireType == v {
					requireType = "[]float64"
					break
				}
			}
			for _, v := range numTypes {
				if paramType == v {
					paramType = "[]float64"
					break
				}
			}
		}

		// If Reply value is an empty slice, it is nil after json.Unmarshal reply buffer into Reply - map[string]interface{}.
		if gstring.StartWith(requireType, "[]") && paramType == "nil" {
			return nil
		}

		// json.Unmarshal could decode interface slice to map[string]interface slice.
		if requireType == "[]interface {}" && paramType == "[]map[string]interface {}" {
			return nil
		}

		if requireType != paramType {
			return gerrors.New("function %s output param %s require type %s, but %s got", function, requireKey, requireType, paramType)
		}
	}
	return nil
}
