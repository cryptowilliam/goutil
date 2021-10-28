package grand

import (
	"github.com/Pallinder/go-randomdata"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"math/rand"
)

type (
	Ranges struct {
		ranges [][2]rune
	}
)

func NewRanges() *Ranges {
	return &Ranges{}
}

func (r *Ranges) Allow(min, max rune) *Ranges {
	if min > max {
		min, max = max, min
	}
	r.ranges = append(r.ranges, [2]rune{min, max})
	return r
}

func (r *Ranges) AllowRune(a rune) *Ranges {
	r.ranges = append(r.ranges, [2]rune{a, a})
	return r
}

func (r *Ranges) Generate(length uint64) string {
	s := ""
	for _, rng := range r.ranges {
		for a := rng[0]; a <= rng[1]; a++ {
			s += string(a)
		}
	}

	res := ""
	for i := 0; i < int(length); i++ {
		res += string(s[RandomInt(0, len(s)-1)])
	}

	return res
}

// Note: min & max are included in random output.
func RandomInt(min, max int) int {
	if min == max {
		return min
	}
	return randomdata.Number(min, max+1)
}

// Returns an int >= min, <= max
func Int(min, max int) int {
	return min + rand.Intn(max-min+1)
}

func RandomBool() bool {
	return RandomInt(0, 3) == 1
}

func RandomString(size int) string {
	if size <= 0 {
		panic(gerrors.New("invalid size %d", size))
	}
	return randomdata.RandStringRunes(size)
}

func RandomFloat(min, max, decimalPoint int) float64 {
	return randomdata.Decimal(min, max, decimalPoint)
}

func RandomSimplePassword(minLen, maxLen int) string {
	pwdchars := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	return RandomPassword(pwdchars, minLen, maxLen)
}

func RandomPassword(chars string, minLen, maxLen int) string {
	if len(chars) == 0 || maxLen <= 0 {
		return ""
	}
	if minLen <= 0 {
		minLen = 1
	}

	resultLen := RandomInt(minLen, maxLen)
	result := ""
	charsLen := len(chars)
	for i := 0; i < resultLen; i++ {
		idx := RandomInt(0, charsLen-1)
		result += string(chars[idx])
	}
	return result
}
