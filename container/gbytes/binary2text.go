package gbytes

import (
	"encoding/base64"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/jbenet/go-base58"
	"strings"
)

const (
	EncodeBase64 = "base64"
	EncodeBase58 = "base58"
)

func Encode(encode string, binary []byte) (text string, err error) {
	encode = strings.ToLower(encode)

	switch encode {
	case "base64":
		return base64.StdEncoding.EncodeToString(binary), nil
	case "base58":
		return base58.Encode(binary), nil
	}

	return "", gerrors.Errorf("unknown encode %s", encode)
}
