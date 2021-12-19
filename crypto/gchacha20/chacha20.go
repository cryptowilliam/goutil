package gchacha20

// Attention should be paid to the fact that not only the same passphrase
// should be used for decryption, but also the same random number (nonce).
// If the random number used for decryption is different from that used for
// encryption, it will not be decrypted correctly.
//
// The nonce must be unique for one key for all time.
// The length of the nonce determines the version of ChaCha20:
// - 8 bytes:  ChaCha20 with a 64 bit nonce and a 2^64 * 64 byte period.
// - 12 bytes: ChaCha20-IETF-Poly1305 with a 96 bit nonce in RFC 7539 and a 2^32 * 64 byte period.
// - 24 bytes: XChaCha20-Poly1305 with a 192 bit nonce in RFC 8439 and a 2^64 * 64 byte period.

import (
	"bufio"
	"bytes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/crypto/gcrypto"
	"github.com/cryptowilliam/goutil/sys/gio"
	"golang.org/x/crypto/chacha20"
	"golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/scrypt"
	"io"
	"io/ioutil"
	"time"
)

type (
	ChaCha20 struct{}

	// ChaCha20Cipher is a modification of Salsa20 published in 2008.
	// It uses a new round function that increases diffusion and increases
	// performance on some architectures.
	ChaCha20Cipher struct {
		key     []byte
		version gcrypto.Cipher
	}

	// ChaCha20Maker is a stream style ChaCha20-poly-1305 codec generator.
	ChaCha20Maker struct {
		key              []byte
		version          gcrypto.Cipher // "chacha20-ietf-poly1305", "xchacha20-poly1305"
		chaR             chacha20.Cipher
		chaW             chacha20.Cipher
		correctNonceSize int
	}

	// ChaCha20RWC is a stream style ChaCha20-poly-1305 codec.
	ChaCha20RWC struct {
		csr *cipher.StreamReader
		csw *cipher.StreamWriter
	}
)

var (
	errUnknownChaChaVer = "unknown chacha20 version '%s'"
)

func NewChaCha20() *ChaCha20 {
	return &ChaCha20{}
}

// CodecWithPassphrase creates ChaCha20 codec with string passphrase.
func (c *ChaCha20) CodecWithPassphrase(passphrase string, version gcrypto.Cipher) (*ChaCha20Cipher, error) {
	if len(passphrase) == 0 {
		return nil, fmt.Errorf("passphrase is required")
	}
	key := passphraseToKey(passphrase)

	return c.CodecWithKey(key, version)
}

// CodecWithKey creates ChaCha20 codec with bytes key.
func (c *ChaCha20) CodecWithKey(key []byte, version gcrypto.Cipher) (*ChaCha20Cipher, error) {
	if len(key) != 32 {
		return nil, fmt.Errorf("invalid key")
	}
	result := &ChaCha20Cipher{}
	result.key = key
	result.version = version
	return result, nil
}

// MakerWithPassphrase creates ChaCha20 stream codec with string passphrase.
func (c *ChaCha20) MakerWithPassphrase(passphrase string, version gcrypto.Cipher) (gcrypto.CipherRWCMaker, error) {

	if len(passphrase) == 0 {
		return nil, fmt.Errorf("passphrase is required")
	}
	key := passphraseToKey(passphrase)

	return c.MakerWithKey(key, version)

}

// MakerWithKey creates ChaCha20 stream codec with bytes key.
func (c *ChaCha20) MakerWithKey(key []byte, version gcrypto.Cipher) (gcrypto.CipherRWCMaker, error) {
	switch version {
	case gcrypto.CipherChaCha20IETFPoly1305, gcrypto.CipherXChaCha20Poly1305:
	default:
		return nil, gerrors.New(errUnknownChaChaVer, version)
	}

	if len(key) != 32 {
		return nil, fmt.Errorf("invalid key")
	}

	correctNonceSize, err := getNonceSize(version)
	if err != nil {
		return nil, err
	}

	result := &ChaCha20Maker{}
	result.key = key
	result.correctNonceSize = correctNonceSize
	result.version = version
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

	block, err := newChaCha20AEAD(c.key, c.version)
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
	nonceSize, err := getNonceSize(c.version)
	if err != nil {
		return nil, err
	}
	if len(b) < nonceSize {
		return nil, fmt.Errorf("invalid cipher text")
	}

	r := bufio.NewReader(bytes.NewReader(b))
	nonce := make([]byte, nonceSize)
	if _, err := io.ReadFull(r, nonce); err != nil {
		return nil, err
	}

	ciphertext, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	block, err := newChaCha20AEAD(c.key, c.version)
	if err != nil {
		return nil, err
	}

	rawData, err := block.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return rawData, nil
}

func (m *ChaCha20Maker) NonceSize() int {
	return m.correctNonceSize
}

// Make wraps io.ReadWriteCloser and generate a CipherRWC.
// genNonce:
// true: generate random nonce and write to other side
// false: read nonce from other side
// timeout:
// it is recommended, used to avoid unexpected block of reader `rwc`.
func (m *ChaCha20Maker) Make(rwc io.ReadWriteCloser, genNonce bool, timeout *time.Duration, nonceCodec gcrypto.EqLenCipher) (gcrypto.CipherRWC, error) {
	if rwc == nil {
		return nil, gerrors.New("nil rwc")
	}
	var nonce []byte
	correctNonceSize := m.NonceSize()
	err := error(nil)

	// read or write nonce
	if genNonce { // generate and write nonce
		nonce, err = generateNonce(m.key, m.version)
		if err != nil {
			return nil, err
		}
		if nonceCodec != nil {
			if err := nonceCodec.Encrypt(nonce); err != nil {
				return nil, err
			}
		}
		n, err := rwc.Write(nonce[:correctNonceSize])
		if err != nil {
			return nil, err
		}
		if n != correctNonceSize {
			return nil, gerrors.New("write nonce size %d != correct nonce size %d", n, correctNonceSize)
		}
	} else { // read nonce from writer side.
		nonce = make([]byte, correctNonceSize)
		_, err = gio.ReadFull(rwc, nonce, timeout) // timeout: avoid unexpected block of reader `rwc`
		if err != nil {
			return nil, err
		}
		if nonceCodec != nil {
			if err = nonceCodec.Decrypt(nonce); err != nil {
				return nil, err
			}
		}
	}

	chaR, err := chacha20.NewUnauthenticatedCipher(m.key, nonce)
	if err != nil {
		return nil, err
	}
	chaW, err := chacha20.NewUnauthenticatedCipher(m.key, nonce)
	if err != nil {
		return nil, err
	}

	s := &ChaCha20RWC{
		csr: &cipher.StreamReader{
			S: chaR,
			R: rwc,
		},
		csw: &cipher.StreamWriter{
			S: chaW,
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

func newChaCha20AEAD(key []byte, version gcrypto.Cipher) (cipher.AEAD, error) {
	switch version {
	case gcrypto.CipherChaCha20IETFPoly1305:
		return chacha20poly1305.New(key)
	case gcrypto.CipherXChaCha20Poly1305:
		return chacha20poly1305.NewX(key)
	default:
		return nil, gerrors.New(errUnknownChaChaVer, version)
	}
}

func getNonceSize(version gcrypto.Cipher) (int, error) {
	switch version {
	case gcrypto.CipherChaCha20IETFPoly1305: // 12
		/*block, err := chacha20poly1305.New(key)
		if err != nil {
			return 0, err
		}
		return block.NonceSize(), nil*/
		return 12, nil
	case gcrypto.CipherXChaCha20Poly1305: // 24
		/*block, err := chacha20poly1305.NewX(key)
		if err != nil {
			return 0, err
		}
		return block.NonceSize(), nil*/
		return 24, nil
	default:
		return 0, gerrors.New(errUnknownChaChaVer, version)
	}

}

func generateNonce(key []byte, version gcrypto.Cipher) ([]byte, error) {
	switch version {
	case gcrypto.CipherChaCha20IETFPoly1305:
		block, err := chacha20poly1305.New(key)
		if err != nil {
			return nil, err
		}
		return randomBytes(uint16(block.NonceSize())), nil
	case gcrypto.CipherXChaCha20Poly1305:
		block, err := chacha20poly1305.NewX(key)
		if err != nil {
			return nil, err
		}
		return randomBytes(uint16(block.NonceSize())), nil
	default:
		return nil, gerrors.New(errUnknownChaChaVer, version)
	}

}
