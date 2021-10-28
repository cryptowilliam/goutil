package gdecimal

import (
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"github.com/cryptowilliam/goutil/container/gnum"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"math"
	"math/big"
	"strconv"
)

/*
Decimal is use to process money.
Money is never a floating point type, float32/float64 are inaccurate data types.

example: save a float64 10.12 into mongodb or mysql, but it may become 10.11999999999999 sometimes when you query from database.


serialize:
BSON/JSON: string
CSV: number
*/

type (
	Decimal decimal.Decimal
)

var (
	Zero = Decimal(decimal.Zero)
	One  = NewFromInt(1)
	N0   = Decimal(decimal.Zero)
	N1   = NewFromInt(1)
)

func NewFromDecimal(d2 decimal.Decimal) Decimal {
	return Decimal(d2)
}

func NewFromFloat32(val float32) Decimal {
	return Decimal(decimal.NewFromFloat32(val))
}

func NewFromFloat64(val float64) Decimal {
	return Decimal(decimal.NewFromFloat(val))
}

// TODO: test required
func NewFromBigInt(val big.Int) Decimal {
	return Decimal(decimal.NewFromBigInt(&val, 0))
}

func NewFromBigFloat(val big.Float) (Decimal, error) {
	origin, err := decimal.NewFromString(val.String())
	if err != nil {
		return Decimal{}, err
	}
	return Decimal(origin), nil
}

func NewFromInt(val int) Decimal {
	return Decimal(decimal.New(int64(val), 0))
}

func NewFromInt8(val int8) Decimal {
	return Decimal(decimal.New(int64(val), 0))
}

func NewFromInt16(val int16) Decimal {
	return Decimal(decimal.New(int64(val), 0))
}

func NewFromInt32(val int32) Decimal {
	return Decimal(decimal.New(int64(val), 0))
}

func NewFromInt64(val int64) Decimal {
	return Decimal(decimal.New(val, 0))
}

func NewFromUint(val uint) Decimal {
	r, _ := decimal.NewFromString(strconv.FormatUint(uint64(val), 10))
	return Decimal(r)
}

func NewFromUint8(val uint8) Decimal {
	r, _ := decimal.NewFromString(strconv.FormatUint(uint64(val), 10))
	return Decimal(r)
}

func NewFromUint16(val uint16) Decimal {
	r, _ := decimal.NewFromString(strconv.FormatUint(uint64(val), 10))
	return Decimal(r)
}

func NewFromUint32(val uint32) Decimal {
	r, _ := decimal.NewFromString(strconv.FormatUint(uint64(val), 10))
	return Decimal(r)
}

func NewFromUint64(val uint64) Decimal {
	r, _ := decimal.NewFromString(strconv.FormatUint(val, 10))
	return Decimal(r)
}

func NewFromString(val string) (Decimal, error) {
	d2, err := decimal.NewFromString(val)
	if err != nil {
		return Zero, err
	}
	return Decimal(d2), nil
}

func NewFromStringEx(val string, prec int) (Decimal, error) {
	f, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return Zero, err
	}
	return NewFromString(strconv.FormatFloat(f, 'f', prec, 64))
}

func (d Decimal) WithPrec(prec int) Decimal {
	d2, _ := NewFromStringEx(d.String(), prec)
	return d2
}

func (d Decimal) raw() decimal.Decimal {
	return decimal.Decimal(d)
}

func (d Decimal) Abs() Decimal {
	if d.IsNegative() {
		return Zero.Sub(d)
	} else {
		return d
	}
}

func (d Decimal) TurnPositiveNegative() Decimal {
	return Zero.Sub(d)
}

func (d Decimal) Add(d2 Decimal) Decimal {
	return Decimal(d.raw().Add(d2.raw()))
}

func (d Decimal) AddInt(n int) Decimal {
	return d.Add(NewFromInt(n))
}

func (d Decimal) AddFloat64(n float64) Decimal {
	return d.Add(NewFromFloat64(n))
}

func (d Decimal) Sub(d2 Decimal) Decimal {
	return Decimal(d.raw().Sub(d2.raw()))
}

func (d Decimal) SubInt(n int) Decimal {
	return d.Sub(NewFromInt(n))
}

func (d Decimal) SubFloat64(n float64) Decimal {
	return d.Sub(NewFromFloat64(n))
}

func (d Decimal) Mul(d2 Decimal) Decimal {
	return Decimal(d.raw().Mul(d2.raw()))
}

func (d Decimal) MulInt(n int) Decimal {
	return d.Mul(NewFromInt(n))
}

func (d Decimal) MulFloat64(n float64) Decimal {
	return d.Mul(NewFromFloat64(n))
}

func (d Decimal) Div(d2 Decimal) Decimal {
	return Decimal(d.raw().Div(d2.raw()))
}

func (d Decimal) DivInt(n int) Decimal {
	return d.Div(NewFromInt(n))
}

func (d Decimal) DivFloat64(n float64) Decimal {
	return d.Div(NewFromFloat64(n))
}

// FIXME： 这个是RoundUp还是RoundDown？
func (d Decimal) DivRound(d2 Decimal, precision int) Decimal {
	return Decimal(d.raw().DivRound(d2.raw(), int32(precision)))
}

func (d Decimal) DivRoundDown(d2 Decimal, precision int) Decimal {
	return Decimal(d.raw().Div(d2.raw()).RoundBank(int32(precision)))
}

func (d Decimal) DivRoundInt(d2 int, precision int) Decimal {
	return d.DivRound(NewFromInt(d2), precision)
}

func (d Decimal) DivRoundFloat64(d2 float64, precision int) Decimal {
	return d.DivRound(NewFromFloat64(d2), precision)
}

func (d Decimal) IsPositive() bool {
	return d.raw().IsPositive()
}

func (d Decimal) IsNegative() bool {
	return d.raw().IsNegative()
}

func (d Decimal) IsZero() bool {
	return d.raw().IsZero()
}

func (d Decimal) GreaterThan(cmp Decimal) bool {
	return d.raw().GreaterThan(cmp.raw())
}

func (d Decimal) GreaterThanInt(cmp int) bool {
	return d.GreaterThan(NewFromInt(cmp))
}

func (d Decimal) GreaterThanFloat64(cmp float64) bool {
	return d.GreaterThan(NewFromFloat64(cmp))
}

func (d Decimal) GreaterThanOrEqual(cmp Decimal) bool {
	return d.raw().GreaterThanOrEqual(cmp.raw())
}

func (d Decimal) GreaterThanOrEqualInt(cmp int) bool {
	return d.GreaterThanOrEqual(NewFromInt(cmp))
}

func (d Decimal) GreaterThanOrEqualFloat64(cmp float64) bool {
	return d.GreaterThanOrEqual(NewFromFloat64(cmp))
}

func (d Decimal) LessThan(cmp Decimal) bool {
	return d.raw().LessThan(cmp.raw())
}

func (d Decimal) LessThanInt(cmp int) bool {
	return d.LessThan(NewFromInt(cmp))
}

func (d Decimal) LessThanFloat64(cmp float64) bool {
	return d.LessThan(NewFromFloat64(cmp))
}

func (d Decimal) LessThanOrEqual(cmp Decimal) bool {
	return d.raw().LessThanOrEqual(cmp.raw())
}

func (d Decimal) LessThanOrEqualInt(cmp int) bool {
	return d.LessThanOrEqual(NewFromInt(cmp))
}

func (d Decimal) LessThanOrEqualFloat64(cmp float64) bool {
	return d.LessThanOrEqual(NewFromFloat64(cmp))
}

func (d Decimal) Equal(cmp Decimal) bool {
	return d.raw().Equal(cmp.raw())
}

func (d Decimal) EqualInt(cmp int) bool {
	return d.Equal(NewFromInt(cmp))
}

func (d Decimal) EqualFloat64(cmp float64) bool {
	return d.Equal(NewFromFloat64(cmp))
}

func (d Decimal) IntPart() int {
	return int(d.raw().IntPart())
}

func (d Decimal) Int64Part() int64 {
	return d.raw().IntPart()
}

func (d Decimal) Float64() float64 {
	r, _ := d.raw().Float64()
	return r
}

func (d Decimal) Float64Ex() (val float64, exact bool) {
	return d.raw().Float64()
}

func (d Decimal) String() string {
	return d.raw().String()
}

func (d *Decimal) SetInt(val int) {
	*d = NewFromInt(val)
}

func (d *Decimal) SetFloat64(val float64) {
	*d = NewFromFloat64(val)
}

// prec就是小数点后最多支持多少位
func (d Decimal) Trunc(prec int, step float64) Decimal {
	if step <= 0 {
		return d
	}
	return NewFromFloat64(math.Trunc(math.Floor(d.Float64()/step)*step*math.Pow10(prec)) / math.Pow10(prec))
}

// min 就是小数点后最多多少位的小数写法，比如，min=0.00001，prec就是5
func (d Decimal) Trunc2(min Decimal, step float64) Decimal {
	if min.Float64() <= 0 {
		return d
	}
	prec := int(math.Log10(N1.Div(min).Float64()))
	//fmt.Println(prec)

	return d.Trunc(prec, step)
}

func Min(first Decimal, args ...Decimal) Decimal {
	min := first

	for _, v := range args {
		if v.LessThan(min) {
			min = v
		}
	}
	return min
}

func Max(first Decimal, args ...Decimal) Decimal {
	max := first

	for _, v := range args {
		if v.GreaterThan(max) {
			max = v
		}
	}
	return max
}

func ToFloat64s(in []Decimal) []float64 {
	var r []float64
	for _, v := range in {
		f64 := v.Float64()
		r = append(r, f64)
	}
	return r
}

func MaxPrec(in []Decimal) int {
	maxPrec := 0
	for _, v := range in {
		prec := v.BitsAfterDecimalPoint()
		if prec > maxPrec {
			maxPrec = prec
		}
	}
	return maxPrec
}

func ToElegantFloat64s(in []Decimal) []gnum.ElegantFloat {
	var r []gnum.ElegantFloat
	maxPrec := 0
	for _, v := range in {
		prec := v.BitsAfterDecimalPoint()
		if prec > maxPrec {
			maxPrec = prec
		}
	}
	for _, v := range in {
		r = append(r, gnum.NewElegantFloat(v.Float64(), maxPrec))
	}
	return r
}

func (d Decimal) BitsAfterDecimalPoint() int {
	/** old implement
	s := d.String()
	if !strings.Contains(s, ".") {
		return 0
	}
	ss := strings.Split(s, ".")
	return len(ss[1])*/
	return int(math.Abs(float64(d.raw().Exponent())))
}

/*
NOTE:
if omitempty required in JSON, use *Decimal in your structure,
no other way to implement this even you change source code MarshalJSON function in shopsprint/decimal
*/
func (d Decimal) MarshalJSON() ([]byte, error) {
	return decimal.Decimal(d).MarshalJSON()
}

// WARNING: this is dangerous for decimals with many digits, since many JSON
// unmarshallers (ex: Javascript's) will unmarshal JSON numbers to IEEE 754
// double-precision floating point numbers, which means you can potentially
// silently lose precision.
func (d Decimal) MarshalJSONWithoutQuotes() ([]byte, error) {
	decimal.MarshalJSONWithoutQuotes = true
	defer func() { decimal.MarshalJSONWithoutQuotes = false }()
	return decimal.Decimal(d).MarshalJSON()
}

func (d *Decimal) UnmarshalJSON(b []byte) error {
	return (*decimal.Decimal)(d).UnmarshalJSON(b)
}

func (d Decimal) MarshalBinary() (data []byte, err error) {
	return decimal.Decimal(d).MarshalBinary()
}

func (d *Decimal) UnmarshalBinary(b []byte) error {
	return (*decimal.Decimal)(d).UnmarshalBinary(b)
}

func (d Decimal) MarshalText() (text []byte, err error) {
	return decimal.Decimal(d).MarshalText()
}

func (d *Decimal) UnmarshalText(text []byte) error {
	return (*decimal.Decimal)(d).UnmarshalText(text)
}

// in BSON, Decimal128 could has 34 numbers after decimal point, it will loss precision if more than 34 numbers after decimal point
// so we serialize Decimal into string, just like in JSON.
// in BSON, omitempty of Decimal works
// reference: https://github.com/hackeryard/configcenter/blob/a9638f554b4cf47fb13c37c5c3611e1e0696fd25/src/common/metadata/time.go#L94
func (d Decimal) MarshalBSONValue() (bsontype.Type, []byte, error) {
	// returned []byte is not []byte(string), but bson value encoded from string
	// so don't use
	// return bsontype.String, []byte(d.String()), nil
	// it is totally error, bson.Unmarshal will report EOF error, it can't recognize []byte(string)

	return bsonx.String(d.String()).MarshalBSONValue()
}

func (d *Decimal) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t != bsontype.String {
		return gerrors.Errorf("bsontype(%s) not allowed in Decimal.UnmarshalBSONValue, only string accept", t.String())
	}
	str, _, ok := bsoncore.ReadString(data)
	if !ok {
		return gerrors.Errorf("decode string, but string not found")
	}
	dec, err := NewFromString(str)
	if err != nil {
		return err
	}
	*d = dec
	return nil
}

// WARN:
// BSON Decimal128 has max 34 number after "."
// but Decimal doesn't have this limit (limited by your memory)
func (d *Decimal) ToBSONDecimal128() (primitive.Decimal128, error) {
	d128, err := primitive.ParseDecimal128(d.String())
	if err != nil {
		return primitive.NewDecimal128(0, 0), err
	}
	return d128, nil
}

func (d *Decimal) ToElegantFloat() gnum.ElegantFloat {
	return gnum.NewElegantFloat(d.Float64(), d.BitsAfterDecimalPoint())
}
