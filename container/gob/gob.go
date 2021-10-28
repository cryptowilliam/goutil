package gob

import (
	"bytes"
	"encoding/gob"
)

func Encode(v interface{}) ([]byte, error) {
	var w bytes.Buffer
	enc := gob.NewEncoder(&w)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func Decode(data []byte, outputStructPtr interface{}) error {
	var w bytes.Buffer
	w.Write(data)
	dec := gob.NewDecoder(&w)
	if err := dec.Decode(outputStructPtr); err != nil {
		return err
	}
	return nil
}
