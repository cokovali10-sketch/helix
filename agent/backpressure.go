// Package agent provides backpressure
package agent

import (
	"sync"
	"time"
)

// Backpressure manages backpressure
type Backpressure struct {
	limit    int
	current  int
	mu       sync.Mutex
	waiters  chan struct{}
}

// NewBackpressure creates backpressure
func NewBackpressure(limit int) *Backpressure {
	return &Backpressure{limit: limit, waiters: make(chan struct{}, limit)}
}

// Acquire acquires slot
func (bp *Backpressure) Acquire() bool {
	bp.mu.Lock()
	defer bp.mu.Unlock()
	if bp.current >= bp.limit {
		return false
	}
	bp.current++
	return true
}

// Release releases slot
func (bp *Backpressure) Release() {
	bp.mu.Lock()
	defer bp.mu.Unlock()
	if bp.current > 0 {
		bp.current--
	}
}

// Current returns current
func (bp *Backpressure) Current() int {
	bp.mu.Lock()
	defer bp.mu.Unlock()
	return bp.current
}

// Limit returns limit
func (bp *Backpressure) Limit() int {
	return bp.limit
}

// Wait waits for slot
func (bp *Backpressure) Wait(timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for {
		if bp.Acquire() {
			return true
		}
		if time.Now().After(deadline) {
			return false
		}
		time.Sleep(time.Millisecond)
	}
}
