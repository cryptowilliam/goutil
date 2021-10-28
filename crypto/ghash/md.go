package ghash

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"io"
	"os"
)

func Md5File(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()
	// files could be pretty big, lets buffer
	r := bufio.NewReader(f)

	h := md5.New()
	_, err = io.Copy(h, r)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil // return hex.EncodeToString(hash.Sum(nil)), nil
}

func Md5Buf(buf []byte) (string, error) {
	if buf == nil {
		return "", gerrors.New("nil input buf in func Md5Buf")
	}
	hash := md5.New()
	hash.Write(buf)
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func Md5Str(str string) (string, error) {
	if len(str) == 0 {
		return "", gerrors.New("empty str")
	}
	m := md5.New()
	m.Write([]byte(str))
	return hex.EncodeToString(m.Sum(nil)), nil
}
