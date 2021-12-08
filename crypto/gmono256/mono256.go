package gmono256

/**
This package implements monoalphabetic cipher for single byte, each byte may contain 256 numbers,
monoalphabetic cipher is also called simple substitution cipher.
Reference: https://en.wikipedia.org/wiki/Substitution_cipher
This is a simple encryption algorithm that can be used in short or low security requirements message transmission.
*/

import (
	"encoding/base64"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/crypto/gcrypto"
	"io"
	"math/rand"
	"strings"
	"time"
)

// length of codec alphabet, its elements are [0, 255]
const alphabetLen = 256

type (
	// codec alphabet
	alphabet [alphabetLen]byte

	// Mono256Cipher is monoalphabetic cipher codec
	Mono256Cipher struct {
		encAlphabet *alphabet // alphabet for encoding
		decAlphabet *alphabet // alphabet for decoding
	}

	Mono256Maker struct {
		cipher *Mono256Cipher
	}
)

func init() {
	// Update random seeds to prevent generating the same random alphabet.
	rand.Seed(time.Now().Unix())
}

// Len implements sort.Interface
func (abt *alphabet) Len() int {
	return alphabetLen
}

// Less implements sort.Interface
func (abt *alphabet) Less(i, j int) bool {
	return abt[i] < abt[j]
}

// Swap implements sort.Interface
func (abt *alphabet) Swap(i, j int) {
	abt[i], abt[j] = abt[j], abt[i]
}

// ToBase64 convert 256 bytes alphabet to base64 string.
func (abt *alphabet) ToBase64() string {
	return base64.StdEncoding.EncodeToString(abt[:])
}

// Base64ToAlphabet convert base64 string to 256 bytes alphabet.
func Base64ToAlphabet(b64s string) (*alphabet, error) {
	b, err := base64.StdEncoding.DecodeString(strings.TrimSpace(b64s))
	if err != nil {
		return nil, err
	}
	if len(b) != alphabetLen {
		return nil, gerrors.New("alphabet length %d != %d", len(b), alphabetLen)
	}
	rst := alphabet{}
	copy(rst[:], b)
	return &rst, nil
}

// Generate a random combination of 256 byte alphabet, which are finally encoded as strings using base64,
// without any duplicate byte, and must consist of 0-255, and all need to be included.
func randAlphabet() string {
	// Generate a random byte array consisting of 0~255
	intArr := rand.Perm(alphabetLen)
	abt := &alphabet{}

	// Copy random int array to byte array.
	for idx, val := range intArr {
		abt[idx] = byte(val)
		// Ensure that all bytes do not have the same index and value, and if they do, regenerate them.
		if idx == val {
			return randAlphabet()
		}
	}

	// Convert to base64 string.
	return abt.ToBase64()
}

// Encrypt plaintext data to ciphertext.
// It implements `EqLenCipher` interface.
func (cipher *Mono256Cipher) Encrypt(b []byte) error {
	for i, v := range b {
		b[i] = cipher.encAlphabet[v]
	}
	return nil
}

// Decrypt from ciphertext data to plaintext.
// It implements `EqLenCipher` interface.
func (cipher *Mono256Cipher) Decrypt(b []byte) error {
	for i, v := range b {
		b[i] = cipher.decAlphabet[v]
	}
	return nil
}

// NewMono256 create new monoalphabetic cipher codec.
func NewMono256(encAlphabet *alphabet) *Mono256Cipher {
	decAlphabet := &alphabet{}
	for i, v := range encAlphabet {
		encAlphabet[i] = v
		decAlphabet[v] = byte(i)
	}
	return &Mono256Cipher{
		encAlphabet: encAlphabet,
		decAlphabet: decAlphabet,
	}
}

// NewRandKeyBase64 generates random key in base64 format.
func NewRandKeyBase64() string {
	return randAlphabet()
}

func NewMono256Maker(b64alphabet string) (gcrypto.CipherRWCMaker, error) {
	encAlphabet, err := Base64ToAlphabet(b64alphabet)
	if err != nil {
		return nil, err
	}

	return &Mono256Maker{cipher: NewMono256(encAlphabet)}, nil
}

func (m *Mono256Maker) Make(rwc io.ReadWriteCloser, readNonce bool, nonceCodec gcrypto.EqLenCipher) (gcrypto.CipherRWC, error) {
	return gcrypto.NewEqLenCipherRWC(m.cipher, rwc), nil
}