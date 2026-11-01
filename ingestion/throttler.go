// Package ingestion provides throttling
package ingestion

import (
	"sync"
	"time"
)

// Throttler throttles ingestion
type Throttler struct {
	rate     int
	interval time.Duration
	tokens   int
	last     time.Time
	mu       sync.Mutex
}

// NewThrottler creates throttler
func NewThrottler(rate int, interval time.Duration) *Throttler {
	return &Throttler{rate: rate, interval: interval, tokens: rate, last: time.Now()}
}

// Allow checks if allowed
func (t *Throttler) Allow() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	now := time.Now()
	elapsed := now.Sub(t.last)
	if elapsed >= t.interval {
		t.tokens = t.rate
		t.last = now
	}
	if t.tokens > 0 {
		t.tokens--
		return true
	}
	return false
}

// Tokens returns tokens
func (t *Throttler) Tokens() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.tokens
}
