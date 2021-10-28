package gaddr

import (
	"bytes"
	"encoding/binary"
	"net"
)

// Converts IP mask to 16 bit unsigned integer.
func MaskToInt(mask net.IPMask) (uint16, error) {
	var i uint16
	buf := bytes.NewReader(mask)
	err := binary.Read(buf, binary.BigEndian, &i)
	return i, err
}
