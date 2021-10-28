package gqueue

/*
This is a fork and modification from https://github.com/eapache/queue/blob/master/queue.go
Fork date: 2017-04-09
Why modify, but not import? Because too many panic!
Modified:
1, Change panic to return nil
2, New() -> newRawQueue()
3, type Queue -> type rawQueue
*/

/*
Package queue provides a fast, ring-buffer queue based on the version suggested by Dariusz Górecki.
Using this instead of other, simpler, queue implementations (slice+append or linked list) provides
substantial memory and time benefits, and fewer GC pauses.
The queue implemented here is as fast as it is for an additional reason: it is *not* thread-safe.
*/

// minQueueLen is smallest capacity that queue may have.
// Must be power of 2 for bitwise modulus: x % n == x & (n - 1).
const minQueueLen = 16

// rawQueue represents a single instance of the queue data structure.
type rawQueue struct {
	buf               []interface{}
	head, tail, count int
}

// New constructs and returns a new rawQueue.
func newRawQueue() *rawQueue {
	return &rawQueue{
		buf: make([]interface{}, minQueueLen),
	}
}

// DataLen returns the number of elements currently stored in the queue.
func (q *rawQueue) Length() int {
	return q.count
}

// resizes the queue to fit exactly twice its current contents
// this can result in shrinking if the queue is less than half-full
func (q *rawQueue) resize() {
	newBuf := make([]interface{}, q.count<<1)

	if q.tail > q.head {
		copy(newBuf, q.buf[q.head:q.tail])
	} else {
		n := copy(newBuf, q.buf[q.head:])
		copy(newBuf[n:], q.buf[:q.tail])
	}

	q.head = 0
	q.tail = q.count
	q.buf = newBuf
}

// Add puts an element on the end of the queue.
func (q *rawQueue) Add(elem interface{}) {
	if q.count == len(q.buf) {
		q.resize()
	}

	q.buf[q.tail] = elem
	// bitwise modulus
	q.tail = (q.tail + 1) & (len(q.buf) - 1)
	q.count++
}

// Peek returns the element at the head of the queue. This call panics
// if the queue is empty.
func (q *rawQueue) Peek() interface{} {
	if q.count <= 0 {
		return nil // panic("queue: Peek() called on empty queue")
	}
	return q.buf[q.head]
}

// Get returns the element at index i in the queue. If the index is
// invalid, the call will panic. This method accepts both positive and
// negative index values. Index 0 refers to the first element, and
// index -1 refers to the last.
func (q *rawQueue) Get(i int) interface{} {
	// If indexing backwards, convert to positive index.
	if i < 0 {
		i += q.count
	}
	if i < 0 || i >= q.count {
		return nil //panic("queue: Get() called with index out of range")
	}
	// bitwise modulus
	return q.buf[(q.head+i)&(len(q.buf)-1)]
}

// Remove removes and returns the element from the front of the queue. If the
// queue is empty, the call will panic.
func (q *rawQueue) Remove() interface{} {
	if q.count <= 0 {
		return nil // panic("queue: Remove() called on empty queue")
	}
	ret := q.buf[q.head]
	q.buf[q.head] = nil
	// bitwise modulus
	q.head = (q.head + 1) & (len(q.buf) - 1)
	q.count--
	// Resize down if buffer 1/4 full.
	if len(q.buf) > minQueueLen && (q.count<<2) == len(q.buf) {
		q.resize()
	}
	return ret
}

// Clean all elements
func (q *rawQueue) Empty() {
	q.count = 0
	q.head = 0
	q.tail = 0
}
