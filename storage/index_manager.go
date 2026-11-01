// Package storage provides index management
package storage

import (
	"fmt"
	"sync"
)

// Index represents an index
type Index struct {
	Name   string   `json:"name"`
	Fields []string `json:"fields"`
	Type   string   `json:"type"`
}

// IndexManager manages indexes
type IndexManager struct {
	indexes map[string]*Index
	mu      sync.RWMutex
}

// NewIndexManager creates manager
func NewIndexManager() *IndexManager {
	return &IndexManager{indexes: make(map[string]*Index)}
}

// Create creates index
func (im *IndexManager) Create(name string, fields []string, indexType string) *Index {
	im.mu.Lock()
	defer im.mu.Unlock()
	idx := &Index{Name: name, Fields: fields, Type: indexType}
	im.indexes[name] = idx
	return idx
}

// Get returns index
func (im *IndexManager) Get(name string) (*Index, bool) {
	im.mu.RLock()
	defer im.mu.RUnlock()
	idx, ok := im.indexes[name]
	return idx, ok
}

// Delete deletes index
func (im *IndexManager) Delete(name string) bool {
	im.mu.Lock()
	defer im.mu.Unlock()
	if _, ok := im.indexes[name]; !ok {
		return false
	}
	delete(im.indexes, name)
	return true
}

// List lists indexes
func (im *IndexManager) List() []*Index {
	im.mu.RLock()
	defer im.mu.RUnlock()
	out := []*Index{}
	for _, idx := range im.indexes {
		out = append(out, idx)
	}
	return out
}

// Exists checks existence
func (im *IndexManager) Exists(name string) bool {
	im.mu.RLock()
	defer im.mu.RUnlock()
	_, ok := im.indexes[name]
	return ok
}

// Count returns count
func (im *IndexManager) Count() int {
	im.mu.RLock()
	defer im.mu.RUnlock()
	return len(im.indexes)
}

// Validate validates index
func (im *IndexManager) Validate(name string) error {
	im.mu.RLock()
	defer im.mu.RUnlock()
	idx, ok := im.indexes[name]
	if !ok {
		return fmt.Errorf("index not found")
	}
	if len(idx.Fields) == 0 {
		return fmt.Errorf("index has no fields")
	}
	return nil
}
