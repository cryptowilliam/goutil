package gcounter

import (
	"github.com/cryptowilliam/goutil/container/gnum"
	"testing"
)

func TestStepRecycleAccumulator_Incr(t *testing.T) {
	sra, err := NewStepRecycleAccumulator(0, 2, 4)
	if err != nil {
		t.Error(err)
		return
	}

	strcorrect := "000011112222"
	strinfact := ""
	for i := 0; i < 12; i++ {
		n := sra.Get()
		strinfact += gnum.FormatInt64(n)
		sra.Incr()
	}

	if strinfact != strcorrect {
		t.Errorf("Correct return is %s, but returns %s", strcorrect, strinfact)
		return
	}
}
