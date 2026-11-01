// Package api provides rate limiting
package api

import (
	"net/http"
	"sync"
	"time"
)

// RateLimiter limits requests per client
type RateLimiter struct {
	limits map[string]*bucket
	mu     sync.Mutex
	rate   int
	window time.Duration
}

type bucket struct {
	count     int
	resetAt   time.Time
}

// NewRateLimiter creates limiter
func NewRateLimiter(rate int, window time.Duration) *RateLimiter {
	return &RateLimiter{limits: make(map[string]*bucket), rate: rate, window: window}
}

// Allow checks if allowed
func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	b, ok := rl.limits[key]
	now := time.Now()
	if !ok || now.After(b.resetAt) {
		rl.limits[key] = &bucket{count: 1, resetAt: now.Add(rl.window)}
		return true
	}
	if b.count >= rl.rate {
		return false
	}
	b.count++
	return true
}

// Middleware wraps handler
func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.RemoteAddr
		if !rl.Allow(key) {
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Reset clears buckets
func (rl *RateLimiter) Reset() {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.limits = make(map[string]*bucket)
}

// Stats returns stats
func (rl *RateLimiter) Stats() map[string]int {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	out := map[string]int{}
	for k, b := range rl.limits {
		out[k] = b.count
	}
	return out
}
