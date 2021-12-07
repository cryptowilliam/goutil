package gchacha20

import (
	"bufio"
	"bytes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"golang.org/x/crypto/chacha20"
	"golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/scrypt"
	"io"
	"io/ioutil"
)

type (

	// ChaCha20Cipher is a modification of Salsa20 published in 2008.
	// It uses a new round function that increases diffusion and increases
	// performance on some architectures.
	ChaCha20Cipher struct {
		key []byte
	}

	// ChaCha20Maker is a stream style ChaCha20-poly-1305 codec generator.
	ChaCha20Maker struct {
		cha *chacha20.Cipher
		key []byte
	}

	// ChaCha20RWC is a stream style ChaCha20-poly-1305 codec.
	ChaCha20RWC struct {
		csr *cipher.StreamReader
		csw *cipher.StreamWriter
	}
)

// NewChaCha20 creates ChaCha20-poly-1305 codec with string passphrase.
func NewChaCha20(passphrase string) (*ChaCha20Cipher, error) {
	if len(passphrase) == 0 {
		return nil, fmt.Errorf("passphrase is required")
	}
	key := passphraseToKey(passphrase)

	return NewChaCha20WithPassKey(key)
}

// NewChaCha20WithPassKey creates ChaCha20-poly-1305 codec with bytes key.
func NewChaCha20WithPassKey(key []byte) (*ChaCha20Cipher, error) {
	if len(key) != 32 {
		return nil, fmt.Errorf("invalid key")
	}
	result := &ChaCha20Cipher{}
	result.key = key
	return result, nil
}


// Encrypt encrypts the given data using ChaCha20-Poly1305 with the given passphrase.
// The passphrase can be a user provided value, and is hashed using scrypt before being used.
// It implements `VarLenCipher` interface.
//
// Will return error if an empty passphrase or data is provided.
func (c ChaCha20Cipher) Encrypt(b []byte) ([]byte, error) {
	if len(b) == 0 {
		return nil, fmt.Errorf("empty plain text")
	}

	block, err := chacha20poly1305.NewX(c.key)
	if err != nil {
		return nil, err
	}

	nonce := randomBytes(uint16(block.NonceSize()))
	ciphertext := block.Seal(nil, nonce, b, nil)

	var writer bytes.Buffer
	if _, err := writer.Write(nonce); err != nil {
		return nil, err
	}
	if _, err := writer.Write(ciphertext); err != nil {
		return nil, err
	}

	return writer.Bytes(), nil
}


// Decrypt decrypts the given encrypted data using ChaCha20-Poly1305 with the given passphrase.
// The passphrase can be a user provided value, and is hashed using scrypt before being used.
// It implements `VarLenCipher` interface.
//
// Will return error if an empty passphrase or data is provided.
func (c ChaCha20Cipher) Decrypt(b []byte) ([]byte, error) {
	if len(b) < 24 {
		return nil, fmt.Errorf("invalid cipher text")
	}

	r := bufio.NewReader(bytes.NewReader(b))
	nonce := make([]byte, 24)
	if _, err := io.ReadFull(r, nonce); err != nil {
		return nil, err
	}

	ciphertext, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	block, err := chacha20poly1305.NewX(c.key)
	if err != nil {
		return nil, err
	}

	rawData, err := block.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return rawData, nil
}

// NewChaCha20Maker creates ChaCha20-poly-1305 stream codec with string passphrase.
func NewChaCha20Maker(passphrase string) (*ChaCha20Maker, error) {
	if len(passphrase) == 0 {
		return nil, fmt.Errorf("passphrase is required")
	}
	key := passphraseToKey(passphrase)

	return NewMakerWithKey(key)
}

// NewMakerWithKey creates ChaCha20-poly-1305 stream codec with bytes key.
func NewMakerWithKey(key []byte) (*ChaCha20Maker, error) {
	if len(key) != 32 {
		return nil, fmt.Errorf("invalid key")
	}

	block, err := chacha20poly1305.NewX(key)
	if err != nil {
		return nil, err
	}
	nonce := randomBytes(uint16(block.NonceSize()))
	cha, err := chacha20.NewUnauthenticatedCipher(key, nonce)
	if err != nil {
		return nil, err
	}

	result := &ChaCha20Maker{}
	result.key = key
	result.cha = cha
	return result, nil
}

func (m *ChaCha20Maker) Make(rwc io.ReadWriteCloser) (*ChaCha20RWC, error) {
	// wraps io.ReadWriteCloser
	if rwc == nil {
		return nil, gerrors.New("nil rwc")
	}

	s := &ChaCha20RWC{
		csr: &cipher.StreamReader{
			S: m.cha,
			R: rwc,
		},
		csw: &cipher.StreamWriter{
			S: m.cha,
			W: rwc,
		},
	}
	return s, nil
}

// Read decrypts cipher data and write plain data into output buffer.
func (s *ChaCha20RWC) Read(b []byte) (int, error) {
	if s.csr == nil {
		return 0, gerrors.New("uninitialized cipher.StreamReader")
	}
	return s.csr.Read(b)
}

// Write encrypts plain data and write cipher data into CipherStreamWriter.
func (s *ChaCha20RWC) Write(b []byte) (int, error) {
	if s.csw == nil {
		return 0, gerrors.New("uninitialized cipher.StreamWriter")
	}
	return s.csw.Write(b)
}

func (s *ChaCha20RWC) Close() error {
	return s.csw.Close()
}

// randomBytes generate random bytes of specified length. Suitable for cryptographical use.
// This may panic if too much data was requested.
func randomBytes(length uint16) []byte {
	randB := make([]byte, length)
	if _, err := rand.Read(randB); err != nil {
		panic(err)
	}
	return randB
}

// passphraseToKey generates a 32-byte key from the given passphrase
func passphraseToKey(passphrase string) []byte {
	key, err := scrypt.Key([]byte(passphrase), nil, 32768, 8, 1, 32)
	if err != nil {
		panic(err)
	}
	if len(key) != 32 {
		panic("invalid key length after hashing")
	}
	return key
}