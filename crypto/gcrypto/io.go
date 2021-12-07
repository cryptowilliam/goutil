package gcrypto

import (
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"io"
)

type (
	// ElcRwc wraps io.ReadWriteCloser and EqLenCipher
	// and it implements io.ReadWriteCloser
	ElcRwc struct {
		cipher EqLenCipher
		rwc io.ReadWriteCloser
	}
)

// Read implements io.ReadWriteCloser
func (cipher *ElcRwc) Read(p []byte) (n int, err error) {
	nRead, errRead := cipher.rwc.Read(p)

	errDecrypt := error(nil)
	if nRead > 0 {
		errDecrypt = cipher.cipher.Decrypt(p[:nRead])
	}
	return nRead, gerrors.Join(errRead, errDecrypt)
}

// Write implements io.ReadWriteCloser
func (cipher *ElcRwc) Write(p []byte) (n int, err error) {
	err = cipher.cipher.Encrypt(p)
	if err != nil {
		return 0, err
	}

	wLen, errWrite := cipher.rwc.Write(p)
	errDecrypt := error(nil)
	if wLen < len(p) {
		errDecrypt = cipher.cipher.Decrypt(p[wLen:]) // restore unsuccessfully written data
	}
	return wLen, gerrors.Join(errWrite, errDecrypt)
}

// Close implements io.ReadWriteCloser.
func (cipher *ElcRwc) Close() error {
	return cipher.rwc.Close()
}

// NewElcRwc create new plaintext ciphertext equal length cipher codec wrapping a `io.ReadWriteCloser`.
func NewElcRwc(elc EqLenCipher, rwc io.ReadWriteCloser) *ElcRwc {
	return &ElcRwc{
		cipher: elc,
		rwc:  rwc,
	}
}