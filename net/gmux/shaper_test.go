package gmux

import (
	"container/heap"
	"testing"
)

func TestShaper(t *testing.T) {
	w1 := writeRequest{prio: 10}
	w2 := writeRequest{prio: 10}
	w3 := writeRequest{prio: 20}
	w4 := writeRequest{prio: 100}
	w5 := writeRequest{prio: (1 << 32) - 1}

	var reqs shaperHeap
	heap.Push(&reqs, w5)
	heap.Push(&reqs, w4)
	heap.Push(&reqs, w3)
	heap.Push(&reqs, w2)
	heap.Push(&reqs, w1)

	var lastPrio = reqs[0].prio
	for len(reqs) > 0 {
		w := heap.Pop(&reqs).(writeRequest)
		if int32(w.prio-lastPrio) < 0 {
			t.Fatal("incorrect shaper priority")
		}

		t.Log("prio:", w.prio)
		lastPrio = w.prio
	}
}
