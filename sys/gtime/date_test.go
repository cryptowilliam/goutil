package gtime

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestDateValid(t *testing.T) {
	if DateValid(2018, 2, 31) {
		t.Error("DateValid() error")
		return
	}
	if DateValid(2018, 11, 31) {
		t.Error("DateValid() error")
		return
	}
	if !DateValid(2018, 1, 31) {
		t.Error("DateValid() error")
		return
	}
}

func TestDate_IsZero(t *testing.T) {
	date := TimeToDate(time.Now(), time.UTC)
	if date.IsZero() {
		t.Error("time.Now() is NOT ZeroDate")
		return
	}
	fmt.Println(ZeroDate.IntYYYYMMDD())
	fmt.Println(int(ZeroDate))
}

func TestDate_StringYYYYMMDD(t *testing.T) {
	date, err := NewDate(2018, 3, 8)
	if err != nil {
		t.Error(err)
		return
	}
	dateString := date.StringYYYYMMDD()
	expected := "20180308"
	if dateString != expected {
		t.Errorf("Correct date string %s, but get %s", expected, date.StringYYYYMMDD())
		return
	}
}

func TestDate_MarshalJSON(t *testing.T) {
	dt := Date(0)
	b, err := json.Marshal(dt)
	if err != nil {
		t.Error(err)
		return
	}
	if string(b) != "\"0000-00-00\"" {
		t.Errorf("Date json.Marshal(Date{}) error, returns %s", string(b))
		return
	}
}

type test_item struct {
	InDate Date `json:"InDate"`
}

func TestDate_UnmarshalJSON(t *testing.T) {
	s := `{"InDate":"2018-05-01"}`
	i := test_item{}
	err := json.Unmarshal([]byte(s), &i)
	if err != nil {
		t.Error(err)
		return
	}
	if i.InDate.String() != "2018-05-01" {
		t.Errorf("Date json.Unmarshal() error, returns '%s'", i.InDate.String())
		return
	}
}

func TestDate_AfterEqual(t *testing.T) {
	/*maxtime, err := ParseDatetimeStringFuzz("2019-01-21 08:00:00.000+08:00")
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(TimeToDate(maxtime).AfterEqual(Yesterday()))*/

	if !Today(time.UTC).AfterEqual(TimeToDate(time.Now(), time.UTC)) {
		t.Error("AfterEqual error1")
	}
	if Yesterday(time.UTC).AfterEqual(TimeToDate(time.Now(), time.UTC)) {
		t.Error("AfterEqual error2")
	}
}

func TestDate_UnixDays(t *testing.T) {
	d := Date(19700102)
	if d.UnixDays() != 1 {
		t.Errorf("UnixDays error")
	}

	d = Date(19691231)
	if d.UnixDays() != -1 {
		t.Errorf("UnixDays error")
	}

	d = Date(19700101)
	if d.UnixDays() != 0 {
		t.Errorf("UnixDays error")
	}
}

func TestDate_String(t *testing.T) {
	d := Date(19691231)
	if d.String() != "1969-12-31" {
		t.Errorf("Date.String() error, returns %s", d.String())
	}

	d = Date(-19691231)
	if d.String() != "-1969-12-31" {
		t.Errorf("Date.String() error, returns %s", d.String())
	}
}
