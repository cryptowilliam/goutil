package gsha

import (
	"bufio"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"os"
)

// sha1 file at path
func Sha1File(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", nil
	}
	defer f.Close()

	// files could be pretty big, lets buffer
	buff := bufio.NewReader(f)
	hash := sha1.New()

	io.Copy(hash, buff)
	return hex.EncodeToString(hash.Sum(nil)), nil
}
