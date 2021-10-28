package gspeed

import (
	"fmt"
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/container/gstring"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// bit size of speed
// 为什么用float64而不用uint64？float64的表示范围比uint64大很多，unint64根本达不到YB级别
type Speed float64

// Map between speed unit of bit and bits size
const (
	_        = iota             // ignore first value by assigning to blank identifier
	Kb Speed = 1 << (10 * iota) // 1 Kb = 1024 bits
	Mb                          // 1 Mb = 1048576 bits
	Gb
	Tb
	Pb
	Eb
	Zb
	Yb
)

// Map between speed unit of byte and bits size
const (
	_        = iota                   // ignore first value by assigning to blank identifier
	KB Speed = 8 * (1 << (10 * iota)) // 1 KB = (8 * 1024) bits
	MB                                // 1 MB = (8 * 1048576) bits
	GB
	TB
	PB
	EB
	ZB
	YB
)

func FromBytesInterval(bytes float64, interval time.Duration) (Speed, error) {
	bytesps := float64(bytes) / (float64(interval) / float64(time.Second))
	return FromBytes(bytesps)
}

func FromBitsInterval(bits float64, interval time.Duration) (Speed, error) {
	bitsps := float64(bits) / (float64(interval) / float64(time.Second))
	return FromBits(bitsps)
}

func FromBytes(size float64) (Speed, error) {
	if size < 0 {
		return Speed(0), gerrors.New("Negative speed error")
	}
	return Speed(size * 8), nil
}

func FromBytesUint64(size uint64) Speed {
	return Speed(float64(size) * 8)
}

func FromBits(size float64) (Speed, error) {
	if size < 0 {
		return Speed(0), gerrors.New("Negative speed error")
	}
	return Speed(size), nil
}

func (s Speed) GetByteSize() float64 {
	return float64(s) / 8.0
}

func (s Speed) GetBitSize() float64 {
	return float64(s)
}

func (s Speed) GreaterThan(s2 Speed) bool {
	return float64(s) > float64(s2)
}

func (s Speed) GreaterThanOrEqual(s2 Speed) bool {
	return float64(s) >= float64(s2)
}

func (s Speed) LessThan(s2 Speed) bool {
	return float64(s) < float64(s2)
}

func (s Speed) LessThanOrEqual(s2 Speed) bool {
	return float64(s) <= float64(s2)
}

func (s Speed) Equals(s2 Speed) bool {
	return float64(s) == float64(s2)
}

func (s Speed) String() string {
	return (&s).StringWithBitUnit()
}

func (s Speed) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", s.String())), nil
}

func (s *Speed) UnmarshalJSON(b []byte) error {
	str := string(b)
	if len(str) <= 1 {
		return gerrors.Errorf("invalid json speed '%s'", s)
	}
	if str[0] != '"' || str[len(str)-1] != '"' {
		return gerrors.Errorf("invalid json speed '%s'", s)
	}
	str = gstring.RemoveHead(str, 1)
	str = gstring.RemoveTail(str, 1)
	speed, err := ParseString(str)
	if err != nil {
		return err
	}
	*s = speed
	return nil
}

func (s *Speed) StringAuto() string {
	b := *s
	switch {
	case b >= YB:
		return fmt.Sprintf("%.2fYB", b/YB)
	case b >= Yb:
		return fmt.Sprintf("%.2fYb", b/Yb)
	case b >= ZB:
		return fmt.Sprintf("%.2fZB", b/ZB)
	case b >= Zb:
		return fmt.Sprintf("%.2fZb", b/Zb)
	case b >= EB:
		return fmt.Sprintf("%.2fEB", b/EB)
	case b >= Eb:
		return fmt.Sprintf("%.2fEb", b/Eb)
	case b >= PB:
		return fmt.Sprintf("%.2fPB", b/PB)
	case b >= Pb:
		return fmt.Sprintf("%.2fPb", b/Pb)
	case b >= TB:
		return fmt.Sprintf("%.2fTB", b/TB)
	case b >= Tb:
		return fmt.Sprintf("%.2fTb", b/Tb)
	case b >= GB:
		return fmt.Sprintf("%.2fGB", b/GB)
	case b >= Gb:
		return fmt.Sprintf("%.2fGb", b/Gb)
	case b >= MB:
		return fmt.Sprintf("%.2fMB", b/MB)
	case b >= Mb:
		return fmt.Sprintf("%.2fMb", b/Mb)
	case b >= KB:
		return fmt.Sprintf("%.2fKB", b/KB)
	case b >= Kb:
		return fmt.Sprintf("%.2fKb", b/Kb)
	}
	return fmt.Sprintf("%.2fB", b)
}

func (s *Speed) StringWithByteUnit() string {
	b := *s
	switch {
	case b >= YB:
		return fmt.Sprintf("%.2fYB", b/YB)
	case b >= ZB:
		return fmt.Sprintf("%.2fZB", b/ZB)
	case b >= EB:
		return fmt.Sprintf("%.2fEB", b/EB)
	case b >= PB:
		return fmt.Sprintf("%.2fPB", b/PB)
	case b >= TB:
		return fmt.Sprintf("%.2fTB", b/TB)
	case b >= GB:
		return fmt.Sprintf("%.2fGB", b/GB)
	case b >= MB:
		return fmt.Sprintf("%.2fMB", b/MB)
	case b >= KB:
		return fmt.Sprintf("%.2fKB", b/KB)
	}
	return fmt.Sprintf("%.2fB", b)
}

func (s *Speed) StringWithBitUnit() string {
	b := *s
	switch {
	case b >= Yb:
		return fmt.Sprintf("%.2fYb", b/Yb)
	case b >= Zb:
		return fmt.Sprintf("%.2fZb", b/Zb)
	case b >= Eb:
		return fmt.Sprintf("%.2fEb", b/Eb)
	case b >= Pb:
		return fmt.Sprintf("%.2fPb", b/Pb)
	case b >= Tb:
		return fmt.Sprintf("%.2fTb", b/Tb)
	case b >= Gb:
		return fmt.Sprintf("%.2fGb", b/Gb)
	case b >= Mb:
		return fmt.Sprintf("%.2fMb", b/Mb)
	case b >= Kb:
		return fmt.Sprintf("%.2fKb", b/Kb)
	}
	return fmt.Sprintf("%.2fb", b)
}

// speed string sample: "2M" "2Mb" "2Mbit" "2Mbits" "2Mbyte" "2Mbytes" "2 Mb" "*/s" "*ps"
func ParseString(speed string) (Speed, error) {
	s := strings.TrimSpace(speed)
	s = strings.Replace(s, " ", "", -1) // Remove space in middle
	if len(s) == 0 {
		return Speed(0), gerrors.New("speed string \"" + speed + "\" is empty")
	}

	// Parse number from head
	nonDigitPos := -1
	for i, v := range s {
		if !unicode.IsDigit(v) && v != '.' {
			nonDigitPos = i
			break
		}
	}
	if nonDigitPos < 1 {
		return Speed(0), gerrors.New("speed string \"" + speed + "\" lack of speed unit part")
	}
	numberStr := s[:nonDigitPos] // NOTICE: end pos char is not included in return string
	number, err := strconv.ParseFloat(numberStr, 64)
	if err != nil {
		return Speed(0), err
	}

	// Parse unit from tail
	unitStr := s[len(numberStr):]
	// Remove per second string from tail
	if gstring.EndWith(strings.ToLower(unitStr), "/s") || gstring.EndWith(strings.ToLower(unitStr), "ps") {
		unitStr = gstring.RemoveTail(unitStr, 2)
	}
	if len(unitStr) == 0 {
		return Speed(0), gerrors.New("speed string \"" + speed + "\" has no unit")
	}
	// Parse "k" "m" "g" "t"...
	unitA := unitStr[0:1] // NOTICE: end pos char is not included in return string
	unitA = strings.ToLower(unitA)
	// Parse "b" "B" "bit" "byte" "bits" "bytes"
	var unitB string
	if len(unitStr) == 1 { // such like "2M", unit string is "M", its length is 1
		unitB = "bit"
	} else if len(unitStr) > 1 {
		unitStr = unitStr[1:]
		if unitStr == "B" || strings.ToLower(unitStr) == "byte" || strings.ToLower(unitStr) == "bytes" {
			unitB = "byte"
		} else if unitStr == "b" || strings.ToLower(unitStr) == "bit" || strings.ToLower(unitStr) == "bits" {
			unitB = "bit"
		} else {
			return Speed(0), gerrors.New("speed string \"" + speed + "\" unit syntax error")
		}
	}

	// Calculate speed
	switch unitA {
	case "k":
		if unitB == "byte" {
			return FromBits(number * KB.GetBitSize())
		} else {
			return FromBits(number * Kb.GetBitSize())
		}
	case "m":
		if unitB == "byte" {
			return FromBits(number * MB.GetBitSize())
		} else {
			return FromBits(number * Mb.GetBitSize())
		}
	case "g":
		if unitB == "byte" {
			return FromBits(number * GB.GetBitSize())
		} else {
			return FromBits(number * Gb.GetBitSize())
		}
	case "t":
		if unitB == "byte" {
			return FromBits(number * TB.GetBitSize())
		} else {
			return FromBits(number * Tb.GetBitSize())
		}
	case "p":
		if unitB == "byte" {
			return FromBits(number * PB.GetBitSize())
		} else {
			return FromBits(number * Pb.GetBitSize())
		}
	case "e":
		if unitB == "byte" {
			return FromBits(number * EB.GetBitSize())
		} else {
			return FromBits(number * Eb.GetBitSize())
		}
	case "z":
		if unitB == "byte" {
			return FromBits(number * ZB.GetBitSize())
		} else {
			return FromBits(number * Zb.GetBitSize())
		}
	case "y":
		if unitB == "byte" {
			return FromBits(number * YB.GetBitSize())
		} else {
			return FromBits(number * Yb.GetBitSize())
		}
	default:
		return Speed(0), gerrors.New("speed string \"" + speed + "\" is invalid speed string")
	}
}
