// Package storage provides cache management
package storage

import (
	"sync"
	"time"
)

// CacheEntry represents cache entry
type CacheEntry struct {
	Key       string
	Value     interface{}
	ExpiresAt time.Time
}

// CacheManager manages cache
type CacheManager struct {
	entries map[string]*CacheEntry
	mu      sync.RWMutex
	ttl     time.Duration
}

// NewCacheManager creates manager
func NewCacheManager(ttl time.Duration) *CacheManager {
	return &CacheManager{entries: make(map[string]*CacheEntry), ttl: ttl}
}

// Set sets value
func (cm *CacheManager) Set(key string, value interface{}) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.entries[key] = &CacheEntry{Key: key, Value: value, ExpiresAt: time.Now().Add(cm.ttl)}
}

// Get gets value
func (cm *CacheManager) Get(key string) (interface{}, bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	e, ok := cm.entries[key]
	if !ok {
		return nil, false
	}
	if time.Now().After(e.ExpiresAt) {
		return nil, false
	}
	return e.Value, true
}

// Delete deletes value
func (cm *CacheManager) Delete(key string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	delete(cm.entries, key)
}

// Clear clears cache
func (cm *CacheManager) Clear() {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.entries = make(map[string]*CacheEntry)
}

// Size returns size
func (cm *CacheManager) Size() int {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return len(cm.entries)
}

// Cleanup removes expired
func (cm *CacheManager) Cleanup() {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	now := time.Now()
	for k, e := range cm.entries {
		if now.After(e.ExpiresAt) {
			delete(cm.entries, k)
		}
	}
}
