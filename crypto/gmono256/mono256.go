package gmono256

/**
This package implements monoalphabetic cipher for single byte, each byte may contain 256 numbers,
monoalphabetic cipher is also called simple substitution cipher.
Reference: https://en.wikipedia.org/wiki/Substitution_cipher
This is a simple encryption algorithm that can be used in applications with low security requirements.
*/

import (
	"encoding/base64"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"math/rand"
	"strings"
	"time"
)

const alphabetLen = 256

type (
	alphabet [alphabetLen]byte

	Mono256Cipher struct {
		encAlphabet *alphabet // alphabet for encoding
		decAlphabet *alphabet // alphabet for decoding
	}
)

func init() {
	// Update random seeds to prevent generating the same random alphabet.
	rand.Seed(time.Now().Unix())
}

func (abt *alphabet) Len() int {
	return alphabetLen
}

func (abt *alphabet) Less(i, j int) bool {
	return abt[i] < abt[j]
}

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

// Encode plaintext data to ciphertext.
func (cipher *Mono256Cipher) Encode(bs []byte) {
	for i, v := range bs {
		bs[i] = cipher.encAlphabet[v]
	}
}

// Decode from ciphertext data to plaintext.
func (cipher *Mono256Cipher) Decode(bs []byte) {
	for i, v := range bs {
		bs[i] = cipher.decAlphabet[v]
	}
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
