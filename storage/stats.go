// Package storage provides stats collection
package storage

import (
	"sync"
	"time"
)

// StorageStats holds storage stats
type StorageStats struct {
	Reads       int64
	Writes      int64
	Deletes     int64
	Errors      int64
	BytesRead   int64
	BytesWritten int64
	LastRead    time.Time
	LastWrite   time.Time
	mu          sync.Mutex
}

// NewStorageStats creates stats
func NewStorageStats() *StorageStats {
	return &StorageStats{}
}

// RecordRead records read
func (s *StorageStats) RecordRead(bytes int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Reads++
	s.BytesRead += bytes
	s.LastRead = time.Now()
}

// RecordWrite records write
func (s *StorageStats) RecordWrite(bytes int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Writes++
	s.BytesWritten += bytes
	s.LastWrite = time.Now()
}

// RecordDelete records delete
func (s *StorageStats) RecordDelete() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Deletes++
}

// RecordError records error
func (s *StorageStats) RecordError() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Errors++
}

// Snapshot returns snapshot
func (s *StorageStats) Snapshot() map[string]interface{} {
	s.mu.Lock()
	defer s.mu.Unlock()
	return map[string]interface{}{
		"reads":         s.Reads,
		"writes":        s.Writes,
		"deletes":       s.Deletes,
		"errors":        s.Errors,
		"bytes_read":    s.BytesRead,
		"bytes_written": s.BytesWritten,
		"last_read":     s.LastRead,
		"last_write":    s.LastWrite,
	}
}

// Reset resets stats
func (s *StorageStats) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Reads = 0
	s.Writes = 0
	s.Deletes = 0
	s.Errors = 0
	s.BytesRead = 0
	s.BytesWritten = 0
}
