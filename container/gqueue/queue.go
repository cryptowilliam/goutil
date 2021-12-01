package gqueue

import "container/list"

type (
	Queue struct {
		list *list.List
	}
)

func NewQueue() *Queue {
	return &Queue{list: list.New()}
}

func (q *Queue) Push(v interface{}) {
	q.list.PushBack(v)
}

func (q *Queue) Pop() interface{} {
	result := q.list.Front().Value
	q.list.Remove(q.list.Front())
	return result
}

func (q *Queue) Len() int {
	return q.list.Len()
}