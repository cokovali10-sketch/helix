// Package agent provides heartbeat
package agent

import (
	"sync"
	"time"
)

// Heartbeat manages heartbeat
type Heartbeat struct {
	interval time.Duration
	lastBeat time.Time
	mu       sync.Mutex
	stop     chan struct{}
}

// NewHeartbeat creates heartbeat
func NewHeartbeat(interval time.Duration) *Heartbeat {
	hb := &Heartbeat{interval: interval, stop: make(chan struct{})}
	go hb.run()
	return hb
}

func (hb *Heartbeat) run() {
	ticker := time.NewTicker(hb.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			hb.Beat()
		case <-hb.stop:
			return
		}
	}
}

// Beat records beat
func (hb *Heartbeat) Beat() {
	hb.mu.Lock()
	defer hb.mu.Unlock()
	hb.lastBeat = time.Now()
}

// LastBeat returns last beat
func (hb *Heartbeat) LastBeat() time.Time {
	hb.mu.Lock()
	defer hb.mu.Unlock()
	return hb.lastBeat
}

// IsAlive checks alive
func (hb *Heartbeat) IsAlive(timeout time.Duration) bool {
	hb.mu.Lock()
	defer hb.mu.Unlock()
	return time.Since(hb.lastBeat) < timeout
}

// Stop stops heartbeat
func (hb *Heartbeat) Stop() {
	close(hb.stop)
}
