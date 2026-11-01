// Package storage provides snapshot management
package storage

import (
	"fmt"
	"sync"
	"time"
)

// SnapshotInfo represents snapshot info
type SnapshotInfo struct {
	ID        string    `json:"id"`
	Path      string    `json:"path"`
	Size      int64     `json:"size"`
	CreatedAt time.Time `json:"created_at"`
	Status    string    `json:"status"`
}

// SnapshotManager manages snapshots
type SnapshotManager struct {
	snapshots map[string]*SnapshotInfo
	mu        sync.RWMutex
}

// NewSnapshotManager creates manager
func NewSnapshotManager() *SnapshotManager {
	return &SnapshotManager{snapshots: make(map[string]*SnapshotInfo)}
}

// Create creates snapshot
func (sm *SnapshotManager) Create(path string) *SnapshotInfo {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	s := &SnapshotInfo{
		ID:        fmt.Sprintf("snap-%d", time.Now().UnixNano()),
		Path:      path,
		CreatedAt: time.Now(),
		Status:    "completed",
		Size:      2048,
	}
	sm.snapshots[s.ID] = s
	return s
}

// Get returns snapshot
func (sm *SnapshotManager) Get(id string) (*SnapshotInfo, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	s, ok := sm.snapshots[id]
	return s, ok
}

// List lists snapshots
func (sm *SnapshotManager) List() []*SnapshotInfo {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	out := []*SnapshotInfo{}
	for _, s := range sm.snapshots {
		out = append(out, s)
	}
	return out
}

// Delete deletes snapshot
func (sm *SnapshotManager) Delete(id string) bool {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if _, ok := sm.snapshots[id]; !ok {
		return false
	}
	delete(sm.snapshots, id)
	return true
}

// Count returns count
func (sm *SnapshotManager) Count() int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return len(sm.snapshots)
}
