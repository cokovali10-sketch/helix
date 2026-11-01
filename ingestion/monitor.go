// Package ingestion provides monitoring
package ingestion

import (
	"sync"
	"time"
)

// Monitor monitors ingestion
type Monitor struct {
	recordsIn   int64
	recordsOut  int64
	errors      int64
	startedAt   time.Time
	mu          sync.Mutex
}

// NewMonitor creates monitor
func NewMonitor() *Monitor {
	return &Monitor{startedAt: time.Now()}
}

// RecordIn records input
func (m *Monitor) RecordIn(n int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.recordsIn += n
}

// RecordOut records output
func (m *Monitor) RecordOut(n int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.recordsOut += n
}

// RecordError records error
func (m *Monitor) RecordError(n int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.errors += n
}

// Stats returns stats
func (m *Monitor) Stats() map[string]interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()
	return map[string]interface{}{
		"records_in":   m.recordsIn,
		"records_out":  m.recordsOut,
		"errors":       m.errors,
		"uptime":       time.Since(m.startedAt).String(),
		"started_at":   m.startedAt,
	}
}

// Reset resets
func (m *Monitor) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.recordsIn = 0
	m.recordsOut = 0
	m.errors = 0
	m.startedAt = time.Now()
}
