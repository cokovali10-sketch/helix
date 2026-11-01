// Package api provides response cache
package api

import (
	"sync"
	"time"
)

// ResponseCache caches responses
type ResponseCache struct {
	items map[string]*cacheItem
	mu    sync.RWMutex
	ttl   time.Duration
}

type cacheItem struct {
	value     interface{}
	expiresAt time.Time
}

// NewResponseCache creates cache
func NewResponseCache(ttl time.Duration) *ResponseCache {
	return &ResponseCache{items: make(map[string]*cacheItem), ttl: ttl}
}

// Get gets cached value
func (rc *ResponseCache) Get(key string) (interface{}, bool) {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	item, ok := rc.items[key]
	if !ok {
		return nil, false
	}
	if time.Now().After(item.expiresAt) {
		return nil, false
	}
	return item.value, true
}

// Set sets cached value
func (rc *ResponseCache) Set(key string, value interface{}) {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	rc.items[key] = &cacheItem{value: value, expiresAt: time.Now().Add(rc.ttl)}
}

// Delete deletes cached value
func (rc *ResponseCache) Delete(key string) {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	delete(rc.items, key)
}

// Clear clears cache
func (rc *ResponseCache) Clear() {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	rc.items = make(map[string]*cacheItem)
}

// Size returns size
func (rc *ResponseCache) Size() int {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	return len(rc.items)
}

// Cleanup removes expired
func (rc *ResponseCache) Cleanup() {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	now := time.Now()
	for k, v := range rc.items {
		if now.After(v.expiresAt) {
			delete(rc.items, k)
		}
	}
}
