package gcrypto

import (
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"io"
)

type (
	// ElcRwc wraps io.ReadWriteCloser and EqualLengthCipher
	// and it implements io.ReadWriteCloser
	ElcRwc struct {
		cipher EqualLengthCipher
		rwc io.ReadWriteCloser
	}
)

// Read implements io.ReadWriteCloser
func (cipher *ElcRwc) Read(p []byte) (n int, err error) {
	cipher.rwc.Read(p)
	err = cipher.cipher.Decrypt(p)
	return len(p), err
}

// Write implements io.ReadWriteCloser
func (cipher *ElcRwc) Write(p []byte) (n int, err error) {
	err = cipher.cipher.Encrypt(p)
	if err != nil {
		return 0, err
	}

	wLen, err := cipher.rwc.Write(p)
	if wLen < len(p) {
		errDecrypt := cipher.cipher.Decrypt(p[wLen:]) // restore unsuccessfully written data
		if errDecrypt != nil {
			if err == nil {
				err = errDecrypt
			} else {
				err = gerrors.Wrap(err, errDecrypt.Error())
			}
		}
	}
	return wLen, err
}

// Close implements io.ReadWriteCloser.
func (cipher *ElcRwc) Close() error {
	return cipher.rwc.Close()
}

// NewElcRwc create new plaintext ciphertext equal length cipher codec wrapping a `io.ReadWriteCloser`.
func NewElcRwc(elc EqualLengthCipher, rwc io.ReadWriteCloser) *ElcRwc {
	return &ElcRwc{
		cipher: elc,
		rwc:  rwc,
	}
}