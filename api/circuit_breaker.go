// Package api provides circuit breaker
package api

import (
	"errors"
	"sync"
	"time"
)

// CircuitBreaker implements circuit breaker
type CircuitBreaker struct {
	failures     int
	threshold    int
	resetTimeout time.Duration
	lastFailure  time.Time
	state        string // closed, open, half-open
	mu           sync.Mutex
}

// NewCircuitBreaker creates breaker
func NewCircuitBreaker(threshold int, resetTimeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{threshold: threshold, resetTimeout: resetTimeout, state: "closed"}
}

// Allow checks if allowed
func (cb *CircuitBreaker) Allow() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	if cb.state == "open" {
		if time.Since(cb.lastFailure) > cb.resetTimeout {
			cb.state = "half-open"
			return true
		}
		return false
	}
	return true
}

// RecordSuccess records success
func (cb *CircuitBreaker) RecordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.failures = 0
	cb.state = "closed"
}

// RecordFailure records failure
func (cb *CircuitBreaker) RecordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.failures++
	cb.lastFailure = time.Now()
	if cb.failures >= cb.threshold {
		cb.state = "open"
	}
}

// State returns state
func (cb *CircuitBreaker) State() string {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.state
}

// Execute executes function
func (cb *CircuitBreaker) Execute(fn func() error) error {
	if !cb.Allow() {
		return errors.New("circuit breaker open")
	}
	err := fn()
	if err != nil {
		cb.RecordFailure()
		return err
	}
	cb.RecordSuccess()
	return nil
}
