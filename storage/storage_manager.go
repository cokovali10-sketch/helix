// Package storage provides storage manager
package storage

import (
	"sync"
)

// StorageManager manages storage components
type StorageManager struct {
	engine *StorageEngine
	cache  *CacheManager
	index  *IndexManager
	mu     sync.RWMutex
}

// NewStorageManager creates manager
func NewStorageManager(engine *StorageEngine, cache *CacheManager, index *IndexManager) *StorageManager {
	return &StorageManager{engine: engine, cache: cache, index: index}
}

// Engine returns engine
func (sm *StorageManager) Engine() *StorageEngine {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.engine
}

// Cache returns cache
func (sm *StorageManager) Cache() *CacheManager {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.cache
}

// Index returns index
func (sm *StorageManager) Index() *IndexManager {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.index
}
