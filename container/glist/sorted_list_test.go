package glist

import "testing"

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("%s != %s", a, b)
	}
}

func TestSortedList_PopMinKey(t *testing.T) {
	list := NewSortedList()
	list.Push(10, "ten")
	list.Push(15, "fifteen")
	list.Push(5, "five")
	list.Push(9, "nine")
	k, v, exist := list.PopMinWeight()
	assertEqual(t, exist, true)
	assertEqual(t, k, int64(5))
	assertEqual(t, v.(string), "five")
	assertEqual(t, list.Len(), 3)

	list.Println()
}

func TestSortedList_PopMaxKey(t *testing.T) {
	list := NewSortedList()
	list.Push(10, "ten")
	list.Push(15, "fifteen")
	list.Push(5, "five")
	list.Push(9, "nine")
	k, v, exist := list.PopMaxWeight()
	assertEqual(t, exist, true)
	assertEqual(t, k, int64(15))
	assertEqual(t, v.(string), "fifteen")
	assertEqual(t, list.Len(), 3)

	list.Println()
}
