package gbase

import "encoding/base64"

func Base64Encode(message []byte) string {
	return base64.RawStdEncoding.EncodeToString(message)
}

func Base64Decode(s string) ([]byte, error) {
	return base64.RawStdEncoding.DecodeString(s)
}
