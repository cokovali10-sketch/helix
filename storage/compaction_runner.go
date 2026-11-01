// Package storage provides compaction runner
package storage

import (
	"sync"
	"time"
)

// CompactionRunner runs compaction
type CompactionRunner struct {
	running bool
	mu      sync.Mutex
	stop    chan struct{}
}

// NewCompactionRunner creates runner
func NewCompactionRunner() *CompactionRunner {
	return &CompactionRunner{stop: make(chan struct{})}
}

// Start starts runner
func (cr *CompactionRunner) Start(fn func()) {
	cr.mu.Lock()
	if cr.running {
		cr.mu.Unlock()
		return
	}
	cr.running = true
	cr.mu.Unlock()
	go func() {
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				fn()
			case <-cr.stop:
				return
			}
		}
	}()
}

// Stop stops runner
func (cr *CompactionRunner) Stop() {
	cr.mu.Lock()
	defer cr.mu.Unlock()
	if cr.running {
		close(cr.stop)
		cr.running = false
		cr.stop = make(chan struct{})
	}
}

// IsRunning returns running
func (cr *CompactionRunner) IsRunning() bool {
	cr.mu.Lock()
	defer cr.mu.Unlock()
	return cr.running
}
