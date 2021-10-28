package grange

import "testing"

func TestRangeFilter_AddInt64(t *testing.T) {
	rf := NewRangeFilter()

	rf.AddInt64(0)
	for i := 10; i <= 1000; i++ {
		rf.AddInt64(int64(i))
	}
	for i := 1234; i <= 2345; i++ {
		rf.AddInt64(int64(i))
	}
	for i := 900; i <= 1001; i++ {
		rf.AddInt64(int64(i))
	}
	if rf.String() != `[0,0] [10,1001] [1234,2345]` {
		t.Fatalf("AddInt error, get %s", rf.String())
		return
	}
}

func TestRangeFilter_AddRange(t *testing.T) {
	rf := NewRangeFilter()

	rf.AddRange(Range{2000, 5678})
	rf.AddRange(Range{5679, 6789})
	if rf.String() != `[2000,6789]` {
		t.Fatalf("AddRange error, get %s", rf.String())
		return
	}
}

func TestRange_IsOverlapEx(t *testing.T) {

}

func TestRangeFilter_SubRange(t *testing.T) {
	rf := NewRangeFilter()

	rf.SubRange(NewRange(1, 2))
	if rf.String() != `` {
		t.Fatalf("SubRange error, get %s", rf.String())
		return
	}

	rf.AddRange(NewRange(1, 100))
	rf.SubRange(NewRange(-10, 1))
	if rf.String() != `[2,100]` {
		t.Fatalf("SubRange error, get %s", rf.String())
		return
	}

	rf.SubRange(NewRange(101, 200))
	if rf.String() != `[2,100]` {
		t.Fatalf("SubRange error, get %s", rf.String())
		return
	}

	rf.SubRange(NewRange(50, 50))
	if rf.String() != `[2,49] [51,100]` {
		t.Fatalf("SubRange error, get %s", rf.String())
		return
	}
}

func TestRangeFilter_MinMax(t *testing.T) {
	rf := NewRangeFilter()
	rf.AddRange(Range{-2000, 5678})
	rf.AddRange(Range{5679, 6789})

	minNum, maxNum, ok := rf.MinMax()
	if !ok {
		t.Fatal("should be ok")
		return
	}
	if minNum != -2000 {
		t.Fatalf("min number should be -2000, but get %d", minNum)
		return
	}
	if maxNum != 6789 {
		t.Fatalf("max number should be 6789, but get %d", maxNum)
		return
	}
}

func TestRangeFilter_Lacks(t *testing.T) {
	rf := NewRangeFilter()
	rf.AddRange(Range{-2000, 5678})
	rf.AddRange(Range{5679, 6789})

	lacks := rf.Lacks()
	if lacks.Len() != 0 {
		t.Fatalf("lacks should be empty, but get %s", lacks.String())
		return
	}

	rf.SubRange(NewRange(100, 200))
	lacks = rf.Lacks()
	minNum, maxNum, ok := lacks.MinMax()
	if !ok {
		t.Fatalf("lacks should be empty, but get %s", lacks.String())
		return
	}
	if minNum != 100 {
		t.Fatalf("min number should be -2000, but get %d", minNum)
		return
	}
	if maxNum != 200 {
		t.Fatalf("max number should be -2000, but get %d", maxNum)
		return
	}
}
