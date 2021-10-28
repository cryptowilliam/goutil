package gnum

import (
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/kenshaw/baseconv"
	"strconv"
)

func baseNumToString(base int) (string, error) {
	switch base {
	case 2:
		return baseconv.DigitsBin, nil
	case 8:
		return baseconv.DigitsOct, nil
	case 10:
		return baseconv.DigitsDec, nil
	case 16:
		return baseconv.DigitsHex, nil
	case 36:
		return baseconv.Digits36, nil
	case 62:
		return baseconv.Digits62, nil
	case 64:
		return baseconv.Digits64, nil
	default:
		return "", gerrors.New("Unsupported base " + strconv.FormatInt(int64(base), 10))
	}
}

func BaseConvert(num string, fromBase, toBase int) (string, error) {
	f, e := baseNumToString(fromBase)
	if e != nil {
		return "", e
	}
	t, e := baseNumToString(toBase)
	if e != nil {
		return "", e
	}
	return baseconv.Convert(num, f, t)
}
