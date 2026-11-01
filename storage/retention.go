// Package storage provides retention management
package storage

import (
	"sync"
	"time"
)

// RetentionPolicy defines retention
type RetentionPolicy struct {
	Name     string
	Duration time.Duration
	MaxSize  int64
}

// RetentionManager manages retention
type RetentionManager struct {
	policies map[string]*RetentionPolicy
	mu       sync.RWMutex
}

// NewRetentionManager creates manager
func NewRetentionManager() *RetentionManager {
	return &RetentionManager{policies: make(map[string]*RetentionPolicy)}
}

// AddPolicy adds policy
func (rm *RetentionManager) AddPolicy(p *RetentionPolicy) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	rm.policies[p.Name] = p
}

// GetPolicy returns policy
func (rm *RetentionManager) GetPolicy(name string) (*RetentionPolicy, bool) {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	p, ok := rm.policies[name]
	return p, ok
}

// RemovePolicy removes policy
func (rm *RetentionManager) RemovePolicy(name string) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	delete(rm.policies, name)
}

// ListPolicies lists policies
func (rm *RetentionManager) ListPolicies() []*RetentionPolicy {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	out := []*RetentionPolicy{}
	for _, p := range rm.policies {
		out = append(out, p)
	}
	return out
}

// ShouldRetain checks if data should be retained
func (rm *RetentionManager) ShouldRetain(ts time.Time, policyName string) bool {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	p, ok := rm.policies[policyName]
	if !ok {
		return true
	}
	return time.Since(ts) < p.Duration
}

// Count returns count
func (rm *RetentionManager) Count() int {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	return len(rm.policies)
}
