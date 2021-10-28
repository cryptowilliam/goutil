package gstruct

// encode struct to bytes array or decode bytes array to struct
// https://github.com/bitgoin/packer
// https://github.com/zhuangsirui/binpacker

// https://github.com/dxhbiz/codec/issues/1#issuecomment-331349108

import (
	"github.com/dxhbiz/codec"
)

func StructPack(v interface{}, p []byte) error {
	p, err := codec.Encode(v)
	if err != nil {
		return err
	}
	return nil
}

func StructUnpack(p []byte, v interface{}) error {
	err := codec.Decode(p, &v)
	if err != nil {
		return err
	}
	return nil
}
