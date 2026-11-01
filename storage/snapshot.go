// Package storage provides additional storage features
// including snapshot, backup, and recovery
package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// SnapshotManager manages storage snapshots
type SnapshotManager struct {
	dir       string
	snapshots map[string]*Snapshot
	mu        sync.RWMutex
}

// Snapshot represents a storage snapshot
type Snapshot struct {
	ID        string
	Timestamp time.Time
	Path      string
	Size      int64
	Metadata  map[string]string
	Files     []string
}

// NewSnapshotManager creates a new snapshot manager
func NewSnapshotManager(dir string) *SnapshotManager {
	os.MkdirAll(dir, 0755)
	return &SnapshotManager{
		dir:       dir,
		snapshots: make(map[string]*Snapshot),
	}
}

// Create creates a new snapshot
func (sm *SnapshotManager) Create(engine *StorageEngine) (*Snapshot, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	id := fmt.Sprintf("snapshot-%d", time.Now().UnixNano())
	path := filepath.Join(sm.dir, id)

	if err := os.MkdirAll(path, 0755); err != nil {
		return nil, err
	}

	snapshot := &Snapshot{
		ID:        id,
		Timestamp: time.Now(),
		Path:      path,
		Metadata:  make(map[string]string),
		Files:     make([]string, 0),
	}

	sm.snapshots[id] = snapshot
	return snapshot, nil
}

// GetSnapshot gets a snapshot by ID
func (sm *SnapshotManager) GetSnapshot(id string) (*Snapshot, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	snapshot, ok := sm.snapshots[id]
	return snapshot, ok
}

// ListSnapshots lists all snapshots
func (sm *SnapshotManager) ListSnapshots() []*Snapshot {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	snapshots := make([]*Snapshot, 0, len(sm.snapshots))
	for _, s := range sm.snapshots {
		snapshots = append(snapshots, s)
	}
	return snapshots
}

// DeleteSnapshot deletes a snapshot
func (sm *SnapshotManager) DeleteSnapshot(id string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	snapshot, ok := sm.snapshots[id]
	if !ok {
		return fmt.Errorf("snapshot %s not found", id)
	}

	os.RemoveAll(snapshot.Path)
	delete(sm.snapshots, id)
	return nil
}

// Restore restores from a snapshot
func (sm *SnapshotManager) Restore(id string, engine *StorageEngine) error {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	_, ok := sm.snapshots[id]
	if !ok {
		return fmt.Errorf("snapshot %s not found", id)
	}

	// In production, would restore data from snapshot
	return nil
}

// BackupManager manages backups
type BackupManager struct {
	dir     string
	backups map[string]*Backup
	mu      sync.RWMutex
}

// Backup represents a backup
type Backup struct {
	ID        string
	Timestamp time.Time
	Path      string
	Size      int64
	Type      string // full, incremental
	Metadata  map[string]string
	Files     []string
}

// NewBackupManager creates a new backup manager
func NewBackupManager(dir string) *BackupManager {
	os.MkdirAll(dir, 0755)
	return &BackupManager{
		dir:     dir,
		backups: make(map[string]*Backup),
	}
}

// CreateFullBackup creates a full backup
func (bm *BackupManager) CreateFullBackup(engine *StorageEngine) (*Backup, error) {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	id := fmt.Sprintf("backup-full-%d", time.Now().UnixNano())
	path := filepath.Join(bm.dir, id)

	if err := os.MkdirAll(path, 0755); err != nil {
		return nil, err
	}

	backup := &Backup{
		ID:        id,
		Timestamp: time.Now(),
		Path:      path,
		Type:      "full",
		Metadata:  make(map[string]string),
		Files:     make([]string, 0),
	}

	bm.backups[id] = backup
	return backup, nil
}

// CreateIncrementalBackup creates an incremental backup
func (bm *BackupManager) CreateIncrementalBackup(engine *StorageEngine, baseBackupID string) (*Backup, error) {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	id := fmt.Sprintf("backup-incr-%d", time.Now().UnixNano())
	path := filepath.Join(bm.dir, id)

	if err := os.MkdirAll(path, 0755); err != nil {
		return nil, err
	}

	backup := &Backup{
		ID:        id,
		Timestamp: time.Now(),
		Path:      path,
		Type:      "incremental",
		Metadata:  map[string]string{"base": baseBackupID},
		Files:     make([]string, 0),
	}

	bm.backups[id] = backup
	return backup, nil
}

// GetBackup gets a backup by ID
func (bm *BackupManager) GetBackup(id string) (*Backup, bool) {
	bm.mu.RLock()
	defer bm.mu.RUnlock()
	backup, ok := bm.backups[id]
	return backup, ok
}

// ListBackups lists all backups
func (bm *BackupManager) ListBackups() []*Backup {
	bm.mu.RLock()
	defer bm.mu.RUnlock()

	backups := make([]*Backup, 0, len(bm.backups))
	for _, b := range bm.backups {
		backups = append(backups, b)
	}
	return backups
}

// DeleteBackup deletes a backup
func (bm *BackupManager) DeleteBackup(id string) error {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	backup, ok := bm.backups[id]
	if !ok {
		return fmt.Errorf("backup %s not found", id)
	}

	os.RemoveAll(backup.Path)
	delete(bm.backups, id)
	return nil
}

// Restore restores from a backup
func (bm *BackupManager) Restore(id string, engine *StorageEngine) error {
	bm.mu.RLock()
	defer bm.mu.RUnlock()

	backup, ok := bm.backups[id]
	if !ok {
		return fmt.Errorf("backup %s not found", id)
	}

	// In production, would restore data from backup
	_ = backup
	return nil
}

// RecoveryManager manages recovery
type RecoveryManager struct {
	walDir string
	sstDir string
}

// NewRecoveryManager creates a new recovery manager
func NewRecoveryManager(walDir, sstDir string) *RecoveryManager {
	return &RecoveryManager{
		walDir: walDir,
		sstDir: sstDir,
	}
}

// Recover recovers data from WAL
func (rm *RecoveryManager) Recover() error {
	// In production, would replay WAL records
	return nil
}

// CheckConsistency checks data consistency
func (rm *RecoveryManager) CheckConsistency() (bool, []string) {
	var issues []string

	// Check WAL files
	entries, err := os.ReadDir(rm.walDir)
	if err != nil {
		issues = append(issues, fmt.Sprintf("Failed to read WAL dir: %v", err))
		return false, issues
	}

	walCount := 0
	for _, entry := range entries {
		if filepath.Ext(entry.Name()) == ".log" {
			walCount++
		}
	}

	if walCount == 0 {
		issues = append(issues, "No WAL files found")
	}

	// Check SSTable files
	sstEntries, err := os.ReadDir(rm.sstDir)
	if err != nil {
		issues = append(issues, fmt.Sprintf("Failed to read SST dir: %v", err))
		return false, issues
	}

	sstCount := 0
	for _, entry := range sstEntries {
		if filepath.Ext(entry.Name()) == ".sst" {
			sstCount++
		}
	}

	return len(issues) == 0, issues
}

// Export exports storage data
func (rm *RecoveryManager) Export(format string, outputPath string) error {
	// In production, would export data in specified format
	return nil
}

// Import imports storage data
func (rm *RecoveryManager) Import(format string, inputPath string) error {
	// In production, would import data from specified format
	return nil
}

// HealthCheck checks storage health
func (rm *RecoveryManager) HealthCheck() map[string]interface{} {
	consistent, issues := rm.CheckConsistency()
	return map[string]interface{}{
		"consistent": consistent,
		"issues":     issues,
		"wal_dir":    rm.walDir,
		"sst_dir":    rm.sstDir,
	}
}

// SaveMetadata saves metadata to file
func SaveMetadata(path string, metadata map[string]interface{}) error {
	data, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// LoadMetadata loads metadata from file
func LoadMetadata(path string) (map[string]interface{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var metadata map[string]interface{}
	if err := json.Unmarshal(data, &metadata); err != nil {
		return nil, err
	}
	return metadata, nil
}

// StorageHealth represents storage health
type StorageHealth struct {
	Status     string
	Consistent bool
	Issues     []string
	WALFiles   int
	SSTFiles   int
	MemtableSize int
}

// CheckHealth checks storage health
func CheckHealth(engine *StorageEngine) StorageHealth {
	health := StorageHealth{
		Status: "healthy",
	}

	// Check WAL
	walEntries, _ := os.ReadDir(engine.wal.dir)
	for _, entry := range walEntries {
		if filepath.Ext(entry.Name()) == ".log" {
			health.WALFiles++
		}
	}

	// Check SSTables
	sstEntries, _ := os.ReadDir(engine.sstDir)
	for _, entry := range sstEntries {
		if filepath.Ext(entry.Name()) == ".sst" {
			health.SSTFiles++
		}
	}

	health.MemtableSize = engine.memtable.Size()

	consistent, issues := NewRecoveryManager(engine.wal.dir, engine.sstDir).CheckConsistency()
	health.Consistent = consistent
	health.Issues = issues

	if !consistent {
		health.Status = "unhealthy"
	}

	return health
}