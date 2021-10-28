package gbytes

import (
	"bytes"
	"encoding/binary"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/container/ginterface"
	"unsafe"
)

func BytesToUint32(b []byte) (uint32, error) {
	if len(b) != 4 {
		return 0, gerrors.New("b len %d != expected len 4", len(b))
	}
	var num uint32
	if err := binary.Read(bytes.NewBuffer(b[:]), binary.LittleEndian, &num); err != nil {
		return 0, err
	}
	return num, nil
}

func BytesToUint64(b []byte) (uint64, error) {
	if len(b) != 8 {
		return 0, gerrors.New("b len %d != expected len 8", len(b))
	}
	var num uint64
	if err := binary.Read(bytes.NewBuffer(b[:]), binary.LittleEndian, &num); err != nil {
		return 0, err
	}
	return num, nil
}

func BytesToInt32(b []byte) (int32, error) {
	if len(b) != 4 {
		return 0, gerrors.New("b len %d != expected len 4", len(b))
	}
	var num int32
	if err := binary.Read(bytes.NewBuffer(b[:]), binary.LittleEndian, &num); err != nil {
		return 0, err
	}
	return num, nil
}

func BytesToInt64(b []byte) (int64, error) {
	if len(b) != 8 {
		return 0, gerrors.New("b len %d != expected len 8", len(b))
	}
	var num int64
	if err := binary.Read(bytes.NewBuffer(b[:]), binary.LittleEndian, &num); err != nil {
		return 0, err
	}
	return num, nil
}

func BytesToNum(b []byte, numTypeSample interface{}) (interface{}, error) {
	if len(b) != 8 {
		return 0, gerrors.New("b len %d != expected len 8", len(b))
	}
	/*var num int64
	if err := binary.Read(bytes.NewBuffer(b[:]), binary.LittleEndian, &num); err != nil {
		return 0, err
	}
	return num, nil*/

	ntsType := ginterface.Type(numTypeSample)
	switch ntsType {
	case "uint8":
		if len(b) != 1 {
			return 0, gerrors.New("b len %d != expected len 1", len(b))
		}
		var num uint8
		if err := binary.Read(bytes.NewBuffer(b[:]), binary.LittleEndian, &num); err != nil {
			return 0, err
		}
		return num, nil
	case "uint16":
		if len(b) != 2 {
			return 0, gerrors.New("b len %d != expected len 2", len(b))
		}
		var num uint16
		if err := binary.Read(bytes.NewBuffer(b[:]), binary.LittleEndian, &num); err != nil {
			return 0, err
		}
		return num, nil
	case "uint32":
		if len(b) != 4 {
			return 0, gerrors.New("b len %d != expected len 4", len(b))
		}
		var num uint32
		if err := binary.Read(bytes.NewBuffer(b[:]), binary.LittleEndian, &num); err != nil {
			return 0, err
		}
		return num, nil
	case "uint64":
		if len(b) != 8 {
			return 0, gerrors.New("b len %d != expected len 8", len(b))
		}
		var num uint64
		if err := binary.Read(bytes.NewBuffer(b[:]), binary.LittleEndian, &num); err != nil {
			return 0, err
		}
		return num, nil
	case "int8":
		if len(b) != 1 {
			return 0, gerrors.New("b len %d != expected len 1", len(b))
		}
		var num int8
		if err := binary.Read(bytes.NewBuffer(b[:]), binary.LittleEndian, &num); err != nil {
			return 0, err
		}
		return num, nil
	case "int16":
		if len(b) != 2 {
			return 0, gerrors.New("b len %d != expected len 2", len(b))
		}
		var num int16
		if err := binary.Read(bytes.NewBuffer(b[:]), binary.LittleEndian, &num); err != nil {
			return 0, err
		}
		return num, nil
	case "int32":
		if len(b) != 4 {
			return 0, gerrors.New("b len %d != expected len 4", len(b))
		}
		var num int32
		if err := binary.Read(bytes.NewBuffer(b[:]), binary.LittleEndian, &num); err != nil {
			return 0, err
		}
		return num, nil
	case "int64":
		if len(b) != 8 {
			return 0, gerrors.New("b len %d != expected len 8", len(b))
		}
		var num int64
		if err := binary.Read(bytes.NewBuffer(b[:]), binary.LittleEndian, &num); err != nil {
			return 0, err
		}
		return num, nil
	case "int":
		if len(b) != int(unsafe.Sizeof(int(1))) {
			return 0, gerrors.New("b len %d != expected len 8", len(b))
		}
		var num int
		if err := binary.Read(bytes.NewBuffer(b[:]), binary.LittleEndian, &num); err != nil {
			return 0, err
		}
		return num, nil
	default:
		return 0, gerrors.New("unsupported type %s", ntsType)
	}
}
