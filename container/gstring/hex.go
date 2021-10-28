package gstring

import (
	"encoding/hex"
)

// Is every letter is hex letter [0-9] [a-f]
func IsHexString(s string) bool {
	if len(s) == 0 {
		return false
	}

	_, err := hex.DecodeString(s)
	if err != nil {
		return false
	}
	return true
}

func IsMd5String(s string) bool {
	if len(s) != 32 {
		return false
	}
	return IsHexString(s)
}
