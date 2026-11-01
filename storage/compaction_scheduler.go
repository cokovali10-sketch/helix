// Package storage provides compaction scheduler
package storage

import (
	"sync"
	"time"
)

// CompactionScheduler schedules compaction
type CompactionScheduler struct {
	interval time.Duration
	stop     chan struct{}
	mu       sync.Mutex
	running  bool
}

// NewCompactionScheduler creates scheduler
func NewCompactionScheduler(interval time.Duration) *CompactionScheduler {
	return &CompactionScheduler{interval: interval, stop: make(chan struct{})}
}

// Start starts scheduler
func (cs *CompactionScheduler) Start(fn func()) {
	cs.mu.Lock()
	if cs.running {
		cs.mu.Unlock()
		return
	}
	cs.running = true
	cs.mu.Unlock()
	go func() {
		ticker := time.NewTicker(cs.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				fn()
			case <-cs.stop:
				return
			}
		}
	}()
}

// Stop stops scheduler
func (cs *CompactionScheduler) Stop() {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	if cs.running {
		close(cs.stop)
		cs.running = false
	}
}

// IsRunning returns running
func (cs *CompactionScheduler) IsRunning() bool {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	return cs.running
}
