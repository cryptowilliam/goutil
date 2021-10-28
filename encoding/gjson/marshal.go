package gjson

import (
	"encoding/json"
	"fmt"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/mailru/easyjson"
	"reflect"
	"strings"
)

func MarshalBuffer(v interface{}, indent bool) ([]byte, error) {
	var buf []byte
	err := error(nil)
	if indent {
		buf, err = json.MarshalIndent(v, "", "\t")
	} else {
		buf, err = json.Marshal(v)
	}
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func MarshalString(v interface{}, indent bool) (string, error) {
	b, err := MarshalBuffer(v, indent)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func MarshalBufferWithErrFmt(v interface{}, indent bool, errFmt string) []byte {
	var buf []byte
	err := error(nil)
	if indent {
		buf, err = json.MarshalIndent(v, "", "\t")
	} else {
		buf, err = json.Marshal(v)
	}
	if err != nil {
		return []byte(fmt.Sprintf(errFmt, err.Error()))
	}
	return buf
}

func MarshalStringWithErrFmt(v interface{}, indent bool, errFmt string) string {
	return string(MarshalBufferWithErrFmt(v, indent, errFmt))
}

func MarshalBufferDefault(v interface{}, indent bool) []byte {
	return MarshalBufferWithErrFmt(v, indent, `{"Error":"%s"}`)
}

func MarshalStringDefault(v interface{}, indent bool) string {
	return MarshalStringWithErrFmt(v, indent, `{"Error":"%s"}`)
}

func marshalFast(v easyjson.Marshaler) ([]byte, error) {
	return easyjson.Marshal(v)
}

// JSONEncode encodes structure data into JSON
func JSONEncode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// JSONDecode decodes JSON data into a structure
func JSONDecode(data []byte, to interface{}) error {
	if !strings.Contains(reflect.ValueOf(to).Type().String(), "*") {
		return gerrors.New("json decode error - memory address not supplied")
	}
	return json.Unmarshal(data, to)
}
