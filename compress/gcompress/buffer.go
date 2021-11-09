package gcompress

import "bytes"

type (
	BufferCompress struct {
		algo CompAlgo
	}
)

func NewBufferCompress(inBuf []byte, outBuf bytes.Buffer, algo CompAlgo) (*BufferCompress, error) {
	return nil, nil
}
