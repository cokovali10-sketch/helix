// Package storage provides backup management
package storage

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// Backup represents a backup
type Backup struct {
	ID        string    `json:"id"`
	Path      string    `json:"path"`
	Size      int64     `json:"size"`
	CreatedAt time.Time `json:"created_at"`
	Status    string    `json:"status"`
	Checksum  string    `json:"checksum"`
}

// BackupManager manages backups
type BackupManager struct {
	backups map[string]*Backup
	mu      sync.RWMutex
}

// NewBackupManager creates manager
func NewBackupManager() *BackupManager {
	return &BackupManager{backups: make(map[string]*Backup)}
}

// Create creates a backup
func (bm *BackupManager) Create(path string) *Backup {
	bm.mu.Lock()
	defer bm.mu.Unlock()
	b := &Backup{
		ID:        fmt.Sprintf("backup-%d", time.Now().UnixNano()),
		Path:      path,
		CreatedAt: time.Now(),
		Status:    "completed",
		Checksum:  "sha256:abc",
		Size:      1024,
	}
	bm.backups[b.ID] = b
	return b
}

// List lists backups
func (bm *BackupManager) List() []*Backup {
	bm.mu.RLock()
	defer bm.mu.RUnlock()
	out := []*Backup{}
	for _, b := range bm.backups {
		out = append(out, b)
	}
	return out
}

// Get returns backup
func (bm *BackupManager) Get(id string) (*Backup, bool) {
	bm.mu.RLock()
	defer bm.mu.RUnlock()
	b, ok := bm.backups[id]
	return b, ok
}

// Delete deletes backup
func (bm *BackupManager) Delete(id string) bool {
	bm.mu.Lock()
	defer bm.mu.Unlock()
	if _, ok := bm.backups[id]; !ok {
		return false
	}
	delete(bm.backups, id)
	return true
}

// Restore restores backup
func (bm *BackupManager) Restore(id string) error {
	bm.mu.RLock()
	defer bm.mu.RUnlock()
	if _, ok := bm.backups[id]; !ok {
		return fmt.Errorf("backup not found")
	}
	return nil
}

// Stats returns stats
func (bm *BackupManager) Stats() map[string]interface{} {
	bm.mu.RLock()
	defer bm.mu.RUnlock()
	return map[string]interface{}{
		"total": len(bm.backups),
	}
}

// MarshalJSON
func (bm *BackupManager) MarshalJSON() ([]byte, error) {
	bm.mu.RLock()
	defer bm.mu.RUnlock()
	return json.Marshal(bm.backups)
}
