// Package api provides throttling
package api

import (
	"sync"
	"time"
)

// Throttler throttles requests
type Throttler struct {
	interval time.Duration
	last     time.Time
	mu       sync.Mutex
}

// NewThrottler creates throttler
func NewThrottler(interval time.Duration) *Throttler {
	return &Throttler{interval: interval}
}

// Allow checks if allowed
func (t *Throttler) Allow() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	now := time.Now()
	if now.Sub(t.last) >= t.interval {
		t.last = now
		return true
	}
	return false
}

// Wait blocks until allowed
func (t *Throttler) Wait() {
	for !t.Allow() {
		time.Sleep(time.Millisecond)
	}
}

// Reset resets
func (t *Throttler) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.last = time.Time{}
}
