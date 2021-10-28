package gtime

import (
	"encoding/json"
	"github.com/cryptowilliam/goutil/encoding/gjson"
	"testing"
	"time"
)

func TestHumanDuration_MarshalJSON(t *testing.T) {
	type S struct {
		D HumanDuration
	}
	s := S{
		D: HumanDuration(Day*8 + (time.Minute * 24)),
	}
	jsonStr := gjson.MarshalStringDefault(s, false)
	if jsonStr != `{"D":"1 week 1 day 24 minutes"}` {
		t.Errorf("HumanDuration MarshalJSON error")
		return
	}
}

func TestHumanDuration_UnmarshalJSON(t *testing.T) {
	jsonStr := `{"D":"1 week 1 day 24 minutes"}`

	type S struct {
		D HumanDuration
	}
	s := &S{}
	if err := json.Unmarshal([]byte(jsonStr), s); err != nil {
		t.Error(err)
		return
	}
	expected := Day*8 + (time.Minute * 24)
	if s.D.ToDuration() != expected {
		t.Errorf("HumanDuration UnmarshalJSON error")
		return
	}
}
