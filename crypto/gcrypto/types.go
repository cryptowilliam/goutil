package gcrypto

type (
	// Cipher is an algorithm for performing encryption or decryption in cryptography.
	Cipher string

	// EqualLengthCipher defines the interface for encryption algorithms where
	// plaintext and ciphertext have the same length.
	EqualLengthCipher interface {
		Encrypt(p []byte) error
		Decrypt(p []byte) error
	}
)

var (
	CipherMono256 Cipher = "mono256"
)