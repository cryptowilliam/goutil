package gbytes

import (
	"bytes"
	"encoding/binary"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/container/ginterface"
)

func Index(p []byte, toSearch byte) int {
	for i, bt := range p {
		if bt == toSearch {
			return i
		}
	}
	return -1
}

func NumToBytes(n interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := error(nil)
	nType := ginterface.Type(n)
	switch nType {
	case "uint8":
		err = binary.Write(buf, binary.LittleEndian, n.(uint8))
	case "uint16":
		err = binary.Write(buf, binary.LittleEndian, n.(uint16))
	case "uint32":
		err = binary.Write(buf, binary.LittleEndian, n.(uint32))
	case "uint64":
		err = binary.Write(buf, binary.LittleEndian, n.(uint64))
	case "int8":
		err = binary.Write(buf, binary.LittleEndian, n.(int8))
	case "int16":
		err = binary.Write(buf, binary.LittleEndian, n.(int16))
	case "int32":
		err = binary.Write(buf, binary.LittleEndian, n.(int32))
	case "int64":
		err = binary.Write(buf, binary.LittleEndian, n.(int64))
	case "int":
		err = binary.Write(buf, binary.LittleEndian, n.(int))
	default:
		return nil, gerrors.New("unsupported type %s", nType)
	}
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
