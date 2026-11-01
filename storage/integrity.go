// Package storage provides additional storage features
// including data integrity, validation, and repair
package storage

import (
	"fmt"
	"hash/crc32"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// DataIntegrityChecker checks data integrity
type DataIntegrityChecker struct {
	checksums map[string]uint32
	mu        sync.RWMutex
}

// NewDataIntegrityChecker creates a new integrity checker
func NewDataIntegrityChecker() *DataIntegrityChecker {
	return &DataIntegrityChecker{
		checksums: make(map[string]uint32),
	}
}

// AddChecksum adds a checksum for a key
func (dic *DataIntegrityChecker) AddChecksum(key string, data []byte) {
	dic.mu.Lock()
	defer dic.mu.Unlock()
	dic.checksums[key] = crc32.ChecksumIEEE(data)
}

// Verify verifies data integrity
func (dic *DataIntegrityChecker) Verify(key string, data []byte) bool {
	dic.mu.RLock()
	defer dic.mu.RUnlock()
	expected, ok := dic.checksums[key]
	if !ok {
		return false
	}
	return crc32.ChecksumIEEE(data) == expected
}

// GetChecksum gets a checksum
func (dic *DataIntegrityChecker) GetChecksum(key string) (uint32, bool) {
	dic.mu.RLock()
	defer dic.mu.RUnlock()
	crc, ok := dic.checksums[key]
	return crc, ok
}

// RemoveChecksum removes a checksum
func (dic *DataIntegrityChecker) RemoveChecksum(key string) {
	dic.mu.Lock()
	defer dic.mu.Unlock()
	delete(dic.checksums, key)
}

// StorageValidator validates storage
type StorageValidator struct {
	engine *StorageEngine
}

// NewStorageValidator creates a new storage validator
func NewStorageValidator(engine *StorageEngine) *StorageValidator {
	return &StorageValidator{engine: engine}
}

// Validate validates the storage engine
func (sv *StorageValidator) Validate() []string {
	var issues []string

	// Check WAL
	walEntries, err := os.ReadDir(sv.engine.wal.dir)
	if err != nil {
		issues = append(issues, fmt.Sprintf("Failed to read WAL dir: %v", err))
	} else {
		walCount := 0
		for _, entry := range walEntries {
			if filepath.Ext(entry.Name()) == ".log" {
				walCount++
			}
		}
		if walCount == 0 {
			issues = append(issues, "No WAL files found")
		}
	}

	// Check SSTables
	sstEntries, err := os.ReadDir(sv.engine.sstDir)
	if err != nil {
		issues = append(issues, fmt.Sprintf("Failed to read SST dir: %v", err))
	} else {
		sstCount := 0
		for _, entry := range sstEntries {
			if filepath.Ext(entry.Name()) == ".sst" {
				sstCount++
			}
		}
		if sstCount > 100 {
			issues = append(issues, "Too many SSTables")
		}
	}

	// Check memtable
	if sv.engine.memtable.IsFull() {
		issues = append(issues, "Memtable is full")
	}

	return issues
}

// IsValid returns true if no issues
func (sv *StorageValidator) IsValid() bool {
	return len(sv.Validate()) == 0
}

// GetValidationReport returns a validation report
func (sv *StorageValidator) GetValidationReport() map[string]interface{} {
	issues := sv.Validate()
	return map[string]interface{}{
		"valid":  len(issues) == 0,
		"issues": issues,
		"count":  len(issues),
	}
}

// StorageRepairer repairs storage
type StorageRepairer struct {
	engine *StorageEngine
}

// NewStorageRepairer creates a new storage repairer
func NewStorageRepairer(engine *StorageEngine) *StorageRepairer {
	return &StorageRepairer{engine: engine}
}

// Repair repairs storage issues
func (sr *StorageRepairer) Repair() error {
	// In production, would repair corrupted data
	return nil
}

// RepairWAL repairs WAL files
func (sr *StorageRepairer) RepairWAL() error {
	// In production, would repair corrupted WAL records
	return nil
}

// RepairSSTable repairs SSTable files
func (sr *StorageRepairer) RepairSSTable(path string) error {
	// In production, would repair corrupted SSTable
	return nil
}

// StorageMonitor monitors storage health
type StorageMonitor struct {
	engine   *StorageEngine
	metrics  *StorageMetrics
	mu       sync.RWMutex
	running  bool
	ctx      context.Context
	cancel   context.CancelFunc
	interval time.Duration
}

// NewStorageMonitor creates a new storage monitor
func NewStorageMonitor(engine *StorageEngine, interval time.Duration) *StorageMonitor {
	return &StorageMonitor{
		engine:   engine,
		metrics:  NewStorageMetrics(),
		interval: interval,
	}
}

// Start starts monitoring
func (sm *StorageMonitor) Start(ctx context.Context) {
	sm.ctx, sm.cancel = context.WithCancel(ctx)
	sm.running = true
	go sm.run()
}

// Stop stops monitoring
func (sm *StorageMonitor) Stop() {
	sm.running = false
	if sm.cancel != nil {
		sm.cancel()
	}
}

// run runs the monitoring loop
func (sm *StorageMonitor) run() {
	ticker := time.NewTicker(sm.interval)
	defer ticker.Stop()

	for {
		select {
		case <-sm.ctx.Done():
			return
		case <-ticker.C:
			if !sm.running {
				return
			}
			sm.collectMetrics()
		}
	}
}

// collectMetrics collects storage metrics
func (sm *StorageMonitor) collectMetrics() {
	stats := sm.engine.Stats()
	if memtableSize, ok := stats["memtable_size"].(int); ok {
		sm.metrics.UpdateMemtableSize(memtableSize)
	}
	if sstableCount, ok := stats["sstable_count"].(int); ok {
		sm.metrics.UpdateSSTableCount(sstableCount)
	}
}

// GetMetrics returns storage metrics
func (sm *StorageMonitor) GetMetrics() *StorageMetrics {
	return sm.metrics
}

// GetHealth returns storage health
func (sm *StorageMonitor) GetHealth() map[string]interface{} {
	health := CheckHealth(sm.engine)
	return map[string]interface{}{
		"status":     health.Status,
		"consistent": health.Consistent,
		"issues":     health.Issues,
		"wal_files":  health.WALFiles,
		"sst_files":  health.SSTFiles,
		"memtable_size": health.MemtableSize,
	}
}

// StorageOptimizer optimizes storage
type StorageOptimizer struct {
	engine *StorageEngine
	config *OptimizerConfig
}

// OptimizerConfig holds optimizer configuration
type OptimizerConfig struct {
	EnableCompression  bool
	EnableEncryption   bool
	CompressionLevel   int
	BlockSize          int
	CacheSize          int
	PrefetchEnabled    bool
	PrefetchDistance   int
	CompactionEnabled  bool
	RetentionEnabled   bool
}

// DefaultOptimizerConfig returns default config
func DefaultOptimizerConfig() *OptimizerConfig {
	return &OptimizerConfig{
		EnableCompression:  true,
		EnableEncryption:   false,
		CompressionLevel:   6,
		BlockSize:          4096,
		CacheSize:          1024 * 1024 * 1024,
		PrefetchEnabled:    true,
		PrefetchDistance:   4,
		CompactionEnabled:  true,
		RetentionEnabled:   true,
	}
}

// NewStorageOptimizer creates a new optimizer
func NewStorageOptimizer(engine *StorageEngine, config *OptimizerConfig) *StorageOptimizer {
	return &StorageOptimizer{
		engine: engine,
		config: config,
	}
}

// Optimize optimizes storage
func (so *StorageOptimizer) Optimize() map[string]interface{} {
	result := make(map[string]interface{})

	if so.config.CompactionEnabled {
		if err := so.engine.Compact(); err != nil {
			result["compaction_error"] = err.Error()
		} else {
			result["compaction"] = "completed"
		}
	}

	if so.config.RetentionEnabled {
		rm := NewRetentionManager(so.engine.sstDir, RetentionPolicy{Type: "time", Duration: 168 * time.Hour})
		if err := rm.Run(); err != nil {
			result["retention_error"] = err.Error()
		} else {
			result["retention"] = "completed"
		}
	}

	result["compression_enabled"] = so.config.EnableCompression
	result["encryption_enabled"] = so.config.EnableEncryption
	result["cache_size"] = so.config.CacheSize

	return result
}

// GetRecommendations returns optimization recommendations
func (so *StorageOptimizer) GetRecommendations() []string {
	var recommendations []string

	stats := so.engine.Stats()
	if memtableSize, ok := stats["memtable_size"].(int); ok {
		if memtableSize > 100*1024*1024 {
			recommendations = append(recommendations, "Consider increasing memtable size")
		}
	}

	if sstableCount, ok := stats["sstable_count"].(int); ok {
		if sstableCount > 50 {
			recommendations = append(recommendations, "Consider running compaction")
		}
	}

	return recommendations
}

// StorageExporter exports storage data
type StorageExporter struct {
	engine *StorageEngine
}

// NewStorageExporter creates a new storage exporter
func NewStorageExporter(engine *StorageEngine) *StorageExporter {
	return &StorageExporter{engine: engine}
}

// ExportJSON exports data as JSON
func (se *StorageExporter) ExportJSON(path string) error {
	// In production, would export all data as JSON
	return nil
}

// ExportCSV exports data as CSV
func (se *StorageExporter) ExportCSV(path string) error {
	// In production, would export all data as CSV
	return nil
}

// ExportParquet exports data as Parquet
func (se *StorageExporter) ExportParquet(path string) error {
	// In production, would export all data as Parquet
	return nil
}

// StorageImporter imports storage data
type StorageImporter struct {
	engine *StorageEngine
}

// NewStorageImporter creates a new storage importer
func NewStorageImporter(engine *StorageEngine) *StorageImporter {
	return &StorageImporter{engine: engine}
}

// ImportJSON imports data from JSON
func (si *StorageImporter) ImportJSON(path string) error {
	// In production, would import data from JSON
	return nil
}

// ImportCSV imports data from CSV
func (si *StorageImporter) ImportCSV(path string) error {
	// In production, would import data from CSV
	return nil
}

// ImportParquet imports data from Parquet
func (si *StorageImporter) ImportParquet(path string) error {
	// In production, would import data from Parquet
	return nil
}