package gcompress

import "bytes"

type (
	BufferCompress struct {
		algo Comp
	}
)

func NewBufferCompress(inBuf []byte, outBuf bytes.Buffer, algo Comp) (*BufferCompress, error) {
	return nil, nil
}
