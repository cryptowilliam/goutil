package gjson

import (
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/tidwall/gjson"
	"reflect"
	"time"
)

type JsonValue gjson.Result

func Get(json, path string) JsonValue {
	return JsonValue(gjson.Get(json, path))
}

func Set(jsonstr string, path []string, val interface{}) (string, error) {
	js, err := simplejson.NewJson([]byte(jsonstr))
	if err != nil {
		return "", err
	}
	js.SetPath(path, val)
	b, err := js.MarshalJSON()
	return string(b), err
}

func (v JsonValue) Exists() bool {
	return gjson.Result(v).Exists()
}

func (v JsonValue) String() string {
	return gjson.Result(v).String()
}

func (v JsonValue) Time() time.Time {
	return gjson.Result(v).Time()
}

func (v JsonValue) Bool() bool {
	return gjson.Result(v).Bool()
}

func (v JsonValue) Float() float64 {
	return gjson.Result(v).Float()
}

func (v JsonValue) Int64() int64 {
	return int64(gjson.Result(v).Float())
}

func MarshalAndPrintln(x interface{}) {
	buf, err := json.Marshal(x)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(buf))
}

type (
	IterFn = func(key string, val interface{}) (newVal interface{}, modified bool, err error)
)

func Iterate(jsonStr *string, indent bool, iterFn IterFn) error {
	if jsonStr == nil {
		return gerrors.New("can't iterate nil jsonStr")
	}

	var jsonIF interface{}
	if err := json.Unmarshal([]byte(*jsonStr), &jsonIF); err != nil {
		return err
	}
	jsonMap := jsonIF.(map[string]interface{})

	if err := iterateMap(&jsonMap, iterFn); err != nil {
		return err
	}

	newJsonStr, err := MarshalString(jsonMap, indent)
	if err != nil {
		return err
	}
	*jsonStr = newJsonStr
	return nil
}

func iterateMap(jsonMap *map[string]interface{}, iterFn IterFn) error {
	for k, v := range *jsonMap {
		if reflect.TypeOf(v) != nil && reflect.TypeOf(v).Kind() == reflect.Map {
			vMap := v.(map[string]interface{})
			if err := iterateMap(&vMap, iterFn); err != nil {
				return err
			}
		} else if reflect.TypeOf(v) != nil && reflect.TypeOf(v).Kind() == reflect.Slice{
			vSlice := v.([]interface{})
			for _, vEntry := range vSlice {
				if reflect.TypeOf(vEntry).Kind() == reflect.Map {
					vMap := vEntry.(map[string]interface{})
					if err := iterateMap(&vMap, iterFn); err != nil {
						return err
					}
				}
			}
		} else {
			newVal, modified, err := iterFn(k, v)
			if err != nil {
				return err
			}
			if modified {
				(*jsonMap)[k] = newVal
			}
		}
	}
	return nil
}