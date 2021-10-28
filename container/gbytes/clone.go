package gbytes

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

// 深拷贝任意对象，完整复制数据
func Clone2(src, dst interface{}) error {
	buff := new(bytes.Buffer)
	enc := gob.NewEncoder(buff)
	dec := gob.NewDecoder(buff)
	if err := enc.Encode(src); err != nil {
		return err
	}
	if err := dec.Decode(dst); err != nil {
		return err
	}
	return nil
}

func DeepCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}

func DeepCopyEx(src interface{}) (interface{}, error) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		fmt.Println("en")
		return nil, err
	}
	var r interface{}
	err := gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(r)
	if err != nil {
		fmt.Println("ex")
		return nil, err
	}
	return r, nil
}
