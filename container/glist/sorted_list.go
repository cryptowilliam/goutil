package glist

import (
	"container/list"
	"fmt"
)

type SortedList struct {
	list *list.List
}

type sortedElement struct {
	weight int64
	value  interface{}
}

func NewSortedList() *SortedList {
	return &SortedList{
		list: list.New(),
	}
}

func (sl *SortedList) Init() {
	sl.list.Init()
}

func (sl *SortedList) Len() int {
	return sl.list.Len()
}

func (sl *SortedList) Push(weight int64, value interface{}) {
	newValue := sortedElement{weight: weight, value: value}

	if sl.list.Len() == 0 {
		sl.list.PushFront(newValue)
	} else if weight <= sl.list.Front().Value.(sortedElement).weight { /* smallest new head */
		sl.list.PushFront(newValue)
	} else if weight >= sl.list.Back().Value.(sortedElement).weight { /* largest new tail */
		sl.list.PushBack(newValue)
	} else {
		for e := sl.list.Front(); e != nil; e = e.Next() {
			if weight <= e.Value.(sortedElement).weight {
				sl.list.InsertBefore(newValue, e)
				break
			}
		}
	}
}

func (sl *SortedList) PopMinWeight() (weight int64, value interface{}, exist bool) {
	if sl.list.Len() == 0 {
		return 0, nil, false
	}

	node := sl.list.Front()
	weight = node.Value.(sortedElement).weight
	value = node.Value.(sortedElement).value
	sl.list.Remove(node)
	return weight, value, true
}

func (sl *SortedList) PopMaxWeight() (weight int64, value interface{}, exist bool) {
	if sl.list.Len() == 0 {
		return 0, nil, false
	}

	node := sl.list.Back()
	weight = node.Value.(sortedElement).weight
	value = node.Value.(sortedElement).value
	sl.list.Remove(node)
	return weight, value, true
}

func (sl *SortedList) Println() {
	for e := sl.list.Front(); e != nil; e = e.Next() {
		fmt.Printf("<-[weight:%d][val:%v]->", e.Value.(sortedElement).weight, e.Value.(sortedElement).value)
	}
}
