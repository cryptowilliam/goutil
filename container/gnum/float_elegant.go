package gnum

import (
	"math"
	"strconv"
)

// ElegantFloat is used to format float to better looking and enough precision string
// it supports NaN when json Marshal/Unmarshal

// original            -> prec: 2       -> humanReadPrec: 2
// 37                  -> 37.00         -> 37.00
// 12237.89374         -> 12237.89      -> 12237.89
// 3.3483300000000003  -> 3.35          -> 3.35
// 0.00883300000000003 -> 0.01          -> 0.0088
// 0.000012800003      -> 0.00          -> 0.000013
type ElegantFloat struct {
	val          float64
	prec         int
	fmtNaNasNull bool // false default. true: format NaN to "NaN"; false: format NaN to null, null is valid in javascript number calc like echarts
}

func NewElegantFloat(val float64, prec int) ElegantFloat {
	return ElegantFloat{val: val, prec: prec}
}

func NewElegantFloatArray(vals []float64, prec int) []ElegantFloat {
	var r []ElegantFloat
	for _, v := range vals {
		r = append(r, NewElegantFloat(v, prec))
	}
	return r
}

func NewElegantFloatPtrArray(vals []*float64, prec int) []*ElegantFloat {
	var r []*ElegantFloat
	for _, v := range vals {
		if v == nil {
			r = append(r, nil)
		} else {
			nf := NewElegantFloat(*v, prec)
			r = append(r, &nf)
		}
	}
	return r
}

func NewElegantFloatPtrArray2(vals []float64, prec int) []*ElegantFloat {
	var r []*ElegantFloat
	for _, v := range vals {
		nf := NewElegantFloat(v, prec)
		r = append(r, &nf)
	}
	return r
}

func NewElegantFloatPtrArray3(vals []float64, prec int, nilVal float64) []*ElegantFloat {
	var r []*ElegantFloat
	for _, v := range vals {
		if v == nilVal {
			r = append(r, nil)
		} else {
			nf := NewElegantFloat(v, prec)
			r = append(r, &nf)
		}
	}
	return r
}

func ElegantFloatArrayToFloatArray(in []ElegantFloat) []float64 {
	var r []float64
	for i := range in {
		r = append(r, in[i].val)
	}
	return r
}

func DetectMaxPrec(vals []ElegantFloat, humanReadPrec int) int {
	r := defaultPrec
	tmp := defaultPrec
	for _, v := range vals {
		tmp = DetectPrecByHumanReadPrec(v.val, humanReadPrec)
		if tmp > r {
			r = tmp
		}
	}
	return r
}

func (t *ElegantFloat) SetPrec(prec int) {
	if prec > invalidPrec {
		t.prec = prec
	}
}

func (t *ElegantFloat) SetHumanReadPrec(humanReadPrec int) {
	t.prec = DetectPrecByHumanReadPrec(t.val, humanReadPrec)
}

func (t *ElegantFloat) Raw() float64 {
	return t.val
}

// TODO: NaN +-Inf会被ParseFloat64认为是错误吗，需要测试一下
// UnmarshalJSON will unmarshal using 2006-01-02T15:04:05+07:00 layout
func (t *ElegantFloat) UnmarshalJSON(b []byte) error {
	val, err := strconv.ParseFloat(string(b), 64)
	if err != nil {
		switch string(b) {
		case "NaN":
			t.val = math.NaN()
			return nil
		case "+Inf":
			t.val = math.Inf(1)
			return nil
		case "-Inf":
			t.val = math.Inf(-1)
			return nil
		default:
			return err
		}
	}

	t.val = val
	t.prec = invalidPrec
	return nil
}

// MarshalJSON will marshal using 2006-01-02T15:04:05+07:00 layout
//
// FIXME
// if t.val = math.NaN and t.JSON(false), json.Marshal(t) will output error:
// json: error calling MarshalJSON for type *gnum.ElegantFloat: invalid character 'N' looking for beginning of value
// why and how to fix?
func (t ElegantFloat) MarshalJSON() ([]byte, error) {
	return t.JSON(t.fmtNaNasNull)
}

// What will happen if value is math.NaN? it will output bytes buffer `"NaN"`
func (t *ElegantFloat) JSON(fmtNaNasNull bool) ([]byte, error) {
	if t.prec <= invalidPrec {
		t.prec = defaultPrec
	}

	if math.IsNaN(t.val) {
		s := ""
		if fmtNaNasNull {
			s = "null"
		} else {
			s = strconv.FormatFloat(t.val, 'f', t.prec, 64)
		}
		//return []byte(`"` + s + `"`), nil
		//fmt.Println("->", s)
		return []byte(s), nil
	} else if math.IsInf(t.val, -1) || math.IsInf(t.val, 1) {
		s := strconv.FormatFloat(t.val, 'f', t.prec, 64)
		return []byte(s), nil
	} else {
		s := strconv.FormatFloat(t.val, 'f', t.prec, 64)
		return []byte(s), nil
	}
}

func (t *ElegantFloat) String() string {
	return strconv.FormatFloat(t.val, 'f', t.prec, 64)
}
