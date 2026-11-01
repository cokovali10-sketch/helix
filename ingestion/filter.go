// Package ingestion provides filtering
package ingestion

import (
	"fmt"
	"sync"
)

// FilterFunc filters records
type FilterFunc func(map[string]interface{}) bool

// Filter filters records
type Filter struct {
	predicates []FilterFunc
	mu         sync.RWMutex
}

// NewFilter creates filter
func NewFilter() *Filter {
	return &Filter{predicates: []FilterFunc{}}
}

// Add adds predicate
func (f *Filter) Add(fn FilterFunc) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.predicates = append(f.predicates, fn)
}

// Apply applies filter
func (f *Filter) Apply(record map[string]interface{}) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	for _, p := range f.predicates {
		if !p(record) {
			return false
		}
	}
	return true
}

// FieldEquals filter
func FieldEquals(field string, value interface{}) FilterFunc {
	return func(r map[string]interface{}) bool {
		v, ok := r[field]
		if !ok {
			return false
		}
		return v == value
	}
}

// FieldContains filter
func FieldContains(field, substr string) FilterFunc {
	return func(r map[string]interface{}) bool {
		v, ok := r[field]
		if !ok {
			return false
		}
		s, ok := v.(string)
		if !ok {
			return false
		}
		return contains(s, substr)
	}
}

// FieldGreaterThan filter
func FieldGreaterThan(field string, threshold float64) FilterFunc {
	return func(r map[string]interface{}) bool {
		v, ok := r[field]
		if !ok {
			return false
		}
		switch val := v.(type) {
		case float64:
			return val > threshold
		case int:
			return float64(val) > threshold
		}
		return false
	}
}

// FieldExists filter
func FieldExists(field string) FilterFunc {
	return func(r map[string]interface{}) bool {
		_, ok := r[field]
		return ok
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || indexOf(s, substr) >= 0)
}

func indexOf(s, substr string) int {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// Count returns predicate count
func (f *Filter) Count() int {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return len(f.predicates)
}

// Clear clears predicates
func (f *Filter) Clear() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.predicates = []FilterFunc{}
}

// String returns summary
func (f *Filter) String() string {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return fmt.Sprintf("Filter: %d predicates", len(f.predicates))
}
