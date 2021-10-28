package ghash

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"hash"
)

type HashType string

const (
	HashTypeSHA1       HashType = "sha1"
	HashTypeSHA256     HashType = "sha256"
	HashTypeSHA512     HashType = "sha512"
	HashTypeSHA512_384 HashType = "sha512_384"
	HashTypeMD5        HashType = "md5"
)

// GetSHA256 returns a SHA256 hash of a byte array
func GetSHA256(message []byte) []byte {
	digest := sha256.New()
	digest.Write(message)
	return digest.Sum(nil)
}

// GetSHA512 returns a SHA512 hash of a byte array
func GetSHA512(message []byte) []byte {
	sha := sha512.New()
	sha.Write(message)
	return sha.Sum(nil)
}

func GetSHA1(message []byte) []byte {
	digest := sha1.New()
	digest.Write(message)
	return digest.Sum(nil)
}

// GetMD5 returns a MD5 hash of a byte array
func GetMD5(message []byte) []byte {
	digest := md5.New()
	digest.Write(message)
	return digest.Sum(nil)
}

func GetHex(message []byte) string {
	return hex.EncodeToString(message)
}

// HexEncodeToString takes in a hexadecimal byte array and returns a string
func HexEncodeToString(message []byte) string {
	return hex.EncodeToString(message)
}

// GetHMAC returns a keyed-hash message authentication code using the desired
// hashtype
func GetHMAC(hashType HashType, plain, secret []byte) []byte {
	var hash func() hash.Hash

	switch hashType {
	case HashTypeSHA1:
		hash = sha1.New
	case HashTypeSHA256:
		hash = sha256.New
	case HashTypeSHA512:
		hash = sha512.New
	case HashTypeSHA512_384:
		hash = sha512.New384
	case HashTypeMD5:
		hash = md5.New
	default:
		panic(gerrors.New("unsupported hash type %s", hashType))
	}

	mac := hmac.New(hash, secret)
	mac.Write(plain)
	return mac.Sum(nil)
}
