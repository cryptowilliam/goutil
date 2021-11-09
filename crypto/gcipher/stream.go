package gcipher

import "io"

type (
	CipherStream struct {

	}
)

func (c *CipherStream) Read(p []byte) (n int, err error) {
}

func (c *CipherStream) Write(p []byte) (n int, err error) {
	n, err = c.w.Write(p)
	err = c.w.Flush()
	return n, err
}

func (c *CipherStream) Close() error {

}

func NewStream(cipherAlgo CipherAlgo, rwc io.ReadWriteCloser) (io.ReadWriteCloser, error) {

}