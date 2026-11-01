// Package storage provides additional storage features
// including advanced caching and prefetching
package storage

import (
	"sync"
	"time"
)

// AdvancedCache provides advanced caching capabilities
type AdvancedCache struct {
	cache      map[string]*CacheEntry
	lru        []string
	capacity   int
	currentSize int
	mu         sync.RWMutex
	stats      CacheStats
	ttl        time.Duration
}

// CacheEntry represents a cache entry
type CacheEntry struct {
	Key       string
	Value     interface{}
	CreatedAt time.Time
	AccessedAt time.Time
	AccessCount int64
	Size      int
	TTL       time.Duration
}

// CacheStats tracks cache statistics
type CacheStats struct {
	Hits       int64
	Misses     int64
	Evictions  int64
	Insertions int64
	Expirations int64
}

// NewAdvancedCache creates a new advanced cache
func NewAdvancedCache(capacity int, ttl time.Duration) *AdvancedCache {
	return &AdvancedCache{
		cache:    make(map[string]*CacheEntry),
		lru:      make([]string, 0),
		capacity: capacity,
		ttl:      ttl,
	}
}

// Get gets a value from cache
func (ac *AdvancedCache) Get(key string) (interface{}, bool) {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	entry, ok := ac.cache[key]
	if !ok {
		ac.stats.Misses++
		return nil, false
	}

	if time.Since(entry.CreatedAt) > ac.ttl {
		delete(ac.cache, key)
		ac.removeFromLRU(key)
		ac.stats.Expirations++
		return nil, false
	}

	entry.AccessedAt = time.Now()
	entry.AccessCount++
	ac.stats.Hits++

	// Move to front of LRU
	ac.moveToFront(key)

	return entry.Value, true
}

// Put puts a value in cache
func (ac *AdvancedCache) Put(key string, value interface{}, size int) {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	// Check if key exists
	if existing, ok := ac.cache[key]; ok {
		ac.currentSize -= existing.Size
	}

	// Evict if necessary
	for ac.currentSize+size > ac.capacity && len(ac.cache) > 0 {
		ac.evictLRU()
	}

	entry := &CacheEntry{
		Key:        key,
		Value:      value,
		CreatedAt:  time.Now(),
		AccessedAt: time.Now(),
		AccessCount: 1,
		Size:       size,
		TTL:        ac.ttl,
	}

	ac.cache[key] = entry
	ac.currentSize += size
	ac.stats.Insertions++
	ac.lru = append([]string{key}, ac.lru...)
}

// Delete deletes a value from cache
func (ac *AdvancedCache) Delete(key string) bool {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	entry, ok := ac.cache[key]
	if !ok {
		return false
	}

	delete(ac.cache, key)
	ac.removeFromLRU(key)
	ac.currentSize -= entry.Size
	return true
}

// Clear clears the cache
func (ac *AdvancedCache) Clear() {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	ac.cache = make(map[string]*CacheEntry)
	ac.lru = make([]string, 0)
	ac.currentSize = 0
}

// Size returns the cache size
func (ac *AdvancedCache) Size() int {
	ac.mu.RLock()
	defer ac.mu.RUnlock()
	return len(ac.cache)
}

// Capacity returns the cache capacity
func (ac *AdvancedCache) Capacity() int {
	return ac.capacity
}

// Utilization returns the cache utilization
func (ac *AdvancedCache) Utilization() float64 {
	ac.mu.RLock()
	defer ac.mu.RUnlock()
	if ac.capacity == 0 {
		return 0
	}
	return float64(ac.currentSize) / float64(ac.capacity)
}

// Stats returns cache statistics
func (ac *AdvancedCache) Stats() CacheStats {
	ac.mu.RLock()
	defer ac.mu.RUnlock()
	return ac.stats
}

// HitRate returns the cache hit rate
func (ac *AdvancedCache) HitRate() float64 {
	ac.mu.RLock()
	defer ac.mu.RUnlock()
	total := ac.stats.Hits + ac.stats.Misses
	if total == 0 {
		return 0
	}
	return float64(ac.stats.Hits) / float64(total)
}

// Contains checks if a key exists
func (ac *AdvancedCache) Contains(key string) bool {
	ac.mu.RLock()
	defer ac.mu.RUnlock()
	entry, ok := ac.cache[key]
	if !ok {
		return false
	}
	return time.Since(entry.CreatedAt) <= ac.ttl
}

// Keys returns all keys
func (ac *AdvancedCache) Keys() []string {
	ac.mu.RLock()
	defer ac.mu.RUnlock()
	keys := make([]string, 0, len(ac.cache))
	for key := range ac.cache {
		keys = append(keys, key)
	}
	return keys
}

// CleanupExpired removes expired entries
func (ac *AdvancedCache) CleanupExpired() int {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	removed := 0
	for key, entry := range ac.cache {
		if time.Since(entry.CreatedAt) > ac.ttl {
			delete(ac.cache, key)
			ac.removeFromLRU(key)
			ac.currentSize -= entry.Size
			ac.stats.Expirations++
			removed++
		}
	}
	return removed
}

func (ac *AdvancedCache) moveToFront(key string) {
	for i, k := range ac.lru {
		if k == key {
			ac.lru = append(ac.lru[:i], ac.lru[i+1:]...)
			ac.lru = append([]string{key}, ac.lru...)
			break
		}
	}
}

func (ac *AdvancedCache) removeFromLRU(key string) {
	for i, k := range ac.lru {
		if k == key {
			ac.lru = append(ac.lru[:i], ac.lru[i+1:]...)
			break
		}
	}
}

func (ac *AdvancedCache) evictLRU() {
	if len(ac.lru) == 0 {
		return
	}

	key := ac.lru[len(ac.lru)-1]
	ac.lru = ac.lru[:len(ac.lru)-1]

	if entry, ok := ac.cache[key]; ok {
		delete(ac.cache, key)
		ac.currentSize -= entry.Size
		ac.stats.Evictions++
	}
}

// Prefetcher prefetches data
type Prefetcher struct {
	cache    *AdvancedCache
	patterns []string
	mu       sync.RWMutex
	running  bool
	ctx      context.Context
	cancel   context.CancelFunc
}

// NewPrefetcher creates a new prefetcher
func NewPrefetcher(cache *AdvancedCache) *Prefetcher {
	return &Prefetcher{
		cache:    cache,
		patterns: make([]string, 0),
	}
}

// Start starts prefetching
func (p *Prefetcher) Start(ctx context.Context) {
	p.ctx, p.cancel = context.WithCancel(ctx)
	p.running = true
	go p.run()
}

// Stop stops prefetching
func (p *Prefetcher) Stop() {
	p.running = false
	if p.cancel != nil {
		p.cancel()
	}
}

// run runs the prefetch loop
func (p *Prefetcher) run() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-p.ctx.Done():
			return
		case <-ticker.C:
			if !p.running {
				return
			}
			p.prefetch()
		}
	}
}

// prefetch prefetches data
func (p *Prefetcher) prefetch() {
	// In production, would prefetch based on patterns
}

// AddPattern adds a prefetch pattern
func (p *Prefetcher) AddPattern(pattern string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.patterns = append(p.patterns, pattern)
}

// GetPatterns returns prefetch patterns
func (p *Prefetcher) GetPatterns() []string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	patterns := make([]string, len(p.patterns))
	copy(patterns, p.patterns)
	return patterns
}

// CacheWarmer warms up the cache
type CacheWarmer struct {
	cache  *AdvancedCache
	loader func(string) (interface{}, int, error)
	mu     sync.Mutex
}

// NewCacheWarmer creates a new cache warmer
func NewCacheWarmer(cache *AdvancedCache, loader func(string) (interface{}, int, error)) *CacheWarmer {
	return &CacheWarmer{
		cache:  cache,
		loader: loader,
	}
}

// Warm warms up the cache with given keys
func (cw *CacheWarmer) Warm(keys []string) {
	for _, key := range keys {
		value, size, err := cw.loader(key)
		if err == nil {
			cw.cache.Put(key, value, size)
		}
	}
}

// WarmAll warms up all keys
func (cw *CacheWarmer) WarmAll(keys []string) {
	cw.Warm(keys)
}

// CacheInvalidator invalidates cache entries
type CacheInvalidator struct {
	cache *AdvancedCache
	rules []InvalidationRule
	mu    sync.RWMutex
}

// InvalidationRule represents an invalidation rule
type InvalidationRule struct {
	Pattern  string
	Action   string // delete, refresh
	Priority int
}

// NewCacheInvalidator creates a new cache invalidator
func NewCacheInvalidator(cache *AdvancedCache) *CacheInvalidator {
	return &CacheInvalidator{
		cache: cache,
		rules: make([]InvalidationRule, 0),
	}
}

// AddRule adds an invalidation rule
func (ci *CacheInvalidator) AddRule(rule InvalidationRule) {
	ci.mu.Lock()
	defer ci.mu.Unlock()
	ci.rules = append(ci.rules, rule)
}

// Invalidate invalidates cache entries matching a pattern
func (ci *CacheInvalidator) Invalidate(pattern string) int {
	ci.mu.RLock()
	defer ci.mu.RUnlock()

	count := 0
	for _, key := range ci.cache.Keys() {
		if matchPattern(key, pattern) {
			ci.cache.Delete(key)
			count++
		}
	}
	return count
}

// InvalidateAll invalidates all cache entries
func (ci *CacheInvalidator) InvalidateAll() {
	ci.cache.Clear()
}

func matchPattern(s, pattern string) bool {
	// Simple pattern matching
	if pattern == "*" {
		return true
	}
	if len(pattern) > 0 && pattern[len(pattern)-1] == '*' {
		prefix := pattern[:len(pattern)-1]
		return len(s) >= len(prefix) && s[:len(prefix)] == prefix
	}
	return s == pattern
}

// CacheMonitor monitors cache performance
type CacheMonitor struct {
	cache    *AdvancedCache
	metrics  map[string]float64
	mu       sync.RWMutex
	interval time.Duration
	ctx      context.Context
	cancel   context.CancelFunc
	running  bool
}

// NewCacheMonitor creates a new cache monitor
func NewCacheMonitor(cache *AdvancedCache, interval time.Duration) *CacheMonitor {
	return &CacheMonitor{
		cache:    cache,
		metrics:  make(map[string]float64),
		interval: interval,
	}
}

// Start starts monitoring
func (cm *CacheMonitor) Start(ctx context.Context) {
	cm.ctx, cm.cancel = context.WithCancel(ctx)
	cm.running = true
	go cm.run()
}

// Stop stops monitoring
func (cm *CacheMonitor) Stop() {
	cm.running = false
	if cm.cancel != nil {
		cm.cancel()
	}
}

// run runs the monitoring loop
func (cm *CacheMonitor) run() {
	ticker := time.NewTicker(cm.interval)
	defer ticker.Stop()

	for {
		select {
		case <-cm.ctx.Done():
			return
		case <-ticker.C:
			if !cm.running {
				return
			}
			cm.collectMetrics()
		}
	}
}

// collectMetrics collects cache metrics
func (cm *CacheMonitor) collectMetrics() {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	stats := cm.cache.Stats()
	cm.metrics["hits"] = float64(stats.Hits)
	cm.metrics["misses"] = float64(stats.Misses)
	cm.metrics["evictions"] = float64(stats.Evictions)
	cm.metrics["insertions"] = float64(stats.Insertions)
	cm.metrics["expirations"] = float64(stats.Expirations)
	cm.metrics["size"] = float64(cm.cache.Size())
	cm.metrics["capacity"] = float64(cm.cache.Capacity())
	cm.metrics["utilization"] = cm.cache.Utilization()
	cm.metrics["hit_rate"] = cm.cache.HitRate()
}

// GetMetrics returns cache metrics
func (cm *CacheMonitor) GetMetrics() map[string]float64 {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	result := make(map[string]float64)
	for k, v := range cm.metrics {
		result[k] = v
	}
	return result
}

// GetStats returns cache statistics
func (cm *CacheMonitor) GetStats() CacheStats {
	return cm.cache.Stats()
}