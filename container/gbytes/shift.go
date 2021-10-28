package gbytes

import (
	"encoding/binary"
)

/* encode 8 bits bool */
func EncodeBool(p []byte, c bool) []byte {
	if c {
		p[0] = 1
	} else {
		p[0] = 0
	}
	return p[1:]
}

/* decode 8 bits bool */
func DecodeBool(p []byte, c *bool) []byte {
	if p[0] != 0 {
		*c = true
	} else {
		*c = false
	}
	return p[1:]
}

/* encode 8 bits unsigned int */
func EncodeUint8(p []byte, c uint8) []byte {
	p[0] = c
	return p[1:]
}

/* decode 8 bits unsigned int */
func DecodeUint8(p []byte, c *byte) []byte {
	*c = p[0]
	return p[1:]
}

/* encode 16 bits unsigned int (lsb) */
func EncodeUint16(p []byte, w uint16) []byte {
	binary.LittleEndian.PutUint16(p, w)
	return p[2:]
}

/* decode 16 bits unsigned int (lsb) */
func DecodeUint16(p []byte, w *uint16) []byte {
	*w = binary.LittleEndian.Uint16(p)
	return p[2:]
}

/* encode 32 bits unsigned int (lsb) */
func EncodeUint32(p []byte, l uint32) []byte {
	binary.LittleEndian.PutUint32(p, l)
	return p[4:]
}

/* decode 32 bits unsigned int (lsb) */
func DecodeUint32(p []byte, l *uint32) []byte {
	*l = binary.LittleEndian.Uint32(p)
	return p[4:]
}

/* encode 64 bits unsigned int (lsb) */
func EncodeUint64(p []byte, l uint64) []byte {
	binary.LittleEndian.PutUint64(p, l)
	return p[8:]
}

/* decode 64 bits unsigned int (lsb) */
func DecodeUint64(p []byte, l *uint64) []byte {
	*l = binary.LittleEndian.Uint64(p)
	return p[8:]
}
