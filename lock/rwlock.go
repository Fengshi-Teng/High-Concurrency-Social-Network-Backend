// Package lock provides an implementation of a read-write lock
// that uses condition variables and mutexes.
package lock

import "sync"

type RWlock struct {
	mu             sync.Mutex
	cond           *sync.Cond
	readers        int
	writer         bool
	waitingWriters int // prevent Wlock starving
}

func NewRWLock() *RWlock {
	l := &RWlock{}
	l.cond = sync.NewCond(&l.mu)
	return l
}

// Lock:
func (l *RWlock) Lock() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.waitingWriters++
	for l.writer || l.readers > 0 {
		l.cond.Wait()
	}
	l.waitingWriters--
	l.writer = true
}

// Unlock:
func (l *RWlock) Unlock() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.writer = false
	l.cond.Broadcast()
}

// RLock:
func (l *RWlock) RLock() {
	l.mu.Lock()
	defer l.mu.Unlock()

	for l.writer || l.waitingWriters > 0 || l.readers >= 32 {
		l.cond.Wait()
	}
	l.readers++
}

// RUnlock:
func (l *RWlock) RUnlock() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.readers--
	if l.readers == 0 {
		l.cond.Broadcast() 
	}  else {
		l.cond.Signal() // Wake one waiter if capacity freed (<32 readers)
	}
}
