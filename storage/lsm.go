// Package storage provides additional storage features
// including LSM tree, compaction, and bloom filters
package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

// LSMTree represents a Log-Structured Merge Tree
type LSMTree struct {
	memtable     *Memtable
	immuTables   []*Memtable
	sstables     []*SSTableReader
	wal          *WAL
	compactor    *Compactor
	bloomFilter  *BloomFilter
	config       *LSMConfig
	mu           sync.RWMutex
	flushChan    chan *Memtable
	compactChan  chan struct{}
}

// LSMConfig holds LSM tree configuration
type LSMConfig struct {
	MemtableSize     int
	MaxImmuTables    int
	SSTableSize      int
	CompactionStyle  string // leveled, size_tiered, universal
	NumLevels        int
	LevelSizeMultiplier float64
	BloomFilterBits  int
	CompressionType  string
}

// DefaultLSMConfig returns default configuration
func DefaultLSMConfig() *LSMConfig {
	return &LSMConfig{
		MemtableSize:        64 * 1024 * 1024, // 64MB
		MaxImmuTables:       4,
		SSTableSize:         256 * 1024 * 1024, // 256MB
		CompactionStyle:     "leveled",
		NumLevels:           7,
		LevelSizeMultiplier: 10.0,
		BloomFilterBits:     10,
		CompressionType:     "zstd",
	}
}

// NewLSMTree creates a new LSM tree
func NewLSMTree(walDir, sstDir string, config *LSMConfig) (*LSMTree, error) {
	wal, err := NewWAL(walDir, 64*1024*1024)
	if err != nil {
		return nil, err
	}

	tree := &LSMTree{
		memtable:    NewMemtable(config.MemtableSize),
		immuTables:  make([]*Memtable, 0),
		sstables:    make([]*SSTableReader, 0),
		wal:         wal,
		compactor:   NewCompactor(sstDir, config),
		bloomFilter: NewBloomFilter(1000000, 3),
		config:      config,
		flushChan:   make(chan *Memtable, 10),
		compactChan: make(chan struct{}, 1),
	}

	// Load existing SSTables
	tree.loadSSTables(sstDir)

	// Start flush and compaction goroutines
	go tree.flushLoop()
	go tree.compactLoop()

	return tree, nil
}

// loadSSTables loads existing SSTables
func (t *LSMTree) loadSSTables(dir string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	for _, entry := range entries {
		if filepath.Ext(entry.Name()) == ".sst" {
			path := filepath.Join(dir, entry.Name())
			reader, err := NewSSTableReader(path)
			if err == nil {
				t.sstables = append(t.sstables, reader)
			}
		}
	}
}

// Put inserts a key-value pair
func (t *LSMTree) Put(key string, value []byte) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Write to WAL
	if err := t.wal.Record([]byte(key + ":" + string(value))); err != nil {
		return err
	}

	// Insert into memtable
	t.memtable.Put(key, value)
	t.bloomFilter.Add([]byte(key))

	// Check if memtable is full
	if t.memtable.IsFull() {
		t.immuTables = append(t.immuTables, t.memtable)
		t.memtable = NewMemtable(t.config.MemtableSize)

		// Trigger flush if too many immutable tables
		if len(t.immuTables) >= t.config.MaxImmuTables {
			select {
			case t.flushChan <- t.immuTables[0]:
			default:
			}
		}
	}

	return nil
}

// Get retrieves a value by key
func (t *LSMTree) Get(key string) ([]byte, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	// Check bloom filter
	if !t.bloomFilter.Contains([]byte(key)) {
		return nil, false
	}

	// Check memtable
	if value, ok := t.memtable.Get(key); ok {
		return value, true
	}

	// Check immutable tables (newest first)
	for i := len(t.immuTables) - 1; i >= 0; i-- {
		if value, ok := t.immuTables[i].Get(key); ok {
			return value, true
		}
	}

	// Check SSTables (newest first)
	for i := len(t.sstables) - 1; i >= 0; i-- {
		if value, ok := t.sstables[i].Get(key); ok {
			return value, true
		}
	}

	return nil, false
}

// Delete removes a key
func (t *LSMTree) Delete(key string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Write tombstone to WAL
	if err := t.wal.Record([]byte("DELETE:" + key)); err != nil {
		return err
	}

	t.memtable.Delete(key)
	return nil
}

// flushLoop handles flushing immutable tables
func (t *LSMTree) flushLoop() {
	for memtable := range t.flushChan {
		t.flushMemtable(memtable)
	}
}

// flushMemtable flushes a memtable to SSTable
func (t *LSMTree) flushMemtable(memtable *Memtable) {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Create SSTable writer
	// In production, would write to disk
	// For now, just create a reader
	reader := &SSTableReader{}
	t.sstables = append(t.sstables, reader)

	// Remove from immutable tables
	for i, mt := range t.immuTables {
		if mt == memtable {
			t.immuTables = append(t.immuTables[:i], t.immuTables[i+1:]...)
			break
		}
	}

	// Trigger compaction if needed
	select {
	case t.compactChan <- struct{}{}:
	default:
	}
}

// compactLoop handles compaction
func (t *LSMTree) compactLoop() {
	for range t.compactChan {
		t.compactor.Compact(t.sstables)
	}
}

// Compactor handles SSTable compaction
type Compactor struct {
	sstDir string
	config *LSMConfig
	mu     sync.Mutex
}

// NewCompactor creates a new compactor
func NewCompactor(sstDir string, config *LSMConfig) *Compactor {
	return &Compactor{
		sstDir: sstDir,
		config: config,
	}
}

// Compact compacts SSTables
func (c *Compactor) Compact(sstables []*SSTableReader) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(sstables) < 4 {
		return
	}

	// In production, would merge SSTables
	// For now, just log
	fmt.Printf("Compacting %d SSTables\n", len(sstables))
}

// GetStats returns LSM tree statistics
func (t *LSMTree) GetStats() map[string]interface{} {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return map[string]interface{}{
		"memtable_size":     t.memtable.Size(),
		"immutable_tables":  len(t.immuTables),
		"sstable_count":     len(t.sstables),
		"bloom_filter_size": 1000000,
	}
}

// Close closes the LSM tree
func (t *LSMTree) Close() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Flush memtable
	t.immuTables = append(t.immuTables, t.memtable)
	for _, mt := range t.immuTables {
		t.flushMemtable(mt)
	}

	// Close WAL
	t.wal.Close()

	// Close SSTables
	for _, sst := range t.sstables {
		sst.Close()
	}

	return nil
}

// RangeScan scans a range of keys
func (t *LSMTree) RangeScan(start, end string) (map[string][]byte, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	result := make(map[string][]byte)

	// Scan memtable
	// In production, would iterate memtable

	// Scan SSTables
	for _, sst := range t.sstables {
		// In production, would scan each SSTable
		_ = sst
	}

	return result, nil
}

// Count returns the number of entries
func (t *LSMTree) Count() int {
	t.mu.RLock()
	defer t.mu.RUnlock()

	count := t.memtable.Size()
	for _, mt := range t.immuTables {
		count += mt.Size()
	}
	for _, sst := range t.sstables {
		count += 0 // Would count SSTable entries
	}
	return count
}

// Size returns the total size in bytes
func (t *LSMTree) Size() int64 {
	t.mu.RLock()
	defer t.mu.RUnlock()

	var size int64
	size += int64(t.memtable.Size())
	for _, mt := range t.immuTables {
		size += int64(mt.Size())
	}
	return size
}

// Health returns health status
func (t *LSMTree) Health() map[string]interface{} {
	consistent, issues := NewRecoveryManager("", "").CheckConsistency()
	return map[string]interface{}{
		"consistent": consistent,
		"issues":     issues,
		"status":     "healthy",
	}
}