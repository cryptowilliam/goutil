package g2fa

/**
copyright: this module is ported from github.com/robbiev/two-factor-auth at 2020/02/18, because a 2FA function needed
*/

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"fmt"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"strings"
	"time"
)

func toBytes(value int64) []byte {
	var result []byte
	mask := int64(0xFF)
	shifts := [8]uint16{56, 48, 40, 32, 24, 16, 8, 0}
	for _, shift := range shifts {
		result = append(result, byte((value>>shift)&mask))
	}
	return result
}

func toUint32(bytes []byte) uint32 {
	return (uint32(bytes[0]) << 24) + (uint32(bytes[1]) << 16) +
		(uint32(bytes[2]) << 8) + uint32(bytes[3])
}

func oneTimePassword(key []byte, value []byte) uint32 {
	// sign the value using HMAC-SHA1
	hmacSha1 := hmac.New(sha1.New, key)
	hmacSha1.Write(value)
	hash := hmacSha1.Sum(nil)

	// We're going to use a subset of the generated hash.
	// Using the last nibble (half-byte) to choose the index to start from.
	// This number is always appropriate as it's maximum decimal 15, the hash will
	// have the maximum index 19 (20 bytes of SHA1) and we need 4 bytes.
	offset := hash[len(hash)-1] & 0x0F

	// get a 32-bit (4-byte) chunk from the hash starting at offset
	hashParts := hash[offset : offset+4]

	// ignore the most significant bit as per RFC 4226
	hashParts[0] = hashParts[0] & 0x7F

	number := toUint32(hashParts)

	// size to 6 digits
	// one million is the first number with 7 digits so the remainder
	// of the division will always return < 7 digits
	pwd := number % 1000000

	return pwd
}

// all []byte in this program are treated as Big Endian
func TwoFactorAuth(inKey string) (password string, secondsRemaining int, err error) {
	if inKey == "" {
		return "", 0, gerrors.New("must specify key to use")
	}

	// decode the key from the first argument
	inputNoSpaces := strings.Replace(inKey, " ", "", -1)
	inputNoSpacesUpper := strings.ToUpper(inputNoSpaces)
	key, err := base32.StdEncoding.DecodeString(inputNoSpacesUpper)
	if err != nil {
		return "", 0, err
	}

	// generate a one-time password using the time at 30-second intervals
	epochSeconds := time.Now().Unix()
	pwd := oneTimePassword(key, toBytes(epochSeconds/30))

	secondsRemaining = int(30 - (epochSeconds % 30))
	return fmt.Sprintf("%06d", pwd), secondsRemaining, nil
}
