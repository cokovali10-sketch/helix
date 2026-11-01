// Package ingestion provides additional monitoring
package ingestion

import (
	"sync"
	"time"
)

// IngestMonitor monitors ingestion
type IngestMonitor struct {
	recordsIn   int64
	recordsOut  int64
	errors      int64
	startedAt   time.Time
	mu          sync.Mutex
}

// NewIngestMonitor creates monitor
func NewIngestMonitor() *IngestMonitor {
	return &IngestMonitor{startedAt: time.Now()}
}

// RecordIn records input
func (m *IngestMonitor) RecordIn(n int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.recordsIn += n
}

// RecordOut records output
func (m *IngestMonitor) RecordOut(n int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.recordsOut += n
}

// RecordError records error
func (m *IngestMonitor) RecordError(n int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.errors += n
}

// Stats returns stats
func (m *IngestMonitor) Stats() map[string]interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()
	return map[string]interface{}{
		"records_in":  m.recordsIn,
		"records_out": m.recordsOut,
		"errors":      m.errors,
		"uptime":      time.Since(m.startedAt).String(),
	}
}

// Reset resets
func (m *IngestMonitor) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.recordsIn = 0
	m.recordsOut = 0
	m.errors = 0
	m.startedAt = time.Now()
}
