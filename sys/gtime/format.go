package gtime

import (
	"github.com/cryptowilliam/goutil/container/gstring"
	"time"
)

var (
	AllLayouts []string
)

var (
	LayoutMM_DD_HH               = registerLayout("15:04:05")
	LayoutYYYY                   = registerLayout("2006")
	LayoutYYYY_MM                = registerLayout("2006-01")
	LayoutYYYY_MM_DD             = registerLayout("2006-01-02")
	LayoutYYYY_MM_DD_HH          = registerLayout("2006-01-02 15")
	LayoutYYYY_MM_DD_HH_mm       = registerLayout("2006-01-02 15:04")
	LayoutYYYY_MM_DD_HH_mm_SS    = registerLayout("2006-01-02 15:04:05")
	LayoutYYYY_MM_DD_HH_mm_SS_NS = registerLayout("2006-01-02 15:04:05.999999999")
	LayoutRFC3339GoExtension     = registerLayout("2006-01-02 15:04:05.999999999 -0700 MST") // Golang custom time layout in time.String(), it is very close to RFC3339.
	LayoutRFC3339Milli           = registerLayout("2006-01-02T15:04:05.999Z")                // RFC 3339 with milliseconds
	LayoutANSIC                  = registerLayout(time.ANSIC)                                // "Mon Jan _2 15:04:05 2006"
	LayoutUnixDate               = registerLayout(time.UnixDate)                             // "Mon Jan _2 15:04:05 MST 2006"
	LayoutRubyDate               = registerLayout(time.RubyDate)                             // "Mon Jan 02 15:04:05 -0700 2006"
	LayoutRFC822                 = registerLayout(time.RFC822)                               // "02 Jan 06 15:04 MST"
	LayoutRFC822Z                = registerLayout(time.RFC822Z)                              // "02 Jan 06 15:04 -0700" RFC822 with numeric zone
	LayoutRFC850                 = registerLayout(time.RFC850)                               // "Monday, 02-Jan-06 15:04:05 MST"
	LayoutRFC1123                = registerLayout(time.RFC1123)                              // "Mon, 02 Jan 2006 15:04:05 MST"
	LayoutRFC1123Z               = registerLayout(time.RFC1123Z)                             // "Mon, 02 Jan 2006 15:04:05 -0700" RFC1123 with numeric zone
	LayoutRFC3339                = registerLayout(time.RFC3339)                              // "2006-01-02T15:04:05Z07:00"
	LayoutRFC3339Nano            = registerLayout(time.RFC3339Nano)                          // "2006-01-02T15:04:05.999999999Z07:00"
	LayoutKitchen                = registerLayout(time.Kitchen)                              // "3:04PM"
	LayoutStamp                  = registerLayout(time.Stamp)                                // "Jan _2 15:04:05"
	LayoutStampMilli             = registerLayout(time.StampMilli)                           // "Jan _2 15:04:05.000"
	LayoutStampMicro             = registerLayout(time.StampMicro)                           // "Jan _2 15:04:05.000000"
	LayoutStampNano              = registerLayout(time.StampNano)                            // "Jan _2 15:04:05.000000000"
)

func registerLayout(s string) string {
	AllLayouts = append(AllLayouts, s)
	return s
}

func (t Time) Std() time.Time {
	return time.Time(t)
}

func (t Time) Elegant(layout string) ElegantTime {
	return NewElegantTime(t.Std(), layout)
}

func TimeToIntYYYYMMDDHHMM(t time.Time) int {
	return (t.Year() * 100000000) + (int(t.Month()) * 1000000) + (t.Day() * 10000) + (t.Hour() * 100) + t.Minute()
}

// ElegantTime is the time.Time with JSON marshal and unmarshal capability
type ElegantTime struct {
	val    time.Time
	layout string
}

func NewElegantTime(tm time.Time, layout string) ElegantTime {
	if layout == "" {
		layout = LayoutRFC3339GoExtension
	}
	return ElegantTime{val: tm, layout: layout}
}

func NewElegantTimeArray(tms []time.Time, layout string) []ElegantTime {
	var r []ElegantTime
	for _, v := range tms {
		r = append(r, NewElegantTime(v, layout))
	}
	return r
}

func (t *ElegantTime) Raw() time.Time {
	return t.val
}

func (t *ElegantTime) SetLayout(layout string) {
	t.layout = layout
}

// UnmarshalJSON will unmarshal using 2006-01-02T15:04:05+07:00 layout
func (t *ElegantTime) UnmarshalJSON(b []byte) error {
	val, err := ParseDatetimeStringFuzz(string(b))
	if err != nil {
		return err
	}

	t.val = val
	t.layout = LayoutRFC3339GoExtension
	return nil
}

// MarshalJSON will marshal using 2006-01-02T15:04:05+07:00 layout
func (t *ElegantTime) MarshalJSON() ([]byte, error) {
	return t.JSON()
}

func (t *ElegantTime) JSONAutoDetect() ([]byte, error) {
	s := t.val.Format(t.DetectBestLayout())
	return []byte(`"` + s + `"`), nil
}

func (t *ElegantTime) JSON() ([]byte, error) {
	if t.layout == "" {
		t.layout = LayoutRFC3339GoExtension
	}
	s := t.val.Format(t.layout)
	return []byte(`"` + s + `"`), nil
}

func (t *ElegantTime) DetectBestLayout() string {
	hasMonth := t.val.Month() != 1
	hasDay := t.val.Day() != 1
	hasHour := t.val.Hour() != 0
	hasMinute := t.val.Minute() != 0
	hasSecond := t.val.Second() != 0
	hasNanosecond := t.val.Nanosecond() != 0
	hasNonUTCTimeZone := t.val.Location() != time.UTC

	// count start with month
	cswNanosecond := 0
	if hasNanosecond {
		cswNanosecond++
	}

	cswSecond := 0
	if hasSecond {
		cswSecond++
	}
	cswSecond += cswNanosecond

	cswMinute := 0
	if hasMinute {
		cswMinute++
	}
	cswMinute += cswSecond

	cswHour := 0
	if hasHour {
		cswHour++
	}
	cswHour += cswMinute

	cswDay := 0
	if hasDay {
		cswDay++
	}
	cswDay += cswHour

	cswMonth := 0
	if hasMonth {
		cswMonth++
	}
	cswMonth += cswDay

	layoutTZ := ""
	if hasNonUTCTimeZone {
		layoutTZ = " -0700 MST"
	}
	if !hasMonth && !hasDay && !hasHour && !hasMinute && !hasSecond && !hasNanosecond {
		return LayoutYYYY + layoutTZ
	} else if !hasDay && !hasHour && !hasMinute && !hasSecond && !hasNanosecond {
		return LayoutYYYY_MM + layoutTZ
	} else if !hasHour && !hasMinute && !hasSecond && !hasNanosecond {
		return LayoutYYYY_MM_DD + layoutTZ
	} else if !hasMinute && !hasSecond && !hasNanosecond {
		return LayoutYYYY_MM_DD_HH + layoutTZ
	} else if !hasSecond && !hasNanosecond {
		return LayoutYYYY_MM_DD_HH_mm + layoutTZ
	} else if !hasNanosecond {
		return LayoutYYYY_MM_DD_HH_mm_SS + layoutTZ
	} else {
		return LayoutYYYY_MM_DD_HH_mm_SS_NS + layoutTZ
	}
}

func DetectBestLayout(in []ElegantTime) string {
	if in == nil || len(in) == 0 {
		return ""
	}

	LAYOUTTIMEZONE := " -0700 MST"

	layoutHead := ""
	layoutTimeZone := ""
	for _, v := range in {
		tmp := v.DetectBestLayout()
		tmpHead := tmp
		tmpTZ := ""
		if gstring.EndWith(tmp, LAYOUTTIMEZONE) {
			tmpHead = gstring.RemoveTail(tmp, len(LAYOUTTIMEZONE))
			tmpTZ = LAYOUTTIMEZONE
		}

		if len(tmpHead) > len(layoutHead) {
			layoutHead = tmpHead
		}
		if len(tmpTZ) > 0 {
			layoutTimeZone = LAYOUTTIMEZONE
		}
	}
	return layoutHead + layoutTimeZone
}

func DetectBestLayoutRaw(in []time.Time) string {
	var inET []ElegantTime
	for _, v := range in {
		inET = append(inET, NewElegantTime(v, ""))
	}
	return DetectBestLayout(inET)
}
