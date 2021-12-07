package gcrypto

import "io"

type (
	// Cipher is an algorithm for performing encryption or decryption in cryptography.
	Cipher string

	// EqLenCipher defines the interface for encryption algorithms where
	// plaintext and ciphertext have the same length.
	EqLenCipher interface {
		Encrypt(b []byte) error
		Decrypt(b []byte) error
	}

	// VarLenCipher defines the interface for encryption algorithms where
	// plaintext and ciphertext have different length.
	VarLenCipher interface {
		Encrypt(b []byte) ([]byte, error)
		Decrypt(b []byte) ([]byte, error)
	}

	// CipherReadWriteCloser defines the interface for encryption algorithms where
	// plaintext and ciphertext have different length.
	CipherReadWriteCloser interface {
		io.ReadWriteCloser
	}
)

var (
	CipherMono256 Cipher = "mono256"
	CipherChaCha20 Cipher = "chacha20"
)