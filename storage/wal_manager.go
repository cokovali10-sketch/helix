// Package storage provides WAL management
package storage

import (
	"fmt"
	"sync"
	"time"
)

// WALEntry represents a WAL entry
type WALEntry struct {
	LSN       uint64    `json:"lsn"`
	Type      string    `json:"type"`
	Data      []byte    `json:"data"`
	Timestamp time.Time `json:"timestamp"`
}

// WALManager manages WAL
type WALManager struct {
	entries []WALEntry
	mu      sync.Mutex
	lsn     uint64
}

// NewWALManager creates manager
func NewWALManager() *WALManager {
	return &WALManager{entries: []WALEntry{}}
}

// Append appends entry
func (w *WALManager) Append(typ string, data []byte) uint64 {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.lsn++
	entry := WALEntry{
		LSN:       w.lsn,
		Type:      typ,
		Data:      data,
		Timestamp: time.Now(),
	}
	w.entries = append(w.entries, entry)
	return w.lsn
}

// Entries returns entries
func (w *WALManager) Entries() []WALEntry {
	w.mu.Lock()
	defer w.mu.Unlock()
	out := make([]WALEntry, len(w.entries))
	copy(out, w.entries)
	return out
}

// LastLSN returns last LSN
func (w *WALManager) LastLSN() uint64 {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.lsn
}

// Truncate truncates up to LSN
func (w *WALManager) Truncate(upTo uint64) {
	w.mu.Lock()
	defer w.mu.Unlock()
	newEntries := []WALEntry{}
	for _, e := range w.entries {
		if e.LSN > upTo {
			newEntries = append(newEntries, e)
		}
	}
	w.entries = newEntries
}

// Count returns count
func (w *WALManager) Count() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	return len(w.entries)
}

// Get returns entry by LSN
func (w *WALManager) Get(lsn uint64) (*WALEntry, bool) {
	w.mu.Lock()
	defer w.mu.Unlock()
	for i := range w.entries {
		if w.entries[i].LSN == lsn {
			return &w.entries[i], true
		}
	}
	return nil, false
}

// String returns summary
func (w *WALManager) String() string {
	return fmt.Sprintf("WALManager: %d entries, lsn=%d", len(w.entries), w.lsn)
}
