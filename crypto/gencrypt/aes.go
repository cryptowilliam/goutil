package gencrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"io"
)

// AesGcm256
// Note:
// When perform multiple times, it returns different results even user don't change input,
// because random number used in encrypt process, please DON'T use it when you need same output for unchanging input.

type (
	AesGcm256 struct{}
)

func NewAesGcm256() *AesGcm256 {
	return &AesGcm256{}
}

// Generate sha256 sum of original key,
// make it 32 bytes long to meet the AES length requirement for the key.
func normalize(key []byte) []byte {
	sum := make([]byte, 0)

	for _, item := range sha256.Sum256(key) {
		sum = append(sum, item)
	}

	return sum
}

// 256-bit AES-GCM with a random nonce
// Reference: https://github.com/gtank/cryptopasta/blob/master/encrypt.go
// Encrypt encrypts data using 256-bit AES-GCM.  This both hides the content of
// the data and provides a check that it hasn't been altered. Output takes the
// form nonce|ciphertext|tag where '|' indicates concatenation.
func (a *AesGcm256) Encrypt(plaintext []byte, key []byte, normalizeWithSha256 bool) (ciphertext []byte, err error) {
	if !normalizeWithSha256 && len(key) != 32 {
		return nil, gerrors.New("AES GCM 256 requires 32 bytes secret")
	}

	if normalizeWithSha256 {
		key = normalize(key)
	}

	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

// 256-bit AES-GCM with a random nonce
// Reference: https://github.com/gtank/cryptopasta/blob/master/encrypt.go
// Decrypt decrypts data using 256-bit AES-GCM.  This both hides the content of
// the data and provides a check that it hasn't been altered. Expects input
// form nonce|ciphertext|tag where '|' indicates concatenation.
func (a *AesGcm256) Decrypt(ciphertext []byte, key []byte, normalizeWithSha256 bool) (plaintext []byte, err error) {
	if !normalizeWithSha256 && len(key) != 32 {
		return nil, gerrors.New("AES GCM 256 requires 32 bytes secret")
	}

	if normalizeWithSha256 {
		key = normalize(key)
	}

	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < gcm.NonceSize() {
		return nil, gerrors.New("malformed ciphertext")
	}

	return gcm.Open(nil,
		ciphertext[:gcm.NonceSize()],
		ciphertext[gcm.NonceSize():],
		nil,
	)
}
