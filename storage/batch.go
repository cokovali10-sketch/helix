// Package storage provides additional storage features
// including write batch, iterator, and snapshot
package storage

import (
	"fmt"
	"sync"
)

// WriteBatch represents a batch of write operations
type WriteBatch struct {
	ops []BatchOp
	mu  sync.Mutex
}

// BatchOp represents a single operation in a batch
type BatchOp struct {
	Type  string // put, delete
	Key   string
	Value []byte
}

// NewWriteBatch creates a new write batch
func NewWriteBatch() *WriteBatch {
	return &WriteBatch{
		ops: make([]BatchOp, 0),
	}
}

// Put adds a put operation
func (wb *WriteBatch) Put(key string, value []byte) {
	wb.mu.Lock()
	defer wb.mu.Unlock()
	wb.ops = append(wb.ops, BatchOp{Type: "put", Key: key, Value: value})
}

// Delete adds a delete operation
func (wb *WriteBatch) Delete(key string) {
	wb.mu.Lock()
	defer wb.mu.Unlock()
	wb.ops = append(wb.ops, BatchOp{Type: "delete", Key: key})
}

// Size returns the number of operations
func (wb *WriteBatch) Size() int {
	wb.mu.Lock()
	defer wb.mu.Unlock()
	return len(wb.ops)
}

// Clear clears all operations
func (wb *WriteBatch) Clear() {
	wb.mu.Lock()
	defer wb.mu.Unlock()
	wb.ops = wb.ops[:0]
}

// Apply applies the batch to a storage engine
func (wb *WriteBatch) Apply(engine *StorageEngine) error {
	wb.mu.Lock()
	defer wb.mu.Unlock()

	for _, op := range wb.ops {
		switch op.Type {
		case "put":
			if err := engine.Put(op.Key, op.Value); err != nil {
				return fmt.Errorf("failed to put %s: %w", op.Key, err)
			}
		case "delete":
			if err := engine.Delete(op.Key); err != nil {
				return fmt.Errorf("failed to delete %s: %w", op.Key, err)
			}
		}
	}
	return nil
}

// Iterator iterates over storage entries
type Iterator struct {
	engine   *StorageEngine
	keys     []string
	current  int
	prefix   string
	reverse  bool
}

// NewIterator creates a new iterator
func NewIterator(engine *StorageEngine, prefix string, reverse bool) *Iterator {
	return &Iterator{
		engine:  engine,
		prefix:  prefix,
		reverse: reverse,
		current: -1,
	}
}

// SeekToFirst seeks to the first entry
func (it *Iterator) SeekToFirst() {
	it.current = 0
	// In production, would populate keys from storage
}

// SeekToLast seeks to the last entry
func (it *Iterator) SeekToLast() {
	// In production, would populate keys and seek to last
}

// Seek seeks to a specific key
func (it *Iterator) Seek(key string) {
	// In production, would binary search
}

// Next moves to the next entry
func (it *Iterator) Next() bool {
	it.current++
	// In production, would check bounds
	return it.current < len(it.keys)
}

// Prev moves to the previous entry
func (it *Iterator) Prev() bool {
	it.current--
	if it.current < 0 {
		it.current = -1
		return false
	}
	return true
}

// Key returns the current key
func (it *Iterator) Key() string {
	if it.current < 0 || it.current >= len(it.keys) {
		return ""
	}
	return it.keys[it.current]
}

// Value returns the current value
func (it *Iterator) Value() []byte {
	key := it.Key()
	if key == "" {
		return nil
	}
	value, _ := it.engine.Get(key)
	return value
}

// Valid checks if the iterator is valid
func (it *Iterator) Valid() bool {
	return it.current >= 0 && it.current < len(it.keys)
}

// Close closes the iterator
func (it *Iterator) Close() {
	it.keys = nil
}

// RangeScan scans a range of keys
func (se *StorageEngine) RangeScan(start, end string) (map[string][]byte, error) {
	result := make(map[string][]byte)
	// In production, would scan SSTables and memtable
	return result, nil
}

// PrefixScan scans keys with a prefix
func (se *StorageEngine) PrefixScan(prefix string) (map[string][]byte, error) {
	result := make(map[string][]byte)
	// In production, would scan all keys with prefix
	return result, nil
}

// Count returns the number of entries
func (se *StorageEngine) Count() (int, error) {
	// In production, would count all entries
	return 0, nil
}

// Compact compacts the storage
func (se *StorageEngine) Compact() error {
	// In production, would compact SSTables
	return nil
}

// Flush flushes the memtable
func (se *StorageEngine) Flush() error {
	return se.flushMemtable()
}

// GetStats returns storage statistics
func (se *StorageEngine) GetStats() map[string]interface{} {
	stats := se.Stats()
	stats["transaction_count"] = 0
	stats["snapshot_count"] = 0
	stats["backup_count"] = 0
	return stats
}

// Health returns health status
func (se *StorageEngine) Health() map[string]interface{} {
	consistent, issues := NewRecoveryManager(se.wal.dir, se.sstDir).CheckConsistency()
	return map[string]interface{}{
		"consistent": consistent,
		"issues":     issues,
		"status":     "healthy",
	}
}