package gstring

import (
	"github.com/cryptowilliam/goutil/basic/gtest"
	"strings"
	"testing"
)

const (
	test_utf8_str = " 你好世界0x3a8F2a0032Dc1dfc38914734EFe21ba27893e8C7  "
)

func TestRemoveIndex(t *testing.T) {
	tcs := gtest.NewCaseList()
	tcs.New().Input("a").Input(0).Expect("")
	tcs.New().Input("ab").Input(0).Expect("b")
	tcs.New().Input("abc").Input(0).Expect("bc")
	tcs.New().Input("abc").Input(2).Expect("ab")
	tcs.New().Input("abc").Input(1).Expect("ac")
	tcs.New().Input("abc").Input(-1).Expect("abc")
	tcs.New().Input("abc").Input(3).Expect("abc")

	for _, v := range tcs.Get() {
		str := v.Inputs[0].(string)
		idx := v.Inputs[1].(int)
		exp := v.Expects[0].(string)
		got := RemoveIndex(str, idx)
		if got != exp {
			gtest.PrintlnExit(t, "RemoveIndex(%s, %d) expect %s but get %s", str, idx, exp, got)
		}
	}
}

func TestEndWith(t *testing.T) {
	if EndWith("kline", "-kline") {
		t.Errorf("EndWith kline test error")
		return
	}
	if !EndWith("kline", "line") {
		t.Errorf("EndWith kline test error")
		return
	}
}

func TestIndexUTF8(t *testing.T) {
	if LenUTF8(test_utf8_str) != 49 {
		t.Errorf("test_utf8_str length in utf-8 should be 49")
	}
}

func TestTrySubstrLenUTF8(t *testing.T) {
	if TrySubstrLenUTF8(test_utf8_str, 5, 42) != "0x3a8F2a0032Dc1dfc38914734EFe21ba27893e8C7" {
		t.Errorf("TrySubstrLenUTF8 error")
	}
}

func TestSplitByLen(t *testing.T) {
	real := SplitByLen("abc123ABC!@#$", 3)
	expected := []string{"abc", "123", "ABC", "!@#", "$"}
	if strings.Join(real, ",") != strings.Join(expected, ",") {
		t.Errorf("SplitByLen error1")
		return
	}

	real = SplitByLen("123456", 3)
	expected = []string{"123", "456"}
	if strings.Join(real, ",") != strings.Join(expected, ",") {
		t.Errorf("SplitByLen error2")
		return
	}
}

func TestSplitChunksAscii(t *testing.T) {
	type item struct {
		src       string
		chunksize int
		fromleft  bool
		expect    []string
	}
	items := []item{
		{src: "123", chunksize: 3, fromleft: true, expect: []string{"123"}},
		{src: "123", chunksize: 3, fromleft: false, expect: []string{"123"}},
		{src: "123", chunksize: 4, fromleft: true, expect: []string{"123"}},
		{src: "123", chunksize: 4, fromleft: false, expect: []string{"123"}},
		{src: "1234567", chunksize: 3, fromleft: true, expect: []string{"123", "456", "7"}},
		{src: "1234567", chunksize: 3, fromleft: false, expect: []string{"1", "234", "567"}},
		{src: "123456", chunksize: 3, fromleft: true, expect: []string{"123", "456"}},
		{src: "123456", chunksize: 3, fromleft: false, expect: []string{"123", "456"}},
	}

	for _, v := range items {
		res := SplitChunksAscii(v.src, v.chunksize, v.fromleft)
		if !Equal(res, v.expect) {
			t.Errorf("expect %s, but get %s", v.expect, res)
			return
		}
	}
}

func TestOnlyFirstLetterUpperCase(t *testing.T) {
	if res := OnlyFirstLetterUpperCase("namebuFFER"); res != "Namebuffer" {
		t.Errorf("TestOnlyFirstLetterUpperCase error %s", res)
		return
	}
}

func TestSortByHex(t *testing.T) {
	s := "722abBCcA"
	correctSorted := "227ABCabc"
	r := SortByHex(s)
	if r != correctSorted {
		t.Errorf("%s after sorted %s, but should be %s", s, r, correctSorted)
	}
}
