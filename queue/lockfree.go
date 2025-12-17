package queue

import (
	"sync/atomic"
)

type Request struct {
	Command   string
	Id        int
	Body      string
	TimeStamp float64
}

type node struct {
	value *Request
	next  atomic.Pointer[node]
}

// LockfreeQueue represents a FIFO structure with operations to enqueue
// and dequeue tasks represented as Request
type LockFreeQueue struct {
	head atomic.Pointer[node]
	tail atomic.Pointer[node]
}

// NewQueue creates and initializes a LockFreeQueue
func NewLockFreeQueue() *LockFreeQueue {
	dummy := &node{}
	q := &LockFreeQueue{}
	q.head.Store(dummy)
	q.tail.Store(dummy)
	return q
}

// Enqueue adds a series of Request to the queue
func (queue *LockFreeQueue) Enqueue(task *Request) {
	n := &node{value: task}
	for {
		tail := queue.tail.Load()
		next := tail.next.Load()

		if next == nil {
			if tail.next.CompareAndSwap(nil, n) {
				queue.tail.CompareAndSwap(tail, n)
				return
			}
		} else {
			// help advance tail if lagging
			queue.tail.CompareAndSwap(tail, next)
		}
	}
}

// Dequeue removes a Request from the queue
func (queue *LockFreeQueue) Dequeue() *Request {
	for {
		head := queue.head.Load()
		tail := queue.tail.Load()
		next := head.next.Load()

		if next == nil {
			return nil // empty
		}

		// help advance tail if lagging
		if head == tail {
			queue.tail.CompareAndSwap(tail, next)
			continue
		}

		val := next.value
		if queue.head.CompareAndSwap(head, next) {
			return val
		}
	}
}


func (q *LockFreeQueue) IsEmpty() bool {
    head := q.head.Load()
    next := head.next.Load()
    return next == nil
}
