// Package catalog provides index management
package catalog

import (
	"sync"
)

// CatalogIndex represents an index
type CatalogIndex struct {
	Name   string
	Fields []string
}

// CatalogIndexManager manages indexes
type CatalogIndexManager struct {
	indexes map[string]*CatalogIndex
	mu      sync.RWMutex
}

// NewCatalogIndexManager creates manager
func NewCatalogIndexManager() *CatalogIndexManager {
	return &CatalogIndexManager{indexes: make(map[string]*CatalogIndex)}
}

// Create creates index
func (cim *CatalogIndexManager) Create(name string, fields []string) *CatalogIndex {
	cim.mu.Lock()
	defer cim.mu.Unlock()
	idx := &CatalogIndex{Name: name, Fields: fields}
	cim.indexes[name] = idx
	return idx
}

// Get returns index
func (cim *CatalogIndexManager) Get(name string) (*CatalogIndex, bool) {
	cim.mu.RLock()
	defer cim.mu.RUnlock()
	idx, ok := cim.indexes[name]
	return idx, ok
}

// Delete deletes index
func (cim *CatalogIndexManager) Delete(name string) bool {
	cim.mu.Lock()
	defer cim.mu.Unlock()
	if _, ok := cim.indexes[name]; !ok {
		return false
	}
	delete(cim.indexes, name)
	return true
}

// List lists indexes
func (cim *CatalogIndexManager) List() []*CatalogIndex {
	cim.mu.RLock()
	defer cim.mu.RUnlock()
	out := []*CatalogIndex{}
	for _, idx := range cim.indexes {
		out = append(out, idx)
	}
	return out
}

// Count returns count
func (cim *CatalogIndexManager) Count() int {
	cim.mu.RLock()
	defer cim.mu.RUnlock()
	return len(cim.indexes)
}
