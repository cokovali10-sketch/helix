// Package storage provides additional storage engine features
// including compression, retention, and compaction
package storage

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Compressor interface
type Compressor interface {
	Compress(data []byte) ([]byte, error)
	Decompress(data []byte) ([]byte, error)
	Name() string
}

// GzipCompressor implements gzip compression
type GzipCompressor struct {
	level int
}

func NewGzipCompressor(level int) *GzipCompressor {
	return &GzipCompressor{level: level}
}

func (g *GzipCompressor) Compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer, err := gzip.NewWriterLevel(&buf, g.level)
	if err != nil {
		return nil, err
	}
	if _, err := writer.Write(data); err != nil {
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (g *GzipCompressor) Decompress(data []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	return io.ReadAll(reader)
}

func (g *GzipCompressor) Name() string {
	return "gzip"
}

// ZlibCompressor implements zlib compression
type ZlibCompressor struct {
	level int
}

func NewZlibCompressor(level int) *ZlibCompressor {
	return &ZlibCompressor{level: level}
}

func (z *ZlibCompressor) Compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer, err := zlib.NewWriterLevel(&buf, z.level)
	if err != nil {
		return nil, err
	}
	if _, err := writer.Write(data); err != nil {
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (z *ZlibCompressor) Decompress(data []byte) ([]byte, error) {
	reader, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	return io.ReadAll(reader)
}

func (z *ZlibCompressor) Name() string {
	return "zlib"
}

// NoOpCompressor is a no-op compressor
type NoOpCompressor struct{}

func (n *NoOpCompressor) Compress(data []byte) ([]byte, error) {
	return data, nil
}

func (n *NoOpCompressor) Decompress(data []byte) ([]byte, error) {
	return data, nil
}

func (n *NoOpCompressor) Name() string {
	return "none"
}

// RetentionManager manages data retention
type RetentionManager struct {
	policy   RetentionPolicy
	sstDir   string
	stats    RetentionStats
	mu       sync.RWMutex
}

// RetentionPolicy defines retention rules
type RetentionPolicy struct {
	Type     string        // time, size, count
	Duration time.Duration // for time-based
	MaxBytes int64         // for size-based
	MaxFiles int           // for count-based
}

// RetentionStats tracks retention statistics
type RetentionStats struct {
	FilesScanned  int
	FilesDeleted  int
	BytesFreed    int64
	OldestRetained time.Time
	NewestRetained time.Time
}

// NewRetentionManager creates a new retention manager
func NewRetentionManager(sstDir string, policy RetentionPolicy) *RetentionManager {
	return &RetentionManager{
		policy: policy,
		sstDir: sstDir,
		stats:  RetentionStats{},
	}
}

// Run runs retention policy
func (rm *RetentionManager) Run() error {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	files, err := os.ReadDir(rm.sstDir)
	if err != nil {
		return err
	}

	var sstFiles []string
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".sst" {
			sstFiles = append(sstFiles, f.Name())
		}
	}

	rm.stats.FilesScanned = len(sstFiles)

	switch rm.policy.Type {
	case "time":
		rm.runTimeBased(sstFiles)
	case "size":
		rm.runSizeBased(sstFiles)
	case "count":
		rm.runCountBased(sstFiles)
	}

	return nil
}

func (rm *RetentionManager) runTimeBased(files []string) {
	cutoff := time.Now().Add(-rm.policy.Duration)
	for _, f := range files {
		info, err := os.Stat(filepath.Join(rm.sstDir, f))
		if err != nil {
			continue
		}
		if info.ModTime().Before(cutoff) {
			os.Remove(filepath.Join(rm.sstDir, f))
			rm.stats.FilesDeleted++
			rm.stats.BytesFreed += info.Size()
		}
	}
}

func (rm *RetentionManager) runSizeBased(files []string) {
	var totalSize int64
	for _, f := range files {
		info, err := os.Stat(filepath.Join(rm.sstDir, f))
		if err != nil {
			continue
		}
		totalSize += info.Size()
	}

	if totalSize <= rm.policy.MaxBytes {
		return
	}

	// Delete oldest files first
	type fileInfo struct {
		name    string
		size    int64
		modTime time.Time
	}
	var infos []fileInfo
	for _, f := range files {
		info, err := os.Stat(filepath.Join(rm.sstDir, f))
		if err != nil {
			continue
		}
		infos = append(infos, fileInfo{f, info.Size(), info.ModTime()})
	}

	// Sort by mod time (oldest first)
	// Simple bubble sort
	for i := 0; i < len(infos); i++ {
		for j := i + 1; j < len(infos); j++ {
			if infos[j].modTime.Before(infos[i].modTime) {
				infos[i], infos[j] = infos[j], infos[i]
			}
		}
	}

	freed := int64(0)
	for _, info := range infos {
		if totalSize-freed <= rm.policy.MaxBytes {
			break
		}
		os.Remove(filepath.Join(rm.sstDir, info.name))
		rm.stats.FilesDeleted++
		rm.stats.BytesFreed += info.size
		freed += info.size
	}
}

func (rm *RetentionManager) runCountBased(files []string) {
	if len(files) <= rm.policy.MaxFiles {
		return
	}

	// Sort by mod time
	type fileInfo struct {
		name    string
		modTime time.Time
	}
	var infos []fileInfo
	for _, f := range files {
		info, err := os.Stat(filepath.Join(rm.sstDir, f))
		if err != nil {
			continue
		}
		infos = append(infos, fileInfo{f, info.ModTime()})
	}

	for i := 0; i < len(infos); i++ {
		for j := i + 1; j < len(infos); j++ {
			if infos[j].modTime.Before(infos[i].modTime) {
				infos[i], infos[j] = infos[j], infos[i]
			}
		}
	}

	toDelete := len(infos) - rm.policy.MaxFiles
	for i := 0; i < toDelete; i++ {
		os.Remove(filepath.Join(rm.sstDir, infos[i].name))
		rm.stats.FilesDeleted++
	}
}

// GetStats returns retention stats
func (rm *RetentionManager) GetStats() RetentionStats {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	return rm.stats
}

// CompactionManager manages SSTable compaction
type CompactionManager struct {
	sstDir     string
	strategy   string
	maxFiles   int
	mu         sync.Mutex
}

// NewCompactionManager creates a new compaction manager
func NewCompactionManager(sstDir, strategy string, maxFiles int) *CompactionManager {
	return &CompactionManager{
		sstDir:   sstDir,
		strategy: strategy,
		maxFiles: maxFiles,
	}
}

// Compact compacts SSTables
func (cm *CompactionManager) Compact() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	files, err := os.ReadDir(cm.sstDir)
	if err != nil {
		return err
	}

	var sstFiles []string
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".sst" {
			sstFiles = append(sstFiles, f.Name())
		}
	}

	if len(sstFiles) < cm.maxFiles {
		return nil
	}

	// In production, would merge SSTables here
	// For now, just log
	fmt.Printf("Compacting %d SSTables in %s\n", len(sstFiles), cm.sstDir)

	return nil
}

// GetSSTableCount returns the number of SSTables
func (cm *CompactionManager) GetSSTableCount() int {
	files, err := os.ReadDir(cm.sstDir)
	if err != nil {
		return 0
	}

	count := 0
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".sst" {
			count++
		}
	}
	return count
}

// StorageMetrics tracks storage metrics
type StorageMetrics struct {
	WritesTotal    int64
	WritesBytes    int64
	ReadsTotal     int64
	ReadsBytes     int64
	CacheHits      int64
	CacheMisses    int64
	MemtableSize   int
	SSTableCount   int
	WALSize        int64
	mu             sync.RWMutex
}

// NewStorageMetrics creates new storage metrics
func NewStorageMetrics() *StorageMetrics {
	return &StorageMetrics{}
}

// RecordWrite records a write operation
func (sm *StorageMetrics) RecordWrite(bytes int64) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.WritesTotal++
	sm.WritesBytes += bytes
}

// RecordRead records a read operation
func (sm *StorageMetrics) RecordRead(bytes int64) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.ReadsTotal++
	sm.ReadsBytes += bytes
}

// RecordCacheHit records a cache hit
func (sm *StorageMetrics) RecordCacheHit() {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.CacheHits++
}

// RecordCacheMiss records a cache miss
func (sm *StorageMetrics) RecordCacheMiss() {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.CacheMisses++
}

// UpdateMemtableSize updates memtable size
func (sm *StorageMetrics) UpdateMemtableSize(size int) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.MemtableSize = size
}

// UpdateSSTableCount updates SSTable count
func (sm *StorageMetrics) UpdateSSTableCount(count int) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.SSTableCount = count
}

// UpdateWALSize updates WAL size
func (sm *StorageMetrics) UpdateWALSize(size int64) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.WALSize = size
}

// GetStats returns metrics as a map
func (sm *StorageMetrics) GetStats() map[string]interface{} {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return map[string]interface{}{
		"writes_total":     sm.WritesTotal,
		"writes_bytes":     sm.WritesBytes,
		"reads_total":      sm.ReadsTotal,
		"reads_bytes":      sm.ReadsBytes,
		"cache_hits":       sm.CacheHits,
		"cache_misses":     sm.CacheMisses,
		"memtable_size":    sm.MemtableSize,
		"sstable_count":    sm.SSTableCount,
		"wal_size":         sm.WALSize,
	}
}

// CacheHitRate returns the cache hit rate
func (sm *StorageMetrics) CacheHitRate() float64 {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	total := sm.CacheHits + sm.CacheMisses
	if total == 0 {
		return 0
	}
	return float64(sm.CacheHits) / float64(total)
}