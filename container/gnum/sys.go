package gnum

import (
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/container/gstring"
	"math"
	"math/big"
)

// Process any custom number systems.
// TODO: 改成strconv中去实现会不会可靠？放到Decimal中是否可行？

const (
	defaultBaseString = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func GetDefaultDigits() []rune {
	return gstring.SplitRunes(defaultBaseString)
}

func parse(n *big.Int, radix int, cpt []int64) []int64 {
	tail := new(big.Int).Mod(n, big.NewInt(int64(radix)))
	cpt = append(cpt, tail.Int64())
	integer := new(big.Int).Div(n, big.NewInt(int64(radix)))
	if integer.Cmp(big.NewInt(0)) > 0 {
		return parse(integer, radix, cpt)
	} else {
		return cpt
	}
}

// Format number to string in specified number system.
// Length of 'numberSystemDigits' is number system length.
// Reference:
// https://github.com/kenticny/numconvert/blob/master/converter.go
func FormatAnySys(numberSystemDigits []rune, n big.Int) string {
	radix := len(numberSystemDigits)
	nAbs := new(big.Int).Abs(&n)
	parsed := parse(nAbs, radix, []int64{})
	res := ""
	for i := len(parsed); i > 0; i-- {
		idx := parsed[i-1]
		res += string(numberSystemDigits[idx])
	}
	if n.Cmp(big.NewInt(0)) >= 0 {
		return res
	} else {
		return "-" + res
	}
}

// Parse string as specified number system.
func ParseAnySys(numberSystemDigits []rune, s string) (*big.Int, error) {
	if len(numberSystemDigits) == 0 {
		return nil, gerrors.New("empty number system digits")
	}
	if len(s) == 0 || s == "-" {
		return nil, gerrors.New("empty number string")
	}

	isPositive := true
	if s[0] == '-' {
		isPositive = false
		s = s[1:]
	}

	base := len(numberSystemDigits)
	res := big.NewInt(0)
	for i := 0; i < len(s); i++ {
		currN := (*uint64)(nil) // What is the value of current byte in specified number system.
		for j := uint64(0); j < uint64(len(numberSystemDigits)); j++ {
			if rune(s[i]) == numberSystemDigits[j] {
				currN = &j
				break
			}
		}
		if currN == nil {
			return nil, gerrors.New("can't parse char '%c'", s[i])
		}
		res = res.Add(res, big.NewInt(int64(math.Pow(float64(base), float64(len(s)-i-1)) * float64(*currN))))
	}

	if !isPositive {
		res = res.Neg(res)
	}
	return res, nil
}
