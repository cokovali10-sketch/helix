// Package storage provides additional storage features
// including advanced query support and optimization
package storage

import (
	"fmt"
	"sync"
	"time"
)

// QueryEngine provides query capabilities for storage
type QueryEngine struct {
	engine *StorageEngine
	cache  *QueryCache
	mu     sync.RWMutex
}

// QueryCache caches query results
type QueryCache struct {
	cache    map[string]*CachedResult
	mu       sync.RWMutex
	maxSize  int
	ttl      time.Duration
}

// CachedResult represents a cached query result
type CachedResult struct {
	Result    interface{}
	Timestamp time.Time
	HitCount  int64
}

// NewQueryEngine creates a new query engine
func NewQueryEngine(engine *StorageEngine) *QueryEngine {
	return &QueryEngine{
		engine: engine,
		cache: &QueryCache{
			cache:   make(map[string]*CachedResult),
			maxSize: 1000,
			ttl:     5 * time.Minute,
		},
	}
}

// Query executes a query
func (qe *QueryEngine) Query(expr string) (interface{}, error) {
	// Check cache
	if cached := qe.cache.Get(expr); cached != nil {
		return cached.Result, nil
	}

	// Execute query
	result, err := qe.execute(expr)
	if err != nil {
		return nil, err
	}

	// Cache result
	qe.cache.Put(expr, result)

	return result, nil
}

// execute executes a query
func (qe *QueryEngine) execute(expr string) (interface{}, error) {
	// In production, would parse and execute query
	return nil, nil
}

// Get queries a specific key
func (qe *QueryEngine) Get(key string) ([]byte, bool) {
	return qe.engine.Get(key)
}

// RangeQuery queries a range of keys
func (qe *QueryEngine) RangeQuery(start, end string) (map[string][]byte, error) {
	return qe.engine.RangeScan(start, end)
}

// PrefixQuery queries keys with a prefix
func (qe *QueryEngine) PrefixQuery(prefix string) (map[string][]byte, error) {
	return qe.engine.PrefixScan(prefix)
}

// Count returns the number of entries
func (qe *QueryEngine) Count() (int, error) {
	return qe.engine.Count()
}

// QueryCache methods

// Get gets a cached result
func (qc *QueryCache) Get(key string) *CachedResult {
	qc.mu.RLock()
	defer qc.mu.RUnlock()

	cached, ok := qc.cache[key]
	if !ok {
		return nil
	}

	if time.Since(cached.Timestamp) > qc.ttl {
		return nil
	}

	cached.HitCount++
	return cached
}

// Put puts a result in cache
func (qc *QueryCache) Put(key string, result interface{}) {
	qc.mu.Lock()
	defer qc.mu.Unlock()

	if len(qc.cache) >= qc.maxSize {
		// Evict oldest
		var oldestKey string
		var oldestTime time.Time
		for k, v := range qc.cache {
			if oldestKey == "" || v.Timestamp.Before(oldestTime) {
				oldestKey = k
				oldestTime = v.Timestamp
			}
		}
		delete(qc.cache, oldestKey)
	}

	qc.cache[key] = &CachedResult{
		Result:    result,
		Timestamp: time.Now(),
		HitCount:  0,
	}
}

// Clear clears the cache
func (qc *QueryCache) Clear() {
	qc.mu.Lock()
	defer qc.mu.Unlock()
	qc.cache = make(map[string]*CachedResult)
}

// Size returns cache size
func (qc *QueryCache) Size() int {
	qc.mu.RLock()
	defer qc.mu.RUnlock()
	return len(qc.cache)
}

// Stats returns cache statistics
func (qc *QueryCache) Stats() map[string]interface{} {
	qc.mu.RLock()
	defer qc.mu.RUnlock()

	var totalHits int64
	for _, v := range qc.cache {
		totalHits += v.HitCount
	}

	return map[string]interface{}{
		"size":       len(qc.cache),
		"max_size":   qc.maxSize,
		"ttl":        qc.ttl.String(),
		"total_hits": totalHits,
	}
}

// StorageAggregator aggregates storage data
type StorageAggregator struct {
	engine *StorageEngine
}

// NewStorageAggregator creates a new aggregator
func NewStorageAggregator(engine *StorageEngine) *StorageAggregator {
	return &StorageAggregator{engine: engine}
}

// Sum aggregates values by key prefix
func (sa *StorageAggregator) Sum(prefix string) (float64, error) {
	result, err := sa.engine.PrefixScan(prefix)
	if err != nil {
		return 0, err
	}

	var sum float64
	for _, v := range result {
		// Parse value as float
		var val float64
		fmt.Sscanf(string(v), "%f", &val)
		sum += val
	}
	return sum, nil
}

// Count counts entries by prefix
func (sa *StorageAggregator) Count(prefix string) (int, error) {
	result, err := sa.engine.PrefixScan(prefix)
	if err != nil {
		return 0, err
	}
	return len(result), nil
}

// Min finds minimum value by prefix
func (sa *StorageAggregator) Min(prefix string) (float64, error) {
	result, err := sa.engine.PrefixScan(prefix)
	if err != nil {
		return 0, err
	}

	min := 0.0
	first := true
	for _, v := range result {
		var val float64
		fmt.Sscanf(string(v), "%f", &val)
		if first || val < min {
			min = val
			first = false
		}
	}
	return min, nil
}

// Max finds maximum value by prefix
func (sa *StorageAggregator) Max(prefix string) (float64, error) {
	result, err := sa.engine.PrefixScan(prefix)
	if err != nil {
		return 0, err
	}

	max := 0.0
	first := true
	for _, v := range result {
		var val float64
		fmt.Sscanf(string(v), "%f", &val)
		if first || val > max {
			max = val
			first = false
		}
	}
	return max, nil
}

// Avg calculates average value by prefix
func (sa *StorageAggregator) Avg(prefix string) (float64, error) {
	result, err := sa.engine.PrefixScan(prefix)
	if err != nil {
		return 0, err
	}

	var sum float64
	count := 0
	for _, v := range result {
		var val float64
		fmt.Sscanf(string(v), "%f", &val)
		sum += val
		count++
	}
	if count == 0 {
		return 0, nil
	}
	return sum / float64(count), nil
}

// StorageIndexer provides indexing for storage
type StorageIndexer struct {
	engine *StorageEngine
	index  map[string]map[string][]string // field -> value -> keys
	mu     sync.RWMutex
}

// NewStorageIndexer creates a new storage indexer
func NewStorageIndexer(engine *StorageEngine) *StorageIndexer {
	return &StorageIndexer{
		engine: engine,
		index:  make(map[string]map[string][]string),
	}
}

// Index indexes a key by field
func (si *StorageIndexer) Index(key, field, value string) {
	si.mu.Lock()
	defer si.mu.Unlock()

	if si.index[field] == nil {
		si.index[field] = make(map[string][]string)
	}
	si.index[field][value] = append(si.index[field][value], key)
}

// Lookup looks up keys by field and value
func (si *StorageIndexer) Lookup(field, value string) []string {
	si.mu.RLock()
	defer si.mu.RUnlock()

	if si.index[field] == nil {
		return nil
	}
	return si.index[field][value]
}

// ListFields lists all indexed fields
func (si *StorageIndexer) ListFields() []string {
	si.mu.RLock()
	defer si.mu.RUnlock()

	fields := make([]string, 0, len(si.index))
	for field := range si.index {
		fields = append(fields, field)
	}
	return fields
}

// ListValues lists all values for a field
func (si *StorageIndexer) ListValues(field string) []string {
	si.mu.RLock()
	defer si.mu.RUnlock()

	if si.index[field] == nil {
		return nil
	}
	values := make([]string, 0, len(si.index[field]))
	for value := range si.index[field] {
		values = append(values, value)
	}
	return values
}

// Remove removes a key from index
func (si *StorageIndexer) Remove(key string) {
	si.mu.Lock()
	defer si.mu.Unlock()

	for field := range si.index {
		for value := range si.index[field] {
			keys := si.index[field][value]
			for i, k := range keys {
				if k == key {
					si.index[field][value] = append(keys[:i], keys[i+1:]...)
					break
				}
			}
		}
	}
}

// Clear clears the index
func (si *StorageIndexer) Clear() {
	si.mu.Lock()
	defer si.mu.Unlock()
	si.index = make(map[string]map[string][]string)
}

// StorageSearcher searches storage
type StorageSearcher struct {
	engine  *StorageEngine
	indexer *StorageIndexer
}

// NewStorageSearcher creates a new storage searcher
func NewStorageSearcher(engine *StorageEngine) *StorageSearcher {
	return &StorageSearcher{
		engine:  engine,
		indexer: NewStorageIndexer(engine),
	}
}

// Search searches for keys matching a query
func (ss *StorageSearcher) Search(query string) ([]string, error) {
	// In production, would use full-text search
	result, err := ss.engine.PrefixScan(query)
	if err != nil {
		return nil, err
	}

	keys := make([]string, 0, len(result))
	for key := range result {
		keys = append(keys, key)
	}
	return keys, nil
}

// SearchByField searches by field value
func (ss *StorageSearcher) SearchByField(field, value string) []string {
	return ss.indexer.Lookup(field, value)
}

// IndexKey indexes a key
func (ss *StorageSearcher) IndexKey(key string, fields map[string]string) {
	for field, value := range fields {
		ss.indexer.Index(key, field, value)
	}
}