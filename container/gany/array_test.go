package gany

import (
	"fmt"
	"testing"
)

func TestSplitByBlockCount(t *testing.T) {
	s := []interface{}{1, 2, 3, 4, 5}
	ss := SplitByBlockCount(s, 0)
	if fmt.Sprintf("%v", ss) != "[]" {
		t.Errorf("SplitByBlockCount error1")
		return
	}

	s = []interface{}{1, 2, 3, 4, 5}
	ss = SplitByBlockCount(s, 1)
	if fmt.Sprintf("%v", ss) != "[[1 2 3 4 5]]" {
		t.Errorf("SplitByBlockCount error2")
		return
	}

	s = []interface{}{1, 2, 3, 4, 5}
	ss = SplitByBlockCount(s, 2)
	if fmt.Sprintf("%v", ss) != "[[1 2 3] [4 5]]" {
		t.Errorf("SplitByBlockCount error3")
		return
	}

	s = []interface{}{1, 2, 3, 4, 5}
	ss = SplitByBlockCount(s, 3)
	if fmt.Sprintf("%v", ss) != "[[1 2] [3 4] [5]]" {
		t.Errorf("SplitByBlockCount error4")
		return
	}

	s = []interface{}{1, 2, 3, 4, 5, 6}
	ss = SplitByBlockCount(s, 3)
	if fmt.Sprintf("%v", ss) != "[[1 2] [3 4] [5 6]]" {
		t.Errorf("SplitByBlockCount error5")
		return
	}

	s = []interface{}{1, 2, 3, 4, 5, 6, 7}
	ss = SplitByBlockCount(s, 4)
	if fmt.Sprintf("%v", ss) != "[[1 2] [3 4] [5 6] [7]]" {
		t.Errorf("SplitByBlockCount error6")
		return
	}
}

func TestSplitByBlockSize(t *testing.T) {
	s := []interface{}{1, 2, 3, 4, 5}
	ss := SplitByBlockSize(s, 0)
	if fmt.Sprintf("%v", ss) != "[[1 2 3 4 5]]" {
		t.Errorf("TestSplitByBlockSize error1")
		return
	}

	s = []interface{}{1, 2, 3, 4, 5}
	ss = SplitByBlockSize(s, 1)
	if fmt.Sprintf("%v", ss) != "[[1] [2] [3] [4] [5]]" {
		t.Errorf("TestSplitByBlockSize error2")
		return
	}

	s = []interface{}{1, 2, 3, 4, 5}
	ss = SplitByBlockSize(s, 2)
	if fmt.Sprintf("%v", ss) != "[[1 2] [3 4] [5]]" {
		t.Errorf("TestSplitByBlockSize error3")
		return
	}

	s = []interface{}{1, 2, 3, 4, 5, 6}
	ss = SplitByBlockSize(s, 2)
	if fmt.Sprintf("%v", ss) != "[[1 2] [3 4] [5 6]]" {
		t.Errorf("TestSplitByBlockSize error4")
		return
	}

	s = []interface{}{1, 2, 3, 4, 5, 6, 7}
	ss = SplitByBlockSize(s, 2)
	if fmt.Sprintf("%v", ss) != "[[1 2] [3 4] [5 6] [7]]" {
		t.Errorf("TestSplitByBlockSize error5")
		return
	}

	s = []interface{}{1, 2, 3, 4, 5, 6, 7, 8}
	ss = SplitByBlockSize(s, 2)
	if fmt.Sprintf("%v", ss) != "[[1 2] [3 4] [5 6] [7 8]]" {
		t.Errorf("TestSplitByBlockSize error6")
		return
	}

	s = []interface{}{1, 2, 3, 4, 5, 6, 7, 8}
	ss = SplitByBlockSize(s, 3)
	if fmt.Sprintf("%v", ss) != "[[1 2 3] [4 5 6] [7 8]]" {
		t.Errorf("TestSplitByBlockSize error7")
		return
	}
}
