// Package catalog provides schema management
package catalog

import (
	"fmt"
	"sync"
)

// Schema represents a schema
type Schema struct {
	Name   string   `json:"name"`
	Fields []string `json:"fields"`
	Type   string   `json:"type"`
}

// SchemaManager manages schemas
type SchemaManager struct {
	schemas map[string]*Schema
	mu      sync.RWMutex
}

// NewSchemaManager creates manager
func NewSchemaManager() *SchemaManager {
	return &SchemaManager{schemas: make(map[string]*Schema)}
}

// Create creates schema
func (sm *SchemaManager) Create(name string, fields []string, typ string) *Schema {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	s := &Schema{Name: name, Fields: fields, Type: typ}
	sm.schemas[name] = s
	return s
}

// Get returns schema
func (sm *SchemaManager) Get(name string) (*Schema, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	s, ok := sm.schemas[name]
	return s, ok
}

// Delete deletes schema
func (sm *SchemaManager) Delete(name string) bool {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if _, ok := sm.schemas[name]; !ok {
		return false
	}
	delete(sm.schemas, name)
	return true
}

// List lists schemas
func (sm *SchemaManager) List() []*Schema {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	out := []*Schema{}
	for _, s := range sm.schemas {
		out = append(out, s)
	}
	return out
}

// Count returns count
func (sm *SchemaManager) Count() int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return len(sm.schemas)
}

// Validate validates schema
func (sm *SchemaManager) Validate(name string) error {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	s, ok := sm.schemas[name]
	if !ok {
		return fmt.Errorf("schema not found")
	}
	if len(s.Fields) == 0 {
		return fmt.Errorf("schema has no fields")
	}
	return nil
}
