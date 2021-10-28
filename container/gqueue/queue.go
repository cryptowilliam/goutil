package gqueue

import (
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"sync"
	"time"
)

type Queue struct {
	inq            *rawQueue
	inqMu          sync.RWMutex
	concurrentSafe bool
	capacity       int
}

// capacity <= 0: unlimited capacity
func NewQueue(concurrentSafe bool, capacity int) *Queue {
	var result Queue
	result.inq = newRawQueue()
	result.concurrentSafe = concurrentSafe
	result.capacity = capacity
	return &result
}

func (q *Queue) PushNoWait(x interface{}) (err error) {
	if q.concurrentSafe {
		q.inqMu.Lock()
	}
	if q.capacity <= 0 /* Unlimited capacity */ || q.inq.Length() < q.capacity {
		q.inq.Add(x)
		err = nil
	} else {
		err = gerrors.New("queue is full")
	}
	if q.concurrentSafe {
		q.inqMu.Unlock()
	}
	return err
}

func (q *Queue) PushWait(x interface{}, timeout time.Duration) error {
	if q.capacity <= 0 { // Unlimited capacity queue -> nowait
		if q.concurrentSafe {
			q.inqMu.Lock()
		}
		q.inq.Add(x)
		if q.concurrentSafe {
			q.inqMu.Unlock()
		}
		return nil
	} else { // Limited capacity queue -> maybe need wait
		ts := time.Now()
		for {
			// Able to push new item
			if q.concurrentSafe {
				q.inqMu.Lock()
			}
			if q.inq.Length() <= q.capacity-1 {
				q.inq.Add(x)
				if q.concurrentSafe {
					q.inqMu.Unlock()
				}
				return nil
			}
			if q.concurrentSafe {
				q.inqMu.Unlock()
			}

			// Unable to push new item
			if timeout <= 0 || time.Now().Sub(ts) < timeout { // Still loop
				time.Sleep(time.Millisecond)
				continue
			} else { // Loop to much time
				return gerrors.New("Push timeout")
			}
		}
	}
}

func (q *Queue) PopNoWait() (x interface{}) {
	if q.concurrentSafe {
		q.inqMu.Lock()
	}
	if q.inq.Length() > 0 {
		x = q.inq.Remove()
	} else {
		x = nil
	}
	if q.concurrentSafe {
		q.inqMu.Unlock()
	}
	return x
}

func (q *Queue) PopWait(timeout time.Duration) interface{} {
	ts := time.Now()
	for {
		if q.concurrentSafe {
			q.inqMu.Lock()
		}
		if q.inq.Length() > 0 {
			x := q.inq.Remove()
			if q.concurrentSafe {
				q.inqMu.Unlock()
			}
			return x
		} else {
			if q.concurrentSafe {
				q.inqMu.Unlock()
			}
			if timeout <= 0 || time.Now().Sub(ts) < timeout {
				time.Sleep(time.Millisecond)
				continue
			} else {
				return nil
			}
		}
	}
}

func (q *Queue) Size() int {
	if q.concurrentSafe {
		q.inqMu.RLock()
	}
	s := q.inq.Length()
	if q.concurrentSafe {
		q.inqMu.RUnlock()
	}
	return s
}

func (q *Queue) Clear() {
	if q.concurrentSafe {
		q.inqMu.Lock()
	}
	q.inq.Empty()
	if q.concurrentSafe {
		q.inqMu.Unlock()
	}
}
