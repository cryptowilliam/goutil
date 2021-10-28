package gnum

import (
	"encoding/binary"
	"math"
	"strconv"
	"strings"
)

var (
	PosInf = math.Inf(1)
	NegInf = math.Inf(-1)
)

const (
	invalidPrec = -2
	defaultPrec = -1
)

const float64EqualityThreshold = 1e-9

func FloatAlmostEqual(a, b float64) bool {
	return math.Abs(a-b) <= float64EqualityThreshold
}

// Converts the float64 into an uint64 without changing the bits, it's the way the bits are interpreted that change.
// big endian
// references:
// https://stackoverflow.com/questions/37758267/golang-float64bits
// https://stackoverflow.com/questions/43693360/convert-float64-to-byte-array
func Float64ToByte(f float64) []byte {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], math.Float64bits(f))
	return buf[:]
}

// little endian
func Float64ToByteLE(f float64) []byte {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], math.Float64bits(f))
	return buf[:]
}

func DetectPrecByHumanReadPrec(val float64, humanReadPrec int) int {
	if humanReadPrec <= invalidPrec {
		return defaultPrec
	}

	// example t.val is 2.001263645807, humanReadPrec is 3
	s := strconv.FormatFloat(val, 'f', -1, 64)
	if strings.Index(s, ".") > 0 {
		s = strings.Split(s, ".")[1] // s is "001263645807"
		for i := range s {
			if s[i] != '0' { // i is 2
				return i + humanReadPrec // returns 5
			}
		}
		return humanReadPrec
	} else {
		return humanReadPrec
	}
}

func DetectMaxPrecRaw(vals []float64, humanReadPrec int) int {
	r := defaultPrec
	tmp := defaultPrec
	for _, v := range vals {
		tmp = DetectPrecByHumanReadPrec(v, humanReadPrec)
		if tmp > r {
			r = tmp
		}
	}
	return r
}
