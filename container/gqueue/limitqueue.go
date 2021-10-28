package gqueue

// 固定大小的queue，达到或者超出大小时，会挤出最旧的元素
// 应用场景：统计最近的信息，太旧的信息没有统计价值的情景。比如速度、访问成功率等。

import "sync"

type LimitQueue struct {
	data  []interface{}
	limit int
	lock  sync.RWMutex
}

func NewLimitQueue(limit int) *LimitQueue {
	return &LimitQueue{limit: limit}
}

func (q *LimitQueue) Add(item interface{}) {
	q.lock.Lock()
	defer q.lock.Unlock()

	if len(q.data) == q.limit {
		q.data = q.data[1:]
	}
	q.data = append(q.data, item)
}

func (q *LimitQueue) Clone() []interface{} {
	q.lock.RLock()
	defer q.lock.RUnlock()
	return q.data
}

func (q *LimitQueue) Size() int {
	q.lock.RLock()
	defer q.lock.RUnlock()
	return len(q.data)
}

func (q *LimitQueue) Clear() {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.data = q.data[:0]
}
