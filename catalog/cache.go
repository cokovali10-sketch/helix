// Package catalog provides cache
package catalog

import (
	"sync"
	"time"
)

// CatalogCache caches catalog data
type CatalogCache struct {
	items map[string]interface{}
	mu    sync.RWMutex
	ttl   time.Duration
}

// NewCatalogCache creates cache
func NewCatalogCache(ttl time.Duration) *CatalogCache {
	return &CatalogCache{items: make(map[string]interface{}), ttl: ttl}
}

// Set sets value
func (cc *CatalogCache) Set(key string, value interface{}) {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	cc.items[key] = value
}

// Get gets value
func (cc *CatalogCache) Get(key string) (interface{}, bool) {
	cc.mu.RLock()
	defer cc.mu.RUnlock()
	v, ok := cc.items[key]
	return v, ok
}

// Delete deletes value
func (cc *CatalogCache) Delete(key string) {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	delete(cc.items, key)
}

// Clear clears cache
func (cc *CatalogCache) Clear() {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	cc.items = make(map[string]interface{})
}

// Size returns size
func (cc *CatalogCache) Size() int {
	cc.mu.RLock()
	defer cc.mu.RUnlock()
	return len(cc.items)
}
