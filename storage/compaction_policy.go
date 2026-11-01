// Package storage provides compaction policy
package storage

import (
	"sync"
	"time"
)

// CompactionPolicy defines policy
type CompactionPolicy struct {
	Name          string
	MinFiles      int
	MaxFiles      int
	SizeThreshold int64
	Interval      time.Duration
}

// CompactionPolicyManager manages policies
type CompactionPolicyManager struct {
	policies map[string]*CompactionPolicy
	mu       sync.RWMutex
}

// NewCompactionPolicyManager creates manager
func NewCompactionPolicyManager() *CompactionPolicyManager {
	return &CompactionPolicyManager{policies: make(map[string]*CompactionPolicy)}
}

// Add adds policy
func (cpm *CompactionPolicyManager) Add(p *CompactionPolicy) {
	cpm.mu.Lock()
	defer cpm.mu.Unlock()
	cpm.policies[p.Name] = p
}

// Get returns policy
func (cpm *CompactionPolicyManager) Get(name string) (*CompactionPolicy, bool) {
	cpm.mu.RLock()
	defer cpm.mu.RUnlock()
	p, ok := cpm.policies[name]
	return p, ok
}

// Remove removes policy
func (cpm *CompactionPolicyManager) Remove(name string) {
	cpm.mu.Lock()
	defer cpm.mu.Unlock()
	delete(cpm.policies, name)
}

// List lists policies
func (cpm *CompactionPolicyManager) List() []*CompactionPolicy {
	cpm.mu.RLock()
	defer cpm.mu.RUnlock()
	out := []*CompactionPolicy{}
	for _, p := range cpm.policies {
		out = append(out, p)
	}
	return out
}

// ShouldCompact checks if should compact
func (cpm *CompactionPolicyManager) ShouldCompact(name string, fileCount int, totalSize int64) bool {
	cpm.mu.RLock()
	defer cpm.mu.RUnlock()
	p, ok := cpm.policies[name]
	if !ok {
		return false
	}
	if fileCount < p.MinFiles {
		return false
	}
	if totalSize < p.SizeThreshold {
		return false
	}
	return true
}

// Count returns count
func (cpm *CompactionPolicyManager) Count() int {
	cpm.mu.RLock()
	defer cpm.mu.RUnlock()
	return len(cpm.policies)
}
