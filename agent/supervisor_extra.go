// Package agent provides supervisor extra
package agent

import (
	"sync"
	"time"
)

// SupervisorExtra supervises components
type SupervisorExtra struct {
	components map[string]bool
	mu         sync.RWMutex
}

// NewSupervisorExtra creates supervisor
func NewSupervisorExtra() *SupervisorExtra {
	return &SupervisorExtra{components: make(map[string]bool)}
}

// Register registers component
func (s *SupervisorExtra) Register(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.components[name] = true
}

// Unregister unregisters
func (s *SupervisorExtra) Unregister(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.components, name)
}

// IsHealthy checks health
func (s *SupervisorExtra) IsHealthy() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, healthy := range s.components {
		if !healthy {
			return false
		}
	}
	return true
}

// SetHealth sets health
func (s *SupervisorExtra) SetHealth(name string, healthy bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.components[name] = healthy
}

// Count returns count
func (s *SupervisorExtra) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.components)
}

// Uptime returns uptime
func (s *SupervisorExtra) Uptime(start time.Time) time.Duration {
	return time.Since(start)
}
