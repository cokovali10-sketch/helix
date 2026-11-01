// Package ingestion provides deduplication
package ingestion

import (
	"sync"
)

// Deduplicator deduplicates records
type Deduplicator struct {
	seen map[string]bool
	mu   sync.Mutex
}

// NewDeduplicator creates deduplicator
func NewDeduplicator() *Deduplicator {
	return &Deduplicator{seen: make(map[string]bool)}
}

// IsDuplicate checks duplicate
func (d *Deduplicator) IsDuplicate(key string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.seen[key] {
		return true
	}
	d.seen[key] = true
	return false
}

// Reset resets
func (d *Deduplicator) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.seen = make(map[string]bool)
}

// Count returns count
func (d *Deduplicator) Count() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.seen)
}
