package gdecimal

import (
	"encoding/json"
	"fmt"
	"github.com/cryptowilliam/goutil/basic/gtest"
	"github.com/cryptowilliam/goutil/encoding/gjson"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

type decimalS1 struct {
	Name  string
	Score Decimal `json:"Score,omitempty"`
}

type decimalS2 struct {
	Name  string
	Score *Decimal `json:"Score,omitempty" bson:"Score,omitempty"`
}

func TestNewFromStringEx(t *testing.T) {
	d, err := NewFromStringEx("12345.1234567", 3)
	if err != nil {
		t.Error(err)
		return
	}
	if d.String() != "12345.123" {
		t.Errorf("NewFromStringEx error1")
		return
	}

	d, err = NewFromStringEx("12345.1234567", 4)
	if err != nil {
		t.Error(err)
		return
	}
	if d.String() != "12345.1235" {
		t.Errorf("NewFromStringEx error2")
		return
	}

	d, err = NewFromStringEx("12345.1234567", 7)
	if err != nil {
		t.Error(err)
		return
	}
	if d.String() != "12345.1234567" {
		t.Errorf("NewFromStringEx error3")
		return
	}

	d, err = NewFromStringEx("12345.1234567", 8)
	if err != nil {
		t.Error(err)
		return
	}
	if d.String() != "12345.1234567" {
		t.Errorf("NewFromStringEx error4")
		return
	}
}

func TestNewDecimalFromInt(t *testing.T) {
	d := NewFromInt(12345)
	if d.String() != "12345" {
		t.Errorf("NewDecimalFromInt error1")
	}

	d = NewFromInt(1234567890123456789)
	if d.String() != "1234567890123456789" {
		t.Errorf("NewDecimalFromInt error2")
	}
}

func TestNewDecimalFromUint(t *testing.T) {
	d := NewFromUint(12345)
	if d.String() != "12345" {
		t.Errorf("NewDecimalFromInt error1")
	}

	d = NewFromUint(1234567890123456789)
	if d.String() != "1234567890123456789" {
		t.Errorf("NewDecimalFromInt error2")
	}
}

func TestDecimal_MarshalJSON(t *testing.T) {
	s := decimalS1{Name: "Bob", Score: Zero}
	if gjson.MarshalStringDefault(s, false) != `{"Name":"Bob","Score":"0"}` {
		t.Errorf("TestDecimal_MarshalJSON error1")
		return
	}

	s2 := decimalS2{Name: "Bob"}
	if gjson.MarshalStringDefault(s2, false) != `{"Name":"Bob"}` {
		t.Errorf("TestDecimal_MarshalJSON error2")
		return
	}
}

func TestDecimal_UnmarshalJSON(t *testing.T) {
	type S struct {
		Name  string
		Score Decimal `json:"Score,omitempty"`
	}

	jsonString := `{"Name":"Tom", "Score":99.9}`
	s := &S{}
	if err := json.Unmarshal([]byte(jsonString), s); err != nil {
		t.Error(err)
		return
	}
	fmt.Println(s)
}

func TestDecimal_Trunc(t *testing.T) {
	d, _ := NewFromString("1.23456789")

	fmt.Println(d.Trunc(8, 0.01))

	if d.Trunc(3, 0.02).String() != "1.22" {
		t.Errorf("Decimal.Trunc error1")
	}
	if d.Trunc(3, 0.03).String() != "1.23" {
		t.Errorf("Decimal.Trunc error2")
	}
	if d.Trunc(3, 0.04).String() != "1.2" {
		t.Errorf("Decimal.Trunc error3")
	}
	if d.Trunc(3, 0.05).String() != "1.2" {
		t.Errorf("Decimal.Trunc error4")
	}
	if d.Trunc(3, 0.06).String() != "1.2" {
		t.Errorf("Decimal.Trunc error5")
	}
	if d.Trunc(3, 0.07).String() != "1.19" {
		t.Errorf("Decimal.Trunc error6")
	}
	if d.Trunc(6, 0.000007).String() != "1.234562" {
		t.Errorf("Decimal.Trunc error7")
	}

	// FIXME 这个例子值得探讨，Trunc是否正确，貌似不太对喔
	d, _ = NewFromString("5141.73181940667768")
	if d.Trunc(8, 0.000001).String() != "5141.73181899" {
		t.Errorf("Decimal.Trunc error8")
	}
}

func TestDecimal_Trunc2(t *testing.T) {
	d, _ := NewFromString("1.23456789")

	fmt.Println(d.Trunc2(NewFromFloat64(0.00001), 0.01))
}

func TestDecimal_MarshalBSON(t *testing.T) {
	d, err := NewFromString("1.23")
	if err != nil {
		t.Error(err)
		return
	}

	type S struct {
		Number Decimal
	}
	s1 := S{
		Number: d,
	}

	buf, err := bson.Marshal(s1)
	if err != nil {
		t.Error(err)
		return
	}

	s2 := new(S)
	if err := bson.Unmarshal(buf, s2); err != nil {
		t.Error(err)
		return
	}
	if s2.Number.Equal(s1.Number) == false {
		t.Errorf("Unmarshal bad result %s vs %s", s1.Number.String(), s2.Number.String())
		return
	}

	s3 := decimalS2{Name: "Bob"}
	if _, err := bson.Marshal(s3); err != nil {
		t.Errorf("TestDecimal_MarshalBSON error2: %s", err.Error())
		return
	}
}

func TestDecimal_MarshalBSON_NilPointer(t *testing.T) {
	s2 := decimalS2{}
	if _, err := bson.Marshal(s2); err != nil {
		t.Errorf("TestDecimal_MarshalBSON_NilPointer error: %s", err.Error())
		return
	}
}

func TestToElegantFloat64s(t *testing.T) {
	var r []Decimal
	r = append(r, NewFromFloat64(0.003))
	r = append(r, NewFromFloat64(0.0004))
	r = append(r, NewFromFloat64(0.00005))
	r = append(r, NewFromFloat64(0.000006))
	r = append(r, NewFromFloat64(0.0000007))
	r = append(r, NewFromFloat64(0.00000008))
	r = append(r, NewFromFloat64(0.000000009))
	r = append(r, NewFromFloat64(0.00000000010))

	efs := ToElegantFloat64s(r)
	for _, v := range efs {
		if len(v.String()) != 12 {
			gtest.PrintlnExit(t, "converted ElegantFloat length should be 12")
		}
	}
}

func TestDecimal_DivRound(t *testing.T) {
	d := NewFromInt(10000)
	fmt.Println(d.DivInt(3).MulInt(3))
	fmt.Println(d.DivRoundInt(3, 30).MulInt(3))
	fmt.Println(d.DivRoundFloat64(0.001436, 30).MulFloat64(0.001436).String())
}

func TestDecimal_Div(t *testing.T) {
	a, _ := decimal.NewFromString("10000")
	b, _ := decimal.NewFromString("0.001436")
	fmt.Println(a.Exponent(), b.Exponent(), a.Div(b).Mul(b).Exponent())
	fmt.Println(a.Div(b).Mul(b).Value())
	fmt.Println(a.Div(b).Mul(b)) // 10000.0000000000000000000308
}

func TestDecimal_BitsAfterDecimalPoint(t *testing.T) {
	cl := gtest.NewCaseList()
	cl.New().Input("1").Expect(0)
	cl.New().Input("12").Expect(0)
	cl.New().Input("123").Expect(0)
	cl.New().Input("1234").Expect(0)
	cl.New().Input("12345").Expect(0)
	cl.New().Input("123456").Expect(0)
	cl.New().Input("1234567").Expect(0)
	cl.New().Input("0.1").Expect(1)
	cl.New().Input("0.12").Expect(2)
	cl.New().Input("0.123").Expect(3)
	cl.New().Input("0.1234").Expect(4)
	cl.New().Input("0.12345").Expect(5)
	cl.New().Input("0.123456").Expect(6)
	cl.New().Input("0.1234567").Expect(7)
	cl.New().Input("0.12345678").Expect(8)
	cl.New().Input("0.123456789").Expect(9)
	cl.New().Input("0.1234567890").Expect(9)
	cl.New().Input("0.1234567891").Expect(10)
	cl.New().Input("0.12345678901").Expect(11)
	cl.New().Input("0.123456789012").Expect(12)
	cl.New().Input("0.1234567890123").Expect(13)
	cl.New().Input("0.12345678901234").Expect(14)
	cl.New().Input("0.123456789012345").Expect(15)
	cl.New().Input("0.1234567890123456").Expect(16)
	cl.New().Input("0.12345678901234567").Expect(17)
	cl.New().Input("0.123456789012345678").Expect(18)
	cl.New().Input("0.1234567890123456789").Expect(19)
	cl.New().Input("0.12345678901234567890").Expect(19)
	cl.New().Input("0.12345678901234567891").Expect(20)
	cl.New().Input("0.123456789012345678901").Expect(21)
	cl.New().Input("0.1234567890123456789012").Expect(22)
	cl.New().Input("0.12345678901234567890123").Expect(23)

	for _, c := range cl.Get() {
		s := c.Inputs[0].(string)
		e := c.Expects[0].(int)
		d, err := NewFromString(s)
		gtest.Assert(t, err)
		if d.BitsAfterDecimalPoint() != e {
			gtest.PrintlnExit(t, "decimal %s expect precision %d, but %d got", d.String(), e, d.BitsAfterDecimalPoint())
		}
	}
}
