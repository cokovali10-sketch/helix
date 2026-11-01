// Package storage provides lock management
package storage

import (
	"sync"
	"time"
)

// Lock represents a lock
type Lock struct {
	Key       string
	Owner     string
	Acquired  time.Time
	ExpiresAt time.Time
}

// LockManager manages locks
type LockManager struct {
	locks map[string]*Lock
	mu    sync.Mutex
}

// NewLockManager creates manager
func NewLockManager() *LockManager {
	return &LockManager{locks: make(map[string]*Lock)}
}

// Acquire acquires lock
func (lm *LockManager) Acquire(key, owner string, ttl time.Duration) bool {
	lm.mu.Lock()
	defer lm.mu.Unlock()
	if l, ok := lm.locks[key]; ok {
		if time.Now().Before(l.ExpiresAt) {
			return false
		}
	}
	lm.locks[key] = &Lock{Key: key, Owner: owner, Acquired: time.Now(), ExpiresAt: time.Now().Add(ttl)}
	return true
}

// Release releases lock
func (lm *LockManager) Release(key, owner string) bool {
	lm.mu.Lock()
	defer lm.mu.Unlock()
	l, ok := lm.locks[key]
	if !ok || l.Owner != owner {
		return false
	}
	delete(lm.locks, key)
	return true
}

// IsLocked checks locked
func (lm *LockManager) IsLocked(key string) bool {
	lm.mu.Lock()
	defer lm.mu.Unlock()
	l, ok := lm.locks[key]
	if !ok {
		return false
	}
	if time.Now().After(l.ExpiresAt) {
		delete(lm.locks, key)
		return false
	}
	return true
}

// Owner returns owner
func (lm *LockManager) Owner(key string) string {
	lm.mu.Lock()
	defer lm.mu.Unlock()
	if l, ok := lm.locks[key]; ok {
		return l.Owner
	}
	return ""
}

// Count returns count
func (lm *LockManager) Count() int {
	lm.mu.Lock()
	defer lm.mu.Unlock()
	return len(lm.locks)
}
