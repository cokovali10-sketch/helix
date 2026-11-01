// Package storage provides additional storage features
// including advanced metrics and monitoring
package storage

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// AdvancedMetrics provides advanced storage metrics
type AdvancedMetrics struct {
	// Write metrics
	WriteCount     int64
	WriteBytes     int64
	WriteErrors    int64
	WriteLatency   time.Duration
	MaxWriteLatency time.Duration

	// Read metrics
	ReadCount      int64
	ReadBytes      int64
	ReadErrors     int64
	ReadLatency    time.Duration
	MaxReadLatency time.Duration
	CacheHits      int64
	CacheMisses    int64

	// Memtable metrics
	MemtableSize    int
	MemtableEntries int
	MemtableFlushes int64

	// SSTable metrics
	SSTableCount    int
	SSTableSize     int64
	Compactions     int64

	// WAL metrics
	WALSize         int64
	WALSegments     int
	WALSyncs        int64

	// Retention metrics
	RetentionFiles  int64
	RetentionBytes  int64

	// General
	Uptime          time.Duration
	StartTime       time.Time

	mu sync.RWMutex
}

// NewAdvancedMetrics creates new advanced metrics
func NewAdvancedMetrics() *AdvancedMetrics {
	return &AdvancedMetrics{
		StartTime: time.Now(),
	}
}

// RecordWrite records a write operation
func (am *AdvancedMetrics) RecordWrite(bytes int64, latency time.Duration) {
	am.mu.Lock()
	defer am.mu.Unlock()
	am.WriteCount++
	am.WriteBytes += bytes
	am.WriteLatency += latency
	if latency > am.MaxWriteLatency {
		am.MaxWriteLatency = latency
	}
}

// RecordWriteError records a write error
func (am *AdvancedMetrics) RecordWriteError() {
	am.mu.Lock()
	defer am.mu.Unlock()
	am.WriteErrors++
}

// RecordRead records a read operation
func (am *AdvancedMetrics) RecordRead(bytes int64, latency time.Duration) {
	am.mu.Lock()
	defer am.mu.Unlock()
	am.ReadCount++
	am.ReadBytes += bytes
	am.ReadLatency += latency
	if latency > am.MaxReadLatency {
		am.MaxReadLatency = latency
	}
}

// RecordReadError records a read error
func (am *AdvancedMetrics) RecordReadError() {
	am.mu.Lock()
	defer am.mu.Unlock()
	am.ReadErrors++
}

// RecordCacheHit records a cache hit
func (am *AdvancedMetrics) RecordCacheHit() {
	am.mu.Lock()
	defer am.mu.Unlock()
	am.CacheHits++
}

// RecordCacheMiss records a cache miss
func (am *AdvancedMetrics) RecordCacheMiss() {
	am.mu.Lock()
	defer am.mu.Unlock()
	am.CacheMisses++
}

// RecordMemtableFlush records a memtable flush
func (am *AdvancedMetrics) RecordMemtableFlush() {
	am.mu.Lock()
	defer am.mu.Unlock()
	am.MemtableFlushes++
}

// RecordCompaction records a compaction
func (am *AdvancedMetrics) RecordCompaction() {
	am.mu.Lock()
	defer am.mu.Unlock()
	am.Compactions++
}

// RecordRetention records retention cleanup
func (am *AdvancedMetrics) RecordRetention(files int64, bytes int64) {
	am.mu.Lock()
	defer am.mu.Unlock()
	am.RetentionFiles += files
	am.RetentionBytes += bytes
}

// UpdateMemtable updates memtable metrics
func (am *AdvancedMetrics) UpdateMemtable(size, entries int) {
	am.mu.Lock()
	defer am.mu.Unlock()
	am.MemtableSize = size
	am.MemtableEntries = entries
}

// UpdateSSTable updates SSTable metrics
func (am *AdvancedMetrics) UpdateSSTable(count int, size int64) {
	am.mu.Lock()
	defer am.mu.Unlock()
	am.SSTableCount = count
	am.SSTableSize = size
}

// UpdateWAL updates WAL metrics
func (am *AdvancedMetrics) UpdateWAL(size int64, segments int, syncs int64) {
	am.mu.Lock()
	defer am.mu.Unlock()
	am.WALSize = size
	am.WALSegments = segments
	am.WALSyncs = syncs
}

// UpdateUptime updates uptime
func (am *AdvancedMetrics) UpdateUptime() {
	am.mu.Lock()
	defer am.mu.Unlock()
	am.Uptime = time.Since(am.StartTime)
}

// GetStats returns metrics as a map
func (am *AdvancedMetrics) GetStats() map[string]interface{} {
	am.mu.RLock()
	defer am.mu.RUnlock()

	avgWriteLatency := time.Duration(0)
	if am.WriteCount > 0 {
		avgWriteLatency = am.WriteLatency / time.Duration(am.WriteCount)
	}

	avgReadLatency := time.Duration(0)
	if am.ReadCount > 0 {
		avgReadLatency = am.ReadLatency / time.Duration(am.ReadCount)
	}

	cacheHitRate := float64(0)
	totalCache := am.CacheHits + am.CacheMisses
	if totalCache > 0 {
		cacheHitRate = float64(am.CacheHits) / float64(totalCache)
	}

	return map[string]interface{}{
		"write_count":       am.WriteCount,
		"write_bytes":       am.WriteBytes,
		"write_errors":      am.WriteErrors,
		"avg_write_latency": avgWriteLatency.String(),
		"max_write_latency": am.MaxWriteLatency.String(),
		"read_count":        am.ReadCount,
		"read_bytes":        am.ReadBytes,
		"read_errors":       am.ReadErrors,
		"avg_read_latency":  avgReadLatency.String(),
		"max_read_latency":  am.MaxReadLatency.String(),
		"cache_hits":        am.CacheHits,
		"cache_misses":      am.CacheMisses,
		"cache_hit_rate":    cacheHitRate,
		"memtable_size":     am.MemtableSize,
		"memtable_entries":  am.MemtableEntries,
		"memtable_flushes":  am.MemtableFlushes,
		"sstable_count":     am.SSTableCount,
		"sstable_size":      am.SSTableSize,
		"compactions":       am.Compactions,
		"wal_size":          am.WALSize,
		"wal_segments":      am.WALSegments,
		"wal_syncs":         am.WALSyncs,
		"retention_files":    am.RetentionFiles,
		"retention_bytes":    am.RetentionBytes,
		"uptime":            am.Uptime.String(),
	}
}

// ToPrometheus converts metrics to Prometheus format
func (am *AdvancedMetrics) ToPrometheus() string {
	am.mu.RLock()
	defer am.mu.RUnlock()

	var output string

	output += fmt.Sprintf("# HELP helix_storage_writes_total Total writes\n")
	output += fmt.Sprintf("# TYPE helix_storage_writes_total counter\n")
	output += fmt.Sprintf("helix_storage_writes_total %d\n", am.WriteCount)

	output += fmt.Sprintf("# HELP helix_storage_write_bytes Total bytes written\n")
	output += fmt.Sprintf("# TYPE helix_storage_write_bytes counter\n")
	output += fmt.Sprintf("helix_storage_write_bytes %d\n", am.WriteBytes)

	output += fmt.Sprintf("# HELP helix_storage_write_errors Total write errors\n")
	output += fmt.Sprintf("# TYPE helix_storage_write_errors counter\n")
	output += fmt.Sprintf("helix_storage_write_errors %d\n", am.WriteErrors)

	output += fmt.Sprintf("# HELP helix_storage_reads_total Total reads\n")
	output += fmt.Sprintf("# TYPE helix_storage_reads_total counter\n")
	output += fmt.Sprintf("helix_storage_reads_total %d\n", am.ReadCount)

	output += fmt.Sprintf("# HELP helix_storage_read_bytes Total bytes read\n")
	output += fmt.Sprintf("# TYPE helix_storage_read_bytes counter\n")
	output += fmt.Sprintf("helix_storage_read_bytes %d\n", am.ReadBytes)

	output += fmt.Sprintf("# HELP helix_storage_read_errors Total read errors\n")
	output += fmt.Sprintf("# TYPE helix_storage_read_errors counter\n")
	output += fmt.Sprintf("helix_storage_read_errors %d\n", am.ReadErrors)

	cacheHitRate := float64(0)
	totalCache := am.CacheHits + am.CacheMisses
	if totalCache > 0 {
		cacheHitRate = float64(am.CacheHits) / float64(totalCache)
	}
	output += fmt.Sprintf("# HELP helix_storage_cache_hit_rate Cache hit rate\n")
	output += fmt.Sprintf("# TYPE helix_storage_cache_hit_rate gauge\n")
	output += fmt.Sprintf("helix_storage_cache_hit_rate %f\n", cacheHitRate)

	output += fmt.Sprintf("# HELP helix_storage_memtable_size Memtable size in bytes\n")
	output += fmt.Sprintf("# TYPE helix_storage_memtable_size gauge\n")
	output += fmt.Sprintf("helix_storage_memtable_size %d\n", am.MemtableSize)

	output += fmt.Sprintf("# HELP helix_storage_memtable_entries Memtable entries\n")
	output += fmt.Sprintf("# TYPE helix_storage_memtable_entries gauge\n")
	output += fmt.Sprintf("helix_storage_memtable_entries %d\n", am.MemtableEntries)

	output += fmt.Sprintf("# HELP helix_storage_sstable_count SSTable count\n")
	output += fmt.Sprintf("# TYPE helix_storage_sstable_count gauge\n")
	output += fmt.Sprintf("helix_storage_sstable_count %d\n", am.SSTableCount)

	output += fmt.Sprintf("# HELP helix_storage_sstable_size SSTable size in bytes\n")
	output += fmt.Sprintf("# TYPE helix_storage_sstable_size gauge\n")
	output += fmt.Sprintf("helix_storage_sstable_size %d\n", am.SSTableSize)

	output += fmt.Sprintf("# HELP helix_storage_compactions Total compactions\n")
	output += fmt.Sprintf("# TYPE helix_storage_compactions counter\n")
	output += fmt.Sprintf("helix_storage_compactions %d\n", am.Compactions)

	output += fmt.Sprintf("# HELP helix_storage_wal_size WAL size in bytes\n")
	output += fmt.Sprintf("# TYPE helix_storage_wal_size gauge\n")
	output += fmt.Sprintf("helix_storage_wal_size %d\n", am.WALSize)

	output += fmt.Sprintf("# HELP helix_storage_wal_segments WAL segments\n")
	output += fmt.Sprintf("# TYPE helix_storage_wal_segments gauge\n")
	output += fmt.Sprintf("helix_storage_wal_segments %d\n", am.WALSegments)

	output += fmt.Sprintf("# HELP helix_storage_uptime Uptime in seconds\n")
	output += fmt.Sprintf("# TYPE helix_storage_uptime gauge\n")
	output += fmt.Sprintf("helix_storage_uptime %f\n", am.Uptime.Seconds())

	return output
}

// StorageProfiler profiles storage performance
type StorageProfiler struct {
	samples   []ProfileSample
	mu        sync.Mutex
	running   bool
	ctx       context.Context
	cancel    context.CancelFunc
	interval  time.Duration
	maxSamples int
}

// ProfileSample represents a profile sample
type ProfileSample struct {
	Timestamp    time.Time
	WriteLatency time.Duration
	ReadLatency  time.Duration
	MemtableSize int
	SSTableCount int
	CacheHitRate float64
}

// NewStorageProfiler creates a new profiler
func NewStorageProfiler(interval time.Duration, maxSamples int) *StorageProfiler {
	return &StorageProfiler{
		samples:    make([]ProfileSample, 0, maxSamples),
		interval:   interval,
		maxSamples: maxSamples,
	}
}

// Start starts profiling
func (sp *StorageProfiler) Start(ctx context.Context) {
	sp.ctx, sp.cancel = context.WithCancel(ctx)
	sp.running = true
	go sp.run()
}

// Stop stops profiling
func (sp *StorageProfiler) Stop() {
	sp.running = false
	if sp.cancel != nil {
		sp.cancel()
	}
}

// run runs the profiling loop
func (sp *StorageProfiler) run() {
	ticker := time.NewTicker(sp.interval)
	defer ticker.Stop()

	for {
		select {
		case <-sp.ctx.Done():
			return
		case <-ticker.C:
			if !sp.running {
				return
			}
			sp.collectSample()
		}
	}
}

// collectSample collects a profile sample
func (sp *StorageProfiler) collectSample() {
	sp.mu.Lock()
	defer sp.mu.Unlock()

	sample := ProfileSample{
		Timestamp: time.Now(),
	}

	sp.samples = append(sp.samples, sample)
	if len(sp.samples) > sp.maxSamples {
		sp.samples = sp.samples[1:]
	}
}

// GetSamples returns collected samples
func (sp *StorageProfiler) GetSamples() []ProfileSample {
	sp.mu.Lock()
	defer sp.mu.Unlock()

	samples := make([]ProfileSample, len(sp.samples))
	copy(samples, sp.samples)
	return samples
}

// GetAverageWriteLatency returns average write latency
func (sp *StorageProfiler) GetAverageWriteLatency() time.Duration {
	sp.mu.Lock()
	defer sp.mu.Unlock()

	if len(sp.samples) == 0 {
		return 0
	}

	var total time.Duration
	for _, s := range sp.samples {
		total += s.WriteLatency
	}
	return total / time.Duration(len(sp.samples))
}

// GetAverageReadLatency returns average read latency
func (sp *StorageProfiler) GetAverageReadLatency() time.Duration {
	sp.mu.Lock()
	defer sp.mu.Unlock()

	if len(sp.samples) == 0 {
		return 0
	}

	var total time.Duration
	for _, s := range sp.samples {
		total += s.ReadLatency
	}
	return total / time.Duration(len(sp.samples))
}

// Export exports profile data
func (sp *StorageProfiler) Export(format string) ([]byte, error) {
	sp.mu.Lock()
	defer sp.mu.Unlock()

	switch format {
	case "json":
		return json.Marshal(sp.samples)
	case "csv":
		var result string
		result += "timestamp,write_latency,read_latency,memtable_size,sstable_count,cache_hit_rate\n"
		for _, s := range sp.samples {
			result += fmt.Sprintf("%d,%v,%v,%d,%d,%.4f\n",
				s.Timestamp.Unix(), s.WriteLatency, s.ReadLatency, s.MemtableSize, s.SSTableCount, s.CacheHitRate)
		}
		return []byte(result), nil
	default:
		return json.Marshal(sp.samples)
	}
}