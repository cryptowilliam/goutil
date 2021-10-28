package guuid

import (
	"github.com/google/uuid"
	"math/big"
	"strings"
)

func NewString(keepDash, upper bool) string {
	s := uuid.New().String()
	if !keepDash {
		s = strings.Replace(s, "-", "", -1)
	}
	if upper {
		return strings.ToUpper(s)
	} else {
		return strings.ToLower(s)
	}
}

// New UUID and convert it to big integer.
func NewBigInt() *big.Int {
	res := new(big.Int)
	res.SetString(strings.Replace(uuid.New().String(), "-", "", -1), 16)
	return res
}
